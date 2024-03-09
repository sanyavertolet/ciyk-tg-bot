package bot

import (
	model "camus/sanyavertolet/bot/pkg/database/model"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func WillComeKeyboard(gameId uint) *tgbotapi.InlineKeyboardMarkup {
	will := tgbotapi.NewInlineKeyboardButtonData("Я приду!", fmt.Sprintf("willCome %d", gameId))
	wont := tgbotapi.NewInlineKeyboardButtonData("Я не приду! :c", fmt.Sprintf("wontCome %d", gameId))
	signRow := tgbotapi.NewInlineKeyboardRow(will, wont)
	toMenu := tgbotapi.NewInlineKeyboardButtonData("К меню", "tomenu")
	toMenuRow := tgbotapi.NewInlineKeyboardRow(toMenu)

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{signRow, toMenuRow}}
}

func SignKeyboard(gameId uint) *tgbotapi.InlineKeyboardMarkup {
	sign := tgbotapi.NewInlineKeyboardButtonData("Записаться", fmt.Sprintf("sign %d", gameId))
	signRow := tgbotapi.NewInlineKeyboardRow(sign)
	toMenu := tgbotapi.NewInlineKeyboardButtonData("К меню", "tomenu")
	toMenuRow := tgbotapi.NewInlineKeyboardRow(toMenu)

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{signRow, toMenuRow}}
}

func UnsignKeyboard(gameId uint) *tgbotapi.InlineKeyboardMarkup {
	unsign := tgbotapi.NewInlineKeyboardButtonData("Выписаться", fmt.Sprintf("unsign %d", gameId))
	unsignRow := tgbotapi.NewInlineKeyboardRow(unsign)
	toMenu := tgbotapi.NewInlineKeyboardButtonData("К меню", "tomenu")
	toMenuRow := tgbotapi.NewInlineKeyboardRow(toMenu)

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{unsignRow, toMenuRow}}
} 

func GoToMainMenuKeyboard() *tgbotapi.InlineKeyboardMarkup {
	toMenu := tgbotapi.NewInlineKeyboardButtonData("К меню", "tomenu")
	toMenuRow := tgbotapi.NewInlineKeyboardRow(toMenu)

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{toMenuRow}}
}

func MainMenuKeyboard() *tgbotapi.InlineKeyboardMarkup {
	future := tgbotapi.NewInlineKeyboardButtonData("Предстоящие игры", "futureGames 0")
	past := tgbotapi.NewInlineKeyboardButtonData("Прошедшие игры", "pastGames 0")

	myFuture := tgbotapi.NewInlineKeyboardButtonData("Мои предстоящие игры", "myFutureGames 0")
	myPast := tgbotapi.NewInlineKeyboardButtonData("Мои прошедшие игры", "myPastGames 0")

	games := tgbotapi.NewInlineKeyboardRow(future, past)
	myGames := tgbotapi.NewInlineKeyboardRow(myFuture, myPast)

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{games, myGames}}
}

func GamesMenuKeyboard(games []model.Game, page int, callbackName string) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for i := range pageSize {
		gameIndex := page * pageSize + i
		if gameIndex < len(games) {
			text := games[gameIndex].Name
			data := "game " + strconv.Itoa(int(games[gameIndex].ID))
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(text, data)))
		}
	}
	
	var paginationButtons []tgbotapi.InlineKeyboardButton
	if page > 0 {
		paginationButtons = append(
			paginationButtons,
			tgbotapi.NewInlineKeyboardButtonData("<", callbackName + " " + strconv.Itoa(page - 1)),
		)
	}
	paginationButtons = append(
		paginationButtons,
		tgbotapi.NewInlineKeyboardButtonData("К меню", "tomenu"),
	)
	if (page + 1) * pageSize < len(games) {
		paginationButtons = append(
			paginationButtons,
			tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf("%s %d", callbackName, page + 1)),
		)
	}
	
	paginationRow := tgbotapi.NewInlineKeyboardRow(paginationButtons...)

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: append(rows, paginationRow)}
}

const pageSize = 5
