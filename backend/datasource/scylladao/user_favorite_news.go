package scylladao

import (
	"github.com/google/uuid"
)

// UserFavoriteNewsRow 分区键: (user_id, news_id)
type UserFavoriteNewsRow struct {
	BaseModel
	UserID uint
	NewsID uuid.UUID
}

func (d *ScyllaDB) SaveUserFavoriteNews(userID, newsID int64) error {
	err := d.DB().Query(
		`INSERT INTO user_favorite_news (user_id, news_id, created_at) VALUES (?, ?, now())`,
		userID, newsID,
	).Exec()
	return err
}

func (d *ScyllaDB) DeleteUserFavoriteNews(userID, newsID int64) error {
	err := d.DB().Query(
		`DELETE FROM user_favorite_news WHERE user_id = ? AND news_id = ?`,
		userID, newsID,
	).Exec()
	return err
}

func (d *ScyllaDB) GetUserFavoriteNewsIDs(userID int64) ([]string, error) {
	iter := d.DB().Query(`SELECT news_id FROM user_favorite_news WHERE user_id = ?`, userID).Iter()

	// 收集收藏的帖子ID
	var newsIDs []string
	var newsID string
	for iter.Scan(&newsID) {
		newsIDs = append(newsIDs, newsID)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return newsIDs, nil
}
