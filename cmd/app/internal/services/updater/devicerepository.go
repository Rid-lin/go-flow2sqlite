package updater

import (
	"go-flow2sqlite/cmd/app/internal/models"
)

type DeviceUpdater interface {
	GetAll(opt *Options) (*models.Devices, error)
	// Create(*models.DeviceOfMikrotik) error
	// Find(int) (*model.User, error)
	// FindByEmail(string) (*model.User, error)
	// DeleteByEmail(string) error
}
