package dbdao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Policy struct {
	BaseModel
	Title       string `gorm:"column:title" json:"title"`
	Content     string `gorm:"column:content" json:"content"`
	Region      string `gorm:"column:region" json:"region"`
	Industry    string `gorm:"column:industry" json:"industry"`
	PublishDate string `gorm:"column:publish_date" json:"publish_date"`
	Source      string `gorm:"column:source" json:"source"`
}

// GetPolicies 获取政策列表
func (d *DB) GetPolicies(region, industry string, page, limit int, needTotal bool) ([]*Policy, int, error) {
	var policies []*Policy
	var total int64

	db := d.DB().Where("region = ? OR ? = ''", region, region)
	if industry != "" {
		db = db.Where("industry = ?", industry)
	}

	err := db.Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Find(&policies).Error
	if err != nil {
		return nil, 0, err
	}

	if needTotal {
		if err := db.Model(&Policy{}).Count(&total).Error; err != nil {
			return nil, 0, err
		}
	}

	return policies, int(total), nil
}

// GetPolicyByID 获取单条政策详情
func (d *DB) GetPolicyByID(id string) (*Policy, error) {
	var policy Policy
	if err := d.DB().First(&policy, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &policy, nil
}

// SearchPolicies 按关键词搜索政策
func (d *DB) SearchPolicies(query string) ([]Policy, error) {
	var policies []Policy
	// 使用模糊匹配作为占位
	if query == "" {
		return []Policy{}, nil
	}
	if err := d.DB().Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%").Find(&policies).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []Policy{}, nil
		}
		return nil, err
	}
	return policies, nil
}

// InsertPolicy 插入单条政策数据
func (d *DB) InsertPolicy(policy *Policy) error {
	if policy == nil {
		return errors.New("policy is empty")
	}
	result := d.DB().Create(policy)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
