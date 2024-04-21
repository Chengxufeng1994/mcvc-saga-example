package model

import "github.com/Chengxufeng1994/go-saga-example/common/model"

// user data model
type User struct {
	model.BaseModel
	Active      bool   `gorm:"default:true"`
	FirstName   string `gorm:"type:varchar(50);not null"`
	LastName    string `gorm:"type:varchar(50);not null"`
	Email       string `gorm:"type:varchar(320);unique;not null"`
	Address     string `gorm:"type:text;not null"`
	PhoneNumber string `gorm:"type:varchar(20);unique;not null"`
	Password    string `gorm:"type:varchar(100);not null"`
}
