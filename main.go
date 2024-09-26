package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// need another field as mongodb stores as bson(binary json)
type Todo struct {
	// omitempty -> dont put in db if empty
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello world!")

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading .env file : ", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal("error connecting to mongodb : ", err)
	}

	defer client.Disconnect(context.Background()) // this happpens when main returns

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal("error pinging mongodb : ", err)
	}

	fmt.Println("Connected to MongoDB Atlas")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + PORT))

}

func getTodos(ctx *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	// use the cursor to iterate over documents returned by db
	if err != nil {
		return err
	}

	// defer -> postpone execution of a function until the surrounding function returns
	// so this will happen after getTodos returns
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return ctx.Status(200).JSON(todos)
}

func createTodo(ctx *fiber.Ctx) error {
	todo := new(Todo)
	// {id:0, completed:false, body:""}

	if err := ctx.BodyParser(todo); err != nil {
		return err
	}
	if todo.Body == "" {
		return ctx.Status(400).SendString("Body is required")
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}
	todo.ID = insertResult.InsertedID.(primitive.ObjectID)
	return ctx.Status(201).JSON(todo)
}

func updateTodo(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid todo id"})
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"Completed": true}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{"success": "true"})
}

func deleteTodo(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid todo id"})
	}

	filter := bson.M{"_id": objectId}
	if _, err = collection.DeleteOne(context.Background(), filter); err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{"success": "true"})
}
