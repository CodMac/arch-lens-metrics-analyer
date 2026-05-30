
---

# 架构度量指标规范：文件复杂度指数 (File Complexity Index - FCI)

## 1. 指标定义

**文件复杂度指数 (File Complexity Index, FCI)** 是指在一个特定源文件的物理边界内，其内部包含的所有函数及类成员方法的圈复杂度（Cyclomatic Complexity, CC）的绝对累加总和。

$$FCI = \sum_{i=1}^{n} CC_i$$

在 `arch-lens` 静态分析引擎对“上帝文件 (God File)”的判定算法中，FCI 是用于量化文件内部“逻辑分支承载极限”的核心综合指标。当一个文件的 $FCI$ 超过设定阈值时，说明该文件内聚集了密集的判定分支，其代码的测试路径和维护成本极高，通常暗示其物理分割不彻底。

---

## 2. 核心价值与局限性

* **价值**：FCI 能够精准捕捉物理文件层面的“逻辑爆炸”。它不仅考虑了代码的体积，更看重代码的控制流复杂度。一个只有 300 行但充斥着密集条件分叉、循环嵌套的文件，其 FCI 值会远高于一个包含 2000 行纯数据结构、赋值或 Getter/Setter 方法的文件。
* **局限性**：FCI 关注的是整个文件控制流分支的**物理总量**。若一个文件内部定义了大量极其简单、但各自圈复杂度为 2 的短小独立函数，其 FCI 也会因为线性累加而破表。因此，FCI 必须配合声明实体数（NDE）和单函数最大圈复杂度联合评估。

---

## 3. 计算方式与规则细节

在 `arch-lens` 的语法分析阶段，算法提取出该文件内的有效函数集合，计算每个函数的圈复杂度（基础分为 1，每遇到一个控制流分支分叉节点则分值 +1），最后将文件中所有函数的圈复杂度进行线性累加。

### 3.1 导致函数圈复杂度 (CC) 增加的 AST 节点元素

1. **条件与循环分支语句**：`if`、`else if`、`for`、`while`、`catch` 块。
2. **多路选择分支**：`switch-case` 结构中的 `case` 节点（注：整个 `switch` 头部节点自身计为 0，仅对具体的 `case` 分支节点计分）。
3. **逻辑运算符**：条件表达式中的逻辑与 `&&`、逻辑或 `||`（因其引入了语言级的短路求值控制流分叉）。

### 3.2 计入与排除的函数边界规则

1. **计入有效函数集合的元素**：顶层独立函数、类/结构体成员方法（如 Java Method 或 Go Receiver Method）、内部私有嵌套类（Inner Class）中包含的所有实现方法。
2. **排除在外的元素**：没有实际函数体的抽象方法声明、接口方法定义（Interface Methods），以及内部不包含任何控制流分支的平凡属性读写器（标准的 Getter/Setter 方法，以防 Java Bean 属性过多造成指标虚高）。

---

## 4. 典型误判场景与约束说明 (Edge Cases)

为防止特定设计模式导致的指标虚高或漏报，静态分析引擎必须应用以下约束：

* **约束 A（大路由表与工厂模式弱化）**：在命令分发器（Dispatcher）或大工厂类（Factory）文件内，通常存在一个包含巨大 `switch-case` 结构的路由函数。这会导致该函数的圈复杂度和文件的 FCI 瞬间破表。引擎应支持通过配置文件或特定注解（如 `@Factory`）为此类特殊职责的文件放宽 FCI 的触发阈值。
* **约束 B（死代码与无用函数拦截）**：若文件中存在被注释掉的函数、或完全无法被外部触达的死代码函数，只要它们依然保留在源文件中并能被解析为有效 AST 函数节点，引擎默认会对其进行复杂度累加。为了提高精准度，建议在 FCI 扫描前配置未引用代码清理插件。
* **约束 C（Lambda 表达式与闭包处理）**：在部分语言中，函数内部可能会嵌入匿名 Lambda 表达式。`arch-lens` 的标准算法是将函数体内定义的匿名 Lambda 或闭包结构所包含的控制流分支，一并算入其所属宿主函数的圈复杂度内，不将其作为独立的顶层函数去重复拆分计数。

