package utils

import (
	keyboards "camus/sanyavertolet/bot/pkg/bot/keyboards"
	"fmt"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func ShowMenu(bot *tgapi.BotAPI, chatId int64) {
	messageText := fmt.Sprintf("Чат-бот команды Камю ин ё Кант Ереван")
	message := tgapi.NewMessage(chatId, messageText)
	message.ReplyMarkup = keyboards.MainMenuKeyboard()

	if _, err := bot.Send(message); err != nil {
		log.Print("Could not send message to user")
		return
	}
}
