package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/db/fixtures"
)

func TestHandleAuthenticateSuccess(t *testing.T) {

	var POST_ROUTE = "/auth"
	ctx := context.TODO()
	tdb := Setup(t, ctx)
	defer tdb.TearDown(t, ctx)

	insertedUser := fixtures.AddUser(tdb.store, "james", "foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store)
	app.Post(POST_ROUTE, authHandler.HandleAuth)

	params := AuthParams{
		Email:    fmt.Sprintf("%s_%s@foo.com", insertedUser.FirstName, insertedUser.LastName),
		Password: fmt.Sprintf("%s_%s", insertedUser.FirstName, insertedUser.LastName),
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
	insertedUser := fixtures.AddUser(tdb.store, "james", "foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store)
	app.Post(POST_ROUTE, authHandler.HandleAuth)

	params := AuthParams{
		Email:    fmt.Sprintf("%s_%s@foo.com", insertedUser.FirstName, insertedUser.LastName),
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
