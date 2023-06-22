package updater

import (
	"go-flow2sqlite/cmd/app/internal/models"
	"time"
)

type DeviceRepository interface {
	GetAll(addr string, retry uint) ([]*models.DeviceType, error)
	RunOnce(addr, interval string, retry uint)
	SetTimerUpdateDevice(IntervalStr string)
	GetTimer() *time.Timer
}
