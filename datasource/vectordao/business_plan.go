package vectordao

import (
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func (d *VectorDB) InsertBusinessPlanVector(ctx context.Context, ids []string, embeddings [][]float32) error {
	_, err := d.Client.Insert(
		ctx,
		"business_plan",
		"",
		entity.NewColumnInt64("id", ),
		entity.NewColumnFloatVector("vector", 768, embeddings),
	)
	return err
}
