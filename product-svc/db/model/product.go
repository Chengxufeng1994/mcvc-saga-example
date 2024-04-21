package model

import "github.com/Chengxufeng1994/go-saga-example/common/model"

// product data model
type Product struct {
	model.BaseModel
	Name        string `gorm:"type:varchar(256);not null"`
	Description string `gorm:"type:text;not null"`
	BrandName   string `gorm:"type:varchar(256);not null"`
	Inventory   int64  `gorm:"not null"`
	Price       int64  `gorm:"not null"`
}

// Idempotency data model
type Idempotency struct {
	model.BaseModel
	ProductID  uint64 `gorm:"primaryKey"`
	Amount     int64  `gorm:"not null"`
	Rollbacked bool   `gorm:"not null"`
}
