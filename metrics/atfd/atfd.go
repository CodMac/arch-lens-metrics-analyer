package atfd

import (
	"strings"

	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
)

// CalculateATFD computes Access to Foreign Data for a class
func CalculateATFD(clsQN string, g *core.Graph) int {
	externalClasses := make(map[string]bool)
	methods := FindContainedElements(clsQN, model.Method, g)

	for _, m := range methods {
		outEdges := g.OutEdges[m.QualifiedName]
		for _, edge := range outEdges {
			// 1. Direct Field Access
			if edge.Type == model.Use && edge.Target.Kind == model.Field {
				owner := g.GetOwnerClass(edge.Target.QualifiedName)
				if owner != "" && owner != clsQN {
					externalClasses[owner] = true
				}
			}
			// 2. Getter Method Call (Heuristic)
			if edge.Type == model.Call && edge.Target.Kind == model.Method {
				if IsGetter(edge.Target.Name) {
					owner := g.GetOwnerClass(edge.Target.QualifiedName)
					if owner != "" && owner != clsQN {
						externalClasses[owner] = true
					}
				}
			}
		}
	}
	return len(externalClasses)
}

func IsGetter(name string) bool {
	return (strings.HasPrefix(name, "get") && len(name) > 3) ||
		(strings.HasPrefix(name, "is") && len(name) > 2)
}

func FindContainedElements(parentQN string, kind model.ElementKind, g *core.Graph) []*model.CodeElement {
	var results []*model.CodeElement
	for _, edge := range g.OutEdges[parentQN] {
		if edge.Type == model.Contain && edge.Target.Kind == kind {
			results = append(results, edge.Target)
		}
	}
	return results
}
