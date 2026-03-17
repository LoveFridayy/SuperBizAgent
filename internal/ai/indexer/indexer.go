package indexer

import (
	embedder2 "SuperBizAgent/internal/ai/embedder"
	"SuperBizAgent/utility/client"
	"SuperBizAgent/utility/common"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func NewMilvusIndexer(ctx context.Context) (*milvus.Indexer, error) {
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		return nil, err
	}
	eb, err := embedder2.DoubaoEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config := &milvus.IndexerConfig{
		Client:            cli,
		Collection:        common.MilvusCollectionName,
		Fields:            fields,
		Embedding:         eb,
		DocumentConverter: documentConverter,
	}
	indexer, err := milvus.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return indexer, nil
}

// documentConverter 将 Document 和 float64 向量转换为 Milvus 行数据
func documentConverter(ctx context.Context, docs []*schema.Document, vectors [][]float64) ([]interface{}, error) {
	rows := make([]interface{}, 0, len(docs))
	for i, doc := range docs {
		// 将 float64 向量转换为 float32
		vec32 := make([]float32, len(vectors[i]))
		for j, v := range vectors[i] {
			vec32[j] = float32(v)
		}
		// 序列化 metadata 为 []byte
		metaBytes, err := json.Marshal(doc.MetaData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		row := map[string]interface{}{
			"id":       doc.ID,
			"vector":   vec32,
			"content":  doc.Content,
			"metadata": metaBytes, // 使用 []byte 而不是 string
		}
		rows = append(rows, row)
	}
	return rows, nil
}

var fields = []*entity.Field{
	{
		Name:     "id",
		DataType: entity.FieldTypeVarChar,
		TypeParams: map[string]string{
			"max_length": "256",
		},
		PrimaryKey: true,
	},
	{
		Name:     "vector",
		DataType: entity.FieldTypeFloatVector, // 改为浮点向量
		TypeParams: map[string]string{
			"dim": "1024", // 阿里云 text-embedding-v3 的维度
		},
	},
	{
		Name:     "content",
		DataType: entity.FieldTypeVarChar,
		TypeParams: map[string]string{
			"max_length": "8192",
		},
	},
	{
		Name:     "metadata",
		DataType: entity.FieldTypeJSON,
	},
}
