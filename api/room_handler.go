package api

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.HotelReservationStore
}

func NewRoomHandler(store *db.HotelReservationStore) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {

	var params types.BookRoomParams
	ctx := c.Context()
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := ctx.Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(AuthErrorResponse{
			Status: http.StatusInternalServerError,
			Msg:    "internal server error",
		})
	}

	ok, err = h.isRoomAvailableForBooking(ctx, params, roomID)
	if err != nil {
		return err
	}

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(types.BookingErrorResponse{
			Type: "error",
			Msg:  "room already booked",
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	insertedBooking, err := h.store.Booking.InsertBooking(ctx, &booking)
	if err != nil {
		return err
	}

	return c.JSON(insertedBooking)
}

func (rh *RoomHandler) isRoomAvailableForBooking(ctx context.Context, params types.BookRoomParams, roomID primitive.ObjectID) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}

	bookings, err := rh.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0
	return ok, nil
}
