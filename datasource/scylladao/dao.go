package scylladao

import (
	"time"
	"vc-go/pkg/gid"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type ScyllaDB gocql.Session

func NewScyllaDB(sess *gocql.Session) *ScyllaDB {
	return (*ScyllaDB)(sess)
}

func (d *ScyllaDB) DB() *gocql.Session {
	return d.DB()
}

type BaseModel struct {
	ID        uuid.UUID
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
