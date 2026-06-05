package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

type FilterCategoryResult struct {
	Category string `json:"category"`
}

func (c *Client) Filter(
	title string,
	budget string,
	description string,
) (string, error) {

	rules, err := LoadResponseRules()
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`
Ниже инструкция классификации заказов.

Верни только JSON:

{
  "category": ""
}

Никаких пояснений.

ИНСТРУКЦИЯ:

%s

ЗАКАЗ

Заголовок:
%s

Бюджет:
%s

Описание:
%s
`,
		rules,
		title,
		budget,
		description,
	)

	resp, err := c.client.Responses.New(
		context.Background(),
		responses.ResponseNewParams{
			Model: shared.ResponsesModel("gpt-4.1-nano"),
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(prompt),
			},
		},
	)
	if err != nil {
		return "", err
	}

	raw := strings.TrimSpace(resp.OutputText())

	var result FilterCategoryResult

	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return "", err
	}

	return result.Category, nil
}
