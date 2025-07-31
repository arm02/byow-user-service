package http

import (
	"net/http"

	"github.com/buildyow/byow-user-service/constants"
	"github.com/buildyow/byow-user-service/dto"
	"github.com/buildyow/byow-user-service/lib"
	"github.com/buildyow/byow-user-service/response"
	"github.com/buildyow/byow-user-service/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Usecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: uc}
}

// @Summary Register user
// @Description Register a new user with avatar
// @Tags Authentication
// @Accept multipart/form-data
// @Produce json
// @Param full_name formData string true "Full name" example(John Doe)
// @Param email formData string true "Email" example(john@example.com)
// @Param password formData string true "Password" example(supersecret)
// @Param phone_number formData string true "Phone number" example(628112123123)
// @Param avatar formData file false "Avatar file"
// @Success 201 {object} dto.UserResponseSwagger
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	// Bind form values to struct
	req.Fullname = c.PostForm("full_name")
	req.Email = c.PostForm("email")
	req.Password = c.PostForm("password")
	req.PhoneNumber = c.PostForm("phone_number")

	err := h.Usecase.RegistrationValidation(c.PostForm("email"), c.PostForm("phone_number"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		response.Error(c, http.StatusBadRequest, constants.FAILED_PARSE_MULTIPART)
		return
	}

	// Upload File
	file, _, err := c.Request.FormFile("avatar")
	if err == nil {
		avatarURL, err := lib.CloudinaryUpload(file)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		req.AvatarUrl = avatarURL
	}

	// Call to usecase or saving to DB
	user, err := h.Usecase.Register(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, dto.UserResponse{
		Fullname:    user.Fullname,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		AvatarUrl:   user.AvatarUrl,
		Verified:    user.Verified,
		OnBoarded:   user.OnBoarded,
	})
}

// @Summary Login user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body dto.LoginRequest true "Email & Password"
// @Success 201 {object} dto.UserResponseSwagger
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.Usecase.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Set cookie
	c.SetCookie("token", user.Token, 3600, "/", "", false, true)

	response.Success(c, http.StatusOK, dto.UserResponse{
		Fullname:    user.Fullname,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		AvatarUrl:   user.AvatarUrl,
		Verified:    user.Verified,
		OnBoarded:   user.OnBoarded,
		Token:       user.Token,
	})
}

// @Summary Logout user
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	response.Success(c, http.StatusOK, constants.LOGOUT_SUCCESSFUL)
}

// @Summary Send OTP Verification
// @Tags Verification
// @Produce plain
// @Param email query string true "Email address"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /verification/users/send-otp [get]
func (h *UserHandler) SendOTPVerification(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.Error(c, http.StatusBadRequest, constants.EMAIL_REQUIRED)
		return
	}
	err := h.Usecase.SendOTP(constants.VERIFICATION, email)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, constants.OTP_SENT)
}

// @Summary Verify OTP
// @Tags Verification
// @Accept json
// @Produce plain
// @Param otp body dto.VerifyOTPRequest true "Email & OTP""
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /verification/users/verify-otp [post]
func (h *UserHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	email := req.Email
	otp := req.OTP

	if email == "" || otp == "" {
		response.Error(c, http.StatusBadRequest, constants.EMAIL_OTP_REQUIRED)
		return
	}

	err := h.Usecase.VerifyOTP(email, otp)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, constants.OTP_VERIFIED)
}

// @Summary Check Logged Account
// @Tags Users
// @Description Check if user is logged in and return user info
// @Produce plain
// @Success 200 {object} dto.UserResponseSwagger
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/me [get]
func (h *UserHandler) UserMe(c *gin.Context) {
	email, _ := c.Get("email")
	userID, _ := c.Get("user_id")
	phone, _ := c.Get("phone")
	response.Success(c, http.StatusOK, gin.H{
		"message": constants.VALID_TOKEN,
		"user": map[string]interface{}{
			"user_id": userID,
			"email":   email,
			"phone":   phone,
		},
	})
}

// @Summary Onboarded User
// @Tags Users
// @Description Onboard user to the system
// @Produce plain
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/onboard [get]
func (h *UserHandler) OnBoard(c *gin.Context) {
	emailIface, _ := c.Get("email")
	email, ok := emailIface.(string)
	if !ok {
		response.Error(c, http.StatusBadRequest, emailIface)
		return
	}
	err := h.Usecase.OnBoard(email)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, constants.ONBOARD_SUCCESSFUL)
}

// @Summary Change Password With OTP
// @Tags Authentication
// @Description Change user password using OTP verification
// @Produce plain
// @Param otp body dto.ChangePasswordRequest true "Email, OTP & New Password""
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/users/change-password-otp [post]
func (h *UserHandler) ChangePasswordWithOTP(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	email := req.Email
	otp := req.OTP

	if email == "" || otp == "" {
		response.Error(c, http.StatusBadRequest, constants.EMAIL_OTP_REQUIRED)
		return
	}

	err := h.Usecase.ChangePasswordWithOTP(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, constants.PASSWORD_CHANGED_SUCCESS)
}

// @Summary Send OTP Forgot Password
// @Tags Authentication
// @Produce plain
// @Param email query string true "Email address"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/users/forgot-password/send-otp [get]
func (h *UserHandler) SendOTPForgotPassword(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.Error(c, http.StatusBadRequest, constants.EMAIL_REQUIRED)
		return
	}
	err := h.Usecase.SendOTP(constants.FORGOT_PASSWORD, email)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, constants.OTP_SENT)
}