---

## 5. 详细代码计算实例

以下是一段经典的 Go 语言业务代码，我们来分别计算它内部各个函数的 **圈复杂度 (CC)** 并最终推导出文件的 **FCI**：

```go
package pipeline // 行 1

import "fmt"     // 行 2
                 // 行 3 - 纯空行
// 函数 1：校验访问权限
func CheckAccess(role string) bool { // 行 4 - 函数 1 入口 (基础分 = 1)
	if role == "ADMIN" || role == "MANAGER" { // 行 5 - 命中 if (+1), 命中 || (+1)
		return true // 行 6
	} // 行 7
	return false // 行 8
} // 行 9

// 结构体声明（非函数节点，FCI 忽略）
type OrderProcessor struct { // 行 10
	Mode string // 行 11
} // 行 12

// 函数 2：核心处理分发
func (op *OrderProcessor) Process(status int) string { // 行 13 - 函数 2 入口 (基础分 = 1)
	switch status { // 行 14 - switch 头部不计分
	case 1: // 行 15 - 命中 case (+1)
		return "Pending" // 行 16
	case 2: // 行 17 - 命中 case (+1)
		if op.Mode == "ASYNC" { // 行 18 - 命中 if (+1)
			return "Async_Processing" // 行 19
		} // 行 20
		return "Sync_Processing" // 行 21
	default: // 行 22 - default 不计分
		return "Unknown" // 行 23
	} // 行 24
} // 行 25

```

### 计算结果对撞表：

| 行号 | 文本内容分类 | 所属函数域 | 触发的 AST 分支算子 | 圈复杂度 (CC) 贡献值 | 说明 |
| --- | --- | --- | --- | --- | --- |
| **1 ~ 3** | Package / Import / 空行 | 全局域 | 无 | 0 | 非函数体内逻辑，直接忽略 |
| **4** | `func CheckAccess(...)` | CheckAccess | 函数入口基础分 | **+1** | 初始化函数默认主路径分值 |
| **5** | `if ... role == "ADMIN" || ...` | CheckAccess | `IfStatement` & `BinaryExpr(||)` | **+2** | 包含 1 个 if 判定与 1 个短路或运算符 |
| **6 ~ 8** | 赋值与返回语句 | CheckAccess | 无 | 0 | 顺序执行流，不增加分支 |
| **9** | 函数结尾右括号 `}` | CheckAccess | 无 | 0 | 函数边界闭合 |
| **10~12** | 结构体定义声明 | 全局域 | 无 | 0 | 静态类型定义，非函数节点 |
| **13** | `func (op *OrderProcessor)...` | Process | 函数入口基础分 | **+1** | 初始化函数默认主路径分值 |
| **14** | `switch status {` | Process | `SwitchStatement` 头部 | 0 | 整个 switch 头部不计分 |
| **15** | `case 1:` | Process | `CaseClause` | **+1** | 条件分叉路径增加 |
| **16** | 顺序返回语句 | Process | 无 | 0 | 顺序流 |
| **17** | `case 2:` | Process | `CaseClause` | **+1** | 条件分叉路径增加 |
| **18** | `if op.Mode == "ASYNC" {` | Process | `IfStatement` | **+1** | case 2 嵌套内部的条件分支路径增加 |
| **19~21** | 嵌套返回与顺序流 | Process | 无 | 0 | 顺序流 |
| **22~25** | default / 结尾花括号 | Process | 无 | 0 | 默认分支不计分，函数边界闭合 |

* **单个函数圈复杂度公式推导**：
* $CC(CheckAccess) = 1 (\text{基础分}) + 1 (\text{行 5 的 if}) + 1 (\text{行 5 的 } ||) = 3$
* $CC(Process) = 1 (\text{基础分}) + 1 (\text{行 15 的 case 1}) + 1 (\text{行 17 的 case 2}) + 1 (\text{行 18 的 if}) = 4$


* **最终文件复杂度指数 (FCI)** = $\sum_{i=1}^{2} CC_i = CC(CheckAccess) + CC(Process) = 3 + 4 =$ **7**

---

**ArchLens 度量指标规范 - #03 文件复杂度指数 (FCI)**

---
