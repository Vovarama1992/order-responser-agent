package main

import (
	"fmt"

	"fl-agent/config"
	"fl-agent/internal/gpt"
)

func main() {
	config.LoadEnv()

	client := gpt.New()

	result, err := client.Review(
		"Backend-разработчик для Telegram mini app",
		"",
		"О нас. Разрабатываем Telegram mini app на TON. Нужен backend-разработчик для развития продукта.",
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("CATEGORY:")
	fmt.Println(result.Category)

	fmt.Println()

	fmt.Println("REPLY:")
	fmt.Println(result.Reply)
}
