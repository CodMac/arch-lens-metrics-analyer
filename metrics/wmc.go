package metrics

import (
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
)

const MethodComplexity = "java.method.metrics.complexity"

// CalculateWMC computes Weighted Methods per Class
func CalculateWMC(clsQN string, g *core.Graph) int {
	wmc := 0
	// Include Methods and ScopeBlocks (initializers)
	methods := FindContainedElements(clsQN, model.Method, g)
	blocks := FindContainedElements(clsQN, model.ScopeBlock, g)

	all := append(methods, blocks...)

	for _, m := range all {
		complexity := 1 // Default complexity
		if m.Extra != nil {
			if val, ok := m.Extra.Mores[MethodComplexity].(float64); ok {
				complexity = int(val)
			} else if val, ok := m.Extra.Mores[MethodComplexity].(int); ok {
				complexity = val
			}
		}
		wmc += complexity
	}
	return wmc
}
