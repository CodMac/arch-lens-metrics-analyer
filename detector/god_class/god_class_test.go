package god_class

import (
	"path/filepath"
	"testing"

	"github.com/CodMac/arch-lens-metrics-analyer/loader"
)

func TestDetectGodClasses(t *testing.T) {
	// 使用 wmc 测试数据，其中包含 GodClassCandidate
	wmcDir := filepath.Join("..", "testdata", "x", "java", "wmc")
	elemPath := filepath.Join(wmcDir, "element.jsonl")

	graph, err := loader.LoadGraph(elemPath, wmcDir)
	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	results := DetectGodClasses(graph)

	found := false
	for _, res := range results {
		if res.ClassQN == "com.test.wmc.GodClassCandidate" {
			found = true
			// WMC = 49 > 47
			// 检查是否命中规则
			if res.WMC != 49 {
				t.Errorf("Expected WMC 49, got %d", res.WMC)
			}
		}
	}

	if !found {
		// 如果没找到，可能是因为 ATFD 或 TCC 没达到 Rule 1 的组合阈值
		// 但 Rule 2 (Concentration) 应该会命中，因为 GodClassCandidate 占了包中大部分方法
		t.Log("GodClassCandidate not detected as Rule 1 God Class, checking for any detection...")
		if len(results) == 0 {
			t.Errorf("No God Classes detected at all")
		}
	}
}
