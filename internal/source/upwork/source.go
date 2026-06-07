package upwork

import "fl-agent/internal/model"

type Source struct{}

func NewSource() *Source {
	return &Source{}
}

func (s *Source) Name() string {
	return "upwork"
}

func (s *Source) Fetch() ([]model.Order, error) {
	return ParseJobs()
}

func (s *Source) Parse(order model.Order) (model.Order, error) {
	return order, nil
}
