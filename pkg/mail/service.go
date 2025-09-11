package mail

import (
	"bytes"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type Service struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailService(host string, port int, username, password, from string) *Service {
	return &Service{
		dialer: gomail.NewDialer(host, port, username, password),
		from:   from,
	}
}

func (s *Service) SendMail(to, subject, templateFile string, data any) error {
	path := filepath.Join("pkg/mail/templates", templateFile)
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	return s.dialer.DialAndSend(m)
}
