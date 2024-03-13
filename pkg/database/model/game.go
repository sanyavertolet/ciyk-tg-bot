package database

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Game struct {
	gorm.Model
	Name               string
	Place              string
	Date               time.Time
	MaxPlayers         int
	IsRegistrationOpen bool
	Registrations      []Registration `gorm:"foreignKey:GameId"`
	Users              []User         `gorm:"many2many:registrations;"`
}

func (game *Game) String() string {
	dateString := game.Date.Format("15:04 02.01")
	return fmt.Sprintf("%s: %s в %s\n\n", game.Name, game.Place, dateString)
}

func (game *Game) StringWithUsers() string {
	var stringBuilder strings.Builder

	dateString := game.Date.Format("15:04 02.01")
	stringBuilder.WriteString(fmt.Sprintf("%s\n%s в %s\n\n", game.Name, game.Place, dateString))

	if !game.IsRegistrationOpen {
		stringBuilder.WriteString("Запись на эту игру откроется в ближейшее к игре воскресенье в 22:00")
	} else if len(game.Users) == 0 {
		stringBuilder.WriteString("На эту игру пока никто не записался.")
	} else {
		for i, user := range game.Users {
			if i == game.MaxPlayers {
				stringBuilder.WriteString("\nРезерв:\n")
			}
			userTag := fmt.Sprintf("@%s", user.Tag)
			if user.Tag == "" {
				// fixMe: tag users without telegram tag
				userTag = user.Name
			}
			stringBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, userTag))
		}
	}

	return stringBuilder.String()
}
