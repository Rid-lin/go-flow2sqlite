package store

import (
	"go-flow2sqlite/cmd/app/internal/models"
)

type StatRepository interface {
	// Create(*models.StatInDB) error
	// FindSizeBetweenDate(string, string) (map[string]models.StatDeviceType, error)
	// DeletingDateData(string) error
	Save(chan models.BMes) error
	// SaveOneLine(*models.BMes) error
}
