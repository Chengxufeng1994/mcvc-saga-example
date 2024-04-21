package db

import (
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/db/model"
	"gorm.io/gorm"
)

// Migrator the database migrator instance
type Migrator struct {
	db *gorm.DB
}

// NewMigrator returns a Migrator
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db}
}

// Migrate database migrate
func (m *Migrator) Migrate() error {
	return m.db.AutoMigrate(&model.User{})
}
