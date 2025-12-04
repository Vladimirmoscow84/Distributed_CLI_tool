package cluster

import (
	"bytes"
	"net/http"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/logger"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/model"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/mygrep"
	"github.com/gin-gonic/gin"
)

type ProcessRequest struct {
	Shard  model.Shard      `json:"shard"`
	Config model.GrepConfig `json:"config"`
}

func StartNode(addr string, log logger.Logger) error {
	if log == nil {
		log = logger.NopLogger{}
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(func(c *gin.Context) {
		log.Info("[cluster] " + c.Request.Method + " " + c.Request.URL.Path)
		c.Next()
	})

	r.POST("/process", func(c *gin.Context) {
		var req ProcessRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			log.Error("[cluster] failed to decode request")
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		log.Info("[cluster] received shard")

		res, err := mygrep.Run(req.Config, bytes.NewReader(req.Shard.Data), log)
		if err != nil {
			log.Error("[cluster] failed to process shard")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "processing failed"})
			return
		}

		result := model.ShardResult{
			ID:    req.Shard.ID,
			Lines: res.Lines,
		}

		c.JSON(http.StatusOK, result)
	})

	log.Info("[cluster] starting node at " + addr)
	return r.Run(addr)
}
