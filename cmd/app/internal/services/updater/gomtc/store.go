package updatergomtc

import (
	"go-flow2sqlite/cmd/app/internal/models"
	updater "go-flow2sqlite/cmd/app/internal/services/updater"
	"time"
)

type Store struct {
	opt               *updater.Options
	deviceRepository  *DeviceRepository
	TimerUpdatedevice *time.Timer
}

func NewUpdater(addrPort string) *Store {
	return &Store{
		opt: &updater.Options{
			Address: addrPort,
		},
	}
}

func (s *Store) Device() updater.DeviceRepository {
	if s.deviceRepository != nil {
		return s.deviceRepository
	}

	s.deviceRepository = &DeviceRepository{
		store:   s,
		devices: make([]*models.DeviceType, 0),
	}
	return s.deviceRepository
}
