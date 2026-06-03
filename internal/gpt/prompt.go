package gpt

import (
	"os"
)

func LoadResponseRules() (string, error) {
	data, err := os.ReadFile("response_rules.md")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
