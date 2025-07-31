package constants

const (
	//RESPONSE
	SUCCESS = "SUCCESS"
	ERROR   = "ERROR"

	// Success Messages (still used in responses)
	LOGOUT_SUCCESSFUL        = "LOGOUT_SUCCESSFUL"
	ONBOARD_SUCCESSFUL       = "ONBOARD_SUCCESSFUL"
	PASSWORD_CHANGED_SUCCESS = "PASSWORD_CHANGED_SUCCESS"
	EMAIL_CHANGED_SUCCESS    = "EMAIL_CHANGED_SUCCESS"
	PHONE_CHANGED_SUCCESS    = "PHONE_CHANGED_SUCCESS"
	OTP_VERIFIED             = "OTP_VERIFIED"
	OTP_SENT                 = "OTP_SENT"
	VALID_TOKEN              = "VALID_TOKEN"

	// Default values
	DefaultPageSize = 20

	// OTP Types (still used for email sending)
	FORGOT_PASSWORD  = "forgot_password"
	VERIFICATION     = "verification"
	EMAIL_CHANGED    = "email_changed"
	PASSWORD_CHANGED = "password_changed"
	PHONE_CHANGED    = "phone_changed"
)
