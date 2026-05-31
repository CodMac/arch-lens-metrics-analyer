package metrics

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/CodMac/arch-lens-dep-analyer/model"
	"github.com/CodMac/arch-lens-metrics-analyer/core"
	"gonum.org/v1/gonum/graph/community"
	"gonum.org/v1/gonum/graph/simple"
)

func CalculateCDC(fileQN string, g *core.Graph) (int, string) {
	dg := simple.NewDirectedGraph()
	nodeToID := make(map[string]int64)
	idToNode := make(map[int64]string)
	var idCounter int64 = 0

	fmt.Println("--- [DEBUG] 1. 注册可用节点 ---")
	for qn, e := range g.Elements {
		if core.IsClassLike(e.Kind) {
			nodeToID[qn] = idCounter
			idToNode[idCounter] = qn
			dg.AddNode(simple.Node(idCounter))
			idCounter++
		}
	}

	fmt.Println("--- [DEBUG] 2. 尝试添加边 (逐条核对) ---")
	edges := make(map[string]map[string]bool)
	for _, rel := range g.Relations {
		if rel.Type != "USE" {
			continue
		}

		srcCls := getEnclosingClass(rel.Source.QualifiedName, g)
		tgtCls := getEnclosingClass(rel.Target.QualifiedName, g)

		fmt.Printf("Rel: %s -> %s | Resolved: %s -> %s\n",
			rel.Source.QualifiedName, rel.Target.QualifiedName, srcCls, tgtCls)

		if srcCls != "" && tgtCls != "" && srcCls != tgtCls {
			// 检查：归约后的类名是否在我们的 nodeToID 中
			_, ok1 := nodeToID[srcCls]
			_, ok2 := nodeToID[tgtCls]

			if !ok1 || !ok2 {
				fmt.Printf(">> 警告: 无法在 nodeToID 中找到归约后的类名 (src:%v, tgt:%v)\n", ok1, ok2)
			} else {
				if _, ok := edges[srcCls]; !ok {
					edges[srcCls] = make(map[string]bool)
				}
				edges[srcCls][tgtCls] = true
				dg.SetEdge(simple.Edge{F: simple.Node(nodeToID[srcCls]), T: simple.Node(nodeToID[tgtCls])})
			}
		}
	}

	fmt.Printf("[DEBUG] 最终图边数: %d\n", dg.Edges().Len())
	if dg.Edges().Len() == 0 {
		return 0, "no-edges"
	}

	// 3. 聚类
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	reduced := community.Modularize(dg, 1.0, src)
	communities := reduced.Communities()

	fmt.Println("--- [DEBUG] 3. 社区划分结果 ---")
	communityMap := make(map[string]int)
	for i, nodes := range communities {
		for _, node := range nodes {
			qn := idToNode[node.ID()]
			communityMap[qn] = i
			fmt.Printf("Community %d: %s\n", i, qn)
		}
	}

	// 4. 计算 CDC
	fileClasses := getClassesInFile(fileQN, g)
	externalCommunities := make(map[int]bool)
	for _, cls := range fileClasses {
		selfComm, ok := communityMap[cls]
		if !ok {
			continue
		}
		for tgt := range edges[cls] {
			if tgtComm, found := communityMap[tgt]; found && tgtComm != selfComm {
				externalCommunities[tgtComm] = true
				fmt.Printf("[DEBUG] 发现跨域依赖: %s(Comm:%d) -> %s(Comm:%d)\n", cls, selfComm, tgt, tgtComm)
			}
		}
	}
	return len(externalCommunities), "louvain-optimized"
}

// 采用图拓扑追溯，不再依赖字符串分割
func getEnclosingClass(qn string, g *core.Graph) string {
	curr := qn
	for {
		if e, ok := g.Elements[curr]; ok && core.IsClassLike(e.Kind) {
			return curr
		}
		foundParent := false
		for _, edge := range g.InEdges[curr] {
			if edge.Type == model.Contain {
				curr = edge.Source.QualifiedName
				foundParent = true
				break
			}
		}
		if !foundParent {
			break
		}
	}
	return ""
}

func getClassesInFile(fileQN string, g *core.Graph) []string {
	var classes []string
	if e, ok := g.Elements[fileQN]; ok && core.IsClassLike(e.Kind) {
		classes = append(classes, fileQN)
	}
	for _, edge := range g.OutEdges[fileQN] {
		if edge.Type == model.Contain && core.IsClassLike(edge.Target.Kind) {
			classes = append(classes, edge.Target.QualifiedName)
		}
	}
	return classes
}
