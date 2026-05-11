package metrics

import (
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
)

// CalculateTCC computes Tight Class Cohesion for a class
func CalculateTCC(clsQN string, g *core.Graph) float64 {
	allMethods := FindContainedElements(clsQN, model.Method, g)

	// Filter: Exclude constructors
	var methods []*model.CodeElement
	for _, m := range allMethods {
		isConstructor := false
		if m.Extra != nil {
			if val, ok := m.Extra.Mores["java.method.is_constructor"].(bool); ok && val {
				isConstructor = true
			}
		}
		if !isConstructor {
			methods = append(methods, m)
		}
	}

	if len(methods) <= 1 {
		return 1.0
	}

	methodFields := make(map[string]map[string]bool)
	for _, m := range methods {
		methodFields[m.QualifiedName] = make(map[string]bool)
		outEdges := g.OutEdges[m.QualifiedName]
		for _, edge := range outEdges {
			if edge.Type == model.Use && edge.Target.Kind == model.Field {
				owner := g.GetOwnerClass(edge.Target.QualifiedName)
				if owner == clsQN {
					methodFields[m.QualifiedName][edge.Target.QualifiedName] = true
				}
			}
		}
	}

	np := len(methods) * (len(methods) - 1) / 2
	ndp := 0
	for i := 0; i < len(methods); i++ {
		for j := i + 1; j < len(methods); j++ {
			if shareField(methodFields[methods[i].QualifiedName], methodFields[methods[j].QualifiedName]) {
				ndp++
			}
		}
	}
	return float64(ndp) / float64(np)
}

func shareField(fields1, fields2 map[string]bool) bool {
	for f := range fields1 {
		if fields2[f] {
			return true
		}
	}
	return false
}
