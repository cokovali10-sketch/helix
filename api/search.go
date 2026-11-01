// Package api provides additional API handlers
// including search, autocomplete, and suggestion handlers
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
)

// SearchHandler handles search requests
type SearchHandler struct {
	index   *SearchIndex
	mu      sync.RWMutex
}

// SearchIndex represents a search index
type SearchIndex struct {
	documents map[string]*SearchDocument
	terms     map[string]map[string]int // term -> docID -> count
	mu        sync.RWMutex
}

// SearchDocument represents a searchable document
type SearchDocument struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	Labels   map[string]string `json:"labels"`
	Score    float64           `json:"score"`
}

// NewSearchHandler creates a new search handler
func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		index: &SearchIndex{
			documents: make(map[string]*SearchDocument),
			terms:     make(map[string]map[string]int),
		},
	}
}

// IndexDocument indexes a document
func (sh *SearchHandler) IndexDocument(doc *SearchDocument) {
	sh.index.mu.Lock()
	defer sh.index.mu.Unlock()

	sh.index.documents[doc.ID] = doc

	// Tokenize and index
	words := tokenize(doc.Title + " " + doc.Content)
	for _, word := range words {
		if sh.index.terms[word] == nil {
			sh.index.terms[word] = make(map[string]int)
		}
		sh.index.terms[word][doc.ID]++
	}
}

// RemoveDocument removes a document from the index
func (sh *SearchHandler) RemoveDocument(id string) {
	sh.index.mu.Lock()
	defer sh.index.mu.Unlock()

	delete(sh.index.documents, id)

	for term := range sh.index.terms {
		delete(sh.index.terms[term], id)
		if len(sh.index.terms[term]) == 0 {
			delete(sh.index.terms, term)
		}
	}
}

