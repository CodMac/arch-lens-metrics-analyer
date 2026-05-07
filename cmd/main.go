package main

import (
	"flag"
	"fmt"
	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/detector"
	"github.com/CodMac/arch-lens-metrics-analyer/loader"
	"log"
)

func main() {
	elemPath := flag.String("elem", "", "path to element.jsonl")
	relPath := flag.String("rel", "", "path to relation.jsonl")
	flag.Parse()

	if *elemPath == "" || *relPath == "" {
		log.Fatal("Usage: arch-lens-metrics -elem <path> -rel <path>")
	}

	graph, err := loader.LoadGraph(*elemPath, *relPath)
	if err != nil {
		log.Fatalf("Failed to load graph: %v", err)
	}

	fmt.Printf("Analyzing %d elements...\n", len(graph.Elements))

	// 1. Detect God Classes
	fmt.Println("\n--- God Class Detection ---")
	godClasses := detector.DetectGodClasses(graph)
	for _, res := range godClasses {
		status := ""
		if res.IsGodFormula {
			status += "[Formula]"
		}
		if res.IsConcentrationWarn {
			status += "[Concentration]"
		}
		fmt.Printf("[GOD CLASS] %s %s\n", res.ClassQN, status)
		fmt.Printf("    Metrics: WMC=%d, ATFD=%d, TCC=%.2f, Density=%.2f\n",
			res.WMC, res.ATFD, res.TCC, res.MethodDensity)
	}

	// 2. Detect Circular Dependencies (Class level)
	fmt.Println("\n--- Circular Dependency Detection (Class Level) ---")
	cycles := detector.DetectCircularDependencies(graph, model.Class)
	for _, comp := range cycles.Components {
		fmt.Printf("[CYCLE] nodes: %v\n", comp)
	}
}
