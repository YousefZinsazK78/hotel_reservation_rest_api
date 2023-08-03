package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/hotel_reservation/api"
	"github.com/yousefzinsazk78/hotel_reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel_reservation"
const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("ListenAddr", ":5000", "the ListenAddr of the api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	//handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)

	apiv1 := app.Group("/api/v1")
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	log.Fatal(app.Listen(*listenAddr))
}
