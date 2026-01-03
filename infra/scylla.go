package infra

import (
	"context"
	"vc-go/config"
	"vc-go/datasource/scylladao"

	"github.com/gocql/gocql"
)

func NewScylla(ctx context.Context, c *config.ScyllaConfig) (*scylladao.ScyllaDB, error) {
	cluster := gocql.NewCluster(c.Host)
	cluster.Keyspace = c.DBName
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: c.User,
		Password: c.Password,
	}
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return scylladao.NewScyllaDB(session), nil
}
