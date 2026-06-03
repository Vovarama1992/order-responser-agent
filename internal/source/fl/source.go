package fl

import "fl-agent/internal/model"

type Source struct{}

func NewSource() *Source {
	return &Source{}
}

func (s *Source) Name() string {
	return "fl"
}

func (s *Source) Fetch() ([]model.Order, error) {
	return ParseProjects()
}
