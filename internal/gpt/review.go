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

reply - готовый отклик. Начинается со слова "Здравствуйте". Категорию заказа внутрь отклика не вставляй. Максимум 3200 символов.

days - рекомендуемый срок выполнения в днях.

price - рекомендуемая стоимость в рублях.

Если заказ попадает в категорию "Не моя специализация", верни:
{
  "category": "Не моя специализация",
  "reply": "",
  "days": 0,
  "price": 0
}

Если заказ вообще не связан с IT, разработкой ПО, автоматизацией, интеграциями, backend, frontend, ботами, AI, парсингом, аналитикой, CRM, сайтами или техническими системами — это "Не моя специализация".

Не пытайся притягивать к IT заказы по ландшафтному дизайну, архитектуре зданий, интерьеру, 3D-рендерам, бухгалтерии, юридическим услугам, отзывам, копирайтингу, презентациям, меню, полиграфии, графическому дизайну и похожим не-IT задачам.

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

	raw := strings.TrimSpace(resp.OutputText())

	fmt.Println("========== GPT RAW ==========")
	fmt.Println(raw)
	fmt.Println("======== END GPT RAW ========")

	var result ReviewResult

	err = json.Unmarshal([]byte(raw), &result)
	if err != nil {
		return nil, fmt.Errorf("json parse error: %w\nraw:\n%s", err, raw)
	}

	return &result, nil
}
