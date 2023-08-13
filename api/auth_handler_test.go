package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/hotel_reservation/db/fixtures"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	//insertUser := insertTestUser(t, tdb.UserStore)
	fixtures.AddUser(tdb.Store, "yousef", "yousef", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "yousef@yousef.com",
		Password: "supersecurepassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}
	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response msg to be invalid credentials but got %s", genResp.Msg)
	}
}

func TestAuthenticate(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	//insertUser := insertTestUser(t, tdb.User)
	insertedUser := fixtures.AddUser(tdb.Store, "yousef", "yousef", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "yousef@yousef.com",
		Password: "yousef_yousef",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}

	//set the encrypted password to an empty string, because we dont return that in any response....
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		fmt.Println(insertedUser)
		fmt.Println(authResp.User)
		t.Fatalf("expected the user to be the inserted user.")
	}

}
