package cluster

import (
	"testing"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/logger"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/model"
)

func TestProcessShardsQuorum(t *testing.T) {
	shards := []model.Shard{
		{ID: 1, Data: []byte("hello world\nfoo\n")},
		{ID: 2, Data: []byte("HELLO gopher\nbar\n")},
		{ID: 3, Data: []byte("nothing here\n")},
	}

	cfg := model.GrepConfig{
		Pattern:    "hello",
		IgnoreCase: true,
	}

	c := Cluster{
		Logger: logger.NopLogger{},
		Quorum: 2, // n/2+1
	}

	results := c.ProcessShards(shards, cfg)

	if len(results) != 2 {
		t.Fatalf("expected 2 shard results for quorum, got %d", len(results))
	}

	// Проверяем, что в результатах действительно есть совпадения
	found := 0
	for _, r := range results {
		found += len(r.Lines)
	}

	if found == 0 {
		t.Fatalf("expected to find at least one matching line, got 0")
	}
}
