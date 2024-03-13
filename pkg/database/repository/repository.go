package database

import (
	"gorm.io/gorm"
	"sync"
)

type Repository struct {
	DB    *gorm.DB
	mutex sync.Mutex
}
