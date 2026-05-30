
---

# 架构度量指标规范：紧密类内聚度 (Tight Class Cohesion - TCC)

## 1. 指标定义

**紧密类内聚度 (Tight Class Cohesion, TCC)** 是指在一个特定类（Class）内部，**直接发生关联（Tight Coupled）的方法对（Method Pairs）的总数，占该类内部所有可能的方法对总数的比例**。

两个方法发生“直接关联”的判定标准是：**它们至少共同访问了该类中的同一个成员变量（Field/Attribute）**。

在 `arch-lens` 静态分析引擎对“上帝类 (God Class)”的判定算法中，TCC 是衡量内聚性丧失的关键分母。当一个类的 $TCC < 0.33$（即只有不到三分之一的方法间存在直接数据关联）且体积超标时，即可判定该类发生了严重的职责蔓延，沦为了纯物理堆砌的“上帝类”。

---

## 2. 核心价值与局限性

* **价值**：TCC 能够从语义和数据流拓扑层面精准识别出“伪内聚”。它可以轻易看穿那些虽然写在同一个类里，但各写各的数据、互不干涉的孤立逻辑块，为架构师重构、拆分类（Extract Class）提供强有力的量化数学依据。
* **局限性**：TCC 仅统计**直接通过变量关联**的方法对（紧密耦合），而忽略了通过“方法 A 调用方法 B，方法 B 访问变量 X”这种间接调用链建立的内聚关联（那是松散类内聚度 LCC 的度量范畴）。因此，针对重度依赖链式逻辑而少直接访问变量的特化类，TCC 可能会表现出非预期的低值。

---

## 3. 计算方式与规则细节

设一个类中所有参与计算的有效方法总数为 $N$。则该类内部所有可能组合的方法对（Method Pairs）的最大可能总数 $NP$（即组合数 $C_N^2$）的计算公式为：

$$NP = \frac{N \times (N - 1)}{2}$$

设该类内部实际被判定为“紧密关联（共享至少一个成员变量）”的方法对总数为 $NDC$。则该类的 TCC 计算公式为：

$$TCC = \frac{NDC}{NP}$$

### 3.1 导致方法间建立紧密关联 (Tight Connection) 的判定条件

1. **共享成员变量**：方法 $M_A$ 和方法 $M_B$ 在其各自的控制流和赋值语句中，显式读取（Read）或写入（Write）了同一个当前类的成员变量 $F_1$。
2. **互访不传递**：若方法 $M_A$ 与 $M_B$ 共享变量 $F_1$，方法 $M_B$ 与 $M_C$ 共享变量 $F_2$，在 TCC 计算中，仅认定 $(M_A, M_B)$ 和 $(M_B, M_C)$ 为直接关联对，$M_A$ 与 $M_C$ 如果没有共同访问任何变量，则 $(M_A, M_C)$ **不计入** $NDC$。

### 3.2 计入与排除的方法/变量边界规则

1. **参与计算的方法集合**：必须是当前类内部实现的、包含有效方法体的公共方法（Public Methods）和私有业务方法。
2. **严格排除的方法**：
* **构造函数（Constructors）**：因其天然初始化所有变量，计入会导致 TCC 虚高，必须排除。
* **平凡属性读写器（标准的 Getter/Setter 方法）**：由于其逻辑过于单一，计入会严重污染方法对的拓扑网络，必须排除。


3. **严格排除的变量**：静态全局常量（`static final` 或 Go 的全局 `const`），因为任何方法都可能读取它，无法代表对象实例的状态内聚。

---

## 4. 典型误判场景与约束说明 (Edge Cases)

为防止特定代码风格或框架机制导致的指标失真，静态分析引擎必须应用以下约束：

* **约束 A（极简类的分母防御）**：若一个类非常小，内部只有 $N = 1$ 个有效方法，则 $NP = 0$。此时公式分母为 0 无法计算。**`arch-lens` 规范要求：当 $N \le 1$ 时，TCC 指标默认直接输出 1.0（代表完全内聚）**。
* **约束 B（委派模式与间接调用链）**：若类内部的方法 $M_1$ 自身没有访问任何成员变量，但它内部第一行就直接调用了方法 $M_2$，而 $M_2$ 访问了变量 $F_1$。此时在标准的 TCC 算子中，$M_1$ 无法与任何方法建立紧密连接。若系统中此类委派方法较多，建议联合参考 LCC（Loose Class Cohesion）指标。
* **约束 C（只写不读的垃圾字段）**：某些遗留代码中可能存在一些只被方法写入却从未被系统真正读取的“死变量”，它们会成为维持 TCC 走高的假象。引擎在分析时，可在预处理阶段对未使用的字段进行标记或剔除。

