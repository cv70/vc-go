package dbdao

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// User 用户实体
type User struct {
	BaseModel
	Phone        string   `db:"phone" json:"phone"`
	Username     string   `db:"username" json:"username"`
	Password     string   `db:"password" json:"-"` // 不返回密码
	DeviceID     string   `db:"device_id" json:"device_id"`
	Role         int8     `db:"role" json:"role"`                          // 角色 0: 普通用户 1: 投资人/投资机构 2: 创业者
	Industry     []string `gorm:"column:industry" json:"industry"`         // 关注行业
	Introduction string   `gorm:"column:introduction" json:"introduction"` // 自我介绍
	IsVerified   bool     `db:"is_verified" json:"is_verified"`
}

func (d *DB) GetUserByPhone(phone string) (*User, error) {
	var user User
	result := d.DB().Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (d *DB) GetUsers(ids ...uuid.UUID) (*User, error) {
	var user User
	result := d.DB().Where("id IN ?", ids).First(&user)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// ExistUsername 检查用户名是否存在
func (d *DB) ExistUsername(username string) (bool, error) {
	var exists bool
	result := d.DB().Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if result.Error != nil {
		return false, result.Error
	}
	return exists, nil
}

func (d *DB) CreateUser(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	err := user.Reset()
	if err != nil {
		return err
	}
	result := d.DB().Create(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("failed to insert user")
	}
	return nil
}
