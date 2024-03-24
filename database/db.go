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
	host     = os.Getenv("PGHOST")
	dbPort   = os.Getenv("PGPORT")
	user     = os.Getenv("PGUSER")
	dbName   = os.Getenv("PGDATABASE")
	password = os.Getenv("PGPASSWORD")
	db       *gorm.DB
	err      error
)

func StartDB() {
	errs := godotenv.Load(".env")
	if errs != nil {
		log.Fatalf("Some error occured. Err: %s", errs)
	}
	config := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host     ,dbPort   ,user     ,dbName   ,password )
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
