// Package ingestion provides schema validation
package ingestion

import (
	"fmt"
	"sync"
)

// SchemaField represents field
type SchemaField struct {
	Name     string
	Type     string
	Required bool
}

// IngestSchema represents schema
type IngestSchema struct {
	Fields []SchemaField
	mu     sync.RWMutex
}

// NewIngestSchema creates schema
func NewIngestSchema() *IngestSchema {
	return &IngestSchema{Fields: []SchemaField{}}
}

// AddField adds field
func (s *IngestSchema) AddField(f SchemaField) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Fields = append(s.Fields, f)
}

// Validate validates record
func (s *IngestSchema) Validate(record map[string]interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.Fields {
		v, ok := record[f.Name]
		if f.Required && !ok {
			return fmt.Errorf("missing required field: %s", f.Name)
		}
		if ok {
			switch f.Type {
			case "string":
				if _, ok := v.(string); !ok {
					return fmt.Errorf("field %s must be string", f.Name)
				}
			case "number":
				switch v.(type) {
				case float64, int, int64:
				default:
					return fmt.Errorf("field %s must be number", f.Name)
				}
			}
		}
	}
	return nil
}

// FieldNames returns field names
func (s *IngestSchema) FieldNames() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []string{}
	for _, f := range s.Fields {
		out = append(out, f.Name)
	}
	return out
}
