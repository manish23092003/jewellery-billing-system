package service

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

// generateSecureToken creates a cryptographically secure random token.
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// normalizePhone removes common characters like spaces, dashes, parentheses.
// It also strips the +91 country code for Indian numbers if present to prevent duplicates.
func normalizePhone(phone string) string {
	if phone == "" {
		return ""
	}
	// Remove all non-digit and non-plus characters
	re := regexp.MustCompile(`[^\d+]`)
	clean := re.ReplaceAllString(phone, "")

	if strings.HasPrefix(clean, "+91") {
		clean = strings.TrimPrefix(clean, "+91")
	} else if strings.HasPrefix(clean, "91") && len(clean) == 12 {
		clean = strings.TrimPrefix(clean, "91")
	}

	return clean
}
