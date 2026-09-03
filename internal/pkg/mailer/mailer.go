// Package mailer — минимальная обёртка над net/smtp, без внешних зависимостей.
package mailer

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

// Config — параметры одной отправки. Приходит из internal/settings уже с
// расшифрованным паролем — mailer в шифрование/хранение не лезет.
type Config struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
	UseTLS      bool // true — STARTTLS обязателен (порт 587 и т.п.), сервер без него — ошибка
}

// Send отправляет одно HTML-письмо. Не использует net/smtp.SendMail: тот
// апгрейдит до STARTTLS оппортунистически (если сервер предложил), не давая
// потребовать TLS явно — а нам нужно честно упасть, если UseTLS=true, а
// сервер STARTTLS не поддерживает, а не тихо уйти в plaintext.
func Send(cfg Config, to, subject, htmlBody string) error {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}
	defer client.Close()

	if err := client.Hello("localhost"); err != nil {
		return fmt.Errorf("smtp hello: %w", err)
	}

	if cfg.UseTLS {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return errors.New("smtp: сервер не поддерживает STARTTLS")
		}
		if err := client.StartTLS(&tls.Config{ServerName: cfg.Host}); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}

	if cfg.Username != "" {
		if ok, _ := client.Extension("AUTH"); ok {
			auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth: %w", err)
			}
		}
	}

	if err := client.Mail(cfg.FromAddress); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}

	from := cfg.FromAddress
	if cfg.FromName != "" {
		from = fmt.Sprintf("%s <%s>", cfg.FromName, cfg.FromAddress)
	}
	if _, err := w.Write(buildMessage(from, to, subject, htmlBody)); err != nil {
		w.Close()
		return fmt.Errorf("smtp write body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}

	return client.Quit()
}

func buildMessage(from, to, subject, htmlBody string) []byte {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString(`Content-Type: text/html; charset="UTF-8"` + "\r\n")
	b.WriteString("\r\n")
	b.WriteString(htmlBody)
	return []byte(b.String())
}
