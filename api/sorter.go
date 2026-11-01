// Package api provides sorting helpers
package api

import (
	"fmt"
	"sort"
	"strings"
)

// SortParams holds sort params
type SortParams struct {
	Field string
	Order string
}

// SortStrings sorts strings
func SortStrings(items []string, order string) {
	if order == "desc" {
		sort.Sort(sort.Reverse(sort.StringSlice(items)))
	} else {
		sort.Strings(items)
	}
}

// SortByField sorts slice by field
func SortByField(items []map[string]interface{}, field, order string) {
	sort.Slice(items, func(i, j int) bool {
		vi := fmt.Sprintf("%v", items[i][field])
		vj := fmt.Sprintf("%v", items[j][field])
		vi = strings.ToLower(vi)
		vj = strings.ToLower(vj)
		if order == "desc" {
			return vi > vj
		}
		return vi < vj
	})
}
