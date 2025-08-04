package db

import (
	"context"
	"database/sql"
	"time"
)

type OtpRow struct {
	Id  int
	Otp string
}

type OtpStore struct {
	Db *sql.DB
}

func NewOtpStore(db *sql.DB) *OtpStore {
	return &OtpStore{
		Db: db,
	}
}

func (s *OtpStore) InsertOtp(otp string) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	statement := `INSERT INTO otp (password) values ($1) RETURNING id`

	var row OtpRow

	err := s.Db.QueryRowContext(ctx, statement, otp).Scan(&row)

	if err != nil {
		return 0
	}

	return row.Id
}
