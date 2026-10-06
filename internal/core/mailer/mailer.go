package core_mailer

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type SMTPMailer struct {
	config Config
}

func NewSMTPMailer(config Config) *SMTPMailer {
	return &SMTPMailer{
		config: config,
	}
}

func (m *SMTPMailer) Send(ctx context.Context, to, subject, body string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	addr := net.JoinHostPort(m.config.Host, strconv.Itoa(m.config.Port))
	auth := smtp.PlainAuth("", m.config.User, m.config.Password, m.config.Host)

	msg := strings.Join([]string{
		"From: " + m.config.From,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"UTF-8\"",
		"",
		body,
	}, "\r\n")

	if err := smtp.SendMail(addr, auth, m.config.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("send email to '%s': %w", to, err)
	}

	return nil
}