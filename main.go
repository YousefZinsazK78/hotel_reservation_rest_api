package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/hotel_reservation/api"
	"log"
)

func main() {
	listenAddr := flag.String("ListenAddr", ":5000", "the ListenAddr of the api server")
	flag.Parse()
	app := fiber.New()

	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user", api.HandleListUser)
	apiv1.Get("/user/:id", api.HandleGetUser)

	log.Fatal(app.Listen(*listenAddr))
}
