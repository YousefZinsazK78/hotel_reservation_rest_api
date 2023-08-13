package api

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/yousefzinsazk78/hotel_reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	dbname := os.Getenv(db.MongoDbNameEnvName)
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	dburi := os.Getenv("MONGO_DB_URL_TEST")
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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
