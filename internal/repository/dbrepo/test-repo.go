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
