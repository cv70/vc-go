package scylladao

import (
	"github.com/google/uuid"
)

// NewsRow 分区键: (status, id)
type NewsRow struct {
	BaseModel
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"` // 分类：政策解读、成功案例、融资动态、活动日历
	Tags     []string `json:"tags"`     // 标签
	Region   string   `json:"region"`   // 地区
	Industry []string `json:"industry"` // 行业
	Status   int32    `json:""`         // 草稿,审核中,审核失败,公开,私密
}

func (d *ScyllaDB) SaveNews(r *NewsRow) error {
	err := r.Reset()
	if err != nil {
		return err
	}
	err = d.DB().Query(
		`INSERT INTO news (id, updated_at, title, content, category, tags, region, industry) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.ID, r.UpdatedAt, r.Title, r.Content, r.Category, r.Tags, r.Region, r.Industry,
	).Exec()
	return err
}

func (d *ScyllaDB) GetNews(newsID uuid.UUID) (*NewsRow, error) {
	r := NewsRow{}
	err := d.DB().Query(
		`SELECT updated_at, title, content, category, tags, region, industry FROM news WHERE id = ?`,
		newsID,
	).Scan(&r.UpdatedAt, &r.Title, &r.Content, &r.Category, &r.Tags, &r.Region, &r.Industry)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (d *ScyllaDB) GetNewsList(newsIDs ...uuid.UUID) ([]*NewsRow, error) {
	var newsList []*NewsRow

	// 构建查询条件
	if len(newsIDs) == 0 {
		// 如果没有提供ID，则返回空列表
		return newsList, nil
	}

	// 批量查询
	iter := d.DB().Query(
		`SELECT id, updated_at, title, content, category, tags, region, industry FROM news WHERE id IN ?`,
		newsIDs,
	).Iter()

	for {
		r := NewsRow{}
		ok := iter.Scan(&r.ID, &r.UpdatedAt, &r.Title, &r.Content, &r.Category, &r.Tags, &r.Region, &r.Industry)
		if !ok {
			break
		}
		newsList = append(newsList, &r)
	}

	return newsList, iter.Close()
}
