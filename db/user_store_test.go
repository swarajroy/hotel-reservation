package db

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"
	"github.com/swarajroy/hotel-reservation/db/mongo"
	"github.com/swarajroy/hotel-reservation/types"
	userfixtures "github.com/swarajroy/hotel-reservation/types/user_fixtures"
)

type UserStoreSuite struct {
	suite.Suite
	userStore       UserStore
	testMongoClient *mongo.TestMongoClient
}

func (suite *UserStoreSuite) SetupSuite() {
	const (
		DB_NAME = "hotel-reservation-test"
	)
	client, err := mongo.NewTestMongoClient(DB_NAME)
	if err != nil {
		suite.T().Error("failed to connect to mongo db container in docker using testcontainers")
	}

	suite.testMongoClient = client
	suite.userStore = NewMongoDbUserStore(suite.testMongoClient.Client, DB_NAME)

}

func (suite *UserStoreSuite) TearDownSuite() {
	suite.testMongoClient.Container.Terminate(context.Background())
}

func (suite *UserStoreSuite) AfterTest() {
	suite.userStore.Drop(context.Background())
}

func (suite *UserStoreSuite) TestInsertUser() {
	//user := nil
	expected, err := userfixtures.Next()
	if err != nil {
		suite.T().Fatalf("error generating user")
	}

	actual, err := suite.userStore.InsertUser(context.Background(), expected)

	suite.Nil(err)
	suite.NotNil(actual)
	suite.False(actual.ID.IsZero())
	suite.Equal(expected, actual)
}

func (suite *UserStoreSuite) TestUpdateUserById() {
	var (
		ctx   = context.Background()
		fName = faker.FirstName()
		lName = faker.LastName()
	)
	user, _ := userfixtures.Next()
	insertedUser, _ := suite.userStore.InsertUser(ctx, user)

	err := suite.userStore.UpdateUserById(ctx, types.UpdateUserParams{FirstName: fName, LastName: lName}, insertedUser.ID.Hex())

	suite.Nil(err)

	retrievedUser, err := suite.userStore.GetUserById(ctx, insertedUser.ID.Hex())

	suite.Nil(err)
	suite.NotNil(retrievedUser)
	suite.NotEqual(retrievedUser.FirstName, user.FirstName)
	suite.NotEqual(retrievedUser.LastName, user.LastName)
	suite.Equal(retrievedUser.FirstName, fName)
	suite.Equal(retrievedUser.LastName, lName)

}

func (suite *UserStoreSuite) TestGetUserById() {
	var (
		ctx = context.Background()
	)
	//user := nil
	expected, err := userfixtures.Next()
	if err != nil {
		suite.T().Fatalf("error generating user")
	}

	insertedUser, _ := suite.userStore.InsertUser(context.Background(), expected)

	retrievedUser, err := suite.userStore.GetUserById(ctx, insertedUser.ID.Hex())

	suite.Nil(err)
	suite.NotNil(retrievedUser)
	suite.Equal(insertedUser, retrievedUser)
}

func (suite *UserStoreSuite) TestGetUsers() {
	var (
		ctx = context.Background()
	)
	expected, err := userfixtures.Next()
	if err != nil {
		suite.T().Fatalf("error generating user")
	}

	insertedUser, _ := suite.userStore.InsertUser(context.Background(), expected)
	retrievedUsers, err := suite.userStore.GetUsers(ctx)

	suite.Nil(err)
	suite.NotEmpty(retrievedUsers)
	suite.Contains(retrievedUsers, insertedUser)

}

func (suite *UserStoreSuite) TestGetUserByEmail() {
	var (
		ctx = context.Background()
	)
	expected, err := userfixtures.Next()
	if err != nil {
		suite.T().Fatalf("error generating user")
	}
	insertedUser, _ := suite.userStore.InsertUser(context.Background(), expected)
	retrievedUser, err := suite.userStore.GetUserByEmail(ctx, expected.Email)

	suite.Nil(err)
	suite.NotNil(retrievedUser)
	suite.Equal(insertedUser, retrievedUser)
}

func (suite *UserStoreSuite) TestDeleteUserById() {
	var (
		ctx = context.Background()
	)
	expected, err := userfixtures.Next()
	if err != nil {
		suite.T().Fatalf("error generating user")
	}
	insertedUser, _ := suite.userStore.InsertUser(context.Background(), expected)

	err = suite.userStore.DeleteUserById(ctx, insertedUser.ID.Hex())

	suite.Nil(err)
}

func TestUserStoreSuite(t *testing.T) {
	suite.Run(t, new(UserStoreSuite))
}
