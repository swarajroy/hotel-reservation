package db

import (
	"context"

	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	USER_COLL = "users"
)

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
}

type MongoDbUserStore struct {
	client   *mongo.Client
	userColl *mongo.Collection
}

func NewMongoDbUserStore(client *mongo.Client) *MongoDbUserStore {
	return &MongoDbUserStore{
		client:   client,
		userColl: client.Database(DBNAME).Collection(USER_COLL),
	}
}

func (s *MongoDbUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.userColl.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoDbUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.userColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return []*types.User{}, nil
	}
	return users, nil
}

func (s *MongoDbUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.userColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
