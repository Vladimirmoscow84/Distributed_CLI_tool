package app

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/cluster"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/logger"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/model"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/mygrep"
)

// Params -  входыне параметры
type Params struct {
	Pattern    string
	IgnoreCase bool
	ShowNumber bool
	Invert     bool
	FilePath   string //если пусто, то читается из  stdin
	UseCluster bool
	Quorum     int
}

func Run(p Params, log logger.Logger) error {
	if log == nil {
		log = logger.NopLogger{}
	}
	log.Info("[app] start")
	reader, closer, err := openinput(p.FilePath)
	if err != nil {
		log.Error("[app] failed to open input")
		return err
	}
	if closer != nil {
		defer closer.Close()
	}
	cfg := model.GrepConfig{
		Pattern:    p.Pattern,
		IgnoreCase: p.IgnoreCase,
		ShowNumber: p.ShowNumber,
		Invert:     p.Invert,
	}

	if !p.UseCluster {
		log.Info("[app] running mygrep")
		result, err := mygrep.Run(cfg, reader, log)
		if err != nil {
			log.Error("[app] failed mygrep")
			return err
		}

		log.Info("[app] writening output")

		err = writeOut(os.Stdout, result)
		if err != nil {
			log.Error("[app] failed to write output")
		}
		log.Info("[app] finished successfully")
		return nil
	}

	log.Info("[app] running in cluster mode")

	scanner := bufio.NewScanner(reader)
	shards := make([]model.Shard, 0)
	id := 1

	for scanner.Scan() {
		shards = append(shards, model.Shard{
			ID:   id,
			Data: append(scanner.Bytes(), '\n'),
		})
		id++
	}

	c := cluster.Cluster{
		Logger: log,
		Quorum: p.Quorum,
	}

	results := c.ProcessShards(shards, cfg)

	final := model.GrepResult{}
	for _, r := range results {
		final.Lines = append(final.Lines, r.Lines...)
	}

	log.Info("[app] writing cluster output")
	err = writeOut(os.Stdout, final)
	if err != nil {
		log.Error("[app] failed to write cluster output")
		return err
	}

	log.Info("[app] finished successfully (cluster)")
	return nil
}

// openInput открывает файл, если нет файла, то возвращает stdin
func openinput(path string) (t io.Reader, closer io.Closer, err error) {
	if path == "" {
		return os.Stdin, nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}

// writeOut печатает результат в writerer
func writeOut(w io.Writer, res model.GrepResult) error {
	for _, line := range res.Lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return err
		}
	}
	return nil
}
