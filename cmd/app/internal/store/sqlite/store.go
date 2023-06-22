package sqlitestore

import (
	"database/sql"
	"go-flow2sqlite/cmd/app/internal/store"
)

type Store struct {
	db             *sql.DB
	statRepository *StatRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Store) Stat() store.StatRepository {
	if s.statRepository != nil {
		return s.statRepository
	}

	s.statRepository = &StatRepository{
		store: s,
	}
	return s.statRepository
}
