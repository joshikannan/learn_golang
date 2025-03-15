// belongs to the db package
// package handles the MongoDB connection
package db

import (
	"context"
	// Provides a context for handling long-running operations (e.g., DB queries). Similar to async/await timeouts in Node.js.
	// Manages timeouts, cancellations, deadlines, long-running operations (DB queries  (MongoDB, SQL), HTTP calls, goroutines),
	"fmt" // Used for printing messages (like console.log in JavaScript).
	"log" // Helps in logging errors (like console.error).
	"os"  // Allows access to environment variables (MONGO_URI).

	"go.mongodb.org/mongo-driver/mongo"
	// The official Go driver for MongoDB. (Like mongoose in Node.js). The MongoDB Go Driver (mongo-driver) provides native support.
	// But The official Go driver does not provide schema enforcement or ODM features. We must manually define structures and write queries.
	"go.mongodb.org/mongo-driver/mongo/options" // Used to configure MongoDB client options (like setting the MongoDB URL).
)

// Declares a global variable Collection, which will hold the MongoDB collection.
// The *mongo.Collection means it’s a pointer to a MongoDB collection (like db.collection('todos') in Node.js).
// This allows other parts of the program to use this collection without reconnecting.
var Collection *mongo.Collection

func ConnectMongoDB() { // which we call in main.go to connect to MongoDB.
	MONGODB_URI := os.Getenv("MONGO_URI")

	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI(MONGODB_URI) // clientOptions holds the configured connection settings.
	client, err := mongo.Connect(context.Background(), clientOptions)

	// mongo.Connect() establishes a connection with MongoDB.
	// context.Background() provides a default context (used for managing long-running operations (Database queries (MongoDB, SQL), API requests Goroutines ,Timeouts & cancellations)).
	// The context.Background() is passed to MongoDB to allow cancellation if needed.
	// This prevents long-running DB queries from blocking the app.
	// context.Background() by itself does not prevent long-running operations. Instead, it is used as a starting point for managing execution flow in your app.
	// If you want to prevent long-running operations, you should use context.WithTimeout() or context.WithCancel() based on context.Background().

	if err != nil {
		log.Fatal(err)
	}

	// Ping to check if MongoDB is reachable
	err = client.Ping(context.Background(), nil)
	// .Ping() checks if MongoDB is reachable.

	if err != nil {
		log.Fatal(err)
	}
	// If MongoDB is down, log.Fatal(err) stops the program.
	// MongoDB may connect but still be unreachable due to network issues.
	// This ensures the connection is actually working before using it.

	fmt.Println("Connected to MongoDB Atlas")

	// Assign the todos collection
	// client.Database("golang_db") selects the database "golang_db".
	// Collection("todos") selects the "todos" collection.
	Collection = client.Database("golang_db").Collection("todos")
}
