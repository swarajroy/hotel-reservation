package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/db/mongo"
	userfixtures "github.com/swarajroy/hotel-reservation/types/user_fixtures"
)

type AuthHandlerSuite struct {
	suite.Suite
	store           *db.HotelReservationStore
	testMongoClient *mongo.TestMongoClient
	authHandler     *AuthHandler
}

func (suite *AuthHandlerSuite) SetupSuite() {
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
	suite.authHandler = NewAuthHandler(suite.store)
}

func (suite *AuthHandlerSuite) TearDownSuite() {
	suite.testMongoClient.Container.Terminate(context.Background())
}

func (suite *AuthHandlerSuite) AfterTest() {
	suite.store.User.Drop(context.Background())
}

func (suite *AuthHandlerSuite) TestHandleAuthenticateSuccess() {

	var (
		POST_ROUTE = "/auth"
		fn         = faker.FirstName()
		ln         = faker.LastName()
		email      = faker.Email()
		password   = faker.Password()
		isAdmin    = false
		ctx        = context.TODO()
	)

	user, err := userfixtures.NextWith(fn, ln, email, password, isAdmin)
	if err != nil {
		suite.T().Fatalf("error generating user")
	}
	insertedUser, _ := suite.store.User.InsertUser(ctx, user)

	app := fiber.New()
	app.Post(POST_ROUTE, suite.authHandler.HandleAuth)

	params := AuthParams{
		Email:    email,
		Password: password,
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", POST_ROUTE, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		suite.T().Fatalf("expected %d status got %d status", http.StatusOK, resp.StatusCode)
	}

	var authResponse AuthResponse
	json.NewDecoder(resp.Body).Decode(&authResponse)

	if authResponse.Token == "" {
		suite.T().Fatalf("token expected in the response got %s", authResponse.Token)
	}

	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResponse.User) {
		suite.T().Fatalf("expected insertedUser to be equal to authResponse User")
	}

}

func (suite *AuthHandlerSuite) TestHandleAuthenticateFailure() {

	var (
		POST_ROUTE = "/auth"
		fn         = faker.FirstName()
		ln         = faker.LastName()
		email      = faker.Email()
		password   = faker.Password()
		isAdmin    = false
		ctx        = context.TODO()
	)

	user, err := userfixtures.NextWith(fn, ln, email, password, isAdmin)
	if err != nil {
		suite.T().Fatalf("error generating user")
	}
	_, err = suite.store.User.InsertUser(ctx, user)
	if err != nil {
		suite.T().Errorf("err occured while inserting user %s", err.Error())
	}

	app := fiber.New()
	app.Post(POST_ROUTE, suite.authHandler.HandleAuth)

	params := AuthParams{
		Email:    email,
		Password: faker.Password(),
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", POST_ROUTE, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		suite.T().Fatalf("expected %d status got %d status", http.StatusBadRequest, resp.StatusCode)
	}

	var authErrorResponse AuthErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&authErrorResponse)
	if err != nil {
		suite.T().Fatalf("decoding the auth error response failed")
	}

	if authErrorResponse.Msg != "Bad Request" {
		suite.T().Fatalf("Msg expected in the response should be 'Bad Request' got %s", authErrorResponse.Msg)
	}
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerSuite))
}
