package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.HotelReservationStore
}

func NewBookingHandler(store *db.HotelReservationStore) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// This needs to be admin authorised
func (bh *BookingHandler) HandleBookings(c *fiber.Ctx) error {
	bookings, err := bh.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return nil
	}
	return c.JSON(bookings)
}

// This needs to be user authorised
func (bh *BookingHandler) HandleBooking(c *fiber.Ctx) error {
	booking, err := bh.store.Booking.GetBooking(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}

	if booking.UserID != user.ID || !user.IsAdmin {
		return c.Status(http.StatusUnauthorized).JSON(AuthErrorResponse{
			Status: http.StatusUnauthorized,
			Msg:    "error",
		})
	}

	return c.JSON(booking)
}
