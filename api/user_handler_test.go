package api

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/hotel_reservation/types"
	"net/http/httptest"
	"testing"
)

func TestPostUser(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	userHandler := NewUserHandler(db.User)
	app := fiber.New()
	app.Post("/", userHandler.HandlePostUser)
	params := types.CreateUserParams{
		FirstName: "joseph",
		LastName:  "zinsaz",
		Email:     "joseph@kashani.com",
		Password:  "1231327878",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User

	_ = json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Error("expected user id to set but")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s\n", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s\n", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s\n", params.Email, user.Email)
	}
}
