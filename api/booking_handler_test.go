package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/swarajroy/hotel-reservation/db/fixtures"
	"github.com/swarajroy/hotel-reservation/types"
)

func TestAdminGetBookings(t *testing.T) {
	ctx := context.TODO()

	db := Setup(t, ctx)
	defer db.TearDown(t, ctx)

	user := fixtures.AddUser(db.store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.store, "bar hotel", "london", nil)
	room := fixtures.AddRoom(db.store, types.SINGLE, 99.99, 99.99, hotel.ID)
	booking := fixtures.AddBooking(db.store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5), time.Time{}, 2)
	fmt.Println(booking)
}
