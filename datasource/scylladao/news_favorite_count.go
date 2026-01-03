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
