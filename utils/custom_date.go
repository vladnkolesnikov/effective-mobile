package utils

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CustomDate struct {
	pgtype.Date
}

const DateLayout = "01-2006"

// Converts date string from MM-YYYY format to Time
func (sd *CustomDate) UnmarshalJSON(b []byte) error {
	plainDate := strings.Trim(string(b), "\"")
	result, err := time.Parse(DateLayout, plainDate)

	if err != nil {
		return err
	}

	sd.Time = result
	sd.Valid = true
	sd.InfinityModifier = pgtype.Finite
	return nil
}

func (sd CustomDate) MarshalJSON() ([]byte, error) {
	if !sd.Valid {
		return []byte("null"), nil
	}

	s := sd.Time.Format(DateLayout)

	return json.Marshal(s)
}

func (sd *CustomDate) ParseQueryDate(date string, isRequred bool) error {
	if date == "" {
		if isRequred {
			return errors.New("date is required")
		}
		sd.Time = time.Now()
		sd.Valid = true
		sd.InfinityModifier = pgtype.Finite
		return nil
	}

	parsedDate, err := time.Parse(DateLayout, date)
	if err != nil {
		return err
	}

	sd.Time = parsedDate
	sd.Valid = true
	sd.InfinityModifier = pgtype.Finite
	return nil
}

func (sd *CustomDate) IsInFuture() bool {
	return time.Now().Before(sd.Time)
}
