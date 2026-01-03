package dbdao

import (
	"time"
	"vc-go/pkg/gid"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DB gorm.DB

func NewDB(db *gorm.DB) *DB {
	return (*DB)(db)
}

func (d *DB) DB() *gorm.DB {
	return (*gorm.DB)(d)
}

type BaseModel struct {
	ID        uuid.UUID `gorm:"primarykey"`
	UpdatedAt time.Time
}

func (m *BaseModel) Reset() (err error) {
	m.ID, err = gid.NewUUID()
	if err != nil {
		return err
	}
	sec, nsec := m.ID.Time().UnixTime()
	m.UpdatedAt = time.Unix(sec, nsec)
	return nil
}
