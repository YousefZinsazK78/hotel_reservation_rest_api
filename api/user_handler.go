package api

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/hotel_reservation/db"
	"github.com/yousefzinsazk78/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *UserHandler) HandlePostUser(ctx *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return ctx.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(ctx.Context(), user)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(ctx *fiber.Ctx) error {
	var (
		//values bson.M
		params types.UpdateUserParams
		userID = ctx.Params("id")
	)
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	if err := h.userStore.UpdateUser(ctx.Context(), filter, params); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"updated": userID})
}

func (h *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if err := h.userStore.DeleteUser(ctx.Context(), userID); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"deleted": userID})
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()
	user, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(ctx.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.JSON(map[string]string{"msg": "not found"})
		}
		return err
	}
	return ctx.JSON(users)
}
