package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID      primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumPersons  int                `bson:"numPersons" json:"numPersons"`
	FromDate    time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate    time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
	CancelledAt time.Time          `bson:"cancelledAt,omitempty" json:"cancelledAt,omitempty"`
}

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

type BookingErrorResponse struct {
	Type string
	Msg  string
}

func (bkp BookRoomParams) Validate() error {
	now := time.Now()
	if now.After(bkp.FromDate) || now.After(bkp.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	return nil
}
