package scylladao

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserLikeNewsRow struct {
	BaseModel
	UserID uuid.UUID
	NewsID uuid.UUID
}

// 保存用户点赞帖子记录
func (d *ScyllaDB) SaveUserLikeNews(userID, newsID uuid.UUID) error {
	r := UserLikeNewsRow{
		UserID: userID,
		NewsID: newsID,
	}
	err := r.Reset()
	if err != nil {
		return err
	}
	err = d.DB().Query(
		`INSERT INTO user_like_news (id, updated_at, user_id, news_id) VALUES (?, ?, ?, ?)`,
		r.ID, r.UpdatedAt, r.UserID, r.NewsID,
	).Exec()
	if err != nil {
		return errors.WithMessage(err, "failed to insert news like")
	}
	return nil
}

// 删除用户点赞帖子记录
func (d *ScyllaDB) DeleteUserLikeNews(userID, newsID uuid.UUID) error {
	err := d.DB().Query(
		`DELETE FROM user_like_news WHERE user_id = ? AND news_id = ?`,
		userID, newsID,
	).Exec()
	if err != nil {
		return errors.WithMessage(err, "failed to insert news like")
	}
	return nil
}
