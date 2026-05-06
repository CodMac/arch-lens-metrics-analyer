package detector

import (
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
	"github.com/CodMac/arch-lens-metrics-analyer/metrics"
)

type GodClassResult struct {
	ClassQN string
	ATFD    int
	TCC     float64
	WMC     int
}

func DetectGodClasses(g *core.Graph) []GodClassResult {
	var results []GodClassResult
	for qn, e := range g.Elements {
		if e.Kind == model.Class {
			atfd := metrics.CalculateATFD(qn, g)
			tcc := metrics.CalculateTCC(qn, g)

			if atfd > 0 {
				// Debug log to confirm calculation is happening
				// fmt.Printf("DEBUG: %s ATFD: %d, TCC: %.2f\n", qn, atfd, tcc)
			}

			if atfd > 5 && tcc < 0.33 {
				results = append(results, GodClassResult{
					ClassQN: qn,
					ATFD:    atfd,
					TCC:     tcc,
				})
			}
		}
	}
	return results
}
