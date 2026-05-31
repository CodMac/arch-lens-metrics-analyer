---

# 架构度量指标规范：跨域耦合度 (Cross-Domain Coupling - CDC)

---

## 1. 指标定义

**跨域耦合度 (Cross-Domain Coupling, CDC)** 是一种基于图形拓扑分析的架构指标，用于量化一个源文件在业务领域层面的分散程度。CDC 统计的是一个文件所直接依赖的、来自不同**业务社区（Business Community）**的外部代码实体的数量。

在 `arch-lens` 静态分析引擎中，CDC 不再依赖传统的路径命名空间（Package/Class Name）去切分业务领域，而是通过对项目依赖图运行**社区发现算法（如 Louvain 聚类）**，自动将代码库划分为若干逻辑内聚的业务社区。

当一个文件的 $CDC > 4$ 时，意味着该文件不仅物理体量大，而且其代码逻辑强行缝合了超过 4 个在语义上相对独立的业务领域。

---

## 2. 核心价值与局限性

* **价值**：能够自动识别出“领域交叉点”。它不关心文件引用了多少个具体的类（扇出 / Fan-Out），而是关注引用对象的“业务归属地”。它能识别出那些无视架构边界、强行组装高危流程的“上帝文件”，是架构治理中识别物理边界违规的精准雷达。
* **局限性**：在某些架构设计模式下（如高度抽象的流程编排器），CDC 偏高可能属于合理的架构现状。因此，必须将 CDC 与定义的实体数（NDE）和圈复杂度（FCI）进行综合对撞，防止误报。

---

## 3. 计算方式与规则细节

CDC 的计算过程分为两个阶段：**自动业务社区划分** 与 **跨域依赖统计**。

### 3.1 阶段一：自动业务社区发现
引擎在加载完项目的全量依赖图（Graph）后，执行以下步骤：
1. **构建依赖图**：提取所有类、接口、函数间的调用、引用、继承等依赖关系。
2. **Louvain 聚类算法**：基于图的模块度（Modularity）优化，将图划分为 $N$ 个不重叠的社区（Community）。每个 Community 代表一个隐式的、自动归纳出的业务领域。
3. **映射建立**：生成 `CommunityMap: QN -> CommunityID`。

### 3.2 阶段二：跨域耦合统计
对于每个源文件 $F_i$：
1. **获取自身社区**：$C_{\text{self}} = CommunityMap[F_i]$。
2. **提取外部依赖**：扫描所有 $F_i$ 直接依赖的外部实体集 $E_{targets}$。
3. **计算 CDC 值**：
   $$CDC = |\{ CommunityMap[e] \mid e \in E_{targets} \land CommunityMap[e] \neq C_{\text{self}} \}|$$
   该算法统计该文件所触及的不同外部社区的数量，实现真正意义上的跨域度量。

---

## 4. 优势与约束说明 (Edge Cases)

本指标设计旨在实现“零配置（Zero-Config）”的架构分析：

* **零配置原则**：彻底摒弃通过“包路径命名规则”去识别领域。无论代码存放于哪个包，算法仅根据其与周围代码的调用频繁程度与逻辑关联度，自动聚类为业务域。
* **抗噪性**：对于通用的基础服务、Utils 库或 Logging 框架，由于其被广泛调用，它们会被算法自动聚类到一个“中心化基础社区”，从而避免对 CDC 指标产生严重的干扰。
* **稳定性**：为了应对社区划分的波动性，在多次度量扫描中，引擎会通过维持稳定的全局图快照来确保 CDC 指标的趋势一致性。

---

## 5. 详细度量计算实例

假设系统内有大量类，算法自动将其归类为以下社区：
* `Community A (Order Domain)`
* `Community B (Inventory Domain)`
* `Community C (User Domain)`
* `Community D (Finance Domain)`

分析文件 `OrderProcessor.java` (归属于社区 A) 的 CDC 值：

| 引用实体 | 目标所在社区 | 是否跨域 | 说明 |
| --- | --- | --- | --- |
| `OtherLocal.java` | Community A | **No** | 域内调用，不计入 |
| `Target2.java` | Community B | **Yes** | 跨域依赖，CDC++ |
| `Target3.java` | Community C | **Yes** | 跨域依赖，CDC++ |
| `Order.java` | Community A | **No** | 域内调用，不计入 |
| `PayAgent.java` | Community D | **Yes** | 跨域依赖，CDC++ |

* **CDC(OrderProcessor) = 3** (依赖了 B, C, D 三个外部社区)

---

**ArchLens 度量指标规范 - #04 跨域耦合度 (CDC)**

---
