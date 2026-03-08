package scylladao

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// UserPublishNewsRow 分区键: (status, user_id)
type UserPublishNewsRow struct {
	BaseModel
	UserID uuid.UUID `json:"user_id"`
	NewsID uuid.UUID `json:"news_id"`
	Status int8      `json:"status"`
}

// 保存用户发布帖子记录
func (d *ScyllaDB) SaveUserPublishNews(userID, newsID uuid.UUID) error {
	r := UserPublishNewsRow{
		UserID: userID,
		NewsID: newsID,
	}
	err := r.Reset()
	if err != nil {
		return err
	}
	err = d.DB().Query(
		`INSERT INTO user_publish_news (id, updated_at, user_id, news_id) VALUES (?, ?, ?, ?)`,
		r.ID, r.UpdatedAt, r.UserID, r.NewsID,
	).Exec()
	if err != nil {
		return errors.WithMessage(err, "failed to insert user publish news")
	}
	return nil
}

// 删除用户发布帖子记录
func (d *ScyllaDB) DeleteUserPublishNews(userID, newsID uuid.UUID, status int8) error {
	err := d.DB().Query(
		`DELETE FROM user_publish_news WHERE status = ? AND user_id = ? AND news_id = ?`,
		status, userID, newsID,
	).Exec()
	if err != nil {
		return errors.WithMessage(err, "failed to delete user publish news")
	}
	return nil
}
