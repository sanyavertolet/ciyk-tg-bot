package database

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Game struct {
	gorm.Model
	
	Name string
	Place string
	Date time.Time
	MaxPlayers int
	Registrations []Registration `gorm:"foreignKey:GameId"`
}

func (game *Game) String() string {
	var stringBuilder strings.Builder
	
	stringBuilder.WriteString(fmt.Sprintf("Что: %s\n", game.Name))
	stringBuilder.WriteString(fmt.Sprintf("Где: %s\n", game.Place))
	stringBuilder.WriteString(fmt.Sprintf("Когда: %s\n", game.Date))
	
	return stringBuilder.String()
}
