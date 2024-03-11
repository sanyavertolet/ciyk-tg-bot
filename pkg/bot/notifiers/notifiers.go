package notifiers

import (
	keyboards "camus/sanyavertolet/bot/pkg/bot/keyboards"
	model "camus/sanyavertolet/bot/pkg/database/model"
	database "camus/sanyavertolet/bot/pkg/database/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func Notify(bot *tgbotapi.BotAPI, messageText string, id int64, keyboard *tgbotapi.InlineKeyboardMarkup) {
	message := tgbotapi.NewMessage(id, messageText)
	message.ReplyMarkup = keyboard

	if _, err := bot.Send(message); err != nil {
		log.Print("Cound not notify user")
	}
}

func NotifyUsers(bot *tgbotapi.BotAPI, userIds []int64, message string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	for i := 0; i < len(userIds); i++ {
		Notify(bot, message, userIds[i], keyboard)
	}
}

func NotifyEverybody(bot *tgbotapi.BotAPI, repo *database.Repository, message string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	ids, err := repo.FindAllUserIds()
	if err != nil {
		log.Panic(err)
	}
	NotifyUsers(bot, ids, message, keyboard)
}

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

	keyboard := keyboards.SignKeyboard(game)

	NotifyEverybody(bot, repo, message, keyboard)
}

func NotifyEverybodyGamesAdded(bot *tgbotapi.BotAPI, repo *database.Repository, games []model.Game) {
	userIds, err := repo.FindAllUserIds()
	if err != nil {
		log.Panic(err)
	}

	for _, game := range games {
		message := fmt.Sprintf(
			"Открыта запись на игру!\n\n * Что: %s\n * Где?: %s\n * Когда: %s",
			game.Name,
			game.Place,
			game.Date,
		)

		keyboard := keyboards.SignKeyboard(game)

		NotifyUsers(bot, userIds, message, keyboard)
	}
}
