package inlinequeries

import (
	database "camus/sanyavertolet/bot/pkg/database/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func Hints(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var articles []interface{}
	article := tgbotapi.NewInlineQueryResultArticleMarkdown(
		update.InlineQuery.ID,
		"no",
		"Чтобы показать информацию об игре, введите \"@ciykbot no\" и выберете подходящую игру.",
	)
	article.Description = "Помощь с командой \"no\""
	articles = append(articles, article)

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       articles,
	}

	if _, err := bot.Request(inlineConf); err != nil {
		log.Panic(err)
	}
}

func AboutGame(bot *tgbotapi.BotAPI, update tgbotapi.Update, repo *database.Repository) {
	games, err := repo.FindNextWeekGames()
	if err != nil {
		return
	}

	var articles []interface{}
	for _, game := range games {
		article := tgbotapi.NewInlineQueryResultArticle(
			update.InlineQuery.ID,
			game.Name,
			game.StringWithUsers(),
		)
		article.Description = fmt.Sprintf("%s в %s", game.Place, game.Date.Format("15:04 02.01"))
		articles = append(articles, article)
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       articles,
	}

	if _, err := bot.Request(inlineConf); err != nil {
		log.Panic(err)
	}
}
