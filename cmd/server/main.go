package main

import (
	"flag"
	"log"
	"os"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/app"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/logger"
)

type stdLogger struct{}

func (stdLogger) Info(msg string)  { log.Println("[INFO]", msg) }
func (stdLogger) Debug(msg string) { log.Println("[DEBUG]", msg) }
func (stdLogger) Error(msg string) { log.Println("[ERROR]", msg) }

func main() {

	pattern := flag.String("pattern", "", "search pattern (required)")
	ignoreCase := flag.Bool("i", false, "ignore case")
	showNumber := flag.Bool("n", false, "show line numbers")
	invert := flag.Bool("v", false, "invert match")
	filePath := flag.String("file", "", "path to input file (optional)")

	flag.Parse()

	if *pattern == "" {
		log.Println("[ERROR] pattern is required")
		flag.Usage()
		os.Exit(1)
	}

	var logg logger.Logger = stdLogger{}

	logg.Info("main: application started")

	params := app.Params{
		Pattern:    *pattern,
		IgnoreCase: *ignoreCase,
		ShowNumber: *showNumber,
		Invert:     *invert,
		FilePath:   *filePath,
	}

	if err := app.Run(params, logg); err != nil {
		logg.Error("main: application failed")
		os.Exit(1)
	}

	logg.Info("main: application finished")
}
