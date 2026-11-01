// Package storage provides additional storage features
// including advanced indexing and search
package storage

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// AdvancedIndex provides advanced indexing capabilities
type AdvancedIndex struct {
	indexes    map[string]*Index
	mu         sync.RWMutex
}

// Index represents an index
type Index struct {
	Name     string
	Type     string // btree, hash, inverted, fulltext
	Fields   []string
	Entries  map[string][]string // value -> keys
	mu       sync.RWMutex
}

// NewAdvancedIndex creates a new advanced index
func NewAdvancedIndex() *AdvancedIndex {
	return &AdvancedIndex{
		indexes: make(map[string]*Index),
	}
}

// CreateIndex creates a new index
func (ai *AdvancedIndex) CreateIndex(name, indexType string, fields []string) error {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if _, exists := ai.indexes[name]; exists {
		return fmt.Errorf("index %s already exists", name)
	}

	ai.indexes[name] = &Index{
		Name:    name,
		Type:    indexType,
		Fields:  fields,
		Entries: make(map[string][]string),
	}
	return nil
}

// DropIndex drops an index
func (ai *AdvancedIndex) DropIndex(name string) error {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if _, exists := ai.indexes[name]; !exists {
		return fmt.Errorf("index %s not found", name)
	}

	delete(ai.indexes, name)
	return nil
}

// GetIndex gets an index by name
func (ai *AdvancedIndex) GetIndex(name string) (*Index, bool) {
	ai.mu.RLock()
	defer ai.mu.RUnlock()
	index, ok := ai.indexes[name]
	return index, ok
}

// ListIndexes lists all indexes
func (ai *AdvancedIndex) ListIndexes() []*Index {
	ai.mu.RLock()
	defer ai.mu.RUnlock()

	indexes := make([]*Index, 0, len(ai.indexes))
	for _, idx := range ai.indexes {
		indexes = append(indexes, idx)
	}
	return indexes
}

// AddEntry adds an entry to an index
func (idx *Index) AddEntry(value string, key string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.Entries[value] = append(idx.Entries[value], key)
}

// RemoveEntry removes an entry from an index
func (idx *Index) RemoveEntry(value string, key string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	keys := idx.Entries[value]
	for i, k := range keys {
		if k == key {
			idx.Entries[value] = append(keys[:i], keys[i+1:]...)
			break
		}
	}
}

// Lookup looks up keys by value
func (idx *Index) Lookup(value string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if keys, ok := idx.Entries[value]; ok {
		result := make([]string, len(keys))
		copy(result, keys)
		return result
	}
	return nil
}

// LookupPrefix looks up keys by prefix
func (idx *Index) LookupPrefix(prefix string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var result []string
	for value, keys := range idx.Entries {
		if strings.HasPrefix(value, prefix) {
			result = append(result, keys...)
		}
	}
	return result
}

// LookupRange looks up keys in a range
func (idx *Index) LookupRange(start, end string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var result []string
	for value, keys := range idx.Entries {
		if value >= start && value <= end {
			result = append(result, keys...)
		}
	}
	return result
}

// GetAllValues returns all indexed values
func (idx *Index) GetAllValues() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	values := make([]string, 0, len(idx.Entries))
	for value := range idx.Entries {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}

// GetEntryCount returns the number of entries
func (idx *Index) GetEntryCount() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	count := 0
	for _, keys := range idx.Entries {
		count += len(keys)
	}
	return count
}

// FullTextIndex provides full-text search
type FullTextIndex struct {
	documents map[string]string
	terms     map[string]map[string]int // term -> doc -> count
	mu        sync.RWMutex
}

// NewFullTextIndex creates a new full-text index
func NewFullTextIndex() *FullTextIndex {
	return &FullTextIndex{
		documents: make(map[string]string),
		terms:     make(map[string]map[string]int),
	}
}

// AddDocument adds a document to the index
func (fti *FullTextIndex) AddDocument(id, content string) {
	fti.mu.Lock()
	defer fti.mu.Unlock()

	fti.documents[id] = content

	// Tokenize and index
	words := tokenize(content)
	for _, word := range words {
		if fti.terms[word] == nil {
			fti.terms[word] = make(map[string]int)
		}
		fti.terms[word][id]++
	}
}