// Search searches for documents
func (sh *SearchHandler) Search(query string, limit int) []*SearchDocument {
	sh.index.mu.RLock()
	defer sh.index.mu.RUnlock()

	words := tokenize(query)
	if len(words) == 0 {
		return nil
	}

	// Score documents
	scores := make(map[string]float64)
	for _, word := range words {
		if docs, ok := sh.index.terms[word]; ok {
			for docID, count := range docs {
				scores[docID] += float64(count)
			}
		}
	}

	// Sort by score
	type scoredDoc struct {
		doc   *SearchDocument
		score float64
	}
	var results []scoredDoc
	for docID, score := range scores {
		if doc, ok := sh.index.documents[docID]; ok {
			doc.Score = score
			results = append(results, scoredDoc{doc, score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	docs := make([]*SearchDocument, len(results))
	for i, r := range results {
		docs[i] = r.doc
	}
	return docs
}

// ServeHTTP serves search endpoint
func (sh *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query required", http.StatusBadRequest)
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	results := sh.Search(query, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
		"total":   len(results),
		"query":   query,
	})
}

// AutocompleteHandler handles autocomplete requests
type AutocompleteHandler struct {
	suggestions map[string][]string
	mu          sync.RWMutex
}

// NewAutocompleteHandler creates a new autocomplete handler
func NewAutocompleteHandler() *AutocompleteHandler {
	return &AutocompleteHandler{
		suggestions: make(map[string][]string),
	}
}

// AddSuggestions adds suggestions for a prefix
func (ah *AutocompleteHandler) AddSuggestions(prefix string, suggestions []string) {
	ah.mu.Lock()
	defer ah.mu.Unlock()
	ah.suggestions[prefix] = suggestions
}

// GetSuggestions gets suggestions for a prefix
func (ah *AutocompleteHandler) GetSuggestions(prefix string, limit int) []string {
	ah.mu.RLock()
	defer ah.mu.RUnlock()

	var result []string
	for p, suggestions := range ah.suggestions {
		if strings.HasPrefix(p, prefix) {
			result = append(result, suggestions...)
		}
	}

	// Also try partial match
	for p, suggestions := range ah.suggestions {
		if strings.Contains(p, prefix) {
			result = append(result, suggestions...)
		}
	}

	// Deduplicate
	seen := make(map[string]bool)
	var unique []string
	for _, s := range result {
		if !seen[s] {
			seen[s] = true
			unique = append(unique, s)
		}
	}

	if limit > 0 && len(unique) > limit {
		unique = unique[:limit]
	}

	return unique
}

// ServeHTTP serves autocomplete endpoint
func (ah *AutocompleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prefix := r.URL.Query().Get("q")
	if prefix == "" {
		http.Error(w, "Query required", http.StatusBadRequest)
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	suggestions := ah.GetSuggestions(prefix, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"suggestions": suggestions,
		"total":       len(suggestions),
	})
}

// SuggestionHandler handles suggestion requests
type SuggestionHandler struct {
	suggestions map[string][]Suggestion
	mu          sync.RWMutex
}

// Suggestion represents a suggestion
type Suggestion struct {
	Text        string `json:"text"`
	Type        string `json:"type"` // metric, function, label, operator
	Description string `json:"description"`
	Score       float64 `json:"score"`
}

// NewSuggestionHandler creates a new suggestion handler
func NewSuggestionHandler() *SuggestionHandler {
	sh := &SuggestionHandler{
		suggestions: make(map[string][]Suggestion),
	}

	// Add default suggestions
	sh.addDefaultSuggestions()

	return sh
}

func (sh *SuggestionHandler) addDefaultSuggestions() {
	// Metric suggestions
	sh.suggestions["metric"] = []Suggestion{
		{Text: "cpu_usage", Type: "metric", Description: "CPU usage percentage", Score: 1.0},
		{Text: "memory_usage", Type: "metric", Description: "Memory usage in bytes", Score: 1.0},
		{Text: "disk_usage", Type: "metric", Description: "Disk usage in bytes", Score: 1.0},
		{Text: "network_rx_bytes", Type: "metric", Description: "Network received bytes", Score: 1.0},
		{Text: "network_tx_bytes", Type: "metric", Description: "Network transmitted bytes", Score: 1.0},
		{Text: "http_requests_total", Type: "metric", Description: "Total HTTP requests", Score: 1.0},
		{Text: "http_request_duration_seconds", Type: "metric", Description: "HTTP request duration", Score: 1.0},
		{Text: "up", Type: "metric", Description: "Target up status", Score: 1.0},
	}

	// Function suggestions
	sh.suggestions["function"] = []Suggestion{
		{Text: "rate", Type: "function", Description: "Calculate per-second rate", Score: 1.0},
		{Text: "irate", Type: "function", Description: "Calculate instant rate", Score: 1.0},
		{Text: "increase", Type: "function", Description: "Calculate increase", Score: 1.0},
		{Text: "sum", Type: "function", Description: "Sum values", Score: 1.0},
		{Text: "avg", Type: "function", Description: "Average values", Score: 1.0},
		{Text: "min", Type: "function", Description: "Minimum value", Score: 1.0},
		{Text: "max", Type: "function", Description: "Maximum value", Score: 1.0},
		{Text: "count", Type: "function", Description: "Count values", Score: 1.0},
		{Text: "topk", Type: "function", Description: "Top K values", Score: 1.0},
		{Text: "bottomk", Type: "function", Description: "Bottom K values", Score: 1.0},
		{Text: "histogram_quantile", Type: "function", Description: "Histogram quantile", Score: 1.0},
		{Text: "abs", Type: "function", Description: "Absolute value", Score: 1.0},
		{Text: "ceil", Type: "function", Description: "Ceiling", Score: 1.0},
		{Text: "floor", Type: "function", Description: "Floor", Score: 1.0},
		{Text: "round", Type: "function", Description: "Round", Score: 1.0},
		{Text: "sqrt", Type: "function", Description: "Square root", Score: 1.0},
	}

	// Operator suggestions
	sh.suggestions["operator"] = []Suggestion{
		{Text: "+", Type: "operator", Description: "Addition", Score: 1.0},
		{Text: "-", Type: "operator", Description: "Subtraction", Score: 1.0},
		{Text: "*", Type: "operator", Description: "Multiplication", Score: 1.0},
		{Text: "/", Type: "operator", Description: "Division", Score: 1.0},
		{Text: "%", Type: "operator", Description: "Modulo", Score: 1.0},
		{Text: "^", Type: "operator", Description: "Power", Score: 1.0},
		{Text: "==", Type: "operator", Description: "Equal", Score: 1.0},
		{Text: "!=", Type: "operator", Description: "Not equal", Score: 1.0},
		{Text: "<", Type: "operator", Description: "Less than", Score: 1.0},
		{Text: ">", Type: "operator", Description: "Greater than", Score: 1.0},
		{Text: "<=", Type: "operator", Description: "Less than or equal", Score: 1.0},
		{Text: ">=", Type: "operator", Description: "Greater than or equal", Score: 1.0},
		{Text: "=~", Type: "operator", Description: "Regex match", Score: 1.0},
		{Text: "!~", Type: "operator", Description: "Regex not match", Score: 1.0},
	}

	// Keyword suggestions
	sh.suggestions["keyword"] = []Suggestion{
		{Text: "by", Type: "keyword", Description: "Group by labels", Score: 1.0},
		{Text: "without", Type: "keyword", Description: "Group without labels", Score: 1.0},
		{Text: "on", Type: "keyword", Description: "Join on labels", Score: 1.0},
		{Text: "ignoring", Type: "keyword", Description: "Join ignoring labels", Score: 1.0},
		{Text: "group_left", Type: "keyword", Description: "Group left join", Score: 1.0},
		{Text: "group_right", Type: "keyword", Description: "Group right join", Score: 1.0},
		{Text: "bool", Type: "keyword", Description: "Boolean modifier", Score: 1.0},
		{Text: "offset", Type: "keyword", Description: "Time offset", Score: 1.0},
	}
}

// GetSuggestions gets suggestions for a query
func (sh *SuggestionHandler) GetSuggestions(query string, limit int) []Suggestion {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		// Return all suggestions
		var all []Suggestion
		for _, suggestions := range sh.suggestions {
			all = append(all, suggestions...)
		}
		if limit > 0 && len(all) > limit {
			all = all[:limit]
		}
		return all
	}

	var result []Suggestion
	for _, suggestions := range sh.suggestions {
		for _, s := range suggestions {
			if strings.HasPrefix(strings.ToLower(s.Text), query) {
				result = append(result, s)
			}
		}
	}

	// Sort by score
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result
}

// ServeHTTP serves suggestion endpoint
func (sh *SuggestionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	suggestions := sh.GetSuggestions(query, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"suggestions": suggestions,
		"total":       len(suggestions),
	})
}

func tokenize(text string) []string {
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
	})
	return words
}