package updater

import "time"

type Store interface {
	Device() DeviceRepository
}
type Options struct {
	Timeout time.Duration
	retry   uint
	Address string
}
