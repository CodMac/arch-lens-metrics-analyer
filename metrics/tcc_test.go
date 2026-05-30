package metrics

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/CodMac/arch-lens-metrics-analyer/loader"
)

func TestCalculateTCC_Scenarios(t *testing.T) {
	tccDir := filepath.Join("..", "testdata", "x", "java", "metrics", "tcc", "dep-out")
	elemPath := filepath.Join(tccDir, "element.jsonl")

	graph, err := loader.LoadGraph(elemPath, tccDir)
	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	testCases := []struct {
		classQN  string
		expected float64
	}{
		{"com.test.tcc.HighCohesion", 0.6666666666666666},
		{"com.test.tcc.MediumCohesion", 0.3333333333333333},
		{"com.test.tcc.LowCohesion", 0.0},
		{"com.test.tcc.SingleMethod", 1.0},
	}

	for _, tc := range testCases {
		t.Run(tc.classQN, func(t *testing.T) {
			tcc := CalculateTCC(tc.classQN, graph)
			if math.Abs(tcc-tc.expected) > 0.0001 {
				t.Errorf("TCC error for %s: expected %.2f, got %.2f", tc.classQN, tc.expected, tcc)
			}
		})
	}
}
