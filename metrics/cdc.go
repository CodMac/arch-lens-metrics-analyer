package metrics

import (
	"strings"

	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
)

// CalculateCDC computes Cross-Domain Coupling and returns the count and identified project root
func CalculateCDC(fileQN string, g *core.Graph) (int, string) {
	// 1. Identify all source packages to find project root LCP
	var sourcePackages []string
	for _, e := range g.Elements {
		if e.Kind == model.Package && e.IsFormSource && e.QualifiedName != "" {
			sourcePackages = append(sourcePackages, e.QualifiedName)
		}
	}

	projectRoot := findLCP(sourcePackages)
	if projectRoot == "" {
		return 0, ""
	}

	// 2. Identify the domain of the current file
	if _, ok := g.Elements[fileQN]; !ok {
		return 0, projectRoot
	}

	// Get the package of the file
	filePkg := g.GetElementPackage(fileQN)
	currentDomain := getDomain(filePkg, projectRoot)

	// 3. Collect distinct domains from IMPORT relations
	externalDomains := make(map[string]bool)

	// CDC usually looks at what the file imports
	// Note: IMPORT relations in our model are Source: FILE -> Target: CLASS/INTERFACE
	for _, edge := range g.OutEdges[fileQN] {
		if edge.Type == model.Import && edge.Target.IsFormSource {
			targetPkg := getPackageFromQN(edge.Target.QualifiedName)
			targetDomain := getDomain(targetPkg, projectRoot)

			if targetDomain != "" && targetDomain != currentDomain {
				externalDomains[targetDomain] = true
			}
		}
	}

	return len(externalDomains), projectRoot
}

func findLCP(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	lcp := strs[0]
	for i := 1; i < len(strs); i++ {
		for !strings.HasPrefix(strs[i], lcp) {
			if len(lcp) == 0 {
				return ""
			}
			idx := strings.LastIndex(lcp, ".")
			if idx == -1 {
				lcp = ""
			} else {
				lcp = lcp[:idx]
			}
		}
	}
	return lcp
}

func getDomain(pkg, root string) string {
	if pkg == root {
		return "root"
	}
	if !strings.HasPrefix(pkg, root) {
		return ""
	}

	suffix := strings.TrimPrefix(pkg, root)
	if strings.HasPrefix(suffix, ".") {
		suffix = suffix[1:]
	}

	parts := strings.Split(suffix, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func getPackageFromQN(qn string) string {
	idx := strings.LastIndex(qn, ".")
	if idx == -1 {
		return ""
	}
	return qn[:idx]
}
