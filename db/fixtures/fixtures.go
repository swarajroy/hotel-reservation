package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.HotelReservationStore, fn, ln string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s_%s@foo.com", fn, ln),
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddHotel(store *db.HotelReservationStore, name, loc string, rooms []primitive.ObjectID) *types.Hotel {
	if rooms == nil {
		rooms = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    rooms,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted hotel = ", insertedHotel)
	return &hotel
}

func AddRoom(store *db.HotelReservationStore, ty types.RoomType, basePrice, price float64, hid primitive.ObjectID) *types.Room {
	room := &types.Room{
		Type:      ty,
		BasePrice: basePrice,
		Price:     price,
		HotelID:   hid,
	}

	insertedRoom, err := store.Room.InsertRoom(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted room = ", insertedRoom)
	return room
}

func AddBooking(store *db.HotelReservationStore, uid, rid primitive.ObjectID, from, till, cancelledAt time.Time, numPersons int) *types.Booking {
	booking := &types.Booking{
		UserID:      uid,
		RoomID:      rid,
		FromDate:    from,
		TillDate:    till,
		NumPersons:  numPersons,
		CancelledAt: cancelledAt,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted booking = ", insertedBooking)
	return insertedBooking
}
