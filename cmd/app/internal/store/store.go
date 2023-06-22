package store

type Store interface {
	Stat() StatRepository
}
