package teststore

import "go-flow2sqlite/cmd/app/internal/models"

type StatRepository struct {
	store *Store
	stats []*models.BMes
}

func (r *StatRepository) Save(inputChannel chan models.BMes) error {
	for bm := range inputChannel {

		r.stats = append(r.stats, &bm)
	}
	return nil
}
func (r *StatRepository) saveOneLine(bm *models.BMes) error {
	r.stats = append(r.stats, bm)

	return nil
}
