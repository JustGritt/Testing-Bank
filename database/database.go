package database

import (
	"github.com/JustGritt/go-grpc/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func Connect() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Account{})
	Database = DbInstance{Db: db}
}

func ConnectTest() {
	db, err := gorm.Open(sqlite.Open("database_test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Account{})
	Database = DbInstance{Db: db}
}
