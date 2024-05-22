package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/db/fixtures"
	"github.com/swarajroy/hotel-reservation/types"
)

func TestAdminUserGetBookingsSuccessful(t *testing.T) {
	ctx := context.TODO()

	db := Setup(t, ctx)
	defer db.TearDown(t, ctx)

	var (
		admin_user     = fixtures.AddUser(db.store, "admin", "admin", true)
		user           = fixtures.AddUser(db.store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.store, "bar hotel", "london", nil)
		room           = fixtures.AddRoom(db.store, types.SINGLE, 99.99, 99.99, hotel.ID)
		booking        = fixtures.AddBooking(db.store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5), time.Time{}, 2)
		app            = fiber.New()
		admin          = app.Group("/", JWTAuthentication(db.store), AdminAuth)
		bookingHandler = NewBookingHandler(db.store)
	)

	_ = booking

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(admin_user))
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response got %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
}

func TestNormalUserGetBookingsFail(t *testing.T) {
	ctx := context.TODO()

	db := Setup(t, ctx)
	defer db.TearDown(t, ctx)

	var (
		user           = fixtures.AddUser(db.store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.store, "bar hotel", "london", nil)
		room           = fixtures.AddRoom(db.store, types.SINGLE, 99.99, 99.99, hotel.ID)
		booking        = fixtures.AddBooking(db.store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5), time.Time{}, 2)
		app            = fiber.New()
		admin          = app.Group("/", JWTAuthentication(db.store), AdminAuth)
		bookingHandler = NewBookingHandler(db.store)
	)

	_ = booking

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}
}

func TestNormalUserGetBookingSuccess(t *testing.T) {
	ctx := context.TODO()

	db := Setup(t, ctx)
	defer db.TearDown(t, ctx)

	var (
		user           = fixtures.AddUser(db.store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.store, "bar hotel", "london", nil)
		room           = fixtures.AddRoom(db.store, types.SINGLE, 99.99, 99.99, hotel.ID)
		booking        = fixtures.AddBooking(db.store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5), time.Time{}, 2)
		app            = fiber.New()
		userRoute      = app.Group("/", JWTAuthentication(db.store))
		bookingHandler = NewBookingHandler(db.store)
	)

	_ = booking

	userRoute.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected a 200 status code got %d", resp.StatusCode)
	}

	var returnedBooking *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&booking); err != nil {
		t.Fatal(err)
	}
	fmt.Println(returnedBooking)
}
