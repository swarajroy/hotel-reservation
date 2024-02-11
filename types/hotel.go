package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json: "name"`
	Location string               `bson:"location" json: "location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json: "rooms"`
}

type RoomType int

const (
	_ RoomType = iota
	SINGLE
	DOUBLE
	DELUXE
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json: "basePrice"`
	Price     float64            `bson:"price" json: "price"`
	HotelID   primitive.ObjectID `bson:"hotelId" json:"hotelId"`
}
