package core

import (
	"strings"

	"github.com/CodMac/arch-lens-dep-analyer/model"
)

type Graph struct {
	Elements  map[string]*model.CodeElement
	Relations []*model.DependencyRelation

	// Fast lookup
	OutEdges map[string][]*model.DependencyRelation
	InEdges  map[string][]*model.DependencyRelation
}

func NewGraph() *Graph {
	return &Graph{
		Elements: make(map[string]*model.CodeElement),
		OutEdges: make(map[string][]*model.DependencyRelation),
		InEdges:  make(map[string][]*model.DependencyRelation),
	}
}

func (g *Graph) AddElement(e *model.CodeElement) {
	if e.QualifiedName == "" && e.Kind == model.File {
		return
	}
	g.Elements[e.QualifiedName] = e
}

func (g *Graph) AddRelation(r *model.DependencyRelation) {
	g.Relations = append(g.Relations, r)

	src := r.Source.QualifiedName
	tgt := r.Target.QualifiedName
	g.OutEdges[src] = append(g.OutEdges[src], r)
	g.InEdges[tgt] = append(g.InEdges[tgt], r)
}

// GetOwnerClass returns the QN of the class-like element containing this element
func (g *Graph) GetOwnerClass(qn string) string {
	curr := qn
	for {
		parent := GetParentQN(curr)
		if parent == "" {
			return ""
		}
		if elem, ok := g.Elements[parent]; ok {
			if IsClassLike(elem.Kind) {
				return parent
			}
		}
		curr = parent
	}
}

// GetElementPackage returns the QN of the package containing this element (especially for FILE)
func (g *Graph) GetElementPackage(qn string) string {
	// For File, check incoming CONTAIN edges from Package
	if e, ok := g.Elements[qn]; ok && e.Kind == model.File {
		for _, edge := range g.InEdges[qn] {
			if edge.Type == model.Contain && edge.Source.Kind == model.Package {
				return edge.Source.QualifiedName
			}
		}
	}
	// For others, use parent QN logic
	return getPackageFromQN(qn)
}

func (g *Graph) GetFileLOC(fileQN string) int {
	if e, ok := g.Elements[fileQN]; ok && e.Kind == model.File {
		if e.Extra != nil {
			if val, ok := e.Extra.Mores["java.file.metrics.loc"].(float64); ok {
				return int(val)
			} else if val, ok := e.Extra.Mores["java.file.metrics.loc"].(int); ok {
				return val
			}
		}
	}
	return 0
}

func getPackageFromQN(qn string) string {
	idx := strings.LastIndex(qn, ".")
	if idx == -1 {
		return ""
	}
	return qn[:idx]
}

func IsClassLike(k model.ElementKind) bool {
	return k == model.Class || k == model.Interface || k == model.Enum || k == model.KAnnotation || k == model.AnonymousClass
}

func GetParentQN(qn string) string {
	idx := strings.LastIndex(qn, ".")
	if idx == -1 {
		return ""
	}
	return qn[:idx]
}
