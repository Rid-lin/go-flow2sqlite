package sqlitestore

import (
	"go-flow2sqlite/cmd/app/internal/store"
)

type Store struct {
	dsn            string
	statRepository *StatRepository
}

func NewStore(dsn string) *Store {
	return &Store{
		dsn: dsn,
	}
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
