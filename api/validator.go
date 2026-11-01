// Package api provides validation helpers
package api

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator validates input
type Validator struct {
	errors []string
}

// NewValidator creates validator
func NewValidator() *Validator {
	return &Validator{errors: []string{}}
}

// Required checks required
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, fmt.Sprintf("%s is required", field))
	}
	return v
}

// MinLength checks min length
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errors = append(v.errors, fmt.Sprintf("%s must be at least %d characters", field, min))
	}
	return v
}

// MaxLength checks max length
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.errors = append(v.errors, fmt.Sprintf("%s must be at most %d characters", field, max))
	}
	return v
}

// Email checks email format
func (v *Validator) Email(field, value string) *Validator {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(value) {
		v.errors = append(v.errors, fmt.Sprintf("%s must be a valid email", field))
	}
	return v
}

// IsValid returns validity
func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

// Errors returns errors
func (v *Validator) Errors() []string {
	return v.errors
}

// Error returns combined error
func (v *Validator) Error() error {
	if v.IsValid() {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(v.errors, "; "))
}
