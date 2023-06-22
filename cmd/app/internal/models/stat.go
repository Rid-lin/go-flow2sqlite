package models

import "time"

type StatInDB struct {
	ID                     int
	Date                   string
	Year, Month, Day, Hour int
	Size                   uint64
	Login, Ipaddress       string
}

type StatType struct {
	PerHour         [24]uint64
	Precent         float64
	SizeOfPrecentil uint64
	Average         uint64
	Size            uint64
	Count           uint32
}

type StatDeviceType struct {
	Mac string
	Ip  string
	StatType
}

type DayStatType struct {
	Ip    string
	Year  int
	Month time.Month
	Day   int
	StatType
}
