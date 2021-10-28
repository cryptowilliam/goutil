package gjson

import (
	"bytes"
	"encoding/json"
)

func FormatIndent(p []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, p, "", "\t")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
