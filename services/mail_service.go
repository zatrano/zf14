package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

type MailService struct {
	host     string
	port     string
	username string
	password string
}

func NewMailService() *MailService {
	return &MailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (m *MailService) SendMail(to, subject, body string) error {
	from := m.username
	header := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n", from, to, subject)
	message := []byte(header + body)

	address := fmt.Sprintf("%s:%s", m.host, m.port)

	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.host,
	})
	if err != nil {
		return fmt.Errorf("tls bağlantısı kurulamadı: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return fmt.Errorf("smtp istemcisi oluşturulamadı: %w", err)
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp kimlik doğrulama başarısız: %w", err)
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("gönderen ayarlanamadı: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("alıcı ayarlanamadı: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("e-posta verisi gönderilemedi: %w", err)
	}
	if _, err := writer.Write(message); err != nil {
		return fmt.Errorf("e-posta yazılamadı: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("e-posta gönderimi tamamlanamadı: %w", err)
	}

	return nil
}
