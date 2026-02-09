package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
)

func printTable(w io.Writer, data any) error {
	rows, err := toRows(data)
	if err != nil {
		return printJSON(w, data)
	}

	if len(rows) == 0 {
		fmt.Fprintln(w, "No results.")
		return nil
	}

	// Collect all keys maintaining a stable order
	keySet := make(map[string]bool)
	for _, row := range rows {
		for k := range row {
			keySet[k] = true
		}
	}
	var headers []string
	for k := range keySet {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	// Header
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	// Separator
	var seps []string
	for _, h := range headers {
		seps = append(seps, strings.Repeat("-", len(h)))
	}
	fmt.Fprintln(tw, strings.Join(seps, "\t"))

	// Rows
	for _, row := range rows {
		var vals []string
		for _, h := range headers {
			vals = append(vals, truncate(fmt.Sprintf("%v", row[h]), 60))
		}
		fmt.Fprintln(tw, strings.Join(vals, "\t"))
	}

	return tw.Flush()
}

// unmarshalWithNumbers uses json.Decoder with UseNumber() to preserve integer precision.
func unmarshalWithNumbers(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(v)
}

func toRows(data any) ([]map[string]any, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Try as array first
	var rows []map[string]any
	if err := unmarshalWithNumbers(b, &rows); err == nil && len(rows) > 0 {
		return rows, nil
	}

	// Try as single object
	var single map[string]any
	if err := unmarshalWithNumbers(b, &single); err == nil {
		return []map[string]any{flattenMap(single, "")}, nil
	}

	// Try reflection for struct slices
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		var result []map[string]any
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			b, _ := json.Marshal(item)
			var m map[string]any
			if unmarshalWithNumbers(b, &m) == nil {
				result = append(result, flattenMap(m, ""))
			}
		}
		if len(result) > 0 {
			return result, nil
		}
	}

	return nil, fmt.Errorf("cannot convert to table rows")
}

func flattenMap(m map[string]any, prefix string) map[string]any {
	result := make(map[string]any)
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]any:
			for fk, fv := range flattenMap(val, key) {
				result[fk] = fv
			}
		default:
			result[key] = val
		}
	}
	return result
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}
