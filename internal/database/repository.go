package database

import "fajntvajb/internal/logger"

func (db *DB) GetRows() ([]string, error) {
	log := logger.Get()
	var res []string
	err := db.connection.Select(&res, "SELECT name FROM test_table")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get rows")
		return nil, err
	}
	return res, nil
}

func (db *DB) InsertRow(name string) error {
	log := logger.Get()
	_, err := db.connection.Exec("INSERT INTO test_table (name) VALUES ($1)", name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert row")
		return err
	}
	return nil
}
