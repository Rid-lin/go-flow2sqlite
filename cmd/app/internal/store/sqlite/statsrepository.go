package sqlitestore

import (
	"go-flow2sqlite/cmd/app/internal/models"
	"time"

	_ "modernc.org/sqlite"
)

type StatRepository struct {
	store *Store
}

func (s *StatRepository) Save(inputChannel chan models.BMes) error {
	for bm := range inputChannel {
		err := s.saveOneLine(&bm)
		if err != nil {
			for r := 3; r <= 0; r-- {
				err := s.saveOneLine(&bm)
				if err == nil {
					break
				}
				time.Sleep(300 * time.Microsecond)
			}
		}

	}
	return nil
}

func (r *StatRepository) saveOneLine(bm *models.BMes) error {
	// 	conn, err := sql.Open("sqlite", r.store.dsn)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer conn.Close()

	// 	insertSQL := `INSERT INTO statistics (date_str, year, month, day, hour, size, login, ipaddress) VALUES `
	// 	query := insertSQL
	// 	values := []interface{}{}
	// 	k := 0
	// 	numFields := 8 // the number of fields you are inserting
	// 	bufItem := int(bufLine) * numFields
	// 	for _, stat := range st {
	// 		n := k * numFields
	// 		if int(bufItem)-n < numFields && len(values) > 0 {
	// 			query = query[:len(query)-1] // remove the trailing comma
	// 			query += ";"
	// 			_, err = conn.Exec(query, values...)
	// 			if err != nil {
	// 				return err
	// 			}

	// 			query = insertSQL
	// 			k = 0
	// 			values = []interface{}{}
	// 		}
	// 		values = append(values,
	// 			stat.Date,
	// 			stat.Year,
	// 			stat.Month,
	// 			stat.Day,
	// 			stat.Hour,
	// 			int(stat.Size),
	// 			stat.Login,
	// 			stat.Ipaddress)

	// 		query += "("
	// 		for j := 0; j < numFields; j++ {
	// 			query += "?" + ","
	// 		}
	// 		query = query[:len(query)-1] + "),"
	// 		k++
	// 	}
	// 	query = query[:len(query)-1] // remove the trailing comma
	// 	query += ";"
	// 	_, err = conn.Exec(query, values...)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return nil

	// ctx := context.Background()
	// s.DB.ExecContext(ctx, "INSERT INTO raw_log VALUES(?,?,?,?,?,?,?,?,?,?)",
	// 	bm.Time,
	// 	bm.Delay,
	// 	bm.IPSrc,
	// 	bm.Protocol,
	// 	bm.SizeOfByte,
	// 	bm.IPDst,
	// 	bm.PortDst,
	// 	bm.MacDst,
	// 	bm.RouterIP,
	// 	bm.PortSRC,
	// )s
	return nil
}
