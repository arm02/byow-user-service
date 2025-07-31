package usecase

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/buildyow/byow-user-service/constants"
	"github.com/buildyow/byow-user-service/domain/entity"
	"github.com/buildyow/byow-user-service/domain/repository"
	"github.com/buildyow/byow-user-service/dto"
	"github.com/buildyow/byow-user-service/infrastructure/jwt"
	"github.com/buildyow/byow-user-service/infrastructure/mailer"
	"github.com/buildyow/byow-user-service/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repo        repository.UserRepository
	JWTSecret   string
	JWTExpire   int
	EmailConfig struct {
		Host string
		Port int
		User string
		Pass string
	}
}

func (u *UserUsecase) RegistrationValidation(email string, phone string) error {
	_, errEmail := u.Repo.FindByEmail(email)
	if errEmail == nil {
		return errors.New(constants.EMAIL_ALREADY_REGISTERED)
	}
	_, errPhoneNumber := u.Repo.FindByPhone(phone)
	if errPhoneNumber == nil {
		return errors.New(constants.PHONE_ALREADY_REGISTERED)
	}
	return nil
}

func (u *UserUsecase) UpdateUserValidation(email string) error {
	_, errEmail := u.Repo.FindByEmail(email)
	if errEmail != nil {
		return errors.New(constants.ERR_NOT_FOUND)
	}
	return nil
}

func (u *UserUsecase) Register(req dto.RegisterRequest) (*entity.User, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	user := &entity.User{
		Fullname:    req.Fullname,
		Email:       req.Email,
		Password:    string(hashed),
		PhoneNumber: req.PhoneNumber,
		AvatarUrl:   req.AvatarUrl,
		Verified:    false,
		OnBoarded:   false,
	}
	err := u.Repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) Login(email, password string) (dto.UserResponse, error) {
	user, err := u.Repo.FindByEmail(email)
	if err != nil {
		return dto.UserResponse{}, errors.New(constants.ERR_NOT_FOUND)
	}
	if !user.Verified {
		return dto.UserResponse{}, errors.New(constants.USER_NOT_VERIFIED)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return dto.UserResponse{}, errors.New(constants.INVALID_CREDENTIALS)
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.PhoneNumber, u.JWTSecret, u.JWTExpire)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{
		Fullname:    user.Fullname,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		AvatarUrl:   user.AvatarUrl,
		Verified:    user.Verified,
		OnBoarded:   user.OnBoarded,
		Token:       token,
	}, nil
}

func (u *UserUsecase) LoginWithoutPassword(email string) (dto.UserResponse, error) {
	user, err := u.Repo.FindByEmail(email)
	if err != nil {
		return dto.UserResponse{}, errors.New(constants.ERR_NOT_FOUND)
	}
	// Generate token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.PhoneNumber, u.JWTSecret, u.JWTExpire)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{
		Fullname:    user.Fullname,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		AvatarUrl:   user.AvatarUrl,
		Verified:    user.Verified,
		OnBoarded:   user.OnBoarded,
		Token:       token,
	}, nil
}

func (u *UserUsecase) SendOTP(otpType, email string) error {
	user, err := u.Repo.FindByEmail(email)
	if err != nil {
		return err
	}
	otp := strconv.Itoa(rand.Intn(900000) + 100000)
	encryptedOTP, err := utils.Encrypt(otp)
	if err != nil {
		return err
	}
	user.OTP = encryptedOTP
	user.OTPType = otpType
	if otpType == constants.VERIFICATION {
		user.OTPExpiresAt = time.Now().Add(5 * time.Minute)
	}
	if otpType == constants.FORGOT_PASSWORD || otpType == constants.EMAIL_CHANGED || otpType == constants.PHONE_CHANGED {
		user.OTPExpiresAt = time.Now().Add(10 * time.Minute)
	}

	if err := u.Repo.Update(user); err != nil {
		return err
	}
	return mailer.SendOTP(email, otp, u.EmailConfig.Host, u.EmailConfig.User, u.EmailConfig.Pass, u.EmailConfig.Port, otpType)
}

func (u *UserUsecase) VerifyOTP(email, otp string) error {
	user, err := u.Repo.FindByEmail(email)
	if err != nil {
		return errors.New(constants.ERR_FETCH_FAILED)
	}
	if time.Now().After(user.OTPExpiresAt) {
		return errors.New(constants.OTP_EXPIRED)
	}

	decryptedOTP, err := utils.Decrypt(user.OTP)
	if err != nil || decryptedOTP != otp {
		return errors.New(constants.OTP_INVALID)
	}

	user.Verified = true
	user.OTP = ""
	user.OTPExpiresAt = time.Time{}
	user.OTPType = ""

	return u.Repo.Update(user)
}

