package db

import (
	"errors"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/db/model"
	"gorm.io/gorm"
)

var ErrInvalidApplication = errors.New("Invalid application name")

// Migrator the database migrator instance
type Migrator struct {
	app string
	db  *gorm.DB
}

// NewMigrator returns a Migrator
func NewMigrator(app string, db *gorm.DB) *Migrator {
	return &Migrator{
		app: app,
		db:  db,
	}
}

// Migrate database migrate
func (m *Migrator) Migrate() error {
	switch m.app {
	case "order":
		return m.db.AutoMigrate(&model.Order{})
	case "payment":
		return m.db.AutoMigrate(&model.Payment{})
	case "product":
		return m.db.AutoMigrate(&model.Product{}, &model.Idempotency{})
	default:
		return ErrInvalidApplication
	}
}
