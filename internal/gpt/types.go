package gpt

type ReviewResult struct {
	Category string `json:"category"`
	Reply    string `json:"reply"`
	Days     int    `json:"days"`
	Price    int    `json:"price"`
}
