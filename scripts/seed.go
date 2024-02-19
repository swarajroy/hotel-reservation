package main

import (
	"context"
	"fmt"
	"log"

	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/db/fixtures"
	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	userStore    db.UserStore
	hotelStore   db.HotelStore
	roomStore    db.RoomStore
	bookingStore db.BookingStore
	store        db.HotelReservationStore
	ctx          = context.Background()
)

func seedHotel(name, location string) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted hotel = ", insertedHotel)
	rooms := []types.Room{
		{
			Type:      types.SINGLE,
			BasePrice: 99.99,
		},
		{
			Type:      types.DOUBLE,
			BasePrice: 199.99,
		},
		{
			Type:      types.DELUXE,
			BasePrice: 129.99,
		},
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("inserted room = ", insertedRoom)
	}
}

func seedUser(c context.Context, fname, lname, email, password string, isAdmin bool) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	res, err := store.User.InsertUser(c, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted User = ", res)
}

func init() {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	userStore = db.NewMongoDbUserStore(client, db.DBNAME)
	hotelStore = db.NewMongoDbHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoDbRoomStore(client, db.DBNAME, hotelStore)
	bookingStore = db.NewMongoDbBookingStore(client, db.DBNAME)
	store = *db.NewHotelReservationStore(userStore, hotelStore, roomStore, bookingStore)
}

func main() {
	fixtures.AddUser(&store, "James", "Foo", "james@foo.com", "supersecurepassword", false)
	fixtures.AddUser(&store, "Alice", "Mclain", "alice@mclain.com", "admin1234", true)
	return
	seedHotel("Bellucia", "France")
	seedHotel("The cozy hotel", "The Netherlands")
	seedHotel("Die another day", "UK")
	seedUser(context.Background(), "James", "Foo", "james@foo.com", "supersecurepassword", false)
	seedUser(context.Background(), "Miceal", "Jordan", "miceal@jordan.com", "airtime23", false)
	seedUser(context.Background(), "Alice", "Mclain", "alice@mclain.com", "admin1234", true)
	fmt.Println("seeded the db")
}
