package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/logger"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/model"
)

func SendShard(addr string, shard model.Shard, cfg model.GrepConfig, log logger.Logger) (model.ShardResult, error) {
	if log == nil {
		log = logger.NopLogger{}
	}

	reqBody := ProcessRequest{
		Shard:  shard,
		Config: cfg,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		log.Error("[cluster client] failed to marshal request")
		return model.ShardResult{}, err
	}

	url := fmt.Sprintf("http://%s/process", addr)
	log.Info("[cluster client] sending shard to " + url)

	httpClient := &http.Client{Timeout: 5 * time.Second}

	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Error("[cluster client] request failed")
		return model.ShardResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("c[cluster client] bad response status")
		return model.ShardResult{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	var result model.ShardResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error("[cluster client] failed to decode response")
		return model.ShardResult{}, err
	}

	log.Info("[cluster client] shard processed successfully")
	return result, nil
}
