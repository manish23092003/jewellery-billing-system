package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"jewellery-billing/internal/config"
)

// BrevoEmailSender sends real emails via Brevo HTTP API.
type BrevoEmailSender struct {
	apiKey string
	from   string
	appURL string
}

func NewBrevoEmailSender(cfg *config.Config) *BrevoEmailSender {
	return &BrevoEmailSender{
		apiKey: cfg.BrevoAPIKey,
		from:   cfg.SMTPFrom,
		appURL: cfg.AppURL,
	}
}

func (s *BrevoEmailSender) SendVerificationEmail(to, name, token string) error {
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

	return s.sendMail(to, name, subject, htmlBody)
}

func (s *BrevoEmailSender) SendPasswordResetEmail(to, name, token string) error {
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

	return s.sendMail(to, name, subject, htmlBody)
}

func (s *BrevoEmailSender) sendMail(toEmail, toName, subject, htmlBody string) error {
	// Brevo API expects this exact payload structure
	payload := map[string]interface{}{
		"sender": map[string]string{
			"email": s.from,
			"name":  "Jewellery Billing System",
		},
		"to": []map[string]string{
			{
				"email": toEmail,
				"name":  toName,
			},
		},
		"subject":     subject,
		"htmlContent": htmlBody,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal brevo payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Brevo uses api-key header
	req.Header.Set("api-key", s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("to", toEmail).Msg("Failed to send email via Brevo")
		return fmt.Errorf("failed to send email via Brevo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Error().Int("status", resp.StatusCode).Str("to", toEmail).Msg("Brevo API returned error")
		return fmt.Errorf("brevo API returned error status: %d", resp.StatusCode)
	}

	log.Info().Str("to", toEmail).Str("subject", subject).Msg("Email sent successfully via Brevo!")
	return nil
}
