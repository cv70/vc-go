package news

import (
	"vc-go/datasource/dbdao"
	"vc-go/datasource/scylladao"
	"vc-go/datasource/vectordao"
)

type NewsDomain struct {
	DB       *dbdao.DB
	Scylla   *scylladao.ScyllaDB
	VectorDB *vectordao.VectorDB
}
