package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/swarajroy/hotel-reservation/api"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/db/fixtures"
	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	userStore    db.UserStore
	hotelStore   db.HotelStore
	roomStore    db.RoomStore
	bookingStore db.BookingStore
	store        *db.HotelReservationStore
	ctx          = context.Background()
)

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
	store = db.NewHotelReservationStore(userStore, hotelStore, roomStore, bookingStore)
}

func main() {
	user := fixtures.AddUser(store, "James", "Foo", false)
	fmt.Println("user -> ", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "Alice", "Mclain", true)
	fmt.Println("admin -> ", api.CreateTokenFromUser(admin))

	hotel := fixtures.AddHotel(store, "Bellucia", "France", nil)
	fmt.Println(hotel)

	room := fixtures.AddRoom(store, types.SINGLE, 99.99, 99.99, hotel.ID)
	fmt.Println(room)

	from := time.Now()
	to := from.AddDate(0, 0, 5)
	booking := fixtures.AddBooking(store, user.ID, room.ID, from, to, time.Time{}, 2)
	fmt.Println(booking)
}
