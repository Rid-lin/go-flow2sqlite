package teststore

import (
	"go-flow2sqlite/cmd/app/internal/models"

	"go-flow2sqlite/cmd/app/internal/store"
)

type Store struct {
	statRepository *StatRepository
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Stat() store.StatRepository {
	if s.statRepository != nil {
		return s.statRepository
	}

	s.statRepository = &StatRepository{
		store: s,
		stats: []*models.BMes{},
	}
	return s.statRepository
}
