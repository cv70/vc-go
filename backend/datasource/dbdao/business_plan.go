package dbdao

import (
	"github.com/pkg/errors"
)

// BusinessPlan 业务计划书结构
type BusinessPlan struct {
	BaseModel
	Title           string    `gorm:"column:title" json:"title"`
	Content         string    `gorm:"column:content" json:"content"`
	Industry        string    `gorm:"column:industry" json:"industry"`
	Region          string    `gorm:"column:region" json:"region"`
	FinancingAmount float64   `gorm:"column:financing_amount" json:"financing_amount"` // 融资金额
	CompanySize     string    `gorm:"column:company_size" json:"company_size"`
	Embedding       []float32 `gorm:"column:embedding" json:"embedding"` // 向量表示
}

// InsertBusinessPlan 插入单条政策数据
func (d *DB) InsertBusinessPlan(plan *BusinessPlan) error {
	if plan == nil {
		return errors.New("business plan is empty")
	}
	result := d.DB().Create(plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *DB) GetBusinessPlan(id uint) (*BusinessPlan, error) {
	plan := BusinessPlan{}
	result := d.DB().Where("id = ?", plan.ID).Find(&plan)
	return &plan, result.Error
}

func (d *DB) GetBusinessPlans(ids ...uint) ([]*BusinessPlan, error) {
	plans := []*BusinessPlan{}
	if len(ids) == 0 {
		return plans, nil
	}

	result := d.DB().Where("id IN ?", ids).Find(&plans)
	return plans, result.Error
}
