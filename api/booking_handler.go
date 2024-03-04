package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
func (bh *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := bh.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return nil
	}
	return c.JSON(bookings)
}

// This needs to be user authorised
func (bh *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
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

func (bh *BookingHandler) HandleDeleteBooking(c *fiber.Ctx) error {
	booking, err := bh.store.Booking.GetBooking(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	log.Info("user = ", user)
	if !ok {
		return fmt.Errorf("not authorized")
	}

	fmt.Printf("user = %+v\n", user)
	fmt.Printf("booking = %+v\n", booking)
	if !user.IsAdmin {
		log.Error("illegal action as user trying to cancel a booking that does not belong to him/her or user is not an admin")
		return c.Status(http.StatusForbidden).JSON(AuthErrorResponse{
			Status: http.StatusForbidden,
			Msg:    "error",
		})
	}
	update := map[string]any{
		"cancelledAt": time.Now(),
	}
	if err := bh.store.Booking.UpdateBookingById(c.Context(), c.Params("id"), update); err != nil {
		return err
	}
	return c.JSON(map[string]string{
		"msg": "updated",
	})
}
