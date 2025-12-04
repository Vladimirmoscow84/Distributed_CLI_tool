package app

import (
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
	FilePath   string
	Cluster    bool
	Quorum     int
	Peers      []string
	Port       int
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

	if p.Cluster {
		log.Info("[app] running in cluster mode")

		c := cluster.Cluster{
			Logger: log,
			Quorum: p.Quorum,
			Peers:  p.Peers,
		}

		// запуск HTTP-сервер для приема shard-ов
		go c.StartServer(p.Port)

		data, err := io.ReadAll(reader)
		if err != nil {
			log.Error("[app] failed to read input")
			return err
		}

		shards := []model.Shard{{ID: 0, Data: data}}

		results := c.ProcessShards(shards, cfg)

		return writeShardsOut(os.Stdout, results)
	} else {
		log.Info("[app] running in local mode")
		result, err := mygrep.Run(cfg, reader, log)
		if err != nil {
			log.Error("[app] failed mygrep")
			return err
		}
		return writeOut(os.Stdout, result)
	}
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
func writeShardsOut(w io.Writer, res []model.ShardResult) error {
	for _, r := range res {
		for _, line := range r.Lines {
			_, err := fmt.Fprintln(w, line)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