func (u *UserUsecase) OnBoard(email string) error {
	user, err := u.Repo.FindByEmail(email)
	if err != nil {
		return err
	}
	user.OnBoarded = true
	if err := u.Repo.Update(user); err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) ChangePasswordWithOTP(req dto.ChangePasswordRequest) error {
	user, err := u.Repo.FindByEmail(req.Email)
	if err != nil {
		return errors.New(constants.ERR_FETCH_FAILED)
	}
	if time.Now().After(user.OTPExpiresAt) {
		return errors.New(constants.OTP_EXPIRED)
	}

	decryptedOTP, err := utils.Decrypt(user.OTP)
	if err != nil || decryptedOTP != req.OTP {
		return errors.New(constants.OTP_INVALID)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	user.Password = string(hashed)
	user.OTP = ""
	user.OTPExpiresAt = time.Time{}
	user.OTPType = ""

	return u.Repo.Update(user)
}

func (u *UserUsecase) ChangePasswordWithOldPassword(email string, req dto.ChangePasswordWithOldPasswordRequest) error {
	user, err := u.Repo.FindByEmail(email)
	if err != nil {
		return errors.New(constants.ERR_FETCH_FAILED)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)) != nil {
		return errors.New(constants.INVALID_ORLD_PASSWORD)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	user.Password = string(hashed)

	return u.Repo.Update(user)
}

func (u *UserUsecase) UpdateUser(req dto.RegisterRequest) (*entity.User, error) {
	user, err := u.Repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New(constants.ERR_FETCH_FAILED)
	}
	if req.AvatarUrl == "" {
		req.AvatarUrl = user.AvatarUrl
	}
	utils.LogWarn("Updating user with email:", req.Email, "and fullname:", req.Fullname)
	err = u.Repo.Update(&entity.User{
		Email:       user.Email,
		Fullname:    req.Fullname,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
		AvatarUrl:   req.AvatarUrl,
		OnBoarded:   user.OnBoarded,
		Verified:    user.Verified,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) UpdateUserByEmail(req dto.ChangeEmailRequest, oldEmail string) error {
	userOldEmail, err := u.Repo.FindByEmail(oldEmail)
	if err != nil {
		return errors.New(constants.ERR_FETCH_FAILED)
	}
	decryptedOTP, err := utils.Decrypt(userOldEmail.OTP)
	if err != nil || decryptedOTP != req.OTP {
		return errors.New(constants.OTP_INVALID)
	}
	if time.Now().After(userOldEmail.OTPExpiresAt) {
		return errors.New(constants.OTP_EXPIRED)
	}

	_, err = u.Repo.FindByEmail(req.NewEmail)
	if err == nil {
		return errors.New(constants.EMAIL_ALREADY_REGISTERED)
	}
	err = u.Repo.UpdateEmail(&entity.User{
		Email:       req.NewEmail,
		Fullname:    userOldEmail.Fullname,
		PhoneNumber: userOldEmail.PhoneNumber,
		Password:    userOldEmail.Password,
		AvatarUrl:   userOldEmail.AvatarUrl,
		OnBoarded:   userOldEmail.OnBoarded,
		Verified:    userOldEmail.Verified,
	}, oldEmail)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) UpdateUserByPhone(req dto.ChangePhoneRequest, oldPhone string) error {
	userOldPhone, err := u.Repo.FindByPhone(oldPhone)
	if err != nil {
		return errors.New(constants.ERR_FETCH_FAILED)
	}
	decryptedOTP, err := utils.Decrypt(userOldPhone.OTP)
	if err != nil || decryptedOTP != req.OTP {
		return errors.New(constants.OTP_INVALID)
	}
	if time.Now().After(userOldPhone.OTPExpiresAt) {
		return errors.New(constants.OTP_EXPIRED)
	}

	_, err = u.Repo.FindByPhone(req.NewPhone)
	if err == nil {
		return errors.New(constants.PHONE_ALREADY_REGISTERED)
	}
	err = u.Repo.UpdatePhone(&entity.User{
		Email:       userOldPhone.Email,
		Fullname:    userOldPhone.Fullname,
		PhoneNumber: req.NewPhone,
		Password:    userOldPhone.Password,
		AvatarUrl:   userOldPhone.AvatarUrl,
		OnBoarded:   userOldPhone.OnBoarded,
		Verified:    userOldPhone.Verified,
	}, oldPhone)
	if err != nil {
		return err
	}
	return nil
}
