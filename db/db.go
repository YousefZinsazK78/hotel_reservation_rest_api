package db

const (
	DBname   = "hotel_reservation"
	DBuri    = "mongodb://localhost:27017"
	UserColl = "users"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
