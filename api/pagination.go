// Package api provides pagination helpers
package api

import (
	"net/http"
	"strconv"
)

// PaginationParams holds pagination params
type PaginationParams struct {
	Page     int
	PageSize int
	Offset   int
	Limit    int
}

// ParsePagination parses pagination from request
func ParsePagination(r *http.Request) PaginationParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
		Limit:    pageSize,
	}
}

// Paginate slices a slice
func Paginate[T any](items []T, p PaginationParams) []T {
	start := p.Offset
	if start > len(items) {
		start = len(items)
	}
	end := start + p.Limit
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

// PaginationResponse wraps paginated response
type PaginationResponse struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginationResponse creates response
func NewPaginationResponse(items interface{}, total, page, pageSize int) PaginationResponse {
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}
	return PaginationResponse{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
