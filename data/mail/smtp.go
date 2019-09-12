package mail

import (
	"crypto/tls"
	"net/smtp"

	"github.com/jordan-wright/email"

	"github.com/eiicon-company/go-core/util/dsn"
)

type (
	smtpMail struct {
		dsn *dsn.MailDSN

		To      []string
		Bcc     []string
		Cc      []string
		From    string
		Subject string
		Text    []byte
		HTML    []byte
	}
)

func (m *smtpMail) Send() error {
	e := email.NewEmail()
	e.To = m.To
	e.Bcc = m.Bcc
	e.Cc = m.Cc
	e.From = m.From
	e.Subject = m.Subject
	if m.Text != nil {
		e.Text = m.Text
	}
	if m.HTML != nil {
		e.HTML = m.HTML
	}

	auth := smtp.PlainAuth("", m.dsn.User, m.dsn.Password, m.dsn.Host)

	if !m.dsn.TLS {
		return e.Send(m.dsn.Addr, auth)
	}

	cfg := &tls.Config{InsecureSkipVerify: true, ServerName: m.dsn.TLSServer}
	return e.SendWithTLS(m.dsn.Addr, auth, cfg)
}
