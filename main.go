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

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := flag.String("listenAddr", ":3000", "The API Servers port")
	flag.Parse()
	app := fiber.New(config)

	apiv1 := app.Group("/api/v1")

	var (
		userStore  = db.NewMongoDbUserStore(client, db.DBNAME)
		hotelStore = db.NewMongoDbHotelStore(client, db.DBNAME)
		roomStore  = db.NewMongoDbRoomStore(client, db.DBNAME, hotelStore)
		store      = &db.HotelReservationStore{
			User:  userStore,
			Hotel: hotelStore,
			Room:  roomStore,
		}
		userHandler  = api.NewUserHandler(store)
		hotelHandler = api.NewHotelHandler(store)
	)

	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/users/:id", userHandler.HandlePutUser)

	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	app.Listen(*listenAddr)
}
