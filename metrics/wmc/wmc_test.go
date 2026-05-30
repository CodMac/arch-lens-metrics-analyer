package wmc

import (
	"path/filepath"
	"testing"

	"github.com/CodMac/arch-lens-metrics-analyer/loader"
)

func TestCalculateWMC_Scenarios(t *testing.T) {
	wmcDir := filepath.Join("..", "..", "testdata", "x", "java", "metrics", "wmc", "dep-out")
	elemPath := filepath.Join(wmcDir, "element.jsonl")

	graph, err := loader.LoadGraph(elemPath, wmcDir)
	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	testCases := []struct {
		classQN  string
		expected int
	}{
		{"com.test.wmc.SimpleClass", 3},
		{"com.test.wmc.GodClassCandidate", 49},
	}

	for _, tc := range testCases {
		t.Run(tc.classQN, func(t *testing.T) {
			wmc := CalculateWMC(tc.classQN, graph)
			if wmc != tc.expected {
				t.Errorf("WMC error for %s: expected %d, got %d", tc.classQN, tc.expected, wmc)
			}
		})
	}
}
