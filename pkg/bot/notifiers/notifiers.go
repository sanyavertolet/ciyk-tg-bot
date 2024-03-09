package notifiers

import (
	keyboards "camus/sanyavertolet/bot/pkg/bot/keyboards"
	model "camus/sanyavertolet/bot/pkg/database/model"
	database "camus/sanyavertolet/bot/pkg/database/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func NotifyHeadOfQueue(bot *tgbotapi.BotAPI, repo *database.Repository, gameId uint) {
	user, game, err := repo.FindFirstQueuingUserByGameId(gameId)
	if err != nil {
		log.Print(err)
		return
	}
	message := fmt.Sprintf("Вы попадаете в основу!\n\n %s", game.String())
	keyboard := keyboards.WillComeKeyboard(game.ID)
	Notify(bot, message, user.ID, keyboard)
}

func NotifyGameAdded(bot *tgbotapi.BotAPI, repo *database.Repository, game model.Game) {
	message := fmt.Sprintf(
		"Новая игра!!\n\n * Что: %s\n * Где?: %s\n * Когда: %s",
		game.Name,
		game.Place,
		game.Date)

	keyboard := keyboards.SignKeyboard(game.ID)

	NotifyEverybody(bot, repo, message, keyboard)
}

func Notify(bot *tgbotapi.BotAPI, messageText string, id int64, keyboard *tgbotapi.InlineKeyboardMarkup) {
	message := tgbotapi.NewMessage(id, messageText)
	message.ReplyMarkup = keyboard
	
	if _, err := bot.Send(message); err != nil {
		log.Print("Cound not notify user")
	}
}

func NotifyEverybody(bot *tgbotapi.BotAPI, repo *database.Repository, message string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	ids, err := repo.FindAllUserIds()
	if err != nil {
		log.Panic("Could not fetch user ids")
		return
	}

	for i := 0; i < len(ids); i++ {
		Notify(bot, message, ids[i], keyboard)
	}
}