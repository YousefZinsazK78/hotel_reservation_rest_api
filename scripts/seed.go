package main

import (
	"context"
	"fmt"
	"github.com/yousefzinsazk78/hotel_reservation/api"
	"github.com/yousefzinsazk78/hotel_reservation/db"
	"github.com/yousefzinsazk78/hotel_reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBuri))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBname).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	dbStore := db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}
	user := fixtures.AddUser(&dbStore, "yousef", "zinsaz", false)
	fmt.Println("yousef->", api.CreateTokenFromUser(user))
	adminUser := fixtures.AddUser(&dbStore, "admin", "admin", true)
	fmt.Println("admin->", api.CreateTokenFromUser(adminUser))
	hotel := fixtures.AddHotel(&dbStore, "negarestan hotel", "iran", 4, nil)
	room := fixtures.AddRoom(&dbStore, "small", true, 60.00, hotel.ID)
	booking := fixtures.AddBooking(&dbStore, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 4))
	fmt.Println("booking->", booking.ID)
}
