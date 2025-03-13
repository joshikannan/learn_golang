package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Todo model
type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

// In-memory store for todos => we we end terminal and reloads , re reun the server by saving it became empty as it is....
var todos = []Todo{} // like a empty array => let todos = [];
// 	todos := []Todo{} → Creates an empty slice (dynamic array,similar to an empty array [] in JavaScript) of Todo items.
// 	// Similar to:
// 	// Go: var todos []Todo
// 	// JS: let todos = [];
// 	// go
// 	// todos := []Todo{
// 	// {ID: 1, Completed: false, Body: "Buy groceries"},
// 	// {ID: 2, Completed: true, Body: "Learn Go"},
// 	//js
// 	// 	let todos = [
// 	//   { id: 1, completed: false, body: "Buy groceries" },
// 	//   { id: 2, completed: true, body: "Learn Go" }
// 	// ];
// 	// }

// ================================================== || Setup Routes || ==================================================
func RoutesWithoutDB(app *fiber.App) {
	app.Get("/api/todos", GetTodos)
	app.Post("/api/todo", CreateTodo)
	app.Patch("/api/todo/:id", UpdateTodo)
	app.Delete("/api/todo/:id", DeleteTodo)
}

// ================================================== || Get All Todos || ==================================================
func GetTodos(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"todos": todos, "success": true})
}

// ================================================== || Create New Todo || ==================================================
func CreateTodo(c *fiber.Ctx) error {
	// similar to app.post("/api/todo", (req, res) => {

	todo := &Todo{}
	// Creates a new Todo object and gets a reference (&Todo) => address of Todo
	// 		Similar to:
	// Go: var todo = &Todo{}
	// JS: let todo = {};

	if err := c.BodyParser(todo); err != nil {
		return err
	}
	// c.BodyParser(todo) → Parses the JSON body from the request and fills the todo struct.
	// If parsing fails (invalid JSON), it returns an error.

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	todo.ID = len(todos) + 1
	todos = append(todos, *todo)
	// *todo => points the value in that addrss
	// append(todos, *todo) → Adds the new todo to the todos slice (array).
	// *todo dereferences the pointer to get the actual struct.
	// js -> todos.push(todo);

	return c.Status(201).JSON(todo)
}

// ================================================== || Update (Patch) Todo || ==================================================
func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	for i, todo := range todos {
		if fmt.Sprint(todo.ID) == id {
			todos[i].Completed = true
			return c.Status(200).JSON(todos[i])
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
}

// ================================================== || Delete Todo || ==================================================
func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	for i, todo := range todos {
		if fmt.Sprint(todo.ID) == id {
			todos = append(todos[:i], todos[i+1:]...)
			// if delete 5 rd todo, in list of 10
			// todos[:i] => 1,2,3,4 => except 5 => before 5th
			// todos[i+1:]... => todo after 5, => 6 to 10 => after th
			// todos = append(todos[:i], todos[i+1:]...) => adds up before 5 and after 5
			return c.Status(200).JSON(fiber.Map{"success": "Todo removed successfully"})
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "Todo not found to delete"})
}
