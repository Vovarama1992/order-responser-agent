package freelance

import "fl-agent/internal/model"

type Source struct{}

func NewSource() *Source {
	return &Source{}
}

func (s *Source) Name() string {
	return "freelance"
}

func (s *Source) Fetch() ([]model.Order, error) {
	return ParseProjects()
}
