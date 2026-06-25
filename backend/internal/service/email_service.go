package service

import (
	"fmt"
	"net/smtp"

	"github.com/rs/zerolog/log"

	"jewellery-billing/internal/config"
)

// EmailSender defines the interface for sending emails.
type EmailSender interface {
	SendVerificationEmail(to, name, token string) error
	SendPasswordResetEmail(to, name, token string) error
}

// ── Console Email Sender (Development) ─────────────────────────────────

// ConsoleEmailSender logs emails to the console instead of sending them.
// Use this for development and testing.
type ConsoleEmailSender struct {
	appURL string
}

func NewConsoleEmailSender(appURL string) *ConsoleEmailSender {
	return &ConsoleEmailSender{appURL: appURL}
}

func (s *ConsoleEmailSender) SendVerificationEmail(to, name, token string) error {
	link := fmt.Sprintf("%s/#/verify-email?token=%s", s.appURL, token)
	log.Info().
		Str("to", to).
		Str("name", name).
		Str("link", link).
		Msg("📧 [DEV] Email Verification — click the link to verify your email")
	return nil
}

func (s *ConsoleEmailSender) SendPasswordResetEmail(to, name, token string) error {
	link := fmt.Sprintf("%s/#/reset-password?token=%s", s.appURL, token)
	log.Info().
		Str("to", to).
		Str("name", name).
		Str("link", link).
		Msg("📧 [DEV] Password Reset — click the link to reset your password")
	return nil
}

// ── SMTP Email Sender (Production) ─────────────────────────────────────

// SMTPEmailSender sends real emails via SMTP.
type SMTPEmailSender struct {
	host     string
	port     string
	username string
	password string
	from     string
	appURL   string
}

func NewSMTPEmailSender(cfg *config.Config) *SMTPEmailSender {
	return &SMTPEmailSender{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUser,
		password: cfg.SMTPPassword,
		from:     cfg.SMTPFrom,
		appURL:   cfg.AppURL,
	}
}

func (s *SMTPEmailSender) SendVerificationEmail(to, name, token string) error {
	link := fmt.Sprintf("%s/#/verify-email?token=%s", s.appURL, token)
	subject := "Verify Your Email — Jewellery Billing"
	body := fmt.Sprintf(`Hello %s,

Welcome to Jewellery Billing! Please verify your email address by clicking the link below:

%s

This link will expire in 24 hours.

If you didn't create an account, you can safely ignore this email.

Best regards,
Jewellery Billing Team`, name, link)

	return s.sendMail(to, subject, body)
}

func (s *SMTPEmailSender) SendPasswordResetEmail(to, name, token string) error {
	link := fmt.Sprintf("%s/#/reset-password?token=%s", s.appURL, token)
	subject := "Reset Your Password — Jewellery Billing"
	body := fmt.Sprintf(`Hello %s,

We received a request to reset your password. Click the link below to set a new password:

%s

This link will expire in 1 hour.

If you didn't request a password reset, you can safely ignore this email.

Best regards,
Jewellery Billing Team`, name, link)

	return s.sendMail(to, subject, body)
}

func (s *SMTPEmailSender) sendMail(to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		s.from, to, subject, body)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	if err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg)); err != nil {
		log.Error().Err(err).Str("to", to).Msg("Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().Str("to", to).Str("subject", subject).Msg("Email sent successfully")
	return nil
}
