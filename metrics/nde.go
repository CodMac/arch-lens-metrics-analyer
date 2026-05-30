package metrics

import (
	"github.com/CodMac/arch-lens-metrics-analyer/core"
)

// CalculateNDE computes Number of Declared Entities (all classes/interfaces in file)
func CalculateNDE(fileQN string, g *core.Graph) int {
	count := 0
	for _, e := range g.Elements {
		if e.Path == fileQN && core.IsClassLike(e.Kind) {
			count++
		}
	}
	return count
}
