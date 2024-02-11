package db

import (
	"context"

	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	HOTEL_COLL = "hotels"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, filter map[string]any, update map[string]any) error
	GetHotels(ctx context.Context, filter map[string]any) ([]*types.Hotel, error)
}

type MongoDbHotelStore struct {
	client    *mongo.Client
	hotelColl *mongo.Collection
}

func NewMongoDbHotelStore(client *mongo.Client, dbname string) *MongoDbHotelStore {
	return &MongoDbHotelStore{
		client:    client,
		hotelColl: client.Database(dbname).Collection(HOTEL_COLL),
	}
}

func (s *MongoDbHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.hotelColl.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoDbHotelStore) UpdateHotel(ctx context.Context, filter map[string]any, update map[string]any) error {
	_, err := s.hotelColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoDbHotelStore) GetHotels(ctx context.Context, filter map[string]any) ([]*types.Hotel, error) {
	resp, err := s.hotelColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := resp.All(ctx, &hotels); err != nil {
		return []*types.Hotel{}, nil
	}
	return hotels, nil
}
