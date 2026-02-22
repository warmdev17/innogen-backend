package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// SendEmail sends an email using SMTP
func SendEmail(to []string, subject, body string) error {
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")
	smtpUser := os.Getenv("MAIL_USER")
	smtpPass := os.Getenv("MAIL_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP environment variables are not set properly")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Format message
	header := make(map[string]string)
	header["From"] = smtpUser
	header["To"] = strings.Join(to, ",")
	header["Subject"] = subject

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	address := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	conn, err := tls.Dial("tcp", address, tlsconfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(smtpUser); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", addr, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create data stream: %w", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data stream: %w", err)
	}

	return nil
}
