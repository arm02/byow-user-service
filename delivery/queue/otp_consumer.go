package queue

import (
	"encoding/json"
	"log"

	"github.com/buildyow/byow-user-service/usecase"
	"github.com/rabbitmq/amqp091-go"
)

type OTPConsumer struct {
	userUsercase usecase.UserUsecase
}

func NewOTPConsumer(userUseCase usecase.UserUsecase) *OTPConsumer {
	return &OTPConsumer{userUsercase: userUseCase}
}

func (c *OTPConsumer) ConsumeOTP(deliveries <-chan amqp091.Delivery) {
	for msg := range deliveries {
		type OTPMessage struct {
			OTPType string `json:"otpType"`
			Email   string `json:"email"`
		}
		var otpMsg OTPMessage
		if err := json.Unmarshal(msg.Body, &otpMsg); err != nil {
			log.Println("Invalid OTP message:", err)
			continue
		}
		log.Println("Processing OTP for", otpMsg.Email)
		_ = c.userUsercase.SendOTP(otpMsg.OTPType, otpMsg.Email)
	}
}
