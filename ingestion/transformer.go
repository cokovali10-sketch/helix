// Package ingestion provides data transformation
package ingestion

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// TransformFunc transforms a record
type TransformFunc func(map[string]interface{}) (map[string]interface{}, error)

// Transformer applies transforms
type Transformer struct {
	transforms []TransformFunc
	mu         sync.RWMutex
}

// NewTransformer creates transformer
func NewTransformer() *Transformer {
	return &Transformer{transforms: []TransformFunc{}}
}

// AddTransform adds a transform
func (t *Transformer) AddTransform(fn TransformFunc) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.transforms = append(t.transforms, fn)
}

// Apply applies all transforms
func (t *Transformer) Apply(record map[string]interface{}) (map[string]interface{}, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	current := record
	for _, fn := range t.transforms {
		next, err := fn(current)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current, nil
}

// RenameField renames a field
func RenameField(from, to string) TransformFunc {
	return func(r map[string]interface{}) (map[string]interface{}, error) {
		if v, ok := r[from]; ok {
			r[to] = v
			delete(r, from)
		}
		return r, nil
	}
}

// AddTimestamp adds timestamp
func AddTimestamp(field string) TransformFunc {
	return func(r map[string]interface{}) (map[string]interface{}, error) {
		r[field] = time.Now().UTC().Format(time.RFC3339)
		return r, nil
	}
}

// DropField drops a field
func DropField(field string) TransformFunc {
	return func(r map[string]interface{}) (map[string]interface{}, error) {
		delete(r, field)
		return r, nil
	}
}

// FilterField keeps only specified fields
func FilterFields(fields []string) TransformFunc {
	return func(r map[string]interface{}) (map[string]interface{}, error) {
		out := map[string]interface{}{}
		for _, f := range fields {
			if v, ok := r[f]; ok {
				out[f] = v
			}
		}
		return out, nil
	}
}

// ScaleValue scales numeric field
func ScaleValue(field string, factor float64) TransformFunc {
	return func(r map[string]interface{}) (map[string]interface{}, error) {
		if v, ok := r[field]; ok {
			switch val := v.(type) {
			case float64:
				r[field] = val * factor
			case int:
				r[field] = float64(val) * factor
			case json.Number:
				f, _ := val.Float64()
				r[field] = f * factor
			}
		}
		return r, nil
	}
}

// ToString converts field to string
func ToString(field string) TransformFunc {
	return func(r map[string]interface{}) (map[string]interface{}, error) {
		if v, ok := r[field]; ok {
			r[field] = fmt.Sprintf("%v", v)
		}
		return r, nil
	}
}

// Count returns number of transforms
func (t *Transformer) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.transforms)
}

// Clear clears transforms
func (t *Transformer) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.transforms = []TransformFunc{}
}
