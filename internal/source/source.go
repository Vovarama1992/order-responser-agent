package source

import "fl-agent/internal/model"

type Source interface {
	Name() string
	Fetch() ([]model.Order, error)
}
