// Package ingestion provides serialization
package ingestion

import (
	"encoding/json"
	"fmt"
)

// Serializer serializes records
type Serializer struct{}

// NewSerializer creates serializer
func NewSerializer() *Serializer {
	return &Serializer{}
}

// ToJSON serializes to JSON
func (s *Serializer) ToJSON(record map[string]interface{}) ([]byte, error) {
	return json.Marshal(record)
}

// FromJSON deserializes from JSON
func (s *Serializer) FromJSON(data []byte) (map[string]interface{}, error) {
	var record map[string]interface{}
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}
	return record, nil
}

// ToJSONString serializes to string
func (s *Serializer) ToJSONString(record map[string]interface{}) (string, error) {
	data, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// BatchToJSON serializes batch
func (s *Serializer) BatchToJSON(records []map[string]interface{}) ([]byte, error) {
	return json.Marshal(records)
}

// BatchFromJSON deserializes batch
func (s *Serializer) BatchFromJSON(data []byte) ([]map[string]interface{}, error) {
	var records []map[string]interface{}
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// Validate validates JSON
func (s *Serializer) Validate(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	return nil
}
