package main

import (
	"e-wallet/config"
	"e-wallet/controllers"
	"e-wallet/routes"
	"e-wallet/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"fmt"
	"log"
	"os"
)

var db *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	var errConn error
	db, errConn = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errConn != nil {
		log.Fatal("failed to connect to the database:", errConn)
	} else {
		log.Println("Database connected successfully")
	}

	walletService := services.NewEWalletSystem(db)

	controllers.InitWalletService(walletService)

	r := gin.Default()
	config.ConnectDB()
	routes.SetupRoutes(r)
	r.Run(":8081")
}
