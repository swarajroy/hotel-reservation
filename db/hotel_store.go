package db

import (
	"context"
	"fmt"

	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	HOTEL_COLL = "hotels"
)

type HotelStore interface {
	Dropper
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, filter map[string]any, update map[string]any) error
	GetHotels(ctx context.Context, filter map[string]any, paginaton *Pagination) ([]*types.Hotel, error)
	GetHotelById(context.Context, string) (*types.Hotel, error)
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

func (s *MongoDbHotelStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping hotel collection ---")
	return s.hotelColl.Drop(ctx)
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

func (s *MongoDbHotelStore) GetHotels(ctx context.Context, filter map[string]any, pag *Pagination) ([]*types.Hotel, error) {
	var (
		skip = (pag.Page - 1) * pag.Limit
	)
	opts := &options.FindOptions{
		Limit: &pag.Limit,
		Skip:  &skip,
	}

	resp, err := s.hotelColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := resp.All(ctx, &hotels); err != nil {
		return []*types.Hotel{}, nil
	}
	return hotels, nil
}

func (s *MongoDbHotelStore) GetHotelById(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, NewDBError(err.Error())
	}
	var hotel types.Hotel
	if err := s.hotelColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}
