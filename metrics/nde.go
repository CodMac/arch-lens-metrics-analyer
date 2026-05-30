package metrics

import (
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
)

// CalculateNDE computes Number of Declared Entities (all top-level classes/interfaces in file)
func CalculateNDE(fileQN string, g *core.Graph) int {
	count := 0
	for _, e := range g.Elements {
		if e.Path == fileQN && core.IsClassLike(e.Kind) {
			// Check if it's top-level (directly contained by the FILE)
			isTopLevel := false
			for _, edge := range g.InEdges[e.QualifiedName] {
				if edge.Type == model.Contain && edge.Source.Kind == model.File {
					isTopLevel = true
					break
				}
			}

			if isTopLevel {
				count++
			}
		}
	}
	return count
}
