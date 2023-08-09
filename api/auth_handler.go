package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yousefzinsazk78/hotel_reservation/db"
	"github.com/yousefzinsazk78/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
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

type genericResp struct {
	Type string `json:"status"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authparams AuthParams
	if err := c.BodyParser(&authparams); err != nil {
		return err
	}

	user, err := h.userStore.GetUsersByEmail(c.Context(), authparams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsPasswordValid(user.EncryptedPassword, authparams.Password) {
		return invalidCredentials(c)
	}

	authResp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}

	return c.JSON(authResp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	invalidAt := &jwt.NumericDate{Time: now.Add(time.Hour * 4)}
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["iat"] = invalidAt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret")
	}
	return tokenString
}
