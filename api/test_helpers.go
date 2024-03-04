package api

import (
	"context"
	"testing"

	"github.com/swarajroy/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	store  *db.HotelReservationStore
}

func Setup(t *testing.T, ctx context.Context) *testdb {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		t.Fatal(err)
	}
	userStore := db.NewMongoDbUserStore(client, db.TEST_DB_NAME)
	hotelStore := db.NewMongoDbHotelStore(client, db.TEST_DB_NAME)
	roomStore := db.NewMongoDbRoomStore(client, db.TEST_DB_NAME, hotelStore)
	bookingStore := db.NewMongoDbBookingStore(client, db.TEST_DB_NAME)
	return &testdb{
		store: &db.HotelReservationStore{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		},
	}

}

func (tdb *testdb) TearDown(t *testing.T, ctx context.Context) {
	if err := tdb.store.User.Drop(ctx); err != nil {
		t.Fatal(err)
	}
}
