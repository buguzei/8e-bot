package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func Start(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var forbiddenWords []string

	wordStage := make(map[int64]bool)

	msgID := make(map[int64]int)

	for update := range updates {
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.From.ID

			switch update.CallbackData() {
			case "cancel":
				delMsg := tgbotapi.NewDeleteMessage(chatID, msgID[chatID])

				_, _ = bot.Send(delMsg)

				wordStage[chatID] = false
			}
			continue
		}

		if update.Message != nil {
			if update.Message.Chat.IsPrivate() {
				chatID := update.Message.Chat.ID

				if wordStage[chatID] {
					forbiddenWords = append(forbiddenWords, strings.ToLower(update.Message.Text))
					wordStage[chatID] = false
					fmt.Println(forbiddenWords)
					continue
				}

				//chatID := update.Message.Chat.ID

				msgID[chatID] = update.Message.MessageID
				go DeleteForbiddenWord(bot, chatID, msgID[chatID], strings.ToLower(update.Message.Text), forbiddenWords)
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
				}
				continue
			}

			/*if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			}*/
		}

	}
}

func DeleteForbiddenWord(bot *tgbotapi.BotAPI, chatID int64, msgID int, text string, forbiddenWords []string) {
	for _, word := range forbiddenWords {
		if len(word) <= len(text) {
			fmt.Println(222)
			for i := 0; i < len(text)-len(word)+1; i++ {
				fmt.Println(i+len(word), len(text))
				if i+len(word) == len(text) {
					if text[i:] == word {

						delMsg := tgbotapi.NewDeleteMessage(chatID, msgID)

						_, err := bot.Send(delMsg)
						if err != nil {
							log.Println(err)
						}
					}
					continue
				}

				if text[i:len(word)+i] == word {

					delMsg := tgbotapi.NewDeleteMessage(chatID, msgID)

					_, err := bot.Send(delMsg)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
