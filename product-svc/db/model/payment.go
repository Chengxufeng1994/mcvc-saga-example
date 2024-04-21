package model

import "github.com/Chengxufeng1994/go-saga-example/common/model"

// Payment data model
type Payment struct {
	model.BaseModel
	UserID       uint64 `gorm:"index;not null"`
	CurrencyCode string `gorm:"not null"`
	Amount       int64  `gorm:"not null"`
}
