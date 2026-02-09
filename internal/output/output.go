package output

import (
	"fmt"
	"io"
)

type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
	FormatYAML  Format = "yaml"
)

// Print outputs data in the specified format.
func Print(w io.Writer, format Format, data any) error {
	switch format {
	case FormatJSON:
		return printJSON(w, data)
	case FormatTable:
		return printTable(w, data)
	case FormatYAML:
		return printYAML(w, data)
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
}
