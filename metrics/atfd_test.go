package metrics

import (
	"path/filepath"
	"testing"

	"github.com/CodMac/arch-lens-metrics-analyer/loader"
)

func TestCalculateATFD_Scenarios(t *testing.T) {
	elemPath := filepath.Join("..", "testdata", "x", "java", "atfd", "element.jsonl")
	relPath := filepath.Join("..", "testdata", "x", "java", "atfd", "relation.jsonl")

	graph, err := loader.LoadGraph(elemPath, relPath)
	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	// 预期值：AtfdTarget 访问了 ForeignData, AnotherForeign, BaseClass
	// 虽然访问了 ForeignData 的多个成员，但 ATFD 应去重计为 3
	clsQN := "com.test.atfd.AtfdTarget"
	atfd := CalculateATFD(clsQN, graph)
	expected := 3

	if atfd != expected {
		t.Errorf("ATFD calculation error for %s: expected %d, got %d", clsQN, expected, atfd)
	}
}
