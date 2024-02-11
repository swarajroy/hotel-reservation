package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/swarajroy/hotel-reservation/db"
)

func JWTAuthentication(store db.HotelReservationStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Api-Token")
		if len(token) == 0 {
			return fmt.Errorf("Unauthorized")
		}

		claims, err := validateToken(token)
		if err != nil {
			fmt.Println("error = ")
			return err
		}

		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		if time.Now().Unix() > expires {
			return fmt.Errorf("token has expired")
		}

		userID := claims["id"].(string)
		user, err := store.User.GetUserById(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("Unauthorized")
		}

		//set the current authenticated user in the context
		c.Context().SetUserValue("user", user)

		return c.Next()
	}

}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	fmt.Println("parseToken = ", token)
	if err != nil {
		fmt.Errorf("failed to parse jwt token = %s", tokenStr)
		return nil, fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("unauthorized")
}
