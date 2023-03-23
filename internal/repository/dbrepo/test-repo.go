package dbrepo

import (
	"errors"
	"time"

	"github.com/aysf/bwago/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}

	return 1, nil
}

// InsertRoomRestriction inserts a room_restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	if r.RoomID == 999 {
		return errors.New("some error")
	}

	return nil
}

// SearchAvailabilityByDates returns true if availability exists for roomID, and false if no availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID == 999 {
		return false, errors.New("some error")
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	layout := "2006-01-02"
	t, _ := time.Parse(layout, "2099-12-31")

	if t.Equal(start) {
		return rooms, errors.New("some error")
	}

	t, _ = time.Parse(layout, "2030-12-30")

	if t.Before(start) {
		rooms = []models.Room{
			{ID: 1, RoomName: "test room", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
	}

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	return models.User{}, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 0, "", nil
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservation []models.Reservation

	return reservation, nil

}

// AllNewReservations returns a slice of all reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservation []models.Reservation

	return reservation, nil

}

// GetReservationByID returns one reservation by ID
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var res models.Reservation

	return res, nil
}

// UpdateReservation updates a reservation in the database
func (m *testDBRepo) UpdateReservation(u models.Reservation) error {

	return nil
}

// DeleteReservation deletes one reservation by id
func (m *testDBRepo) DeleteReservation(id int) error {

	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}
