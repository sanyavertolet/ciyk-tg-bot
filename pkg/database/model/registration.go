package database

import "gorm.io/gorm"

type Registration struct {
	gorm.Model
	
	IsQueuing bool
	
	UserId int64
	User User `gorm:"foreignKey:UserId"`
	
	GameId uint
	Game Game `gorm:"foreignKey:GameId"`
}