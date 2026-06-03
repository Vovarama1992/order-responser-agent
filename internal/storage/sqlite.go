package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			source TEXT NOT NULL,
			order_id TEXT NOT NULL,
			url TEXT NOT NULL,
			PRIMARY KEY(source, order_id)
		)
	`)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Exists(source, orderID string) (bool, error) {
	var exists int

	err := s.db.QueryRow(`
		SELECT 1
		FROM orders
		WHERE source = ?
		  AND order_id = ?
		LIMIT 1
	`, source, orderID).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Storage) Save(source, orderID, url string) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO orders(
			source,
			order_id,
			url
		)
		VALUES (?, ?, ?)
	`, source, orderID, url)

	return err
}
