package mygrep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/logger"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/model"
)

//Run - функция use-case, которая компилирует регулярное выражение, читает построчно входной поток, применяет правила совпадения и возвращает результат

func Run(cfg model.GrepConfig, r io.Reader, log logger.Logger) (model.GrepResult, error) {
	if log == nil {
		log = logger.NopLogger{}
	}

	log.Info("[mygrep] start processing")

	var result model.GrepResult

	pattern := cfg.Pattern
	if cfg.IgnoreCase {
		pattern = "(?i)" + pattern
		log.Debug("[mygrep] ignore case enabled")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Error("[mygrep] regex compile error")
		return result, err
	}

	scanner := bufio.NewScanner(r)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()

		match := re.MatchString(line)
		if cfg.Invert {
			match = !match
		}

		if match {
			if cfg.ShowNumber {
				line = fmt.Sprintf("%d:%s", lineNum, line)
			}
			result.Lines = append(result.Lines, line)
		}

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		log.Error("[mygrep] scanner error")
		return result, err
	}

	log.Info("[mygrep] processing finished")
	log.Debug("[mygrep] found lines count = " + strconv.Itoa(len(result.Lines)))

	return result, nil
}
