
---

# 架构度量指标规范：声明实体数 (Number of Declared Entities - NDE)

## 1. 指标定义

**声明实体数 (Number of Declared Entities, NDE)** 是指在一个特定源文件（File）的物理边界内，所**显式定义的顶层元素（Top-level Elements）的绝对数量总和**。这些顶层元素包括但不限于：类（Classes）、接口（Interfaces）、结构体（Structs）、全局函数/独立函数（Functions）以及全局变量（Global Variables）。

在 `arch-lens` 静态分析引擎中，NDE 是用于判定“上帝文件 (God File)”的决定性指标之一。当一个文件的 $NDE$ 超过设定阈值（例如 $NDE > 10$）时，说明该文件缺乏基本的物理分割，开发人员在不断向同一个文件追加不相关的顶层声明，导致该物理文件演变成了一个边界模糊、维护风险极高的“代码垃圾场”。

---

## 2. 核心价值与局限性

* **价值**：NDE 能从物理和组织视角精准暴露出系统的“模块封装失败”。在提倡单文件单职责的现代语言规范中，高 NDE 值直接对应着极高的物理耦合度与混乱的维护边界。
* **局限性**：NDE 仅关注顶层声明的**数量（广度）**，而不关注单个声明内部的**复杂度（深度）**。例如，一个文件定义了 15 个极其简单的全局常量（NDE = 15），和定义了 2 个上千行的核心业务类（NDE = 2），前者的 NDE 远高于后者，但后者的重构优先级通常更高。因此，NDE 必须配合代码行数（LOC）和文件耦合度共同使用。

---

## 3. 计算方式与规则细节

在 `arch-lens` 对源文件进行 AST 树状解析时，算法仅遍历**根节点（Root Node）下的第一层直接子节点（Direct Children）**。通过对这些子节点的类型进行识别，累加计算 NDE 计数。

### 3.1 严格计入 NDE 的元素（顶层声明）

1. **顶层类型声明**：显式定义的顶层类（Class）、接口（Interface）、结构体（Struct）、枚举（Enum）。
2. **顶层函数/方法声明**：非内嵌于任何类或结构体内部的独立函数（例如 Go 中的 `func Foo()`，或 C++ 中的全局函数）。
3. **顶层全局变量/常量**：在文件全局域声明的变量（Global Variable）或常量（Constant Block）。

### 3.2 严格排除在 NDE 之外的元素（非顶层/辅助元素）

1. **类内部的成员**：类或结构体内部包含的成员变量、属性、局部方法（Methods）或构造函数。这些属于内聚成员，**绝不计入**文件的 NDE。
2. **局部变量与匿名函数**：在函数体或方法体内部声明的局部变量、局部常量、局部内嵌函数或 Lambda 表达式。
3. **语言级基础设施头部**：包声明语句（如 `package main`）、模块导入语句（如 `import` 或 `#include`）。

---

## 4. 典型误判场景与约束说明 (Edge Cases)

为了让 `arch-lens` 在实际的多语言项目（Go/Java/C++）审计中不发生误报，引擎在统计 NDE 时必须遵循以下约束逻辑：

* **约束 A（Go 语言的常量/变量块 `const/var`）**：
  在 Go 语言中，开发者习惯使用括号将一组常量或变量组合在一起：
```go
const (
    StatusPending = 1
    StatusActive  = 2
)

```


**`arch-lens` 审计准则**：在 AST 解析中，若这一组变量属于同一个 `ValueSpec` 顶层节点，则对 NDE 的贡献值**仅计为 1**；若属于多个分散的顶层生命周期，则按实际独立节点计数。
* **约束 B（内部类 / 嵌套类型）**：
  在 Java 中，类内部可能包含有 `private class InnerHelper`（嵌套内部类）。由于该类属于主类的子节点，而非文件的第一层顶层节点，因此**不计入文件的 NDE**。
* **约束 C（现代语言的单文件规范差异）**：
* **Java**：由于语法强制要求（通常一个文件只能有一个 public 类），正常项目的 NDE 普遍 $\le 2$。若出现 NDE 极高的情况，通常是因为开发者在同一个文件尾部堆砌了大量的包级私有类（Package-private Classes）。
* **Go / C++**：语言天然允许单文件内存在大量的全局函数和变量，因此针对 Go/C++ 的上帝文件审计，NDE 的判定阈值应当比 Java 适当放宽 2~3 倍。



---

## 5. 详细代码计算实例

以下是一段经典的 Go 语言业务源文件（`account_hub.go`），我们以此来演示 `arch-lens` 静态分析引擎是如何对其顶层 AST 节点进行精确的 NDE 计算的：

```go
package account // 排除：包声明头部

import (
	"fmt" // 排除：导入语句
	"time"
)

// 顶层声明 1：全局配置变量 (NDE +1)
var GlobalTimeout = 5 * time.Second

// 顶层声明 2：常量组合块 (作为一个整体顶层节点，NDE +1)
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// 顶层声明 3：结构体 Account (NDE +1)
type Account struct {
	ID       string    // 结构体内部成员，排除
	Username string    // 结构体内部成员，排除
	CreatedAt time.Time
}

// 结构体绑定的方法：属于 Account 内部成员，非顶层独立元素，排除！
func (a *Account) IsValid() bool {
	return a.ID != ""
}

// 顶层声明 4：独立接口声明 (NDE +1)
type TokenVerifier interface {
	Verify(token string) bool
}

// 顶层声明 5：顶层全局独立函数 (NDE +1)
func GenerateID() string {
	prefix := "ACC_" // 函数内部局部变量，排除
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
}

```

### AST 顶层节点扫描与 NDE 计数对撞表：

| AST 第一层子节点 | 节点语法类型 (Node Type) | 是否计入 NDE | 判定原理解析 |
| --- | --- | --- | --- |
| `package account` | `PackageClause` | **No** | 属于语言级物理声明头部，过滤 |
| `import (...)` | `GenDecl (IMPORT)` | **No** | 属于依赖导入声明头部，过滤 |
| `var GlobalTimeout = ...` | `GenDecl (VAR)` | **Yes (1)** | 属于物理文件顶层的独立全局变量定义 |
| `const ( RoleAdmin... )` | `GenDecl (CONST)` | **Yes (2)** | 顶层组合常量块，作为单个顶层定义节点计入 |
| `type Account struct` | `GenDecl (TYPE)` | **Yes (3)** | 属于物理文件顶层的独立结构体类型定义 |
| `func (a *Account) IsValid` | `FuncDecl (Method)` | **No** | **关键点**：这是接收者方法，在语义上挂载于 Account 内部，非独立顶层元素 |
| `type TokenVerifier interface` | `GenDecl (TYPE)` | **Yes (4)** | 属于物理文件顶层的独立接口类型定义 |
| `func GenerateID()` | `FuncDecl (Function)` | **Yes (5)** | 属于物理文件顶层的无接收者全局独立函数 |

* **最终原始代码行数 (Raw LOC)** = 36 行
* **最终声明实体数 (Number of Declared Entities - NDE)** = **5**

---

**ArchLens 度量指标规范 - #02 声明实体数 (NDE)**

---