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

func TestUserStoreSuite(t *testing.T) {
	suite.Run(t, new(UserStoreSuite))
}
