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

func TestUserRepository_FindSizeBetweenDate(t *testing.T) {
	db, teardown := sqlitestore.TestDB(t, dsn)
	defer teardown("statistics")
	// db, _ := sqlitestore.TestDB(t, dsn)

	s := sqlitestore.NewStore(db)

	stat := &models.StatInDB{
		Date:      "2021-11-17",
		Year:      2021,
		Month:     11,
		Day:       17,
		Hour:      4,
		Size:      123456789,
		Login:     "10.61.129.241",
		Ipaddress: "10.61.129.241",
	}
	ds, err := s.Stat().FindSizeBetweenDate(stat.Date, stat.Date)
	assert.NoError(t, err)
	dst := map[string]models.StatDeviceType{}
	assert.Equal(t, dst, ds)

	assert.NoError(t, s.Stat().Create(stat))

	d, err := s.Stat().FindSizeBetweenDate(stat.Date, stat.Date)
	assert.NoError(t, err)
	assert.NotNil(t, d)
}

// func TestUserRepository_SaveSample(t *testing.T) {
// 	type Sample struct {
// 		Foo, Bar, Baz string
// 	}

// 	samples := []*Sample{{"a", "b", "c"}, {"d", "e", "f"}, {"g", "h", "i"}, {"j", "k", "k"}}

// 	insertsql := `insert into samples (foo, bar, baz) values `
// 	query := insertsql
// 	bufLine := 6
// 	values := []interface{}{}
// 	k := 0
// 	for _, s := range samples {
// 		numFields := 3
// 		n := k * numFields
// 		if int(bufLine)-n < numFields {
// 			query = query[:len(query)-1] // remove the trailing comma
// 			fmt.Println(query)
// 			fmt.Println(values)
// 			query = insertsql
// 			k = 0
// 			n = k * numFields
// 			values = []interface{}{}
// 		}
// 		values = append(values, s.Foo, s.Bar, s.Baz)

// 		query += `(`
// 		for j := 0; j < numFields; j++ {
// 			query += `$` + strconv.Itoa(n+j+1) + `,`
// 		}
// 		query = query[:len(query)-1] + `),`

// 		k++
// 	}
// 	query = query[:len(query)-1]
// 	fmt.Println(k)
// 	fmt.Println(query)
// 	fmt.Println(values)

// }

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

func TestUserRepository_DeletingDateData(t *testing.T) {
	db, teardown := sqlitestore.TestDB(t, dsn)
	defer teardown("statistics")

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

	err := s.Stat().Save(samples, uint16(len(samples)))
	assert.NoError(t, err)
	d, err := s.Stat().FindSizeBetweenDate(samples[0].Date, samples[0].Date)
	assert.NoError(t, err)
	assert.NotNil(t, d)
	err = s.Stat().DeletingDateData(samples[0].Date)
	assert.NoError(t, err)
	d, err = s.Stat().FindSizeBetweenDate(samples[0].Date, samples[0].Date)
	assert.NoError(t, err)
	assert.Empty(t, d)

}
