package metrics

import (
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
)

// CalculateFCI computes File Complexity Index (sum of CC of all methods in file)
func CalculateFCI(fileQN string, g *core.Graph) int {
	fci := 0
	for _, e := range g.Elements {
		if e.Path == fileQN && (e.Kind == model.Method) {
			complexity := 1

			if e.Extra != nil {
				if val, ok := e.Extra.Mores[MethodComplexity].(float64); ok {
					complexity = int(val)
				} else if val, ok := e.Extra.Mores[MethodComplexity].(int); ok {
					complexity = val
				}
			}
			fci += complexity
		}
	}
	return fci
}
