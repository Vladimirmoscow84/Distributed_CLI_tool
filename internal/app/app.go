package app

import (
	"errors"
	"fmt"
	"io"
	"os"

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
}

func Run(p Params) error {
	reader, closer, err := openinput(p.FilePath)
	if err != nil {
		return nil
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

	result, err := mygrep.Run(cfg, reader)
	if err != nil {
		return err
	}

	return writeOut(os.Stdout, result)
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

// вывод в stdout
func writeOut(w io.Writer, res model.GrepResult) error {
	if w == nil {
		return errors.New("nil writer")
	}
	for _, line := range res.Lines {
		fmt.Fprintln(w, line)
	}
	return nil
}
