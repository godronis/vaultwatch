package alert

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/yourusername/vaultwatch/internal/config"
)

// EmailAlerter sends alert notifications via SMTP email.
type EmailAlerter struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       []string
}

// NewEmailAlerter creates a new EmailAlerter from the provided config.
// Returns an error if required fields are missing.
func NewEmailAlerter(cfg config.EmailConfig) (*EmailAlerter, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("email alerter: SMTP host is required")
	}
	if cfg.From == "" {
		return nil, fmt.Errorf("email alerter: sender address is required")
	}
	if len(cfg.To) == 0 {
		return nil, fmt.Errorf("email alerter: at least one recipient is required")
	}
	port := cfg.Port
	if port == 0 {
		port = 587
	}
	return &EmailAlerter{
		host:     cfg.Host,
		port:     port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
		to:       cfg.To,
	}, nil
}

// Send delivers an alert email for the given secret path and expiry information.
func (e *EmailAlerter) Send(path string, payload map[string]string) error {
	subject := fmt.Sprintf("[VaultWatch] Secret expiring: %s", path)
	body := buildEmailBody(path, payload)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.from,
		strings.Join(e.to, ", "),
		subject,
		body,
	))

	addr := fmt.Sprintf("%s:%d", e.host, e.port)
	var auth smtp.Auth
	if e.username != "" {
		auth = smtp.PlainAuth("", e.username, e.password, e.host)
	}

	if err := smtp.SendMail(addr, auth, e.from, e.to, msg); err != nil {
		return fmt.Errorf("email alerter: failed to send mail: %w", err)
	}
	return nil
}

func buildEmailBody(path string, payload map[string]string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Secret path: %s\n", path))
	for k, v := range payload {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	return sb.String()
}
