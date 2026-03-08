package scylladao

func (d *ScyllaDB) IncrUserFavoriteNewsCount(newsID int64) error {
	err := d.DB().Query(
		`UPDATE news_favorite_count SET count = count + 1 WHERE news_id = ?`,
		newsID,
	).Exec()
	return err
}

func (d *ScyllaDB) DecrUserFavoriteNewsCount(newsID int64) error {
	err := d.DB().Query(
		`UPDATE news_favorite_count SET count = count - 1 WHERE news_id = ?`,
		newsID,
	).Exec()
	return err
}

func (d *ScyllaDB) GetUserFavoriteNewsCount(newsID int64) (int64, error) {
	var count int64
	err := d.DB().Query(
		`SELECT count FROM news_favorite_count WHERE news_id = ?`,
		newsID,
	).Scan(&count)
	return count, err
}
