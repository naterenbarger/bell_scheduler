package services

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewEmailService(host string, port int, username, password, from string) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *EmailService) SendPasswordResetEmail(to, token string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		You have requested to reset your password.
		Please click the following link to reset your password:
		
		%s/reset-password?token=%s
		
		If you did not request this, please ignore this email.
		`, "http://localhost:3000", token)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", s.from, to, subject, body)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
