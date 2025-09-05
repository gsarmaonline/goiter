package sms

import (
	"os"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type (
	SmsRequest struct {
		FromPhoneNumber string
		ToPhoneNumber   string
		Message         string
	}
)

func NewSMS(fromPhoneNumber string) *SmsRequest {
	return &SmsRequest{
		FromPhoneNumber: fromPhoneNumber,
	}
}

func SendSms(smsReq *SmsRequest) (err error) {

	client := twilio.NewRestClient()

	params := &api.CreateMessageParams{}
	params.SetBody(smsReq.Message)
	params.SetFrom(os.Getenv("SMS_FROM_PHONE_NUMBER"))
	params.SetTo(smsReq.ToPhoneNumber)

	_, err = client.Api.CreateMessage(params)
	if err != nil {
		return
	}

	return
}
