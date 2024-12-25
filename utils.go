package gobase

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(word string) string {
	snake := matchFirstCap.ReplaceAllString(word, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

func serializeStruct[T any](structData T) (string, error) {
	jsonData, err := json.Marshal(structData)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func deserializeStruct[T any](structData string) (T, error) {
	var result T
	// First, validate that the JSON matches expected structure
	var rawJSON map[string]interface{}
	if err := json.Unmarshal([]byte(structData), &rawJSON); err != nil {
		return result, err
	}
	// Try to unmarshal into the target type
	decoder := json.NewDecoder(strings.NewReader(structData))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&result); err != nil {
		return result, errors.New("invalid type structure")
	}

	return result, nil
}
