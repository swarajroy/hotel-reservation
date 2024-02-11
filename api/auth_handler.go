package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	store *db.HotelReservationStore
}

func NewAuthHandler(store *db.HotelReservationStore) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type AuthErrorResponse struct {
	Status int
	Msg    string
}

func invalidCreds(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(AuthErrorResponse{
		Status: 400,
		Msg:    "Bad Request",
	})
}
func (auth *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var params AuthParams

	if err := c.BodyParser(&params); err != nil {
		log.Errorf("error occurred err = ", err)
		return invalidCreds(c)
	}

	user, err := auth.store.User.GetUserByEmail(c.Context(), params.Email)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCreds(c)
		}
		return err
	}

	if !types.IsValisPassword(user.EncryptedPassword, params.Password) {
		return invalidCreds(c)
	}

	log.Info("authenticated user = ", user)

	authResp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}

	return c.JSON(authResp)
}

func createTokenFromUser(u *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 3).Unix()
	claims := jwt.MapClaims{
		"id":      u.ID,
		"email":   u.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret = ", err)
	}
	return tokenStr
}
