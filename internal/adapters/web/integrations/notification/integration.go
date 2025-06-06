package notification

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"

	"github.com/tapiaw38/auth-api-be/internal/platform/config"
)

type (
	Integration interface {
		SendEmail(SendEmailInput) error
	}

	integration struct {
		appName  string
		smtpHost string
		smtpPort string
		username string
		password string
	}

	SendEmailInput struct {
		To           string            `json:"to"`
		Subject      string            `json:"subject"`
		TemplateName string            `json:"template_name"`
		Variables    map[string]string `json:"variables"`
	}
)

func NewIntegration(cfg *config.ConfigurationService) Integration {
	return &integration{
		appName:  cfg.AppName,
		smtpHost: cfg.Notification.Email.Host,
		smtpPort: cfg.Notification.Email.Port,
		username: cfg.Notification.Email.Username,
		password: cfg.Notification.Email.Password,
	}
}

func (i *integration) SendEmail(input SendEmailInput) error {
	from := mail.Address{
		Name:    i.appName,
		Address: i.username,
	}
	to := mail.Address{
		Name:    "",
		Address: input.To,
	}

	encodedSubject := fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(input.Subject)))

	headers := map[string]string{
		"From":         from.String(),
		"To":           to.String(),
		"Subject":      encodedSubject,
		"Content-Type": "text/html; charset=UTF-8",
	}

	message := bytes.Buffer{}
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")

	tmpl, err := template.ParseFiles("templates/" + input.TemplateName + ".html")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&message, input.Variables); err != nil {
		return err
	}

	return i.sendSMTPEmail(input.To, i.username, message.String())
}

func (i *integration) sendSMTPEmail(toEmail, fromEmail, message string) error {

	host := i.smtpHost
	servername := i.smtpHost + ":" + i.smtpPort

	auth := smtp.PlainAuth("", i.username, i.password, host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(fromEmail); err != nil {
		return err
	}

	if err = client.Rcpt(toEmail); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = client.Quit()
	if err != nil {
		return err
	}

	return nil
}
