package updater

import "time"

type Options struct {
	timeout time.Duration
	retry   uint
	address string
}

type Option func(*Options)

func Timeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.timeout = timeout
	}
}

func Retry(retry uint) Option {
	return func(options *Options) {
		options.retry = retry
	}
}

type Store interface {
	Device() DeviceUpdater
}
