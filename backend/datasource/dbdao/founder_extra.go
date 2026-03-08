package dbdao

import (
	"github.com/pkg/errors"
)

// FounderExtra 创业者额外信息结构
type FounderExtra struct {
	BaseModel
	UserID                uint64   `gorm:"column:user_id" json:"user_id"`
	CompanySize           string   `gorm:"column:company_size" json:"company_size"`                       // 公司规模
	FundingStage          string   `gorm:"column:funding_stage" json:"funding_stage"`                     // 融资阶段
	Experience            string   `gorm:"column:experience" json:"experience"`                           // 创业经验
	Skills                []string `gorm:"column:skills" json:"skills"`                                   // 技能
	InvestmentStage       string   `gorm:"column:investment_stage" json:"investment_stage"`               // 寻求投资阶段
	InvestmentAmountLower uint64   `gorm:"column:investment_amount_lower" json:"investment_amount_lower"` // 寻求投资金额下限
	InvestmentAmountUpper uint64   `gorm:"column:investment_amount_upper" json:"investment_amount_upper"` // 寻求投资金额上限
	Status                int8     `gorm:"column:status" json:"status"`                                   // 状态
}

func (d *DB) InsertFounder(founder *FounderExtra) error {
	if founder == nil {
		return errors.New("founder extra is empty")
	}
	result := d.DB().Create(founder)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *DB) GetFounders(ids ...uint) ([]*FounderExtra, error) {
	founders := []*FounderExtra{}
	if len(ids) == 0 {
		return founders, nil
	}

	result := d.DB().Where("id IN ?", ids).Find(&founders)
	return founders, result.Error
}
