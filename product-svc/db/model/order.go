package model

import "github.com/Chengxufeng1994/go-saga-example/common/model"

// Order data model
type Order struct {
	model.BaseModel
	ProductID uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"not null"`
	Amount    int64  `gorm:"not null"`
}
