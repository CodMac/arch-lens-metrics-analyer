package metrics

import (
	"path/filepath"
	"testing"

	"github.com/CodMac/arch-lens-metrics-analyer/loader"
)

func TestCalculateNDE_Scenarios(t *testing.T) {
	ndeDir := filepath.Join("..", "testdata", "x", "java", "metrics", "nde", "dep-out")
	elemPath := filepath.Join(ndeDir, "element.jsonl")

	graph, err := loader.LoadGraph(elemPath, ndeDir)
	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	tests := []struct {
		fileQN   string
		expected int
	}{
		{"SingleClass.java", 1},
		{"MultipleTopLevel.java", 5},
		{"InnerClasses.java", 1},
		{"AnonymousClasses.java", 1},
		{"MixedTypes.java", 3},
		{"JustImports.java", 1},
		{"Empty.java", 0},
	}

	for _, tt := range tests {
		t.Run(tt.fileQN, func(t *testing.T) {
			nde := CalculateNDE(tt.fileQN, graph)
			if nde != tt.expected {
				t.Errorf("NDE calculation error for %s: expected %d, got %d", tt.fileQN, tt.expected, nde)
			}
		})
	}
}
