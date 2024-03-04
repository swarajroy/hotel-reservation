package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/api"
	"github.com/swarajroy/hotel-reservation/api/middleware"
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

	var (
		userStore    = db.NewMongoDbUserStore(client, db.DBNAME)
		hotelStore   = db.NewMongoDbHotelStore(client, db.DBNAME)
		roomStore    = db.NewMongoDbRoomStore(client, db.DBNAME, hotelStore)
		bookingStore = db.NewMongoDbBookingStore(client, db.DBNAME)
		store        = &db.HotelReservationStore{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(store)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		authHandler    = api.NewAuthHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(store))
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
	)

	// auth handlers
	auth.Post("/auth", authHandler.HandleAuth)
	// user handlers
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/users/:id", userHandler.HandlePutUser)

	// hotel handler
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	// room handler
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// bookings handler - admin route
	admin.Get("/bookings", bookingHandler.HandleGetBookings)
	admin.Delete("/bookings/:id", bookingHandler.HandleDeleteBooking)
	// bookings handler - user route
	apiv1.Get("/bookings/:id", bookingHandler.HandleGetBooking)

	app.Listen(*listenAddr)
}
