package db

const (
	DBNAME       = "hotel-reservation"
	TEST_DB_NAME = "hotel-reservation-test"
	DB_URI       = "mongodb://127.0.0.1:27017"
)

type HotelReservationStore struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}

func NewHotelReservationStore(user UserStore, hotel HotelStore, room RoomStore) *HotelReservationStore {
	return &HotelReservationStore{
		User:  user,
		Hotel: hotel,
		Room:  room,
	}
}
