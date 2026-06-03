package gpt

import (
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Client struct {
	client openai.Client
}

func New() *Client {
	client := openai.NewClient(
		option.WithAPIKey(
			os.Getenv("KEY"),
		),
	)

	return &Client{
		client: client,
	}
}
