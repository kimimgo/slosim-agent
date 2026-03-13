package tools

import (
	"encoding/json"
	"strconv"
	"strings"
)

// UnmarshalToolInput is a lenient JSON unmarshaler for LLM tool call inputs.
// Some models (e.g., LLaMA, Gemma) send numbers as quoted strings: "4" instead of 4.
// This function tries standard unmarshal first, and on type-mismatch failure,
// normalizes string values that look like numbers before retrying.
func UnmarshalToolInput(input string, v interface{}) error {
	err := json.Unmarshal([]byte(input), v)
	if err == nil {
		return nil
	}

	// Only attempt normalization for type-mismatch errors
	if !strings.Contains(err.Error(), "cannot unmarshal string") {
		return err
	}

	// Parse into raw map and convert string-numbers to bare numbers
	var raw map[string]json.RawMessage
	if jsonErr := json.Unmarshal([]byte(input), &raw); jsonErr != nil {
		return err
	}

	changed := false
	for key, val := range raw {
		s := string(val)
		if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
			numStr := strings.TrimSpace(s[1 : len(s)-1])
			if _, parseErr := strconv.ParseFloat(numStr, 64); parseErr == nil {
				raw[key] = json.RawMessage(numStr)
				changed = true
			}
		}
	}

	if !changed {
		return err
	}

	normalized, marshalErr := json.Marshal(raw)
	if marshalErr != nil {
		return err
	}

	return json.Unmarshal(normalized, v)
}
