package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/db/mongo"
	"github.com/swarajroy/hotel-reservation/types"
	userfixtures "github.com/swarajroy/hotel-reservation/types/user_fixtures"
)

type UserHandlerSuite struct {
	suite.Suite
	store           *db.HotelReservationStore
	userHandler     *UserHandler
	testMongoClient *mongo.TestMongoClient
}

func (suite *UserHandlerSuite) SetupSuite() {
	const (
		DB_NAME = "hotel-reservation-test"
	)
	client, err := mongo.NewTestMongoClient(DB_NAME)
	if err != nil {
		suite.T().Error("failed to connect to mongo db container in docker using testcontainers")
	}

	suite.testMongoClient = client
	userStore := db.NewMongoDbUserStore(suite.testMongoClient.Client, DB_NAME)
	hotelStore := db.NewMongoDbHotelStore(suite.testMongoClient.Client, DB_NAME)
	roomStore := db.NewMongoDbRoomStore(suite.testMongoClient.Client, DB_NAME, hotelStore)
	bookingStore := db.NewMongoDbBookingStore(suite.testMongoClient.Client, DB_NAME)
	store := db.NewHotelReservationStore(userStore, hotelStore, roomStore, bookingStore)
	suite.store = store

	suite.testMongoClient = client
	suite.userHandler = NewUserHandler(store)
}

func (suite *UserHandlerSuite) TearDownSuite() {
	suite.testMongoClient.Container.Terminate(context.Background())
}

func (suite *UserHandlerSuite) AfterTest() {
	suite.store.User.Drop(context.Background())
}

func (suite *UserHandlerSuite) TestPostUser() {
	var POST_ROUTE = "/users"

	app := fiber.New()
	app.Post(POST_ROUTE, suite.userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Password:  faker.Password(),
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", POST_ROUTE, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	var user *types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		suite.T().Errorf("expecting userId to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		suite.T().Errorf("expecting no encrypted password to be set in the response")
	}
	if user.FirstName != params.FirstName {
		suite.T().Errorf("expected firstname %s got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		suite.T().Errorf("expected lastname %s got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		suite.T().Errorf("expected email %s got %s", params.Email, user.Email)
	}
}

func (suite *UserHandlerSuite) TestGetByID() {

	var (
		GET_BY_ID_ROUTE = "/users/:id"
		fn              = faker.FirstName()
		ln              = faker.LastName()
		email           = faker.Email()
		password        = faker.Password()
		isAdmin         = false
		ctx             = context.TODO()
	)

	newUser, err := userfixtures.NextWith(fn, ln, email, password, isAdmin)
	if err != nil {
		suite.T().Fatalf("error generating user")
	}
	expected, err := suite.store.User.InsertUser(ctx, newUser)
	if err != nil {
		suite.T().Errorf("err occured while inserting user %s", err.Error())
	}

	app := fiber.New()
	app.Get(GET_BY_ID_ROUTE, suite.userHandler.HandleGetUser)

	req := httptest.NewRequest("GET", strings.Replace(GET_BY_ID_ROUTE, ":id", expected.ID.Hex(), -1), nil)
	resp, _ := app.Test(req)

	var actual *types.User
	json.NewDecoder(resp.Body).Decode(&actual)

	if actual.ID.Hex() != expected.ID.Hex() {
		suite.T().Errorf("get by id failed expected %s got %s", expected.ID.Hex(), actual.ID.Hex())
	}
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerSuite))
}
