package main

import (
	"database/sql"
	"net/http"
	"time"
)

//noinspection GoUnusedParameter
func calSetup(w http.ResponseWriter, r *http.Request) {
	//blocave to be filled by either json decode or from forms
	//start / end also need to be pulled from the request body
	//blocks := []Timeblock{}
	w.WriteHeader(http.StatusCreated)
}

func setAndInsertBlock(t *sql.Tx, block TimeBlock, newDate time.Time) error {
	block.setDay(newDate)
	_, err := block.insert()
	if err != nil {
		t.Rollback()
		return err
	}
	return nil
}

//noinspection GoUnusedFunction
func applyTemplate(startDate time.Time, stopDate time.Time, blocks []TimeBlock) error {
	Tx, err := db.Begin()
	if err != nil {
		return err
	}
	for startDate.Before(stopDate) {
		for _, block := range blocks {
			switch block.Start.Weekday() {
			case time.Monday:
				tempDate := startDate.AddDate(0, 0, 0)
				if err := setAndInsertBlock(Tx, block, tempDate); err != nil {
					return err
				}
				break
			case time.Tuesday:
				tempDate := startDate.AddDate(0, 0, 1)
				if err := setAndInsertBlock(Tx, block, tempDate); err != nil {
					return err
				}
				break
			case time.Wednesday:
				tempDate := startDate.AddDate(0, 0, 2)
				if err := setAndInsertBlock(Tx, block, tempDate); err != nil {
					return err
				}
				break
			case time.Thursday:
				tempDate := startDate.AddDate(0, 0, 3)
				if err := setAndInsertBlock(Tx, block, tempDate); err != nil {
					return err
				}
				break
			case time.Friday:
				tempDate := startDate.AddDate(0, 0, 4)
				if err := setAndInsertBlock(Tx, block, tempDate); err != nil {
					return err
				}
				break
			}
		}
		startDate = startDate.AddDate(0, 0, 7)
	}
	Tx.Commit()
	return nil
}

// DO NOT USE THIS UNLESS 100% REQUIRED
//noinspection GoUnusedParameter
func undoSetup(w http.ResponseWriter, r *http.Request) {
	//NEEDS AUTH CHECK
	//REQUIRES NOTIFICATIONS FOR ALL USERS ALREADY WITH BLOCKS SCHEDULED
	deleteTime := time.Now()

	q := `DELETE FROM booking WHERE booking_start > ($1)`
	q2 := `DELETE FROM time_block WHERE block_start > ($1)`

	_, err := db.Exec(q, deleteTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Fatal("Could not delete current bookings")
	}
	_, err = db.Exec(q2, deleteTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Fatal("Could not delete time blocks")
	}
	w.WriteHeader(http.StatusGone)
}
