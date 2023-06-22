package sqlitestore

import (
	"fmt"
	"go-flow2sqlite/cmd/app/internal/models"

	"database/sql"

	_ "modernc.org/sqlite"
)

type StatRepository struct {
	store *Store
}

func (r *StatRepository) Create(stat *models.StatInDB) error {
	conn, err := sql.Open(r.store.dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.Exec(
		`INSERT INTO statistics (
			date_str, year, month, day, hour, size, login, ipaddress
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8)`,
		stat.Date, stat.Year, stat.Month, stat.Day, stat.Hour, int(stat.Size), stat.Login, stat.Ipaddress,
	); err != nil {
		return err
	}
	stat.ID = int(conn.LastInsertRowID())
	return nil
}

// FindOnDate возвращает map ключом которой является Mac-адрес устройства
func (r *StatRepository) FindSizeBetweenDate(from, to string) (map[string]models.StatDeviceType, error) {
	conn, err := sqlite.Open(r.store.dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	devStats := map[string]models.StatDeviceType{}

	stmt, err := conn.Prepare(`SELECT ipaddress, login, sum(size) as summ, hour 
	FROM statistics
	WHERE date(date_str) BETWEEN date($1) AND date($2)
	GROUP BY login, ipaddress, hour
	ORDER BY sum(size) DESC;`, from, to)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for {
		// Use Scan to access column data from a row
		var hour int
		var size int
		var ip, mac string
		var pHour [24]uint64

		hasRow, err := stmt.Step()
		if err != nil {
			return nil, err
		}

		if !hasRow {
			break
		}

		err = stmt.Scan(&ip, &mac, &size, &hour)
		if err != nil {
			return nil, err
		}
		stat, ok := devStats[mac]
		if !ok {
			pHour[hour] = uint64(size)
			stat = models.StatDeviceType{
				Ip:  ip,
				Mac: mac,
				StatType: models.StatType{
					PerHour: pHour,
					Size:    uint64(size),
				},
			}
		} else {
			stat.PerHour[hour] += uint64(size)
			stat.Size += uint64(size)
		}

		devStats[mac] = stat
	}
	return devStats, nil
}

func (r *StatRepository) DeletingDateData(date string) error {
	conn, err := sqlite.Open(r.store.dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Exec("delete from statistics where date_str = $1", date)
	if err != nil {
		return err
	}
	n := conn.Changes()
	if n == 0 {
		return fmt.Errorf("nothing delete(%v)", err)
	}
	return nil
}

func (r *StatRepository) Save(st []*models.StatInDB, bufLine uint16) error {
	conn, err := sqlite.Open(r.store.dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	insertSQL := `INSERT INTO statistics (date_str, year, month, day, hour, size, login, ipaddress) VALUES `
	query := insertSQL
	values := []interface{}{}
	k := 0
	numFields := 8 // the number of fields you are inserting
	bufItem := int(bufLine) * numFields
	for _, stat := range st {
		n := k * numFields
		if int(bufItem)-n < numFields && len(values) > 0 {
			query = query[:len(query)-1] // remove the trailing comma
			query += ";"
			err = conn.Exec(query, values...)
			if err != nil {
				return err
			}

			query = insertSQL
			k = 0
			values = []interface{}{}
		}
		values = append(values,
			stat.Date,
			stat.Year,
			stat.Month,
			stat.Day,
			stat.Hour,
			int(stat.Size),
			stat.Login,
			stat.Ipaddress)

		query += "("
		for j := 0; j < numFields; j++ {
			query += "?" + ","
		}
		query = query[:len(query)-1] + "),"
		k++
	}
	query = query[:len(query)-1] // remove the trailing comma
	query += ";"
	err = conn.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}
