package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello World!")
	app := fiber.New()

	todos := []Todo{}

	/**
	function gets a pointer to fiber context
	get all todos
	*/
	app.Get("/api/todos", func(context *fiber.Ctx) error {
		return context.Status(200).JSON(todos)
	})

	// created a todo
	app.Post("/api/todos", func(context *fiber.Ctx) error {
		todo := &Todo{} //{id:0,completed:false,body:""}
		// body of context will be used for todo
		if err := context.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return context.Status(400).JSON(fiber.Map{"error": "todo body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return context.Status(201).JSON(todo)
	})

	// update a todo
	app.Patch("/api/todos/:id", func(context *fiber.Ctx) error {
		id := context.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return context.Status(200).JSON(todos[i])
			}
		}
		return context.Status(404).JSON(fiber.Map{"error": "todo not found"})
	})

	// delete a todo
	app.Delete("/api/todos/delete/:id", func(context *fiber.Ctx) error {
		id := context.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				// make todos a list of all todos before and after this excluding this, unpack values using ... operator
				todos = append(todos[:i], todos[i+1:]...)
				return context.Status(200).JSON(fiber.Map{"success": "true"})
			}
		}
		return context.Status(404).JSON(fiber.Map{"error": "todo not found"})
	})

	log.Fatal(app.Listen(":9000")) // listen at port 4000, log.fatal will print any errors and exit

}
