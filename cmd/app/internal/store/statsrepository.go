package store

import "go-flow2sqlite/cmd/app/internal/models"

type StatRepository interface {
	Create(*models.StatInDB) error
	FindSizeBetweenDate(string, string) (map[string]models.StatDeviceType, error)
	DeletingDateData(string) error
	Save([]*models.StatInDB, uint16) error
}
