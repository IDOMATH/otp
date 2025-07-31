package db

import (
	"context"
	"database/sql"
	"time"
)

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

	statement := `INSERT INTO otp (password) values ($1)`

	return 0
}
