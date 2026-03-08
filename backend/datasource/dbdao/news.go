package dbdao

import (
	"strings"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// News 帖子结构
type News struct {
	BaseModel
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Category    string   `json:"category"` // 分类：政策解读、成功案例、融资动态、活动日历
	Source      string   `json:"source"`   // 来源
	PublishDate string   `json:"publish_date"`
	Tags        []string `json:"tags"`     // 标签
	Region      string   `json:"region"`   // 地区
	Industry    string   `json:"industry"` // 行业
}

func (d *DB) ListNews(ids ...uint) ([]*News, error) {
	newsList := []*News{}
	err := d.DB().Where("id IN ?", ids).Find(&newsList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return newsList, nil
}

func (d *DB) DetailNews(id uint) (*News, error) {
	news := News{}
	err := d.DB().Where("id = ?", id).First(&news).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &news, nil
}

func (d *DB) InsertNews(news *News) error {
	if news == nil {
		return errors.New("insert news is nil")
	}
	return d.DB().Save(news).Error
}

// GetNews 获取帖子列表
func (d *DB) GetNews(category, region, industry string, tags []string, page, limit int) ([]*News, int64, error) {
	newsList := []*News{}
	var total int64

	db := d.DB().Model(&News{})

	query := strings.Builder{}
	args := []interface{}{}
	if category != "" {
		if query.Len() > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString("category = ?")
		args = append(args, category)
	}

	if region != "" {
		if query.Len() > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString("region = ?")
		args = append(args, region)
	}

	if industry != "" {
		if query.Len() > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString("industry = ?")
		args = append(args, industry)
	}

	// 标签筛选（简化实现，实际项目中可能需要更复杂的标签匹配）
	if len(tags) > 0 {
		for _, tag := range tags {
			if query.Len() > 0 {
				query.WriteString(" AND ")
			}
			query.WriteString("JSON_CONTAINS(tags, ?)")
			args = append(args, tag)
		}
		// 这里简化处理，实际需要根据数据库类型调整查询语句
	}

	err := db.Order("publish_date desc").Offset((page - 1) * limit).Limit(limit).Find(&newsList).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return newsList, total, nil
}

// SearchNews 搜索帖子
func (d *DB) SearchNews(query string, page, limit int) ([]*News, int, error) {
	newsList := []*News{}
	var total int64

	// 使用全文搜索（这里简化实现，实际项目中应使用ES或其他全文搜索引擎）
	err := d.DB().Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%").Find(&newsList).Error
	if err != nil {
		return nil, 0, err
	}

	err = d.DB().Model(&News{}).Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return newsList, int(total), nil
}
