package db

import (
	"context"

	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	BOOKING_COLL = "bookings"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context, filter map[string]any) ([]*types.Booking, error)
	GetBooking(ctx context.Context, id string) (*types.Booking, error)
}

type MongoDbBookingStore struct {
	client      *mongo.Client
	bookingColl *mongo.Collection
}

func NewMongoDbBookinglStore(client *mongo.Client, dbname string) *MongoDbBookingStore {
	return &MongoDbBookingStore{
		client:      client,
		bookingColl: client.Database(dbname).Collection(BOOKING_COLL),
	}
}

func (s *MongoDbBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.bookingColl.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoDbBookingStore) GetBookings(ctx context.Context, filter map[string]any) ([]*types.Booking, error) {
	resp, err := s.bookingColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := resp.All(ctx, &bookings); err != nil {
		return []*types.Booking{}, nil
	}
	return bookings, nil
}

func (s *MongoDbBookingStore) GetBooking(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	if err := s.bookingColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}
