package main

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

// создание и подключение SQLite бд
func connectToSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// InitDb инициализация бд (в т.ч. миграции)
func InitDb() (db *gorm.DB) {

	db, err := connectToSQLite()
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&MoneyMovementType{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&Category{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&MoneyMovement{})
	if err != nil {
		return nil
	}
	mmTypes := []MoneyMovementType{
		MoneyMovementType{Name: "ДОХОД"},
		MoneyMovementType{Name: "РАСХОД"},
		MoneyMovementType{Name: "КОПИЛКА"},
		MoneyMovementType{Name: "-КОПИЛКА"},
	}
	db.Create(mmTypes)
	return
}

// User methods

// CreateUser создание пользователя в бд
func CreateUser(db *gorm.DB, telegramId int64, balance int, username string) (result *gorm.DB) {
	user := User{TelegramID: telegramId, Balance: balance, Username: username}
	result = db.Create(&user)
	return
}

// GetUserByTelegramId получение пользователя по телеграм ID
func GetUserByTelegramId(db *gorm.DB, telegramId int64) User {
	var user User
	db.First(&user, "telegram_id = ?", telegramId)
	return user
}

// GetUserById получение пользователя по ID
func GetUserById(db *gorm.DB, id uint) User {
	var user User
	db.First(&user, id)
	return user
}

// GetUserByUsername получение пользователя по его username
func GetUserByUsername(db *gorm.DB, username string) User {
	var user User
	db.First(&user, "username = ?", username)
	return user
}

// UpdateUserBalance инкремент баланса пользователя (текущий баланс + добавляемый баланс),
// либо декремент, для этого нужно использовать отрицательный параметр balance
func UpdateUserBalance(db *gorm.DB, user User, balance int) User {
	db.Model(&user).Update("Balance", user.Balance+balance)
	return user
}

// DeleteUser удаление пользователя
func DeleteUser(db *gorm.DB, id uint) *gorm.DB {
	var user User
	result := db.Delete(&user, id)
	return result
}

// Category methods

// CreateCategory создание категории в бд
func CreateCategory(
	db *gorm.DB,
	user User,
	name string,
	moneyMovementType MoneyMovementType,
) (result *gorm.DB, category Category) {
	category = Category{User: user, Name: name, MoneyMovementTypeID: moneyMovementType.ID}
	result = db.Create(&category)
	return
}

// GetCategoryById получение категории по ID
func GetCategoryById(db *gorm.DB, id uint) Category {
	var category Category
	db.First(&category, id)
	return category
}

// GetCategoryByName получение категории по его username
func GetCategoryByName(db *gorm.DB, name string) Category {
	var category Category
	result := db.First(&category, "name = ?", name)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return category
	}
	return category
}

// DeleteCategory удаление категории
func DeleteCategory(db *gorm.DB, id uint) *gorm.DB {
	var category Category
	result := db.Delete(&category, id)
	return result
}

// MoneyMovement methods

// CreateMoneyMovement создание денежного передвижения в бд
func CreateMoneyMovement(
	db *gorm.DB,
	user User,
	quantity int64,
	category Category,
) (result *gorm.DB) {
	moneyMovement := MoneyMovement{User: user, Quantity: quantity, CategoryID: category.ID}
	result = db.Create(&moneyMovement)
	return
}

// GetMoneyMovementById получение денежного передвижения по ID
func GetMoneyMovementById(db *gorm.DB, id uint) MoneyMovement {
	var moneyMovement MoneyMovement
	db.First(&moneyMovement, id)
	return moneyMovement
}

// DeleteMoneyMovement удаление денежного передвижения
func DeleteMoneyMovement(db *gorm.DB, id uint) *gorm.DB {
	var moneyMovement MoneyMovement
	result := db.Delete(&moneyMovement, id)
	return result
}

// MoneyMovementType methods

// GetMoneyMovementTypeByName получение категории по его username
func GetMoneyMovementTypeByName(db *gorm.DB, name string) MoneyMovementType {
	var moneyMovementType MoneyMovementType
	result := db.First(&moneyMovementType, "name = ?", name)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return moneyMovementType
	}
	return moneyMovementType
}
