package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var validUsername = regexp.MustCompile(`^[A-Z0-9_]+$`)

func validateUsername(username string) error {
	if !validUsername.MatchString(username) {
		return fmt.Errorf("username only allows alphanumeric, uppercase and underscore")
	}
	return nil
}

func SanitizeAndValidateUsername(raw string) (string, error) {
	if raw == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	username := strings.TrimSpace(raw)
	if username == "" {
		return "", fmt.Errorf("username cannot be only spaces")
	}

	username = strings.ToUpper(username)
	if err := validateUsername(username); err != nil {
		return "", err
	}
	return username, nil
}

func SanitizeUsernameWithoutError(username string) string {
	var removeUsernameIllegalRegex = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	username = removeUsernameIllegalRegex.ReplaceAllString(username, "")
	return strings.ToUpper(username)
}
