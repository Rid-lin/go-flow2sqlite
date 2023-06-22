package teststore

import "go-flow2sqlite/cmd/app/internal/models"

type StatRepository struct {
	store *Store
	stats []*models.StatInDB
}

func (r *StatRepository) Create(stat *models.StatInDB) error {
	r.stats = append(r.stats, stat)
	stat.ID = len(r.stats)
	return nil
}

func (r *StatRepository) FindSizeBetweenDate(from, to string) (map[string]models.StatDeviceType, error) {
	devStats := map[string]models.StatDeviceType{}

	for _, v := range r.stats {

		var pHour [24]uint64
		hour := v.Hour
		size := v.Size
		ip := v.Ipaddress
		mac := v.Login
		stat, ok := devStats[mac]
		if !ok {
			pHour[v.Hour] = v.Size
			stat = models.StatDeviceType{
				Ip:  ip,
				Mac: mac,
				StatType: models.StatType{
					PerHour: pHour,
					Size:    size,
				},
			}
		} else {
			stat.PerHour[hour] += size
			stat.Size += size
		}

		devStats[mac] = stat
	}
	return devStats, nil
}

func (r *StatRepository) DeletingDateData(date string) error {
	lenRStat := len(r.stats)
	for i := 0; i < lenRStat; i++ {
		if r.stats[i].Date == date {
			r.stats = append(r.stats[:i], r.stats[i:]...)
		}
	}
	return nil
}

func (r *StatRepository) Save(st []*models.StatInDB, bufLine uint16) error {
	r.stats = append(r.stats, st...)

	return nil
}
