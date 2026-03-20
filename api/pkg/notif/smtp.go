package notif

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTPProvider struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

func NewSMTPProvider(host, port, user, password, from string) *SMTPProvider {
	return &SMTPProvider{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
	}
}

func (p *SMTPProvider) Kirim(_ context.Context, tujuan, subjek, pesan string) error {
	if p.host == "" {
		return nil
	}

	auth := smtp.PlainAuth("", p.user, p.password, p.host)

	header := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		p.from, tujuan, subjek, pesan,
	)

	addr := p.host + ":" + p.port
	err := smtp.SendMail(addr, auth, p.from, []string{tujuan}, []byte(header))
	if err != nil {
		return fmt.Errorf("smtp kirim email: %w", err)
	}

	return nil
}
