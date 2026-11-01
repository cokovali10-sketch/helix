// Package catalog provides validation
package catalog

import (
	"fmt"
	"strings"
)

// ValidateName validates name
func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name required")
	}
	if len(name) > 64 {
		return fmt.Errorf("name too long")
	}
	return nil
}

// ValidateEmail validates email
func ValidateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email")
	}
	return nil
}
