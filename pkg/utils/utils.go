package utils

import (
	"encoding/json"
	"io"
	"regexp"
)

func DecodeBody[customType any](body io.ReadCloser) (*customType, error) {
	var data customType
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func IsUUID(s string) bool {
	guidRegex := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	ok, _ := regexp.MatchString(guidRegex, s)
	return ok
}
