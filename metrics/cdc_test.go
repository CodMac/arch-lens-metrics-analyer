package metrics

import (
	"path/filepath"
	"testing"

	"github.com/CodMac/arch-lens-metrics-analyer/loader"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCDC_Scenarios(t *testing.T) {
	cdcDir := filepath.Join("..", "testdata", "x", "java", "metrics", "cdc", "dep-out")
	elemPath := filepath.Join(cdcDir, "element.jsonl")

	graph, err := loader.LoadGraph(elemPath, cdcDir)
	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	fileQN := "com.example.metrics.CrossDomainProcessor"
	cdc, _ := CalculateCDC(fileQN, graph)

	assert.Equal(t, 3, cdc, "CDC 应该精确识别出 B, C, D 三个社区")
}
