package upwork

import (
	"encoding/json"
	"os/exec"

	"fl-agent/internal/model"
)

func ParseJobs() ([]model.Order, error) {
	out, err := exec.Command(
		"node",
		"internal/source/upwork/parser.js",
	).Output()

	if err != nil {
		return nil, err
	}

	var orders []model.Order

	if err := json.Unmarshal(out, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
