package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		//ID:        "10001",
		FirstName: "Swaraj",
		LastName:  "Roy",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	u := types.User{
		ID:        "10001",
		FirstName: "Swaraj",
		LastName:  "Roy",
	}
	return c.JSON(u)
}
