package models_test

import (
	"go-flow2sqlite/cmd/app/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSumm(t *testing.T) {
	stat := models.StatType{
		PerHour:         [24]uint64{1234567890, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Precent:         50.0,
		SizeOfPrecentil: 50.0,
		Average:         100000,
		Size:            1234567890,
		// VolumePerCheck:  1234567890,
		Count: 20,
	}
	sd := models.StatDevice{
		PerHour: [24]uint64{0, 123456789, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Mac:     "10.61.129.241",
		Ip:      "10.61.129.241",
		Size:    123456789,
		// VolumePerCheck: 123456789,
		Count: 20,
		Year:  2021,
		Month: time.Month(11),
		Day:   17,
		Hour:  2,
	}
	statAfterSumm := models.StatType{
		PerHour:         [24]uint64{1234567890, 123456789, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Precent:         50.0,
		SizeOfPrecentil: 50.0,
		Average:         100000,
		Size:            1358024679,
		// VolumePerCheck:  1234567890,
		Count: 21,
	}
	stat.Summ(&sd)
	assert.Equal(t, statAfterSumm, stat)
}
