package mygrep

import (
	"bufio"
	"io"
	"regexp"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/model"
)

//Run - функция use-case, которая компилирует регулярное выражение, читает построчно входной поток, применяет правила совпадения и возвращает результат

func Run(cfg model.GrepConfig, r io.Reader) (model.GrepResult, error) {
	reg, err := compilePattern(cfg.Pattern, cfg.IgnoreCase)
	if err != nil {
		return model.GrepResult{}, err
	}
	scanner := bufio.NewScanner(r)

	lineNo := 0
	result := model.GrepResult{}

	for scanner.Scan() {
		lineNo++
		line := scanner.Text()

		matched := reg.MatchString(line)
		if cfg.Invert {
			matched = !matched
		}
		if matched {
			if cfg.ShowNumber {
				result.Lines = append(result.Lines, formatStringWithNumber(lineNo, line))
			} else {
				result.Lines = append(result.Lines, line)
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		return model.GrepResult{}, err
	}
	return result, nil

}

// compilePattern компилирует регулярку с учетом ignore-case
func compilePattern(pattern string, ignoreCase bool) (*regexp.Regexp, error) {
	if ignoreCase {
		return regexp.Compile("(?!)" + pattern)
	}
	return regexp.Compile(pattern)
}

// formatStringWithNumber форматирует строку с номером
func formatStringWithNumber(number int, line string) string {
	return numToString(number) + ":" + line
}

// numToString
func numToString(n int) string {
	if n == 0 {
		return "0"
	}

	buf := make([]byte, 0, 10)
	for n > 0 {
		d := n % 10
		buf = append([]byte{byte('0' + d)}, buf...)
		n /= 10
	}
	return string(buf)
}
