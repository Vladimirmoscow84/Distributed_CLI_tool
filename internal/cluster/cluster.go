package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/logger"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/model"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/mygrep"
	"github.com/gin-gonic/gin"
)

type Cluster struct {
	Logger logger.Logger
	Quorum int
	Peers  []string // адреса других узлов
}

// ProcessShards распределяет шарды по локалке и сети, ждёт кворум
func (c *Cluster) ProcessShards(shards []model.Shard, cfg model.GrepConfig) []model.ShardResult {
	results := make([]model.ShardResult, 0, len(shards))

	for _, shard := range shards {
		// локальная обработка
		res, _ := mygrep.Run(cfg, bytes.NewReader(shard.Data), c.Logger)
		results = append(results, model.ShardResult{ID: shard.ID, Lines: res.Lines})

		// отправка на peers
		for _, peer := range c.Peers {
			c.Logger.Info(fmt.Sprintf("[cluster client] sending shard to http://%s/process", peer))
			if err := c.sendShard(peer, shard); err != nil {
				c.Logger.Error(fmt.Sprintf("[cluster client] failed to send shard to %s: %v", peer, err))
			}
		}
	}

	return results
}

func (c *Cluster) sendShard(peer string, shard model.Shard) error {
	data, _ := json.Marshal(shard)
	resp, err := http.Post(fmt.Sprintf("http://%s/process", peer), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// StartServer запускает Gin HTTP сервер для приема shards
func (c *Cluster) StartServer(port int) {
	r := gin.Default()

	r.POST("/process", func(ctx *gin.Context) {
		var shard model.Shard
		if err := ctx.BindJSON(&shard); err != nil {
			c.Logger.Error("failed to decode shard")
			ctx.Status(400)
			return
		}

		c.Logger.Info(fmt.Sprintf("[cluster server] received shard %d", shard.ID))
		cfg := model.GrepConfig{Pattern: "Spartak"} // можно прокинуть через тело запроса
		res, _ := mygrep.Run(cfg, bytes.NewReader(shard.Data), c.Logger)

		ctx.JSON(200, model.ShardResult{ID: shard.ID, Lines: res.Lines})
	})

	c.Logger.Info(fmt.Sprintf("[cluster server] listening on :%d", port))
	r.Run(fmt.Sprintf(":%d", port))
}
