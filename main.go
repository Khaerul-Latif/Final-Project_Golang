package main

import (
	"final-project/database"
	_ "final-project/docs"
	"final-project/router"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// @title Final Project
// @version 1.0
// @description Documentation Final Project
// @termOfService http://swagger.io/terms/
// @schemes					http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorizations
// @description				Type "Bearer" followed by a space and JWT token.

func main() {
	errs := godotenv.Load(".env")
	if errs != nil {
		log.Fatalf("Some error occured. Err: %s", errs)
	}
	database.StartDB()
	r := router.StartApp()
	var PORT = os.Getenv("PORT")
	r.Run(":" + PORT)
}
