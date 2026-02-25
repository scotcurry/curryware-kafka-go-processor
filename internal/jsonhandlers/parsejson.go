package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/logging"
	"encoding/base64"
	"encoding/json"
)

// ParseJSON decodes a base64-encoded JSON string into a value of type T.
func ParseJSON[T any](encoded string) (T, error) {
	var result T

	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		logging.LogError("Error decoding base64: ", err)
		return result, err
	}

	err = json.Unmarshal(decodedBytes, &result)
	if err != nil {
		logging.LogError("Error parsing JSON: ", err)
		return result, err
	}

	return result, nil
}
