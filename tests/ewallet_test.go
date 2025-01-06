// tests/ewallet_test.go

package tests

import (
	"e-wallet/models"
	"e-wallet/services"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupDB() *gorm.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	return db
}

func TestConcurrentCredit(t *testing.T) {
	db := setupDB()

	ws := services.NewEWalletSystem(db)

	user := models.User{Username: "roy", Balance: 0}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	userID := user.ID

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				balance, transactionID, err := ws.Credit(userID, 1000) //user value sum tertentu
				if err != nil {
					t.Errorf("Error crediting: %v", err)
				}
				t.Logf("Transaction ID: %d, New Balance: %f", transactionID, balance)
			}
		}()
	}

	wg.Wait()

	var balance float64
	if err := db.Table("users").Where("id = ?", userID).Select("balance").Scan(&balance).Error; err != nil {
		t.Errorf("Error fetching balance: %v", err)
	}

	expectedBalance := 10000000.0
	if balance != expectedBalance {
		t.Errorf("Expected balance %f, got %f", expectedBalance, balance)
	}
}