---

## 5. 详细代码计算实例

以下是一段经典的 Java 业务代码（`OrderService`），我们来对其内部的成员变量访问关系进行静态拓扑展开，并精密推导其 **TCC** 值：

```java
package com.archlens.demo;

import java.util.List;

public class OrderService {
    // 参与计算的成员变量
    private List<String> orderCache; 
    private double taxRate;         

    // 排除项：构造函数（不参与 TCC 计算）
    public OrderService() {} 

    // 方法 1：加载缓存
    public void loadCache() {
        this.orderCache.add("Order_01"); // 访问变量: orderCache
    }

    // 方法 2：清除缓存
    public void clearCache() {
        this.orderCache.clear();         // 访问变量: orderCache
    }

    // 方法 3：调整税率
    public void updateTax(double rate) {
        this.taxRate = rate;             // 访问变量: taxRate
    }

    // 方法 4：计算总价（孤立业务，未访问任何内部成员变量）
    public double calculateTotal(double price) {
        return price * 1.12;             // 访问变量: 无
    }
}

```

### 5.1 基础元数据提取表

| 方法名称 | 是否参与计算 | 实际访问的当前类成员变量集合 | 说明 |
| --- | --- | --- | --- |
| `OrderService()` | **No** | `{orderCache, taxRate}` | 构造函数，严格排除 |
| `loadCache` | **Yes** | `{"orderCache"}` | 有效业务方法 1 |
| `clearCache` | **Yes** | `{"orderCache"}` | 有效业务方法 2 |
| `updateTax` | **Yes** | `{"taxRate"}` | 有效业务方法 3 |
| `calculateTotal` | **Yes** | `{}` (空集) | 有效业务方法 4（虽然没用字段，但也算有效方法） |

由此得到，参与计算的有效方法总数 **$N = 4$**（分别是 `loadCache`, `clearCache`, `updateTax`, `calculateTotal`）。

---

### 5.2 方法对 (Method Pairs) 全矩阵拓扑对撞表

根据 $N=4$，全系统一共有 $NP = \frac{4 \times 3}{2} = 6$ 个可能的方法对。我们逐一对撞求出实际关联值：

| 方法对组合 | 方法 A 变量集 | 方法 B 变量集 | 交集是否非空 | 是否属于紧密关联对 (NDC) | 共享的变量 |
| --- | --- | --- | --- | --- | --- |
| **(loadCache, clearCache)** | `{"orderCache"}` | `{"orderCache"}` | **{"orderCache"}** | **Yes** | `orderCache` |
| **(loadCache, updateTax)** | `{"orderCache"}` | `{"taxRate"}` | `{\}` (空) | **No** | 无 |
| **(loadCache, calculateTotal)** | `{"orderCache"}` | `{\}` (空) | `{\}` (空) | **No** | 无 |
| **(clearCache, updateTax)** | `{"orderCache"}` | `{"taxRate"}` | `{\}` (空) | **No** | 无 |
| **(clearCache, calculateTotal)** | `{"orderCache"}` | `{\}` (空) | `{\}` (空) | **No** | 无 |
| **(updateTax, calculateTotal)** | `{"taxRate"}` | `{\}` (空) | `{\}` (空) | **No** | 无 |

由此得到，实际发生紧密关联的方法对总数 **$NDC = 1$**（仅有 `(loadCache, clearCache)` 这一对）。

---

### 5.3 TCC 最终公式落地

将上述拓扑元数据带入 `arch-lens` 的标准 TCC 判定公式中：

* 最大可能的方法对总数：

$$NP = \frac{4 \times (4 - 1)}{2} = 6$$


* 实际紧密关联的方法对总数：

$$NDC = 1$$


* 最终紧密类内聚度：

$$TCC = \frac{NDC}{NP} = \frac{1}{6} \approx 0.167$$


* **分析的有效方法总数 ($N$)**：4 个
* **最大可能方法对总数 ($NP$)**：6 对
* **实际紧密关联方法对总数 ($NDC$)**：1 对
* **最终紧密类内聚度 ($TCC$)**：**0.167** （远低于 0.33 阈值，内聚性极差，坏味道明显）

---

**ArchLens 度量指标规范 - #04 紧密类内聚度 (TCC)**

---
