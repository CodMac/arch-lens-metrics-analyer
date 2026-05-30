package god_class

import (
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
	"github.com/CodMac/arch-lens-metrics-analyer/detector"
	"github.com/CodMac/arch-lens-metrics-analyer/metrics/atfd"
	"github.com/CodMac/arch-lens-metrics-analyer/metrics/tcc"
	"github.com/CodMac/arch-lens-metrics-analyer/metrics/wmc"
)

type GodClassResult struct {
	ClassQN             string
	ATFD                int
	TCC                 float64
	WMC                 int
	MethodDensity       float64
	IsGodFormula        bool
	IsConcentrationWarn bool
}

func DetectGodClasses(g *core.Graph) []GodClassResult {
	var results []GodClassResult

	// Helper to get package name for any element
	getPkg := func(qn string) string {
		curr := qn
		if e, ok := g.Elements[curr]; ok && e.Kind == model.Method {
			curr = g.GetOwnerClass(curr)
		}
		if curr == "" {
			return ""
		}
		return detector.getPackageName(curr)
	}

	// Pre-calculate package method counts for Rule 2
	pkgMethodCounts := make(map[string]int)
	for _, e := range g.Elements {
		if e.Kind == model.Method {
			pkg := getPkg(e.QualifiedName)
			if pkg != "" {
				pkgMethodCounts[pkg]++
			}
		}
	}

	for qn, e := range g.Elements {
		if e.Kind == model.Class {
			atfd := atfd.CalculateATFD(qn, g)
			tcc := tcc.CalculateTCC(qn, g)
			wmc := wmc.CalculateWMC(qn, g)

			// Rule 1: God Formula
			isGodFormula := wmc > 47 && atfd > 5 && tcc < 0.33

			// Rule 2: Concentration
			classMethods := atfd.FindContainedElements(qn, model.Method, g)
			pkg := getPkg(qn)
			pkgMethods := pkgMethodCounts[pkg]

			methodDensity := 0.0
			if pkgMethods > 0 {
				methodDensity = float64(len(classMethods)) / float64(pkgMethods)
			}
			isConcentrationWarn := methodDensity > 0.33

			if isGodFormula || isConcentrationWarn {
				results = append(results, GodClassResult{
					ClassQN:             qn,
					ATFD:                atfd,
					TCC:                 tcc,
					WMC:                 wmc,
					MethodDensity:       methodDensity,
					IsGodFormula:        isGodFormula,
					IsConcentrationWarn: isConcentrationWarn,
				})
			}
		}
	}
	return results
}
