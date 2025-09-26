package mail

import (
	"bytes"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type Provider struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailProvider(host string, port int, username, password, from string) *Provider {
	return &Provider{
		dialer: gomail.NewDialer(host, port, username, password),
		from:   from,
	}
}

func (p *Provider) SendMail(to, subject, templateFile string, data any) error {
	path := filepath.Join("internal/provider/mail/templates", templateFile)
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", p.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	return p.dialer.DialAndSend(m)
}
