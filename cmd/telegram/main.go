package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if chatID == "" {
		log.Fatal("TELEGRAM_CHAT_ID environment variable is not set")
	}
	tgBot := NewTelegramBot(token, chatID)

	err := tgBot.SendMessage("yeah baby yeah!!")
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	fmt.Println("message sent successfully!")

}
