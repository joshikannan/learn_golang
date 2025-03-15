// main package is special in Go—it's the entry point of a Go program.
package main

import (
	"fmt"              // Used for printing output to the console.
	"log"              //  Used for logging errors and other important information.
	"myproject/db"     //  db folder -> manages database connection
	"myproject/routes" // routes folder -> route handlers for our API.
	"os"               // Used to read environment variables (like PORT).

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// main() function inside this package will be executed when you run the program.
// where execution starts here
func main() {
	fmt.Println("Starting server...") // similar to console log, Logs a message but does not stop execution.

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	// log.Fatal => used to log error messages and immediately terminate the program.
	// It works similarly to log.Println(), but after printing the message, it calls os.Exit(1),
	// which immediately stops the program with a non-zero exit status.

	// Connect to MongoDB
	// Calls ConnectMongoDB() from the db package (db/connect.go).
	// This function establishes a connection to MongoDB and assigns the collection reference.
	db.ConnectMongoDB()

	// Creates a new instance of the Fiber web framework.
	// This instance (app) inherits all the methods of fiber.App, access all the methods of fiber
	// app := fiber.New() creates an instance of Fiber.
	// app (instance) can access all methods of fiber.App.
	// This works just like in Express.js where app is an instance of Express.
	app := fiber.New()

	// Register API routes
	routes.MongoDbRoutes(app)   // MongoDB-based routes =>  Registers API routes that interact with MongoDB.
	routes.RoutesWithoutDB(app) // In-memory (local) routes => Registers API routes that store data in memory (without a database).

	// Set port from environment variable or default to 4000
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	// Starts the Fiber server on the specified port.
	// If the server fails to start, log.Fatal() logs the error and exits the program.
	log.Fatal(app.Listen(":" + PORT))
}