// RemoveDocument removes a document from the index
func (fti *FullTextIndex) RemoveDocument(id string) {
	fti.mu.Lock()
	defer fti.mu.Unlock()

	delete(fti.documents, id)

	for term := range fti.terms {
		delete(fti.terms[term], id)
		if len(fti.terms[term]) == 0 {
			delete(fti.terms, term)
		}
	}
}

// Search searches for documents
func (fti *FullTextIndex) Search(query string) []string {
	fti.mu.RLock()
	defer fti.mu.RUnlock()

	words := tokenize(query)
	if len(words) == 0 {
		return nil
	}

	// Find documents matching all words
	docScores := make(map[string]int)
	for _, word := range words {
		if docs, ok := fti.terms[word]; ok {
			for docID, count := range docs {
				docScores[docID] += count
			}
		}
	}

	// Sort by score
	type docScore struct {
		id    string
		score int
	}
	var results []docScore
	for id, score := range docScores {
		results = append(results, docScore{id, score})
	}

	// Simple sort
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].score > results[i].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	ids := make([]string, len(results))
	for i, r := range results {
		ids[i] = r.id
	}
	return ids
}

// GetDocument gets a document by ID
func (fti *FullTextIndex) GetDocument(id string) (string, bool) {
	fti.mu.RLock()
	defer fti.mu.RUnlock()
	content, ok := fti.documents[id]
	return content, ok
}

// GetTermCount returns the number of terms
func (fti *FullTextIndex) GetTermCount() int {
	fti.mu.RLock()
	defer fti.mu.RUnlock()
	return len(fti.terms)
}

// GetDocumentCount returns the number of documents
func (fti *FullTextIndex) GetDocumentCount() int {
	fti.mu.RLock()
	defer fti.mu.RUnlock()
	return len(fti.documents)
}

// tokenize tokenizes text into words
func tokenize(text string) []string {
	// Simple tokenization
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
	})
	return words
}

// SpatialIndex provides spatial indexing
type SpatialIndex struct {
	points map[string]Point
	mu     sync.RWMutex
}

// Point represents a spatial point
type Point struct {
	ID  string
	Lat float64
	Lon float64
}

// NewSpatialIndex creates a new spatial index
func NewSpatialIndex() *SpatialIndex {
	return &SpatialIndex{
		points: make(map[string]Point),
	}
}

// AddPoint adds a point to the index
func (si *SpatialIndex) AddPoint(id string, lat, lon float64) {
	si.mu.Lock()
	defer si.mu.Unlock()
	si.points[id] = Point{ID: id, Lat: lat, Lon: lon}
}

// RemovePoint removes a point from the index
func (si *SpatialIndex) RemovePoint(id string) {
	si.mu.Lock()
	defer si.mu.Unlock()
	delete(si.points, id)
}

// FindNearby finds points within a radius
func (si *SpatialIndex) FindNearby(lat, lon, radiusKm float64) []Point {
	si.mu.RLock()
	defer si.mu.RUnlock()

	var result []Point
	for _, p := range si.points {
		dist := haversine(lat, lon, p.Lat, p.Lon)
		if dist <= radiusKm {
			result = append(result, p)
		}
	}
	return result
}

// FindNearest finds the nearest point
func (si *SpatialIndex) FindNearest(lat, lon float64) *Point {
	si.mu.RLock()
	defer si.mu.RUnlock()

	var nearest *Point
	minDist := 0.0
	first := true

	for _, p := range si.points {
		dist := haversine(lat, lon, p.Lat, p.Lon)
		if first || dist < minDist {
			p := p
			nearest = &p
			minDist = dist
			first = false
		}
	}

	return nearest
}

