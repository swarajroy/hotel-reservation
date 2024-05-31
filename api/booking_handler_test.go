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
	"github.com/stretchr/testify/suite"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/db/fixtures"
	"github.com/swarajroy/hotel-reservation/db/mongo"
	"github.com/swarajroy/hotel-reservation/types"
)

type BookingHandlerSuite struct {
	suite.Suite
	bookingHandler  *BookingHandler
	store           *db.HotelReservationStore
	testMongoClient *mongo.TestMongoClient
}

func (suite *BookingHandlerSuite) SetupSuite() {
	const (
		DB_NAME = "hotel-reservation-test"
	)
	client, err := mongo.NewTestMongoClient(DB_NAME)
	if err != nil {
		suite.T().Error("failed to connect to mongo db container in docker using testcontainers")
	}

	suite.testMongoClient = client
	userStore := db.NewMongoDbUserStore(suite.testMongoClient.Client, DB_NAME)
	hotelStore := db.NewMongoDbHotelStore(suite.testMongoClient.Client, DB_NAME)
	roomStore := db.NewMongoDbRoomStore(suite.testMongoClient.Client, DB_NAME, hotelStore)
	bookingStore := db.NewMongoDbBookingStore(suite.testMongoClient.Client, DB_NAME)
	store := db.NewHotelReservationStore(userStore, hotelStore, roomStore, bookingStore)
	suite.store = store
	suite.bookingHandler = NewBookingHandler(store)
}

func (suite *BookingHandlerSuite) TearDownSuite() {
	suite.testMongoClient.Container.Terminate(context.Background())
}

func (suite *BookingHandlerSuite) AfterTest() {
	suite.store.Booking.Drop(context.Background())
}

func TestBookingHandlerSuite(t *testing.T) {
	suite.Run(t, new(BookingHandlerSuite))
}

func (suite *BookingHandlerSuite) TestAdminUserGetBookingsSuccessful() {

	var (
		admin_user     = fixtures.AddUser(suite.store, "admin", "admin", true)
		user           = fixtures.AddUser(suite.store, "james", "foo", false)
		hotel          = fixtures.AddHotel(suite.store, "bar hotel", "london", nil)
		room           = fixtures.AddRoom(suite.store, types.SINGLE, 99.99, 99.99, hotel.ID)
		booking        = fixtures.AddBooking(suite.store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5), time.Time{}, 2)
		app            = fiber.New()
		admin          = app.Group("/", JWTAuthentication(suite.store), AdminAuth)
		bookingHandler = NewBookingHandler(suite.store)
	)

	_ = booking

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(admin_user))
	resp, err := app.Test(req)

	if err != nil {
		suite.T().Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		suite.T().Fatalf("non 200 response got %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		suite.T().Fatal(err)
	}

	if len(bookings) != 1 {
		suite.T().Fatalf("expected 1 booking got %d", len(bookings))
	}
}

func (suite *BookingHandlerSuite) TestNormalUserGetBookingsFail() {
	var (
		user           = fixtures.AddUser(suite.store, "james", "foo", false)
		hotel          = fixtures.AddHotel(suite.store, "bar hotel", "london", nil)
		room           = fixtures.AddRoom(suite.store, types.SINGLE, 99.99, 99.99, hotel.ID)
		booking        = fixtures.AddBooking(suite.store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5), time.Time{}, 2)
		app            = fiber.New()
		admin          = app.Group("/", JWTAuthentication(suite.store), AdminAuth)
		bookingHandler = NewBookingHandler(suite.store)
	)

	_ = booking

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)

	if err != nil {
		suite.T().Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		suite.T().Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}
}

func (suite *BookingHandlerSuite) TestNormalUserGetBookingSuccess() {
	var (
		user           = fixtures.AddUser(suite.store, "james", "foo", false)
		hotel          = fixtures.AddHotel(suite.store, "bar hotel", "london", nil)
		room           = fixtures.AddRoom(suite.store, types.SINGLE, 99.99, 99.99, hotel.ID)
		booking        = fixtures.AddBooking(suite.store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5), time.Time{}, 2)
		app            = fiber.New()
		userRoute      = app.Group("/", JWTAuthentication(suite.store))
		bookingHandler = NewBookingHandler(suite.store)
	)

	_ = booking

	userRoute.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)

	if err != nil {
		suite.T().Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		suite.T().Fatalf("expected a 200 status code got %d", resp.StatusCode)
	}

	var returnedBooking *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&booking); err != nil {
		suite.T().Fatal(err)
	}
	fmt.Println(returnedBooking)
}
