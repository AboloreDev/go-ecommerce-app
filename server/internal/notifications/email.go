package notifications

import (
	"fmt"
	"net"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SimpleEmail struct {
	Body     string
	Subject  string
	To string
}

type EmailNotifier struct {
	cfg *SMTPConfig
}

func NewEmailNotifier(cfg *SMTPConfig) *EmailNotifier {
	return &EmailNotifier{cfg: cfg}
}

func (e *EmailNotifier) SendEmail(email *SimpleEmail) error {
	addr := fmt.Sprintf("%d:%s", e.cfg.Host, e.cfg.Port)

	// Connect directly without TLS
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.cfg.Host)
	if err != nil {
		return err
	}
	defer client.Close()
	defer client.Quit()

	if e.cfg.Username != "" || e.cfg.Password != "" {
		auth := smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)
		err := client.Auth(auth)
		if err != nil {
			return err
		}
	}

	// Set Sender
	err = client.Mail(e.cfg.From)
	if err != nil {
		return err
	}

	// set recipient
	err = client.Rcpt(email.To)
	if err != nil {
		return nil
	}

	// Send message
	w, err := client.Data()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(
		"From: %s\r\n To: %s\r\n Subject: %s\r\n\r\n%s",
		e.cfg.From, email.To, email.Subject, email.Body,
	)

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	return w.Close()
}

func (e *EmailNotifier) SendOrderNotificationMail(userEmail, userName string, userOrder interface{}) error {
	email := &SimpleEmail{
		To: userEmail,
		Subject:  "Order Email Notification",
		Body: fmt.Sprintf(`Hello %s,

	You have successfully places an order.

	Its getting ready and it will be shipped in few days`,
		userOrder,

	`Best regards,
	The Shop Team`, userName),
	}

	return e.SendEmail(email)
}