// haversine calculates the great-circle distance
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Earth radius in km
	dLat := (lat2 - lat1) * 3.14159265359 / 180.0
	dLon := (lon2 - lon1) * 3.14159265359 / 180.0
	a := sin(dLat/2)*sin(dLat/2) + cos(lat1*3.14159265359/180.0)*cos(lat2*3.14159265359/180.0)*sin(dLon/2)*sin(dLon/2)
	c := 2 * atan2(sqrt(a), sqrt(1-a))
	return R * c
}

func sin(x float64) float64 {
	// Simple sin approximation
	return x - x*x*x/6 + x*x*x*x*x/120
}

func cos(x float64) float64 {
	// Simple cos approximation
	return 1 - x*x/2 + x*x*x*x/24
}

func sqrt(x float64) float64 {
	// Simple sqrt
	if x <= 0 {
		return 0
	}
	guess := x / 2
	for i := 0; i < 10; i++ {
		guess = (guess + x/guess) / 2
	}
	return guess
}

func atan2(y, x float64) float64 {
	// Simple atan2
	if x > 0 {
		return atan(y / x)
	} else if x < 0 && y >= 0 {
		return atan(y/x) + 3.14159265359
	} else if x < 0 && y < 0 {
		return atan(y/x) - 3.14159265359
	} else if x == 0 && y > 0 {
		return 3.14159265359 / 2
	} else if x == 0 && y < 0 {
		return -3.14159265359 / 2
	}
	return 0
}

func atan(x float64) float64 {
	// Simple atan approximation
	return x - x*x*x/3 + x*x*x*x*x/5
}

// TimeSeriesIndex provides time-series indexing
type TimeSeriesIndex struct {
	series map[uint64]*SeriesInfo
	mu     sync.RWMutex
}

// SeriesInfo represents series information
type SeriesInfo struct {
	ID        uint64
	Labels    map[string]string
	MinTime   int64
	MaxTime   int64
	Count     int64
	LastValue float64
}

// NewTimeSeriesIndex creates a new time series index
func NewTimeSeriesIndex() *TimeSeriesIndex {
	return &TimeSeriesIndex{
		series: make(map[uint64]*SeriesInfo),
	}
}

// AddSeries adds a series
func (tsi *TimeSeriesIndex) AddSeries(id uint64, labels map[string]string) {
	tsi.mu.Lock()
	defer tsi.mu.Unlock()

	tsi.series[id] = &SeriesInfo{
		ID:     id,
		Labels: labels,
		MinTime: 0,
		MaxTime: 0,
		Count:  0,
	}
}

// UpdateSeries updates a series
func (tsi *TimeSeriesIndex) UpdateSeries(id uint64, timestamp int64, value float64) {
	tsi.mu.Lock()
	defer tsi.mu.Unlock()

	if series, ok := tsi.series[id]; ok {
		if series.MinTime == 0 || timestamp < series.MinTime {
			series.MinTime = timestamp
		}
		if timestamp > series.MaxTime {
			series.MaxTime = timestamp
		}
		series.Count++
		series.LastValue = value
	}
}

// GetSeries gets a series by ID
func (tsi *TimeSeriesIndex) GetSeries(id uint64) (*SeriesInfo, bool) {
	tsi.mu.RLock()
	defer tsi.mu.RUnlock()
	series, ok := tsi.series[id]
	return series, ok
}

// FindSeries finds series matching labels
func (tsi *TimeSeriesIndex) FindSeries(labels map[string]string) []uint64 {
	tsi.mu.RLock()
	defer tsi.mu.RUnlock()

	var result []uint64
	for id, series := range tsi.series {
		match := true
		for k, v := range labels {
			if series.Labels[k] != v {
				match = false
				break
			}
		}
		if match {
			result = append(result, id)
		}
	}
	return result
}

// GetSeriesCount returns the number of series
func (tsi *TimeSeriesIndex) GetSeriesCount() int {
	tsi.mu.RLock()
	defer tsi.mu.RUnlock()
	return len(tsi.series)
}

// RemoveSeries removes a series
func (tsi *TimeSeriesIndex) RemoveSeries(id uint64) {
	tsi.mu.Lock()
	defer tsi.mu.Unlock()
	delete(tsi.series, id)
}