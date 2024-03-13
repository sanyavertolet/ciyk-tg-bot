package database

import (
	model "camus/sanyavertolet/bot/pkg/database/model"
	database "camus/sanyavertolet/bot/pkg/database/repository"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(databaseFileName string) (*database.Repository, error) {
	dsn := fmt.Sprintf("%s?_busy_timeout=5000", databaseFileName)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(
		&model.User{},
		&model.Game{},
		&model.Registration{},
		&model.Checkpoint{},
	); err != nil {
		return nil, err
	}
	return &database.Repository{DB: db}, nil
}
