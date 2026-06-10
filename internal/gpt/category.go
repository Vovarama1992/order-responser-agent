package gpt

import "strings"

var categoryMarkers = []string{
	"благородного дона",
	"дона с бодуна",
	"понурого дона",
	"хлам",
	"не моя специализация",
}

func IsKnownCategory(category string) bool {
	category = strings.ToLower(strings.TrimSpace(category))

	for _, marker := range categoryMarkers {
		if strings.Contains(category, marker) {
			return true
		}
	}

	return false
}
