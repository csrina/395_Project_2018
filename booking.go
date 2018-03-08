package main

import (
	"errors"
	"time"
)

// bookingBlock - struct matching a booking record
//noinspection GoUnusedType
type bookingBlock struct {
	BookingID int       `db:"booking_id" json:"bookingId"`
	BlockID   int       `db:"block_id" json:"blockId"`
	FamilyID  int       `db:"family_id" json:"familyId"`
	UserID    int       `db:"user_id" json:"userID"`
	Start     time.Time `db:"booking_start" json:"start"`
	End       time.Time `db:"booking_end" json:"end"`
}

// inserts a booking into the database
func (b *bookingBlock) insertBooking() error {
	q := `INSERT INTO booking (block_id, family_id, user_id, 
			booking_start, booking_end) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING booking_id`

	err := db.QueryRow(q, b.BlockID, b.FamilyID, b.UserID, b.Start, b.End).Scan(&b.BookingID)
	if err != nil {
		return err
	}
	return nil
}

// updates an existing booking in the database
func (b *bookingBlock) updateBooking() error {
	q := `UPDATE booking SET booking_start = $2, booking_end = $3
			WHERE booking_id = $3`

	results, err := db.Exec(q, b.BlockID, b.Start, b.End)
	if err != nil {
		return err
	}
	count, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		err = errors.New("Booking not updated")
		return err
	}
	return nil
}

// deletes an existing booking in the database
func (b *bookingBlock) deleteBooking() error {
	q := `DELETE FROM booking WHERE bookind_id = $1`

	results, err := db.Exec(q, b.BlockID, b.Start, b.End)
	if err != nil {
		return err
	}
	count, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		err = errors.New("Booking not deleted")
		return err
	}
	return nil
}

/* Joined relation of time_block and booking_block */
type Booking struct {
	BookingID 	int       	`db:"booking_id" json:"bookingId"`
	BlockID   	int       	`db:"block_id" json:"blockId"`
	FamilyID  	int       	`db:"family_id" json:"familyId"`
	UserID    	int       	`db:"user_id" json:"userID"`
	Start     	time.Time 	`db:"booking_start" json:"start"`
	End       	time.Time 	`db:"booking_end" json:"end"`
	BlockStart 	time.Time 	`db:"block_start" json:"endBlock"`
	BlockEnd   	time.Time 	`db:"block_end" json:"endBlock"`
	RoomID		int			`db:"room_id" json:"room_id"`
	Modifier    int         `db:"modifier" json:"modifier"`
	Note     	string    	`db:"note" json:"note"`
}

/* WHY */
func getUserBookings(start time.Time, end time.Time, UID int) ([]Booking, error) {
	/* Get all bookings in range start-now */
	q := `SELECT * FROM booking NATURAL JOIN time_block
			WHERE (time_block.block_id = booking.block_id
					AND booking.user_id = $1 
					AND time_block.block_start >= $2 
					AND time_block.block_end <= $3)
			ORDER BY time_block.block_start`

	var bookBlocks []Booking
	err := db.Select(&bookBlocks, q, UID, start, end)
	logger.Println("uid ", UID, "start ", start, " end ", end, "\nblocks: ", bookBlocks)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return bookBlocks, nil
}
