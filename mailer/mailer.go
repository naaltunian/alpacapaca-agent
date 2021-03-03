package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

var (
	from     string
	pass     string
	to       string
	smtpHost string
	smtpPort string
)

func init() {
	from = os.Getenv("PACAEMAIL")
	pass = os.Getenv("PACAEMAILPASS")
	to = os.Getenv("RECEIVINGEMAIL")
	smtpHost = os.Getenv("SMTPHOST")
	smtpPort = os.Getenv("SMTPPORT")
}

// Notify emails agent owners of any issues
func Notify(msg string) {
	// for reference: https://www.loginradius.com/blog/async/sending-emails-with-golang/

	// Message
	message := []byte("To: test@test.com\r\n" +
		"Subject: Urgent Trading Notification!\r\n" +
		"\r\n" +
		msg,
	)
	receivingEmail := []string{to}

	// Authentication.
	auth := smtp.PlainAuth("", from, pass, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, receivingEmail, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
