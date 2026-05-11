package loader

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
)

func LoadGraph(elementPath, relationPath string) (*core.Graph, error) {
	g := core.NewGraph()

	if err := loadElements(elementPath, g); err != nil {
		return nil, err
	}

	if err := loadRelations(relationPath, g); err != nil {
		return nil, err
	}

	return g, nil
}

func loadElements(path string, g *core.Graph) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var e model.CodeElement
		if err := json.Unmarshal([]byte(scanner.Text()), &e); err != nil {
			continue
		}
		g.AddElement(&e)
	}
	return scanner.Err()
}

func loadRelations(path string, g *core.Graph) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var r model.DependencyRelation
		if err := json.Unmarshal([]byte(scanner.Text()), &r); err != nil {
			continue
		}

		// 修正：优先使用已加载的完整 Element，不覆盖
		if fullSource, ok := g.Elements[r.Source.QualifiedName]; ok {
			r.Source = fullSource
		} else {
			g.AddElement(r.Source)
		}

		if fullTarget, ok := g.Elements[r.Target.QualifiedName]; ok {
			r.Target = fullTarget
		} else {
			g.AddElement(r.Target)
		}

		g.AddRelation(&r)
	}
	return scanner.Err()
}
