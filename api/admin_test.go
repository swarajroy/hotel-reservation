package api

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/swarajroy/hotel-reservation/types"
	"github.com/valyala/fasthttp"
)

func TestShouldFailWhenUserIsNotAnAdminSetInContext(t *testing.T) {
	app := fiber.New()

	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	c.Context().SetUserValue("user", &types.User{
		IsAdmin: false,
	})

	err := AdminAuth(c)

	if err != nil {
		assert.ErrorIs(t, err, ErrUnAuthorized())
	}
}

func TestShouldFailWhenUserNotSetInContext(t *testing.T) {
	app := fiber.New()

	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	err := AdminAuth(c)

	if err != nil {
		assert.ErrorIs(t, err, ErrUnAuthorized())
	}
}
