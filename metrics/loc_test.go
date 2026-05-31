package metrics

import (
	"path/filepath"
	"testing"

	"github.com/CodMac/arch-lens-metrics-analyer/loader"
)

func TestCalculateFileLOC_Scenarios(t *testing.T) {
	locDir := filepath.Join("..", "testdata", "x", "java", "metrics", "loc", "dep-out")
	elemPath := filepath.Join(locDir, "element.jsonl")

	graph, err := loader.LoadGraph(elemPath, locDir)
	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	tests := []struct {
		fileQN   string
		expected int
	}{
		{"Standard.java", 10},
		{"Mixed.java", 8},
		{"OnlyComments.java", 0},
		{"OnlyEmpty.java", 0},
		{"MultiLineStrings.java", 8},
	}

	for _, tt := range tests {
		t.Run(tt.fileQN, func(t *testing.T) {
			loc := CalculateFileLOC(tt.fileQN, graph)
			if loc != tt.expected {
				t.Errorf("LOC calculation error for %s: expected %d, got %d", tt.fileQN, tt.expected, loc)
			}
		})
	}
}
