# 源码架构分析协议：上帝文件 (God File)

## 1. 缺陷定义
**上帝文件** 是指那些由于缺乏物理分割，导致其职责边界模糊、维护风险极高的源代码文件。它不仅违反了单一职责原则（SRP），还因为文件过大导致编译器解析缓慢、版本控制冲突（Merge Conflict）频发，以及开发者的认知负荷过重。

---

## 2. 典型场景与代码示例

### 2.1 混合职责的巨型文件 (The Hybrid Blob)
一个文件里既定义了数据模型（Model），又定义了业务逻辑（Service），甚至还包含了工具函数（Utils）和常量（Constants）。

### 2.2 多重类嵌套文件 (Nested Class Overflow)
在 Java 或 C# 中，一个公共类文件内部嵌套了 10 几个私有内部类（Inner Classes），逻辑错综复杂。

---

## 3. 抽象度量指标：深度量化与计算方式 (Metrics)

上帝文件的检测侧重于“物理规模”与“职责分布”的交叉比对。

| 全称 | 简名 | 计算方式与定义 | 阈值参考 |
| :--- | :--- | :--- | :--- |
| **Lines of Code** | **LOC** | 统计文件中的总行数（排除空行和纯注释行）。 | $> 1000$ |
| **Number of Declared Entities** | **NDE** | 统计文件内定义的**顶层元素**数量，包括类、接口、函数。 | $> 10$ |
| **File Complexity Index** | **FCI** | 统计文件中所有方法/函数的圈复杂度（CC）的总和。 | $> 100$ |
| **Cross-Domain Coupling** | **CDC** | 统计文件直接依赖的外部独立业务社区（Clusters）的数量。 | $> 4$ |

---

## 4. 缺陷命中规则 (Detection Rules)

判定上帝文件的核心在于：**物理规模巨大**且**逻辑领域发散**。

### 规则 1：超大规模文件判定 (Hyper-Scale)
当一个文件的物理规模和逻辑复杂度双双触顶时命中。
$$Rule_{Hyper} = (LOC > 1000) \land (FCI > 100)$$

### 规则 2：逻辑散乱判定 (Scattered Logic)
如果文件规模中等，但定义了太多的独立实体且跨越多个自动聚类的业务领域。
$$Rule_{Scattered} = (NDE > 15) \land (CDC > 4)$$

---

## 5. 检测算法逻辑 (Pseudo Code)

```python
def detect_god_file(project_files, graph):
    # 1. 执行自动聚类，建立实体到业务社区的映射
    community_map = run_louvain_community_detection(graph)
    
    for file in project_files:
        # 2. 统计常规指标 (LOC, NDE, FCI)
        loc = count_logical_loc(file)
        nde = count_top_level_entities(file)
        fci = sum(calculate_cc(m) for m in file.methods)
            
        # 3. 统计跨域耦合 (CDC)
        # 获取文件依赖的所有外部实体
        external_entities = get_dependencies(file)
        # 统计其涉及的唯一业务社区数量
        communities = {community_map[e] for e in external_entities if community_map[e] != community_map[file]}
        cdc = len(communities)
        
        # 4. 命中判定
        if (loc > 1000 and fci > 100) or (nde > 15 and cdc > 4):
            report_issue(file, "God File", {"LOC": loc, "NDE": nde, "FCI": fci, "CDC": cdc})
```

---

## 6. 治理建议与决策矩阵

| 指标表现 | 根本原因 | 建议重构动作 |
| :--- | :--- | :--- |
| **LOC > 2000** | 纯粹的逻辑堆砌 | **最高优先级**：按功能进行物理切分，每个文件不超过 500 行。 |
| **NDE > 20** | 文件变成了“垃圾袋” | **中优先级**：按类型归类，将 Helper、Constant、Entity 分别归档。 |
| **CDC 较高** | 编排逻辑过于耦合 | **低优先级**：引入事件驱动，减少硬编码的跨域依赖引用。 |
