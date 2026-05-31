package loc

import (
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-dep-analyer/x/java"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
)

func CalculateFileLOC(fileQN string, g *core.Graph) int {
	if e, ok := g.Elements[fileQN]; ok && e.Kind == model.File {
		if e.Extra != nil {
			if val, ok := e.Extra.Mores[java.FileLOC].(float64); ok {
				return int(val)
			} else if val, ok := e.Extra.Mores[java.FileLOC].(int); ok {
				return val
			}
		}
	}
	return 0
}
