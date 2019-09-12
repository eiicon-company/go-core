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
	}
)

func (m *smtpMail) Send(data *Data) error {
	e := email.NewEmail()
	e.To = data.To
	e.Bcc = data.Bcc
	e.Cc = data.Cc
	e.From = data.From
	e.Subject = data.Subject
	if data.Text != nil {
		e.Text = data.Text
	}
	if data.HTML != nil {
		e.HTML = data.HTML
	}

	auth := smtp.PlainAuth("",
		m.dsn.User, m.dsn.Password, m.dsn.Host,
	)

	if !m.dsn.TLS {
		return e.Send(m.dsn.Addr, auth)
	}

	cfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.dsn.TLSServer,
	}
	return e.SendWithTLS(m.dsn.Addr, auth, cfg)
}
