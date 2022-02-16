package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"path"
	"strings"
	"text/template"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
)

var MailRequests = make(chan Mail)

type Mail struct {
	EmailTo   []string    `json:"email_to"`
	EmailFrom string      `json:"email_from"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Subject   string      `json:"subject"`
	Msg       string      `json:"message"`
	Body      []byte
}

func (m *Mail) Send() error {
	return smtp.SendMail(
		fmt.Sprintf("%s:%s", config.GetConfig().Mail.SmtpHost, config.GetConfig().Mail.SmtpPort),
		smtp.PlainAuth("", config.GetConfig().Mail.SmtpPass, config.GetConfig().Mail.SmtpPwd, config.GetConfig().Mail.SmtpHost),
		config.GetConfig().Mail.EmaiFrom,
		m.EmailTo,
		[]byte(m.Msg),
	)
}

func (m *Mail) ParseTemplate() error {
	t, ok := MailTemplates[m.Type]
	if !ok {
		return fmt.Errorf("unsupported type %s", m.Type)
	}

	tmplName := fmt.Sprintf("email_%s.html", t)
	tmplFilepath := path.Clean(fmt.Sprintf("%s/%s", config.GetConfig().Directory.MailTemplates, tmplName))
	tmpl, err := template.ParseFiles(tmplFilepath)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, m.Data); err != nil {
		return err
	}
	m.Body = tpl.Bytes()
	return nil
}

func (m *Mail) BuildMessage() error {
	if err := m.ParseTemplate(); err != nil {
		return err
	}
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", m.EmailFrom)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(m.EmailTo, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", m.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", m.Body)

	m.Msg = msg
	return nil
}

func ProcessEmail(m *Mail) error {
	if err := m.BuildMessage(); err != nil {
		return err
	}

	return m.Send()
}
