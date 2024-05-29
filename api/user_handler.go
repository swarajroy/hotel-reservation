package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/swarajroy/hotel-reservation/db"
	"github.com/swarajroy/hotel-reservation/types"
)

type UserHandler struct {
	store *db.HotelReservationStore
}

func NewUserHandler(store *db.HotelReservationStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	user, err := h.store.User.GetUserById(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	log.Info("Enter HandleGetUsers")
	users, err := h.store.User.GetUsers(c.Context())
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
		return ErrBadRequest()
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

	insertedUser, err := h.store.User.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	log.Info("insertedUser = ", &insertedUser)

	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.store.User.DeleteUserById(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"Deleted": userID})
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var params types.UpdateUserParams
	userID := c.Params("id")

	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	if err := h.store.User.UpdateUserById(c.Context(), params, userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"Updated": userID})
}
