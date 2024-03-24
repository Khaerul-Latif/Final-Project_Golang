package database

import (
	"final-project/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	errs := godotenv.Load(".env")
	if errs != nil {
		log.Fatalf("Some error occured. Err: %s", errs)
	}
	config := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"),os.Getenv("DBNAME"),os.Getenv("PASSWORD"), os.Getenv("SSLMODE"))
	con := config

	db, err = gorm.Open(postgres.Open(con), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	fmt.Println("Successfully connected to database")
	db.Debug().AutoMigrate(&models.User{}, &models.Photo{}, &models.SocialMedia{}, &models.Comment{})
}

func GetDB() *gorm.DB {
	return db
}
