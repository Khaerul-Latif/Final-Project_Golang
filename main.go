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
// @description MyGram is a free photo sharing app written in Go. People can share, view, and comment photos by everyone. Anyone can create an account by registering an email address and selecting a username.
// @termOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description					        Description for what is this security definition being used
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
