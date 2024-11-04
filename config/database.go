package config

import (
	"EcommerceSederhana/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/go_ecom?parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// db.Migrator().DropTable(&models.User{}, &models.Product{}, &models.Transaction{}, &models.Cart{}, &models.TransactionItem{})

	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Transaction{}, &models.Cart{}, &models.TransactionItem{})

	log.Println("Connect Database")

	DB = db
}
