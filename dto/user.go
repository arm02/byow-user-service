package dto

type LoginRequest struct {
	Email    string `json:"email" example:"arm.adrian02@gmail.com"`
	Password string `json:"password" example:"masukaja123"`
}

type RegisterRequest struct {
	Fullname    string `json:"full_name" example:"John Doe"`
	Email       string `json:"email" example:"john@example.com"`
	Password    string `json:"password" example:"supersecret"`
	PhoneNumber string `json:"phone_number" example:"628112123123"`
	AvatarUrl   string `json:"avatar_url"`
}

type UserResponse struct {
	Fullname    string `json:"full_name" example:"John Doe"`
	Email       string `json:"email" example:"john@example.com"`
	PhoneNumber string `json:"phone_number" example:"628112123123"`
	AvatarUrl   string `json:"avatar_url" example:"https://assets/images/img.jpg"`
	Verified    bool   `json:"verified" example:"false"`
	OnBoarded   bool   `json:"on_boarded" example:"false"`
	Token       string `json:"token,omitempty" example:"token"`
	CreatedAt   string `json:"created_at,omitempty" example:"2024-01-15T10:30:00Z"`
}

type UserResponseSwagger struct {
	Status string       `json:"status" example:"SUCCESS"`
	Code   int          `json:"code" example:"200"`
	Data   UserResponse `json:"data"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" example:"john@example.com"`
	OTP   string `json:"otp" example:"000000"`
}

type ChangePasswordRequest struct {
	Email    string `json:"email" example:"john@example.com"`
	OTP      string `json:"otp" example:"000000"`
	Password string `json:"password" example:"newpassword"`
}

type ChangePasswordWithOldPasswordRequest struct {
	OldPassword string `json:"old_password" example:"oldpassword"`
	NewPassword string `json:"new_password" example:"newpassword"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email" example:"john.doe@example.com"`
	OTP      string `json:"otp" example:"000000"`
}

type ChangePhoneRequest struct {
	NewPhone string `json:"new_phone" example:"628112123123"`
	OTP      string `json:"otp" example:"000000"`
}
