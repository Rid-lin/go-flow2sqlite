package models_test

import (
	"testing"

	"go-flow2sqlite/cmd/app/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestTesting_TestStat(t *testing.T) {
	assert.Equal(t, models.TestStat(t),
		&models.StatInDB{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      4,
			Size:      123456789,
			Login:     "10.61.129.53",
			Ipaddress: "10.61.129.53",
		})
}
