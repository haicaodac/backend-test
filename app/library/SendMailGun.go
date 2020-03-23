package library

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

var yourDomain string = "email.hanyny.com" // e.g. mg.yourcompany.com
var privateAPIKey string = "c3338c2334c7c41a4eb4e1f7b3525649-de7062c6-5850118f"

// SendMailGun ...
func SendMailGun(title string, body string, email string) error {

	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)
	m := mg.NewMessage(
		"Hanyny <hotro@hanyny.com>",
		title,
		"",
		email,
	)
	m.SetHtml(body)
	// m.AddAttachment("files/test.jpg")
	// m.AddAttachment("files/test.txt")
	m.SetReplyTo("hotro@hanyny.com")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	_, _, err := mg.Send(ctx, m)
	if err != nil {
		return err
	}
	return nil
}
