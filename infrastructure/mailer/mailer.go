package mailer

import (
	"fmt"

	"github.com/buildyow/byow-user-service/constants"
	"gopkg.in/gomail.v2"
)

func SendOTP(email, otp, host, user, pass string, port int, otpType string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your OTP Code")
	m.SetBody("text/plain", fmt.Sprintf("Your OTP for %s is: %s expired in %d minutes", otpType, otp, getOTPLifetime(otpType)))

	d := gomail.NewDialer(host, port, user, pass)
	return d.DialAndSend(m)
}

func getOTPLifetime(otpType string) int {
	switch otpType {
	case constants.FORGOT_PASSWORD, constants.EMAIL_CHANGED, constants.PHONE_CHANGED:
		return 10 // 10 minutes for forgot password, email changed, and phone changed
	case constants.VERIFICATION:
		return 5 // 5 minutes for verification
	default:
		return 1 // default to 1 minutes
	}
}
