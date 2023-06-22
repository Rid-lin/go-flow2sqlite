package sqlitestore_test

import (
	"go-flow2sqlite/cmd/app/internal/models"
	"testing"

	"go-flow2sqlite/cmd/app/internal/store/sqlitestore"

	"github.com/stretchr/testify/assert"
)

func TestStatRepository_Create(t *testing.T) {
	db, teardown := sqlitestore.TestDB(t, dsn)
	defer teardown("statistics")

	s := sqlitestore.NewStore(db)

	d := models.TestStat(t)

	assert.NoError(t, s.Stat().Create(d))
	assert.NotNil(t, d)

}

func TestStatRepository_Save(t *testing.T) {
	db, teardown := sqlitestore.TestDB(t, dsn)
	defer teardown("statistics")

	// db, _ := sqlitestore.TestDB(t, dsn)

	samples := []*models.StatInDB{
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      0,
			Size:      1000,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      1,
			Size:      1001,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      2,
			Size:      1002,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      3,
			Size:      1003,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      4,
			Size:      1004,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      5,
			Size:      1005,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      6,
			Size:      1006,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      7,
			Size:      1007,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      8,
			Size:      1008,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
		{
			Date:      "2021-11-17",
			Year:      2021,
			Month:     11,
			Day:       17,
			Hour:      9,
			Size:      1009,
			Login:     "10.61.129.241",
			Ipaddress: "10.61.129.241",
		},
	}

	s := sqlitestore.NewStore(db)

	for i := range samples {
		err := s.Stat().Save(samples, uint16(i+1))
		assert.NoError(t, err)
	}
	d, err := s.Stat().FindSizeBetweenDate(samples[0].Date, samples[0].Date)
	assert.NoError(t, err)
	assert.NotNil(t, d)
}
