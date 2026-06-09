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

	result.Category = NormalizeCategory(result.Category)

	return result.Category, nil

}

func NormalizeCategory(category string) string {
	category = strings.TrimSpace(category)

	switch {
	case strings.HasPrefix(category, "а"):
		return "а) Чистая архитектура для благородного дона"

	case strings.HasPrefix(category, "б"):
		return "б) Для дона с бодуна"

	case strings.HasPrefix(category, "в"):
		return "в) Для понурого дона"

	case strings.HasPrefix(category, "г"):
		return "г) Хлам"

	case strings.HasPrefix(category, "д"):
		return "д) Не моя специализация"

	default:
		return category
	}
}
