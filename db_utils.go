package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectToDb(config Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", config.DbConString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CleanupDb(config Config) error {
	db, err := ConnectToDb(config)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"TRUNCATE TABLE sensors_data;")
	return err
}
