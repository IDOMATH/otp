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

	statement := `INSERT INTO otp (password, expires_at) values ($1, $2) RETURNING id`

	var row OtpRow

	err := s.Db.QueryRowContext(ctx, statement, otp, time.Now().Add(time.Hour).Unix()).Scan(&row)

	if err != nil {
		return 0
	}

	return row.Id
}

func (s *OtpStore) GetOtp(id int) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	statement := `SELECT password FROM otp WHERE id = ($1)`

	var pass string

	err := s.Db.QueryRowContext(ctx, statement, id).Scan(&pass)

	if err != nil {
		return ""
	}

	return pass
}

func (s *OtpStore) DeleteOtp(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	statement := `DELETE FROM otp WHERE id = ($1)`

	_, err := s.Db.ExecContext(ctx, statement, id)
	return err
}

func (s *OtpStore) CleanupExpired() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	statement := `DELETE FROM otp WHERE expires_at < ($1)`
	_, err := s.Db.ExecContext(ctx, statement, time.Now().Unix())
	return err
}

func (s *OtpStore) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		_, ok := <-ticker.C
		if !ok {
			break
		}

		s.CleanupExpired()

	}
}
