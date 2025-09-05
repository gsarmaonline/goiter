package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type (
	MailerRequest struct {
		From        string
		To          []string
		Subject     string
		PlainText   string
		HtmlContent string
	}
)

func (s *MailerRequest) SendEmail() (err error) {
	s.From = os.Getenv("SENDGRID_FROM_EMAIL")
	message := mail.NewSingleEmail(mail.NewEmail(s.From, s.From), s.Subject, mail.NewEmail(s.From, s.From), s.PlainText, s.HtmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	if _, err = client.Send(message); err != nil {
		return
	}
	return
}
