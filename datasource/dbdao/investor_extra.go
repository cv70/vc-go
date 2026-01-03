package dbdao

import "github.com/pkg/errors"

// InvestorExtra 投资人额外信息结构
type InvestorExtra struct {
	BaseModel
	UserID                uint64 `gorm:"column:user_id" json:"user_id"`
	Company               string `gorm:"column:company" json:"company"`
	InvestmentStage       string `gorm:"column:investment_stage" json:"investment_stage"`               // 投资阶段
	InvestmentAmountLower uint64 `gorm:"column:investment_amount_lower" json:"investment_amount_lower"` // 投资金额下限
	InvestmentAmountUpper uint64 `gorm:"column:investment_amount_upper" json:"investment_amount_upper"` // 投资金额上限
	Status                int8   `gorm:"column:status" json:"status"`                                   // 状态
}

func (d *DB) InsertInvestor(investor *InvestorExtra) error {
	if investor == nil {
		return errors.New("investor extra is empty")
	}
	err := investor.Reset()
	if err != nil {
		return err
	}
	result := d.DB().Create(investor)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *DB) GetInvestorExtras(ids ...uint) ([]*InvestorExtra, error) {
	investors := []*InvestorExtra{}
	if len(ids) == 0 {
		return investors, nil
	}

	result := d.DB().Where("id IN ?", ids).Find(&investors)
	return investors, result.Error
}
