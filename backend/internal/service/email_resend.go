package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"jewellery-billing/internal/config"
)

// ResendEmailSender sends real emails via Resend HTTP API.
type ResendEmailSender struct {
	apiKey string
	from   string
	appURL string
}

func NewResendEmailSender(cfg *config.Config) *ResendEmailSender {
	// If the user hasn't specified a verified domain, use Resend's default testing domain
	from := cfg.SMTPFrom
	if from == "" {
		from = "onboarding@resend.dev"
	}

	return &ResendEmailSender{
		apiKey: cfg.ResendAPIKey,
		from:   from,
		appURL: cfg.AppURL,
	}
}

func (s *ResendEmailSender) SendVerificationEmail(to, name, token string) error {
	link := fmt.Sprintf("%s/verify-email?token=%s", s.appURL, token)
	subject := "Verify Your Email — Jewellery Billing"
	htmlBody := fmt.Sprintf(`
		<div style="font-family: sans-serif; padding: 20px;">
			<h2>Hello %s,</h2>
			<p>Welcome to Jewellery Billing! Please verify your email address by clicking the button below:</p>
			<a href="%s" style="background-color: #c6a962; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block; margin-top: 10px;">Verify Email</a>
			<p style="margin-top: 20px; font-size: 12px; color: #666;">This link will expire in 24 hours.</p>
		</div>
	`, name, link)

	return s.sendMail(to, subject, htmlBody)
}

func (s *ResendEmailSender) SendPasswordResetEmail(to, name, token string) error {
	link := fmt.Sprintf("%s/reset-password?token=%s", s.appURL, token)
	subject := "Reset Your Password — Jewellery Billing"
	htmlBody := fmt.Sprintf(`
		<div style="font-family: sans-serif; padding: 20px;">
			<h2>Hello %s,</h2>
			<p>We received a request to reset your password. Click the button below to set a new password:</p>
			<a href="%s" style="background-color: #c6a962; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block; margin-top: 10px;">Reset Password</a>
			<p style="margin-top: 20px; font-size: 12px; color: #666;">This link will expire in 1 hour.</p>
		</div>
	`, name, link)

	return s.sendMail(to, subject, htmlBody)
}

func (s *ResendEmailSender) sendMail(to, subject, htmlBody string) error {
	payload := map[string]interface{}{
		"from":    s.from,
		"to":      []string{to},
		"subject": subject,
		"html":    htmlBody,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal resend payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("to", to).Msg("Failed to send email via Resend")
		return fmt.Errorf("failed to send email via Resend: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Error().Int("status", resp.StatusCode).Str("to", to).Msg("Resend API returned error")
		return fmt.Errorf("resend API returned error status: %d", resp.StatusCode)
	}

	log.Info().Str("to", to).Str("subject", subject).Msg("Email sent successfully via Resend!")
	return nil
}
