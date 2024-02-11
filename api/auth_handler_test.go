package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/types"
)

func InsertTestUser(t *testing.T, fname, lname, email, password string, store *db.HotelReservationStore, c context.Context) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		t.Fatal(err)
	}
	res, err := store.User.InsertUser(c, user)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func TestHandleAuthenticateSuccess(t *testing.T) {

	var POST_ROUTE = "/auth"
	ctx := context.TODO()
	tdb := Setup(t, ctx)
	defer tdb.TearDown(t, ctx)
	insertedUser := InsertTestUser(t, "James", "Foo", "james@foo.com", "supersecurepassword", tdb.store, ctx)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store)
	app.Post(POST_ROUTE, authHandler.HandleAuth)

	params := AuthParams{
		Email:    insertedUser.Email,
		Password: "supersecurepassword",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", POST_ROUTE, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d status got %d status", http.StatusOK, resp.StatusCode)
	}

	var authResponse AuthResponse
	json.NewDecoder(resp.Body).Decode(&authResponse)

	if authResponse.Token == "" {
		t.Fatalf("token expected in the response got %s", authResponse.Token)
	}

	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResponse.User) {
		t.Fatalf("expected insertedUser to be equal to authResponse User")
	}

}

func TestHandleAuthenticateFailure(t *testing.T) {

	var POST_ROUTE = "/auth"
	ctx := context.TODO()
	tdb := Setup(t, ctx)
	defer tdb.TearDown(t, ctx)
	insertedUser := InsertTestUser(t, "James", "Foo", "james@foo.com", "supersecurepassword", tdb.store, ctx)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store)
	app.Post(POST_ROUTE, authHandler.HandleAuth)

	params := AuthParams{
		Email:    insertedUser.Email,
		Password: "notsupersecurepassword",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", POST_ROUTE, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d status got %d status", http.StatusBadRequest, resp.StatusCode)
	}

	var authErrorResponse AuthErrorResponse
	json.NewDecoder(resp.Body).Decode(&authErrorResponse)

	if authErrorResponse.Msg != "Bad Request" {
		t.Fatalf("Msg expected in the response should be 'Bad Request' got %s", authErrorResponse.Msg)
	}
}
