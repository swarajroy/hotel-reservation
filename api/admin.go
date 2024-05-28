package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnAuthenticated()
	}
	if !user.IsAdmin {
		return ErrUnAuthorized()
	}
	if err := c.Next(); err != nil {
		return err
	}
	return nil
}
