package retriever

import (
	"SuperBizAgent/internal/ai/embedder"
	"SuperBizAgent/utility/client"
	"SuperBizAgent/utility/common"
	"context"

	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func NewMilvusRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		return nil, err
	}
	eb, err := embedder.DoubaoEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	// 配置 AUTOINDEX 搜索参数
	sp, err := entity.NewIndexAUTOINDEXSearchParam(1)
	if err != nil {
		return nil, err
	}
	r, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      cli,
		Collection:  common.MilvusCollectionName,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:       1,
		Embedding:  eb,
		MetricType: entity.COSINE,
		Sp:         sp,
		// 将 float64 向量转换为 float32
		VectorConverter: func(ctx context.Context, vectors [][]float64) ([]entity.Vector, error) {
			result := make([]entity.Vector, len(vectors))
			for i, vec := range vectors {
				vec32 := make([]float32, len(vec))
				for j, v := range vec {
					vec32[j] = float32(v)
				}
				result[i] = entity.FloatVector(vec32)
			}
			return result, nil
		},
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
