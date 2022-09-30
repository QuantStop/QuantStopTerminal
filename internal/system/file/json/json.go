package json

import (
	"encoding/json"
	"io"
)

// PrettyEncodeJson writes data to .json file in a "pretty", indented and readable way.
// Returns error if there was a problem encoding the data.
func PrettyEncodeJson(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}
