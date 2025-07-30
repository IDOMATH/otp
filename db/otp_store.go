package db

import "database/sql"

type OtpStore struct {
	Db *sql.DB
}

func NewOtpStore(db *sql.DB) *OtpStore {
	return &OtpStore{
		Db: db,
	}
}
