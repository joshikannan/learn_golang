package main

// this is from u tube channel free code camp
// url => https://youtu.be/lNd7XlXwlho

// to me => go is very easy just understand the concept and syntax and flow
// fiber framework as similar to express
//

import (
	"log"
	"myproject/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error on loading .env file")
	}

	PORT := os.Getenv("PORT")

	routes.RoutesWithoutDB(app) // api routes without db => url => https://youtu.be/lNd7XlXwlho

	// 	 In Go, log.Fatal is used to log an error message and immediately exit the program.
	// 	It’s similar to console.error in Node.js but with one key difference—it stops execution.

	// log.Fatal(app.Listen(":4000"))

	log.Fatal(app.Listen(":" + PORT))
}
