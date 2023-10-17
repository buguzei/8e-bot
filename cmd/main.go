package main

import (
	"github.com/buguzei/8e-bot/internal"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	admin, _ := strconv.Atoi(os.Getenv("ADMIN"))
	admin1, _ := strconv.Atoi(os.Getenv("ADMIN1"))
	admin2, _ := strconv.Atoi(os.Getenv("ADMIN2"))

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	internal.Start(bot, int64(admin), int64(admin1), int64(admin2))
}
