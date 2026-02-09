package output

import (
	"io"

	"gopkg.in/yaml.v3"
)

func printYAML(w io.Writer, data any) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer enc.Close()
	return enc.Encode(data)
}
