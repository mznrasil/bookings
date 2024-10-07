package dbrepo

import (
	"context"
	"time"

	"github.com/mznrasil/bookings/internal/models"
)

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `
		INSERT INTO reservations
			(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO room_restrictions
			(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) SearchRoomAvailabilityByDates(
	startDate, endDate time.Time,
	roomID int,
) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT COUNT(id)
		FROM room_restrictions
		WHERE
			room_id = $1 AND
			$2 < end_date AND $3 > start_date
		`

	var count int
	rows := m.DB.QueryRowContext(ctx, query, roomID, startDate, endDate)
	err := rows.Scan(&count)
	if err != nil {
		return false, err
	}

	// if count is 0, then there is availability
	return count == 0, nil
}

func (m *postgresDBRepo) SearchAvailabilityForAllRooms(
	startDate, endDate time.Time,
) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT r.id, r.room_name
		FROM rooms r
		WHERE r.id NOT IN (
			SELECT rr.room_id
			FROM room_restrictions rr
			WHERE $1 < end_date AND $2 > start_date
		)
	`

	var rooms []models.Room

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.RoomName); err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}
