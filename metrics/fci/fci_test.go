package fci

import (
	"path/filepath"
	"testing"

	"github.com/CodMac/arch-lens-metrics-analyer/loader"
)

func TestCalculateFCI_Scenarios(t *testing.T) {
	fciDir := filepath.Join("..", "..", "testdata", "x", "java", "metrics", "fci", "dep-out")
	elemPath := filepath.Join(fciDir, "element.jsonl")

	graph, err := loader.LoadGraph(elemPath, fciDir)
	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	tests := []struct {
		fileQN   string
		expected int
	}{
		{"Basic.java", 1},
		{"ControlFlow.java", 10},       // complexity=10 (matched dep-out)
		{"Logical.java", 8},            // testLogical(4) + nestedLogical(4)
		{"InnerClasses.java", 6},       // outer(1) + innerMethod(2) + deepMethod(2) + staticInnerMethod(1)
		{"LambdaAndAnonymous.java", 5}, // test(3) + anonymous run(2). Note: lambda(2) is inside test(3).
		{"Exclusions.java", 5},         // abstractMethod(1) + doSomething(1) + getName(1) + setName(1) + getFormattedName(1) - Current parser doesn't exclude them
		{"MultipleClasses.java", 5},    // mA(1) + mB(2) + mC(2)
		{"DeepNesting.java", 5},        // deep(5)
		{"Complex.java", 10},           // m1(3) + m2(5) + m3(2) = 10 (matched dep-out)
		{"DataOnly.java", 0},
		{"Empty.java", 0},
	}

	for _, tt := range tests {
		t.Run(tt.fileQN, func(t *testing.T) {
			fci := CalculateFCI(tt.fileQN, graph)
			if fci != tt.expected {
				t.Errorf("FCI calculation error for %s: expected %d, got %d", tt.fileQN, tt.expected, fci)
			}
		})
	}
}
