package scylladao

import "github.com/google/uuid"

func (d *ScyllaDB) IncrNewsVisitCount(newsID uuid.UUID) error {
	err := d.DB().Query(
		`UPDATE news_visit_count SET count = count + 1 WHERE news_id = ?`,
		newsID,
	).Exec()
	return err
}

func (d *ScyllaDB) GetNewsVisitCount(newsID uuid.UUID) (int64, error) {
	var count int64
	err := d.DB().Query(
		`SELECT count FROM news_visit_count WHERE news_id = ?`,
		newsID,
	).Scan(&count)
	return count, err
}
