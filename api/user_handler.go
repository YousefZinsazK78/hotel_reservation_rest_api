package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/hotel_reservation/types"
)

func HandleListUser(ctx *fiber.Ctx) error {
	u := types.User{
		ID:        "0",
		FirstName: "yousef",
		LastName:  "zinsaz",
	}
	return ctx.JSON(u)
}

func HandleGetUser(ctx *fiber.Ctx) error {
	return ctx.JSON("james")
}