// @Summary Update User
// @Description Update user information
// @Tags Users
// @Accept json
// @Produce json
// @Param full_name formData string true "Full name" example(John Doe)
// @Param email formData string true "Email" example(john@example.com)
// @Param avatar formData file false "Avatar file"
// @Success 201 {object} dto.UserResponseSwagger
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/update [post]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req dto.RegisterRequest
	// Bind form values to struct
	req.Fullname = c.PostForm("full_name")
	req.Email = c.PostForm("email")
	req.Password = c.PostForm("password")
	req.PhoneNumber = c.PostForm("phone_number")

	err := h.Usecase.UpdateUserValidation(c.PostForm("email"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		response.Error(c, http.StatusBadRequest, constants.FAILED_PARSE_MULTIPART)
		return
	}

	// Upload File
	file, _, err := c.Request.FormFile("avatar")
	if err == nil {
		avatarURL, err := lib.CloudinaryUpload(file)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		req.AvatarUrl = avatarURL
	}

	// Call to usecase or saving to DB
	user, err := h.Usecase.UpdateUser(req)
	if req.AvatarUrl == "" {
		req.AvatarUrl = user.AvatarUrl
	}
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, dto.UserResponse{
		Fullname:    req.Fullname,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		AvatarUrl:   req.AvatarUrl,
		OnBoarded:   user.OnBoarded,
		Verified:    user.Verified,
	})
}

// @Summary Change Email With OTP
// @Tags Users
// @Description Change user email using OTP verification
// @Produce plain
// @Param otp body dto.ChangeEmailRequest true "OTP & New Email"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/change-email [post]
func (h *UserHandler) ChangeEmail(c *gin.Context) {
	var req dto.ChangeEmailRequest
	oldEmail, _ := c.Get("email")
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.OTP == "" || req.NewEmail == "" {
		response.Error(c, http.StatusBadRequest, constants.EMAIL_OTP_REQUIRED)
		return
	}
	err := h.Usecase.UpdateUserByEmail(req, oldEmail.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	c.SetCookie("token", "", -1, "/", "", false, true) // REMOVE OLD TOKEN
	newLogged, err := h.Usecase.LoginWithoutPassword(req.NewEmail)
	c.SetCookie("token", newLogged.Token, 3600, "/", "", false, true) // SET NEW TOKEN
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{
		"message": constants.EMAIL_CHANGED_SUCCESS,
		"data":    newLogged,
	})
}

// @Summary Send OTP Change Email
// @Tags Users
// @Produce plain
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/change-email/send-otp [get]
func (h *UserHandler) SendOTPEmailChange(c *gin.Context) {
	oldEmail, _ := c.Get("email")
	if oldEmail == "" {
		response.Error(c, http.StatusBadRequest, constants.EMAIL_REQUIRED)
		return
	}
	err := h.Usecase.SendOTP(constants.EMAIL_CHANGED, oldEmail.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, constants.OTP_SENT)
}

// @Summary Change Phone With OTP Email
// @Tags Users
// @Description Change user phone using OTP verification
// @Produce plain
// @Param otp body dto.ChangePhoneRequest true "OTP & New Email"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/change-phone [post]
func (h *UserHandler) ChangePhone(c *gin.Context) {
	oldPhone, _ := c.Get("phone")
	email, _ := c.Get("email")
	if oldPhone == "" {
		response.Error(c, http.StatusBadRequest, constants.PHONE_REQUIRED)
		return
	}
	var req dto.ChangePhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.OTP == "" || req.NewPhone == "" {
		response.Error(c, http.StatusBadRequest, constants.EMAIL_OTP_REQUIRED)
		return
	}
	err := h.Usecase.UpdateUserByPhone(req, oldPhone.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	c.SetCookie("token", "", -1, "/", "", false, true) // REMOVE OLD TOKEN
	newLogged, err := h.Usecase.LoginWithoutPassword(email.(string))
	c.SetCookie("token", newLogged.Token, 3600, "/", "", false, true) // SET NEW TOKEN
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{
		"message": constants.PHONE_CHANGED_SUCCESS,
		"data":    newLogged,
	})
}

// @Summary Send OTP Change Email
// @Tags Users
// @Produce plain
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/change-phone/send-otp [get]
func (h *UserHandler) SendOTPPhoneChange(c *gin.Context) {
	oldEmail, _ := c.Get("email")
	if oldEmail == "" {
		response.Error(c, http.StatusBadRequest, constants.EMAIL_REQUIRED)
		return
	}
	err := h.Usecase.SendOTP(constants.PHONE_CHANGED, oldEmail.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, constants.OTP_SENT)
}

// @Summary Change Password With Old Password
// @Tags Users
// @Description Change user password using old password
// @Produce plain
// @Param otp body dto.ChangePasswordWithOldPasswordRequest true "Email, Old Password & New Password"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/change-password-old [post]
func (h *UserHandler) ChangePasswordWithOldPassword(c *gin.Context) {
	email, _ := c.Get("email")
	if email == "" {
		response.Error(c, http.StatusBadRequest, constants.EMAIL_REQUIRED)
		return
	}
	var req dto.ChangePasswordWithOldPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		response.Error(c, http.StatusBadRequest, constants.ALL_FIELD_REQUIRED)
		return
	}

	err := h.Usecase.ChangePasswordWithOldPassword(email.(string), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, constants.PASSWORD_CHANGED_SUCCESS)
}
