// Package storage provides validation
package storage

import (
	"fmt"
)

// ValidateKey validates key
func ValidateKey(key string) error {
	if key == "" {
		return fmt.Errorf("key required")
	}
	return nil
}

// ValidateValue validates value
func ValidateValue(value []byte) error {
	if len(value) == 0 {
		return fmt.Errorf("value required")
	}
	return nil
}
