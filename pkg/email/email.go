package email

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer struct {
	Client *sendgrid.Client
}

func NewClient(apiKey string) *Mailer {
	return &Mailer{
		Client: sendgrid.NewSendClient(apiKey),
	}
}

func (mailer *Mailer) sendEmail(subject, toName, toEmailAddress, plainTextContent, htmlContent string) error {
	from := mail.NewEmail("Abanoub CEO", "logybyvy@lyft.live")
	toInfo := &mail.Email{
		Name:    toName,
		Address: toEmailAddress,
	}
	message := mail.NewSingleEmail(from, subject, toInfo, plainTextContent, htmlContent)
	_, err := mailer.Client.Send(message)
	if err != nil {
		log.Println("Error While sending emails", err)
		return err
	}

	return nil
}

func (mailer *Mailer) SendWelcomEmail(userName, emailAddress string) error {
	return mailer.sendEmail("welcome to our wonderful app", userName, emailAddress, "Hello our user please visit our site https://www.rescounts.com", "<h1>You are welcome here!</h1>")
}
