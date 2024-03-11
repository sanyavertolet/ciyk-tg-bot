package database

import (
	model "camus/sanyavertolet/bot/pkg/database/model"
	database "camus/sanyavertolet/bot/pkg/database/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(databaseFileName string) (*database.Repository, error) {
	db, err := gorm.Open(sqlite.Open(databaseFileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.User{}, &model.Game{}, &model.Registration{}, &model.Checkpoint{})
	if err != nil {
		return nil, err
	}
	return &database.Repository{DB: db}, nil
}
