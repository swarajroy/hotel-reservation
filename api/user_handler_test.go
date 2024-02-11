package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	store *db.HotelReservationStore
}

func Setup(t *testing.T, ctx context.Context) *testdb {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		t.Fatal(err)
	}
	userStore := db.NewMongoDbUserStore(client, db.TEST_DB_NAME)
	return &testdb{
		store: &db.HotelReservationStore{
			User: userStore,
		},
	}

}

func (tdb *testdb) TearDown(t *testing.T, ctx context.Context) {
	if err := tdb.store.User.Drop(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	var POST_ROUTE = "/users"
	ctx := context.TODO()
	tdb := Setup(t, ctx)
	defer tdb.TearDown(t, ctx)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store)
	app.Post(POST_ROUTE, userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Swaraj",
		LastName:  "Roy",
		Email:     "sroy@golang.org",
		Password:  "sedfdsas",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", POST_ROUTE, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expecting userId to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting no encrypted password to be set in the response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s got %s", params.Email, user.Email)
	}
}

/* func TestGetByID(t *testing.T) {
	ctx := context.TODO()
	tdb := Setup(t, ctx)
	defer tdb.TearDown(t, ctx)

	expected, _ := tdb.UserStore.InsertUser(ctx, &types.User{FirstName: "Swaraj", LastName: "Roy", Email: "s.roy@gmail.com", EncryptedPassword: "asdsadsadsadf"})
	var GET_BY_ID_ROUTE = "/users/" + (expected.ID).Hex()
	t.Log(GET_BY_ID_ROUTE)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Get(GET_BY_ID_ROUTE, userHandler.HandleGetUser)

	req := httptest.NewRequest("GET", GET_BY_ID_ROUTE, nil)

	resp, _ := app.Test(req)

	var actual types.User
	json.NewDecoder(resp.Body).Decode(&actual)

	if actual.ID != expected.ID {
		t.Errorf("get by id failed expected %s got %s", expected.ID.String(), actual.ID)
	}
} */
