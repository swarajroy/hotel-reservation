package db

import (
	"context"

	"github.com/swarajroy/hotel-reservation/db/mongo/utils"
	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	USER_COLL = "users"
)

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUserById(context.Context, string) error
	UpdateUserById(ctx context.Context, params types.UpdateUserParams, id string) error
	GetUserByEmail(context.Context, string) (*types.User, error)
}

type MongoDbUserStore struct {
	client   *mongo.Client
	userColl *mongo.Collection
}

func (s *MongoDbUserStore) Drop(ctx context.Context) error {
	return s.userColl.Drop(ctx)
}

func NewMongoDbUserStore(client *mongo.Client, dbname string) *MongoDbUserStore {
	return &MongoDbUserStore{
		client:   client,
		userColl: client.Database(dbname).Collection(USER_COLL),
	}
}

func (s *MongoDbUserStore) DeleteUserById(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	//TODO : what to do when we want to delete a user by ID but the delete fails... Log it may be
	_, err = s.userColl.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return utils.ResourceNotFound("user", "id", err)
	}
	return nil
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
	var user *types.User
	if err := s.userColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, utils.ResourceNotFound("user", id, err)
	}
	return user, nil
}

func (s *MongoDbUserStore) UpdateUserById(ctx context.Context, params types.UpdateUserParams, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": oid,
	}

	update := bson.M{
		"$set": bson.M{
			"firstName": params.FirstName,
			"lastName":  params.LastName,
		},
	}

	if _, err := s.userColl.UpdateOne(ctx, filter, update); err != nil {
		return utils.ResourceNotFound("user", id, err)
	}
	return nil
}

func (s *MongoDbUserStore) GetUserByEmail(c context.Context, email string) (*types.User, error) {
	filter := bson.M{
		"email": email,
	}
	result := s.userColl.FindOne(c, filter)
	var user *types.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}
