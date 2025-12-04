package cluster

import (
	"bytes"
	"sync"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/logger"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/model"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/mygrep"
)

type Cluster struct {
	Logger    logger.Logger
	Quorum    int
	ShardChan chan model.ShardResult
}

// ProcessShards запускает параллельную обработку шардов через mygrep
func (c *Cluster) ProcessShards(shards []model.Shard, cfg model.GrepConfig) []model.ShardResult {
	if c.Logger == nil {
		c.Logger = logger.NopLogger{}
	}

	c.ShardChan = make(chan model.ShardResult, len(shards))
	var wg sync.WaitGroup

	for _, shard := range shards {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := mygrep.Run(cfg, bytes.NewReader(shard.Data), c.Logger)
			if err != nil {
				c.Logger.Error("[cluster] failed to process shard")
				return
			}
			c.ShardChan <- model.ShardResult{ID: shard.ID, Lines: res.Lines}
		}()
	}

	go func() {
		wg.Wait()
		close(c.ShardChan)
	}()

	// Сбор результатов до достижения кворума
	results := make([]model.ShardResult, 0)
	for r := range c.ShardChan {
		results = append(results, r)
		if len(results) >= c.Quorum {
			c.Logger.Info("cluster: quorum reached")
			break
		}
	}

	return results
}
