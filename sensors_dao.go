package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var ErrNotInitialized = errors.New("DB not initialized")

type SensorsDao interface {
	Init() error
	WriteMeasures(timestamp time.Time, measures []MeasurePayload) error
}

func NewSensorsDao(config Config) SensorsDao {
	return &sensorsDao{config: config}
}

type sensorsDao struct {
	config Config
	db     *sqlx.DB
}

func (s *sensorsDao) Init() error {
	db, err := ConnectToDb(s.config)
	if err != nil {
		return err
	}
	s.db = db

	migrator := NewMigrator()
	if err := migrator.Migrate(s.config); err != nil {
		return err
	}

	return nil
}

func (s *sensorsDao) WriteMeasures(timestamp time.Time, measures []MeasurePayload) error {
	if err := s.checkIsInitialized(); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return errors.Wrap(err, "error in begin transaction")
	}

	stmt, err := tx.Prepare("INSERT INTO sensors_data(t, sensor, value) VALUES ($1, $2, $3);")
	if err != nil {
		return errors.Wrap(err, "error in query preparation")
	}

	for _, measure := range measures {
		if _, err := stmt.Exec(timestamp, measure.Sensor, measure.Value); err != nil {
			return errors.Wrap(err, "error in execute statement")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error in commit transaction")
	}
	if err := stmt.Close(); err != nil {
		return errors.Wrap(err, "error in close statement")
	}

	return nil
}

func (s *sensorsDao) checkIsInitialized() error {
	if !s.isInitialized() {
		return ErrNotInitialized
	}

	return nil
}

func (s *sensorsDao) isInitialized() bool {
	return s.db != nil
}
