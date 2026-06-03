package gpt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

func (c *Client) Review(
	title string,
	budget string,
	description string,
) (*ReviewResult, error) {

	rules, err := LoadResponseRules()
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`
Ниже будет инструкция по классификации заказов и написанию откликов.

Строго следуй этой инструкции.

В ответ верни ТОЛЬКО валидный JSON без markdown, без пояснений и без лишнего текста.

Формат ответа:

{
  "category": "",
  "reply": "",
  "days": 0,
  "price": 0
}

Где:

category - категория заказа согласно инструкции.

reply - готовый отклик. Начинается со слова "Здравствуйте". Категорию заказа внутрь отклика не вставляй.

days - рекомендуемый срок выполнения в днях.

price - рекомендуемая стоимость в рублях.

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
			Model: shared.ResponsesModel("gpt-4.1"),
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(prompt),
			},
		},
	)

	if err != nil {
		return nil, err
	}

	var result ReviewResult

	err = json.Unmarshal(
		[]byte(resp.OutputText()),
		&result,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
