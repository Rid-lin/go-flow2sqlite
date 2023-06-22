package models

import (
	"testing"
)

func TestStat(t *testing.T) *StatInDB {
	return &StatInDB{
		Date:      "2021-11-17",
		Year:      2021,
		Month:     11,
		Day:       17,
		Hour:      4,
		Size:      123456789,
		Login:     "10.61.129.53",
		Ipaddress: "10.61.129.53",
	}
}
