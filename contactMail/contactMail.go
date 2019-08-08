package contactMail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

// MAIL_HOST := os.Getenv("MAIL_HOST")
// MAIL_USERNAME := os.Getenv("MAIL_USERNAME") // TODO: When project over
// MAIL_PASSWORD := os.Getenv("MAIL_PASSWORD")

var (
	host             = "smtp.gmail.com"
	MAIL_USERNAME    = "basketcourtcontact@gmail.com"
	MAIL_PASSWORD    = "udbnbexsqaltunpr"
	port             = "587"
	recipientAddress = []string{"farkhaniyoussef@gmail.com"}
)

type Contact struct {
	Name    string
	Subject string
	Email   string
	Message string
}

func SendMail(contact Contact) error {
	auth := smtp.PlainAuth("", MAIL_USERNAME, MAIL_PASSWORD, host)

	t, err := template.ParseFiles("contactMail/email-template.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer

	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("To: %s\r\n"+"Subject: %s\n%s\n\n", recipientAddress[0], contact.Subject, headers)))

	t.Execute(&body, contact)
	// fmt.Println(body.String())
	err = smtp.SendMail(host+":"+port, auth, MAIL_USERNAME, recipientAddress, body.Bytes())
	if err != nil {
		err = emailErr{
			description: "Mail couldn't be sent",
		}
		return err
	}
	return nil
}

type emailErr struct {
	description string
}

func (err emailErr) Error() string {
	return err.description
}
