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

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	listenAddr := flag.String("ListenAddr", ":5000", "the ListenAddr of the api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBuri))
	if err != nil {
		log.Fatal(err)
	}

	//handlers initialization
	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookStore(client)
		store        = &db.Store{
			User:    userStore,
			Room:    roomStore,
			Hotel:   hotelStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", api.AdminAuth)
	)

	//auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	///user handler
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	///hotel handler
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)

	///room handler
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	//bookings handler
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBookings)

	//admin handler
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	log.Fatal(app.Listen(*listenAddr))
}
