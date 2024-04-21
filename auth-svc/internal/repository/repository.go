package repository

import "gorm.io/gorm"

type GormOption func(*gorm.DB) error
