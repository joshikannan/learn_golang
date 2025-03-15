// this file belongs to the routes package.
// Since it's not main, it's meant to be imported into other files.
package routes

import (
	"context"          // Used for cancellation and timeouts in MongoDB operations.
	"myproject/db"     // Imports the database connection.
	"myproject/models" // Imports the models (structs) for MongoDB.

	"github.com/gofiber/fiber/v2"                // Imports Fiber, the web framework used for building APIs.
	"go.mongodb.org/mongo-driver/bson"           // Allows working with MongoDB BSON (Binary JSON) format.
	"go.mongodb.org/mongo-driver/bson/primitive" // Provides ObjectID handling, required for MongoDB _id fields.
)

// MongoDbRoutes registers all routes
func MongoDbRoutes(app *fiber.App) {
	app.Get("/api/todo", GetTodos)
	app.Post("/api/todo", CreateTodo)
	app.Patch("/api/todo/:id", UpdateTodo)
	app.Delete("/api/todo/:id", DeleteTodo)
}

// GetTodos retrieves all todos
func GetTodos(c *fiber.Ctx) error {
	var todos []models.Todo // todos is a slice (array) of models.Todo, where all fetched todos will be stored.
	cursor, err := db.Collection.Find(context.Background(), bson.M{})
	// 	db.Collection.Find(...) → Fetches all documents from the MongoDB collection.
	// context.Background() → Provides no timeout (not ideal; we can improve this later).
	// bson.M{} → An empty filter, meaning we fetch all todos.
	// 	Returns:
	// cursor → A pointer to the query result.
	// err → Holds any error that occurred.
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	// 	The defer keyword schedules it to run at the end of the function, even if an error occurs.
	// Prevents memory leaks by closing unused database connections.

	// cursor.Next(context.Background())
	// 	This iterates over the cursor (like a loop) to process each document returned from MongoDB.
	// cursor.Next(ctx) moves to the next document in the result set.
	// If there are no more documents, the loop stops.

	for cursor.Next(context.Background()) {
		var todo models.Todo
		// Declares a variable todo of type models.Todo.
		// It will store one MongoDB document at a time as we iterate through the cursor.
		if err := cursor.Decode(&todo); err != nil {
			// Takes the current document and decodes it into the todo struct.
			// If Decode() fails (due to a schema mismatch), it returns an error.
			return err
		}
		todos = append(todos, todo)
		// Adds each decoded todo to the todos slice (array).
	}

	return c.Status(200).JSON(fiber.Map{"todos": todos, "success": true})
	//	✔️ cursor.Next(ctx) → Moves to the next MongoDB document.
	// ✔️ cursor.Decode(&todo) → Converts MongoDB JSON into a Go struct (models.Todo).
	// ✔️ append(todos, todo) → Adds each document to a slice (array).
	// ✔️ defer cursor.Close(ctx) → Closes the cursor after function execution.
}

// CreateTodo adds a new todo
func CreateTodo(c *fiber.Ctx) error {
	todo := new(models.Todo)
	// todo is a new instance of the models.Todo struct.
	// new(models.Todo) initializes an empty todo object.
	if err := c.BodyParser(todo); err != nil {
		// c.BodyParser(todo) converts the incoming JSON request body into the todo struct.
		// If the request isn't valid JSON, an error is returned.
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	insertResult, err := db.Collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}
	// db.Collection.InsertOne(...) inserts todo into MongoDB.
	// context.Background() → Creates a default context (can be improved using timeouts).

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)
	// 	MongoDB automatically generates an _id for new documents.
	// insertResult.InsertedID returns this generated _id, but as an interface{}.
	// (primitive.ObjectID) type assertion converts it into an ObjectID.
	// When inserting into MongoDB, the _id isn't automatically added to the todo object in Go.
	// This step ensures our Go object contains the correct MongoDB _id.
	return c.Status(201).JSON(todo)
}

// UpdateTodo marks a todo as completed
func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	// MongoDB stores _id as ObjectID, not a string.
	// primitive.ObjectIDFromHex(id) converts the string id into an ObjectID.
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo id"})
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": true}}
	// Defines the update operation using MongoDB's $set operator.
	// It updates the completed field to true.

	_, err = db.Collection.UpdateOne(context.Background(), filter, update)
	// 	Calls UpdateOne() on MongoDB to update the matching todo.
	// Arguments:
	// context.Background() → Default context (can be improved with a timeout).
	// filter → Finds the todo based on _id.
	// update → Specifies the update (completed = true).
	// The function returns:
	// _ (ignored) → Update result (e.g., matched & modified count).
	// err → Error if update fails.
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}

// DeleteTodo removes a todo
func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo id"})
	}

	filter := bson.M{"_id": objectID}
	_, err = db.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": "Todo deleted successfully"})
}
