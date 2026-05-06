package detector

import (
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
	"strings"
)

type CircularDependencyResult struct {
	Level      string
	Components [][]string
}

func DetectCircularDependencies(g *core.Graph, kind model.ElementKind) *CircularDependencyResult {
	nodes := make([]string, 0)
	for qn, e := range g.Elements {
		if e.Kind == kind {
			nodes = append(nodes, qn)
		}
	}

	adj := make(map[string]map[string]bool)
	for _, rel := range g.Relations {
		if !isDependencyRelation(rel.Type) {
			continue
		}

		srcOwner := getOwnerAtLevel(rel.Source, kind, g)
		tgtOwner := getOwnerAtLevel(rel.Target, kind, g)

		if srcOwner != "" && tgtOwner != "" && srcOwner != tgtOwner {
			if _, ok := adj[srcOwner]; !ok {
				adj[srcOwner] = make(map[string]bool)
			}
			adj[srcOwner][tgtOwner] = true
		}
	}

	return &CircularDependencyResult{
		Level:      string(kind),
		Components: tarjan(nodes, adj),
	}
}

func getOwnerAtLevel(e *model.CodeElement, targetKind model.ElementKind, g *core.Graph) string {
	if e.Kind == targetKind {
		return e.QualifiedName
	}

	if targetKind == model.Package {
		return getPackageName(e.QualifiedName)
	}

	if targetKind == model.Class || targetKind == model.Interface {
		return g.GetOwnerClass(e.QualifiedName)
	}

	return ""
}

func getPackageName(qn string) string {
	idx := strings.LastIndex(qn, ".")
	if idx == -1 {
		return ""
	}
	return qn[:idx]
}

func isDependencyRelation(t model.DependencyType) bool {
	return t == model.Call || t == model.Use || t == model.Extend || t == model.Implement || t == model.Create
}

func tarjan(nodes []string, adj map[string]map[string]bool) [][]string {
	index := 0
	indices := make(map[string]int)
	lowlink := make(map[string]int)
	onStack := make(map[string]bool)
	stack := []string{}
	result := [][]string{}

	var strongconnect func(v string)
	strongconnect = func(v string) {
		indices[v] = index
		lowlink[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for w := range adj[v] {
			if _, exists := indices[w]; !exists {
				strongconnect(w)
				lowlink[v] = min(lowlink[v], lowlink[w])
			} else if onStack[w] {
				lowlink[v] = min(lowlink[v], indices[w])
			}
		}

		if lowlink[v] == indices[v] {
			component := []string{}
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				component = append(component, w)
				if w == v {
					break
				}
			}
			if len(component) > 1 {
				result = append(result, component)
			}
		}
	}

	for _, v := range nodes {
		if _, exists := indices[v]; !exists {
			strongconnect(v)
		}
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
