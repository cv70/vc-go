package vectordao

import (
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"k8s.io/klog/v2"
)

func (d *VectorDB) SearchSimilarInvestors(ctx context.Context, vector []float32, topk int) (map[string]float32, error) {
	sp, err := entity.NewIndexIvfFlatSearchParam(16)
	if err != nil {
		return nil, err
	}

	ids := map[string]float32{}
	searchResult, err := d.Client.Search(ctx, "investor", nil, "", []string{"vector"}, []entity.Vector{entity.FloatVector(vector)}, "text", entity.COSINE, topk, sp)
	for _, result := range searchResult {
		idColumn, ok := result.Fields.GetColumn("id").(*entity.ColumnString)
		if !ok {
			continue
		}
		for i := 0; i < result.ResultCount; i++ {
			score := result.Scores[i]
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				klog.Errorf("id column value by idx error: %v", err)
				continue
			}
			ids[id] = score
		}
	}
	return ids, nil
}
