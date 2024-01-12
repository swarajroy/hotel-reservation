package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found!"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	log.Info("Enter HandleGetUsers")
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		log.Error("error occurred")
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	log.Info("Enter HandlePostUser")
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	log.Info("Before NewUserFromParams")
	user, err := types.NewUserFromParams(params)
	log.Info("user = ", &user)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	log.Info("insertedUser = ", &insertedUser)

	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.userStore.DeleteUserById(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"Deleted": userID})
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var params types.UpdateUserParams
	userID := c.Params("id")

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := h.userStore.UpdateUserById(c.Context(), params, userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"Updated": userID})
}
