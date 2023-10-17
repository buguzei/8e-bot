package internal

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
)

func Start(bot *tgbotapi.BotAPI, admin, admin1, admin2 int64) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var forbiddenWords []string

	wordStage := make(map[int64]bool)

	msgID := make(map[int64]int)

	file, err := os.ReadFile("ForbiddenWords.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(file, &forbiddenWords)
	for update := range updates {
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.From.ID

			switch update.CallbackData() {
			case "cancel":
				delMsg := tgbotapi.NewDeleteMessage(chatID, msgID[chatID])

				_, err := bot.Send(delMsg)
				if err != nil {
					log.Println(err)
				}

				wordStage[chatID] = false
			}
			continue
		}

		if update.Message != nil {
			if update.Message.From.ID == admin || update.Message.From.ID == admin1 || update.Message.From.ID == admin2 {
				chatID := update.Message.Chat.ID

				if wordStage[chatID] {
					forbiddenWords = append(forbiddenWords, strings.ToLower(update.Message.Text))
					wordStage[chatID] = false
					file, err = json.MarshalIndent(forbiddenWords, "", "    ")
					if err != nil {
						log.Println(err)
					}
					err = os.WriteFile("ForbiddenWords.json", file, 0644)
					if err != nil {
						log.Println(err)
					}
					msg := tgbotapi.NewMessage(update.Message.From.ID, "Сообщение успешно добавленно в список запрещенных слов")
					_, _ = bot.Send(msg)
					continue
				}

				switch update.Message.Text {
				case "/newword":
					msg := tgbotapi.NewMessage(chatID, "Введите новое запрещенное слово:")
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel")),
					)

					msgConfig, err := bot.Send(msg)
					if err != nil {
						log.Println(err)
					}

					msgID[chatID] = msgConfig.MessageID

					wordStage[chatID] = true
					continue
				}
			}

			if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
				chatID := update.Message.Chat.ID
				msgID[chatID] = update.Message.MessageID
				go DeleteForbiddenWord(bot, chatID, msgID[chatID], strings.ToLower(update.Message.Text), forbiddenWords)
				if update.Message.Text == "/myid" {
					msg := tgbotapi.NewMessage(chatID, strconv.FormatInt(update.Message.From.ID, 10))
					_, _ = bot.Send(msg)
					continue
				}
			}
		}

	}
}

func DeleteForbiddenWord(bot *tgbotapi.BotAPI, chatID int64, msgID int, text string, forbiddenWords []string) {
	for _, word := range forbiddenWords {
		if len(word) <= len(text) {
			for i := 0; i < len(text)-len(word)+1; i++ {
				if i+len(word) == len(text) {
					if text[i:] == word {

						delMsg := tgbotapi.NewDeleteMessage(chatID, msgID)

						_, _ = bot.Send(delMsg)
					}
					continue
				}

				if text[i:len(word)+i] == word {

					delMsg := tgbotapi.NewDeleteMessage(chatID, msgID)

					_, _ = bot.Send(delMsg)
				}
			}
		}
	}
}
