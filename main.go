package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/api"
	"github.com/swarajroy/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//DB_URI = "mongodb+srv://flashcardsage:VCB36qjOQwfsfpqY@cluster0.agr8mpl.mongodb.net/?retryWrites=true&w=majority"
	DB_URI = "mongodb://127.0.0.1:27017"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := flag.String("listenAddr", ":3000", "The API Servers port")
	flag.Parse()
	app := fiber.New(config)

	apiv1 := app.Group("/api/v1")

	userHandler := api.NewUserHandler(db.NewMongoDbUserStore(client))

	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/users/:id", userHandler.HandlePutUser)

	app.Listen(*listenAddr)
}
