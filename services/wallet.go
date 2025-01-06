package services

import (
	"errors"
	"log"
	"sync"

	"e-wallet/models"

	"gorm.io/gorm"
)

var DB *gorm.DB

type EWalletSystem struct {
	DB *gorm.DB
	mu sync.Mutex //thread-safe
}

func GetDB() *gorm.DB {
	return DB
}

func NewEWalletSystem(db *gorm.DB) *EWalletSystem {
	log.Println("Creating new EWalletSystem")
	return &EWalletSystem{DB: db}
}

func (ew *EWalletSystem) Credit(userID int, amount float64) (float64, int, error) {
	if amount <= 0 {
		return 0, 0, errors.New("invalid amount")
	}

	// thread-safe operation
	ew.mu.Lock()
	defer ew.mu.Unlock()

	tx := ew.DB.Begin()
	if tx.Error != nil {
		return 0, 0, tx.Error
	}

	err := tx.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	var balance float64
	err = tx.Model(&models.User{}).Where("id = ?", userID).Select("balance").Scan(&balance).Error
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	err = tx.Exec("INSERT INTO transactions (user_id, amount, type) VALUES (?, ?, ?)", userID, amount, "credit").Error
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	var lastTransaction struct {
		ID int `json:"id"`
	}
	err = tx.Raw("SELECT id FROM transactions WHERE user_id = ? ORDER BY id DESC LIMIT 1", userID).Scan(&lastTransaction).Error
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	tx.Commit()

	return balance, lastTransaction.ID, nil
}

func (ew *EWalletSystem) Debit(userID int, amount float64) (float64, int, error) {
	if amount <= 0 {
		return 0, 0, errors.New("invalid amount")
	}

	// thread-safe operation
	ew.mu.Lock()
	defer ew.mu.Unlock()

	tx := ew.DB.Begin()
	if tx.Error != nil {
		return 0, 0, tx.Error
	}

	var balance float64
	err := tx.Model(&models.User{}).Where("id = ?", userID).Select("balance").Scan(&balance).Error
	if err != nil || balance < amount {
		tx.Rollback()
		return 0, 0, errors.New("insufficient funds")
	}

	err = tx.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	err = tx.Model(&models.User{}).Where("id = ?", userID).Select("balance").Scan(&balance).Error
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	err = tx.Exec("INSERT INTO transactions (user_id, amount, type) VALUES (?, ?, ?)", userID, amount, "debit").Error
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	var lastTransaction struct {
		ID int `json:"id"`
	}
	err = tx.Raw("SELECT id FROM transactions WHERE user_id = ? ORDER BY id DESC LIMIT 1", userID).Scan(&lastTransaction).Error
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	tx.Commit()

	return balance, lastTransaction.ID, nil
}
