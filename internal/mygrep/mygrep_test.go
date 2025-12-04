package mygrep

import (
	"strings"
	"testing"

	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/logger"
	"github.com/Vladimirmoscow84/Distributed_CLI_tool/internal/model"
)

func TestBasic(t *testing.T) {
	input := "spartak\nchampion\nspartak forever"
	cfg := model.GrepConfig{
		Pattern:    "spartak",
		IgnoreCase: false,
		ShowNumber: false,
		Invert:     false,
	}

	response, err := Run(cfg, strings.NewReader(input), logger.NopLogger{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(response.Lines))
	}
}

func TestIgnoreCase(t *testing.T) {
	input := "Spartak\nchan\nSPARTAK"

	cfg := model.GrepConfig{
		Pattern:    "spartak",
		IgnoreCase: true,
	}

	res, err := Run(cfg, strings.NewReader(input), logger.NopLogger{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
}

func TestInvert(t *testing.T) {
	input := "spartak\narsenal\norel"

	cfg := model.GrepConfig{
		Pattern: "spartak",
		Invert:  true,
	}

	res, err := Run(cfg, strings.NewReader(input), logger.NopLogger{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
}

func TestShowNumber(t *testing.T) {
	input := "spartak\ncska\nspartak"

	cfg := model.GrepConfig{
		Pattern:    "spartak",
		ShowNumber: true,
	}

	res, err := Run(cfg, strings.NewReader(input), logger.NopLogger{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Lines[0][0:2] != "1:" {
		t.Fatalf("expected line number prefix, got %s", res.Lines[0])
	}
}

func TestInvalidRegex(t *testing.T) {
	cfg := model.GrepConfig{
		Pattern: "([",
	}

	_, err := Run(cfg, strings.NewReader("test"), logger.NopLogger{})
	if err == nil {
		t.Fatalf("expected regex error, got nil")
	}
}
