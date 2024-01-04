package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramID int64  `gorm:"size:255;uniqueIndex"`
	Balance    int    // User balance
	Username   string // A regular string field

}

type MoneyMovementType struct {
	gorm.Model
	Name string `gorm:"size:255;uniqueIndex"`
}

type Category struct {
	gorm.Model

	UserID uint
	User   User

	Name string `gorm:"size:255;uniqueIndex"`

	MoneyMovementTypeID uint
	MoneyMovementType   MoneyMovementType
}

type MoneyMovement struct {
	gorm.Model

	UserID uint
	User   User

	Quantity int64

	CategoryID uint
	Category   Category
}
