package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var validUsername = regexp.MustCompile(`^[A-Z0-9_]+$`)

func SanitizeAndValidateUsername(raw string) (string, error) {
	if raw == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	username := strings.TrimSpace(raw)
	if username == "" {
		return "", fmt.Errorf("username cannot be spaces")
	}

	username = strings.ToUpper(username)
	if !validUsername.MatchString(username) {
		return "", fmt.Errorf("username can only be alphanumeric and underscore")
	}
	return username, nil
}

func SanitizeUsernameWithoutError(username string) string {
	username = validUsername.ReplaceAllString(username, "")
	return strings.ToUpper(username)
}
