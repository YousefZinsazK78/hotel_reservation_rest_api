package api

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/hotel_reservation/db/fixtures"
	"github.com/yousefzinsazk78/hotel_reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "jimmy", "watercooler", false)
		user           = fixtures.AddUser(db.Store, "yousef", "yousef", false)
		hotel          = fixtures.AddHotel(db.Store, "test hotel", "iran", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 44.44, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		bookingHandler = NewBookingHandler(db.Store)
		app            = fiber.New()
		route          = app.Group("/", JWTAuthentication(db.User))
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	have := bookingResp
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
	}

	///test not admin user
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "yousef", "yousef", false)
		hotel          = fixtures.AddHotel(db.Store, "test hotel", "iran", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 44.44, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		bookingHandler = NewBookingHandler(db.Store)
		app            = fiber.New(
			fiber.Config{
				ErrorHandler: ErrorHandler,
			},
		)
		admin = app.Group("/", JWTAuthentication(db.User), AdminAuth)
	)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
	}

	///test not admin user
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized but got %d", resp.StatusCode)
	}
}
