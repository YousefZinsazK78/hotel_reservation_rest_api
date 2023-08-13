package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yousefzinsazk78/hotel_reservation/db"
	"net/http"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrUnAuthorized()
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		iat := claims["iat"].(float64)
		iatInt := int64(iat)
		if time.Now().Unix() > iatInt {
			return NewError(http.StatusUnauthorized, "token Expired!")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrUnAuthorized()
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to pase JWT token:", err)
		return nil, ErrUnAuthorized()
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, ErrUnAuthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnAuthorized()
	}
	return claims, nil
}
