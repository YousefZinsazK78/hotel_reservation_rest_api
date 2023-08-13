package api

import (
	"context"
	"github.com/yousefzinsazk78/hotel_reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

const (
	dburi      = "mongodb://localhost:27017"
	testDbName = "hotel-reservation-test"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		t.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookStore(client),
		},
	}
}
