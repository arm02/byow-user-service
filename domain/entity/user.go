package entity

import "time"

type User struct {
	ID           string    `bson:"_id,omitempty"`
	Fullname     string    `bson:"full_name"`
	Email        string    `bson:"email"`
	Password     string    `bson:"password"`
	PhoneNumber  string    `bson:"phone_number"`
	AvatarUrl    string    `bson:"avatar_url"`
	OnBoarded    bool      `bson:"on_boarded"`
	OTP          string    `bson:"otp,omitempty"`
	OTPType      string    `bson:"otp_type,omitempty"`
	OTPExpiresAt time.Time `bson:"otp_expires_at,omitempty"`
	Verified     bool      `bson:"verified"`
	CreatedAt    time.Time `bson:"created_at"`
}
