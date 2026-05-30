package detector

import (
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
	"github.com/CodMac/arch-lens-metrics-analyer/metrics"
	"github.com/CodMac/arch-lens-metrics-analyer/metrics/fci"
	"github.com/CodMac/arch-lens-metrics-analyer/metrics/nde"
)

type GodFileResult struct {
	FileQN       string
	LOC          int
	NDE          int
	FCI          int
	CDC          int
	ProjectRoot  string
	IsHyperScale bool
	IsScattered  bool
}

func DetectGodFiles(g *core.Graph) []GodFileResult {
	var results []GodFileResult

	for qn, e := range g.Elements {
		if e.Kind == model.File {
			loc := metrics.CalculateFileLOC(qn, g)
			nde := nde.CalculateNDE(qn, g)
			fci := fci.CalculateFCI(qn, g)
			cdc, root := metrics.CalculateCDC(qn, g)

			// Rule 1: Hyper-Scale
			// (LOC > 1000) AND (FCI > 100)
			isHyperScale := loc > 1000 && fci > 100

			// Rule 2: Scattered Logic
			// (NDE > 15) AND (CDC > 4)
			isScattered := nde > 15 && cdc > 4

			if isHyperScale || isScattered {
				results = append(results, GodFileResult{
					FileQN:       qn,
					LOC:          loc,
					NDE:          nde,
					FCI:          fci,
					CDC:          cdc,
					ProjectRoot:  root,
					IsHyperScale: isHyperScale,
					IsScattered:  isScattered,
				})
			}
		}
	}
	return results
}
