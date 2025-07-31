package usecase

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"time"

	"github.com/buildyow/byow-user-service/constants"
	"github.com/buildyow/byow-user-service/domain/entity"
	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/buildyow/byow-user-service/domain/repository"
	"github.com/buildyow/byow-user-service/dto"
	"github.com/buildyow/byow-user-service/infrastructure/jwt"
	"github.com/buildyow/byow-user-service/infrastructure/mailer"
	"github.com/buildyow/byow-user-service/infrastructure/validation"
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
		return appErrors.ErrEmailAlreadyExists
	}
	_, errPhoneNumber := u.Repo.FindByPhone(phone)
	if errPhoneNumber == nil {
		return appErrors.ErrPhoneAlreadyExists
	}
	return nil
}

func (u *UserUsecase) UpdateUserValidation(email string) error {
	_, errEmail := u.Repo.FindByEmail(email)
	if errEmail != nil {
		return appErrors.ErrUserNotFound
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
		return dto.UserResponse{}, appErrors.ErrUserNotFound
	}
	if !user.Verified {
		return dto.UserResponse{}, appErrors.ErrUserNotVerified
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return dto.UserResponse{}, appErrors.ErrInvalidCredentials
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
		return dto.UserResponse{}, appErrors.ErrUserNotFound
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
	// Generate secure random OTP
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}
	otp := strconv.Itoa(int(n.Int64()) + 100000)
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
		return appErrors.ErrUserNotFound
	}
	if time.Now().After(user.OTPExpiresAt) {
		return appErrors.ErrExpiredOTP
	}

	decryptedOTP, err := utils.Decrypt(user.OTP)
	if err != nil || decryptedOTP != otp {
		return appErrors.ErrInvalidOTP
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
	// Validate password strength first
	if valid, message := validation.ValidatePassword(req.Password); !valid {
		return appErrors.NewValidationError(message)
	}

	user, err := u.Repo.FindByEmail(req.Email)
	if err != nil {
		return appErrors.ErrUserNotFound
	}
	if time.Now().After(user.OTPExpiresAt) {
		return appErrors.ErrExpiredOTP
	}

	decryptedOTP, err := utils.Decrypt(user.OTP)
	if err != nil || decryptedOTP != req.OTP {
		return appErrors.ErrInvalidOTP
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return appErrors.NewInternalError("Failed to hash password")
	}
	
	user.Password = string(hashed)
	user.OTP = ""
	user.OTPExpiresAt = time.Time{}
	user.OTPType = ""

	return u.Repo.Update(user)
}

func (u *UserUsecase) ChangePasswordWithOldPassword(email string, req dto.ChangePasswordWithOldPasswordRequest) error {
	// Validate new password strength first
	if valid, message := validation.ValidatePassword(req.NewPassword); !valid {
		return appErrors.NewValidationError(message)
	}

	user, err := u.Repo.FindByEmail(email)
	if err != nil {
		return appErrors.ErrUserNotFound
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)) != nil {
		return appErrors.ErrInvalidOldPassword
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return appErrors.NewInternalError("Failed to hash password")
	}
	
	user.Password = string(hashed)

	return u.Repo.Update(user)
}

func (u *UserUsecase) UpdateUser(req dto.RegisterRequest) (*entity.User, error) {
	user, err := u.Repo.FindByEmail(req.Email)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}
	if req.AvatarUrl == "" {
		req.AvatarUrl = user.AvatarUrl
	}
	utils.LogWarn("Updating user with email:", req.Email, "and fullname:", req.Fullname)
	
	// Update existing user object to preserve all fields including CreatedAt
	user.Fullname = req.Fullname
	user.AvatarUrl = req.AvatarUrl
	
	err = u.Repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) UpdateUserByEmail(req dto.ChangeEmailRequest, oldEmail string) error {
	userOldEmail, err := u.Repo.FindByEmail(oldEmail)
	if err != nil {
		return appErrors.ErrUserNotFound
	}
	decryptedOTP, err := utils.Decrypt(userOldEmail.OTP)
	if err != nil || decryptedOTP != req.OTP {
		return appErrors.ErrInvalidOTP
	}
	if time.Now().After(userOldEmail.OTPExpiresAt) {
		return appErrors.ErrExpiredOTP
	}

	_, err = u.Repo.FindByEmail(req.NewEmail)
	if err == nil {
		return appErrors.ErrEmailAlreadyExists
	}
	
	// Update existing user object to preserve all fields including CreatedAt
	userOldEmail.Email = req.NewEmail
	userOldEmail.OTP = ""
	userOldEmail.OTPExpiresAt = time.Time{}
	userOldEmail.OTPType = ""
	
	err = u.Repo.UpdateEmail(userOldEmail, oldEmail)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecase) UpdateUserByPhone(req dto.ChangePhoneRequest, oldPhone string) error {
	userOldPhone, err := u.Repo.FindByPhone(oldPhone)
	if err != nil {
		return appErrors.ErrUserNotFound
	}
	decryptedOTP, err := utils.Decrypt(userOldPhone.OTP)
	if err != nil || decryptedOTP != req.OTP {
		return appErrors.ErrInvalidOTP
	}
	if time.Now().After(userOldPhone.OTPExpiresAt) {
		return appErrors.ErrExpiredOTP
	}

	_, err = u.Repo.FindByPhone(req.NewPhone)
	if err == nil {
		return appErrors.ErrPhoneAlreadyExists
	}
	
	// Update existing user object to preserve all fields including CreatedAt
	userOldPhone.PhoneNumber = req.NewPhone
	userOldPhone.OTP = ""
	userOldPhone.OTPExpiresAt = time.Time{}
	userOldPhone.OTPType = ""
	
	err = u.Repo.UpdatePhone(userOldPhone, oldPhone)
	if err != nil {
		return err
	}
	return nil
}
