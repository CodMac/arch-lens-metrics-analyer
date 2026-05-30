在“上帝类 (God Class)”的缺陷判定中，WMC 是用来衡量类内部逻辑总体积与总复杂度的决定性基础指标。如果说 TCC 衡量的是内聚性的丧失（分母），那么 WMC 衡量的就是复杂度的恶性膨胀（分子）。

根据你 arch-lens 仓库中关于上帝类的架构审计协议，WMC 的计算摒弃了泛化的“只数方法个数”的低精度做法，而是将类中所有方法的圈复杂度（Cyclomatic Complexity, CC）进行累加。

以下是严格按照统一格式为你输出的 WMC 指标规范文档：

---

# 架构度量指标规范：加权方法数 (Weighted Method Count - WMC)

## 1. 指标定义

**加权方法数 (Weighted Method Count, WMC)** 是指在一个特定类（Class / Struct）的物理边界内，其内部包含的所有有效成员方法的复杂度权重的绝对累加总和。

在 `arch-lens` 静态分析引擎中，方法的权重采用其真实的**圈复杂度 (Cyclomatic Complexity, CC)** 进行度量。其数学计算公式为：

$$WMC = \sum_{i=1}^{n} CC(m_i)$$

在对“上帝类 (God Class)”的判定算法中，WMC 是最核心的体积与复杂度触发阈值。当一个类的 $WMC \ge 47$（工业级标准阈值）时，说明该类内部承载了过多的业务判定分支，代码的认知负荷已经溢出，极度需要进行职责剥离或重构。

---

## 2. 核心价值与局限性

* **价值**：WMC 能够精准识别“隐蔽的巨型类”。如果一个类只有 3 个方法，但每个方法内部都包含极其复杂的深层 `if-else` 或 `switch-case` 判定（每个方法 CC = 20），传统的“方法计数”指标只有 3，会将其判定为健康；而 WMC 的值会飙升至 60，精准揪出这种由于局部逻辑爆炸导致的架构坏味道。
* **局限性**：WMC 关注的是类内部复杂度的**绝对值**。对于一些天然需要大量分支、但职责高度单一的特定模式类（例如纯路由转发器、大工厂类），WMC 也会轻易破表。因此，WMC 必须与低类内聚度（TCC）、高外部访问（ATFD）联合使用来确诊上帝类。

---

## 3. 计算方式与规则细节

在 `arch-lens` 的 AST 拓扑解析阶段，算法首先提取出当前类下的所有方法节点，基于控制流图（CFG）计算每个独立方法的圈复杂度 $CC(m_i)$（基础分为 1，每遇到一个控制流分支分叉节点则分值 +1），最后进行线性求和。

### 3.1 参与 WMC 累加的有效方法边界

1. **类成员方法**：Java 中的普通实例方法、静态方法（Static Methods），以及 Go 中的结构体接收者方法（Receiver Methods）。
2. **私有内嵌类方法**：若当前主类内部声明了私有嵌套类（Inner Class），嵌套类内部的方法复杂度**一并累加**到主类的 WMC 中。

### 3.2 严格排除在 WMC 之外的元素

1. **无主声明**：不包含实际函数体的抽象方法声明、接口方法定义（Interface Methods）。
2. **特殊辅助方法**：
* 内部无任何控制流分支的平凡属性读写器（标准的 Getter/Setter 方法）。
* 语言级脚手架重写（如 Java 中自动生成的 `toString()`、`hashCode()`、`equals()`），以防污染核心业务复杂度。



---

## 4. 典型误判场景与约束说明 (Edge Cases)

为防止开发框架特化或代码风格差异导致的指标失真，静态分析引擎必须应用以下约束：

* **约束 A（匿名内部类与 Lambda 闭包）**：在方法内部声明的匿名内部类或 Lambda 表达式，其内部如果包含控制流分支，**其复杂度直接算作其所属宿主方法的局部复杂度**，一同累加进 WMC，不作为独立方法分开统计。
* **约束 B（超大单体方法防御）**：若一个类的 WMC 极高（如 80），但检查发现其 90% 的复杂度都来自于同一个单体方法（该方法 CC = 72），这种情况的本质坏味道是“长方法 (Long Method)”，而非“上帝类 (God Class)”。重构方向应优先通过“提炼方法 (Extract Method)”平摊类复杂度。
* **约束 C（继承方法污染拦截）**：WMC **仅统计当前类中显式编写/重写的代码方法**。父类中定义且当前类直接继承、未做任何代码覆盖（Override）的方法，其复杂度绝不计入当前类的 WMC 池中。

---

## 5. 详细代码计算实例

以下是一段经典的 Java 业务代码（`PaymentProcessor`），我们来对其各个成员方法进行 AST 分支扫描，并精密推导其 **WMC** 值：

```java
package com.archlens.demo;

public class PaymentProcessor {
    private String env;

    // 排除项：标准的 Getter 方法，不参与 WMC 计算
    public String getEnv() { return this.env; } 

    // 方法 1：路由分发校验
    public boolean checkRoute(String channel) { // 入口基础分 = 1
        if (channel == null || channel.isEmpty()) { // 命中 if (+1), 命中 || (+1)
            return false;
        }
        return true;
    } // 方法 1 最终 CC = 3

    // 方法 2：执行扣款
    public void executePay(double amount, String type) { // 入口基础分 = 1
        switch (type) { // switch 头部不计分
            case "WECHAT": // 命中 case (+1)
                System.out.println("WeChat Pay");
                break;
            case "ALIPAY": // 命中 case (+1)
                if (amount > 1000) { // 命中嵌套 if (+1)
                    System.out.println("Large Alipay");
                }
                break;
            default: // default 不计分
                break;
        }
    } // 方法 2 最终 CC = 4
}

```

### 计算结果对撞表：

| 行号 | 文本内容分类 | 所属方法域 | 触发的 AST 分支算子 | 圈复杂度 (CC) 贡献值 | 说明 |
| --- | --- | --- | --- | --- | --- |
| **1~4** | Package / 字段定义 | 全局域 | 无 | 0 | 非方法内逻辑，直接忽略 |
| **6** | `getEnv()` 读写器 | getEnv | 属性 Getter 算子 | 0 | 平凡读写器，直接排除 |
| **9** | `func checkRoute(...)` | checkRoute | 方法入口基础分 | **+1** | 初始化方法默认主路径分值 |
| **10** | `if (channel == null || ...)` | checkRoute | `IfStatement` & `BinaryExpr(||)` | **+2** | 包含 1 个 if 判定与 1 个短路或运算符 |
| **11~14** | 顺序流与结尾花括号 | checkRoute | 无 | 0 | 边界闭合，checkRoute 累计 CC = 3 |
| **17** | `func executePay(...)` | executePay | 方法入口基础分 | **+1** | 初始化方法默认主路径分值 |
| **18** | `switch (type) {` | executePay | `SwitchStatement` 头部 | 0 | 整个 switch 头部不计分 |
| **19** | `case "WECHAT":` | executePay | `CaseClause` | **+1** | 分支路径增加 |
| **22** | `case "ALIPAY":` | executePay | `CaseClause` | **+1** | 分支路径增加 |
| **23** | `if (amount > 1000) {` | executePay | `IfStatement` | **+1** | case 内部嵌套的条件分支路径增加 |
| **26~29** | default / 结尾花括号 | executePay | 无 | 0 | 边界闭合，executePay 累计 CC = 4 |

* **单个方法圈复杂度公式推导**：
* $CC(checkRoute) = 1 (\text{基础分}) + 1 (\text{if}) + 1 (||) = 3$
* $CC(executePay) = 1 (\text{基础分}) + 1 (\text{case 1}) + 1 (\text{case 2}) + 1 (\text{嵌套 if}) = 4$


* **最终加权方法数 (WMC)** = $\sum_{i=1}^{2} CC_i = CC(checkRoute) + CC(executePay) = 3 + 4 =$ **7**

---

**ArchLens 度量指标规范 - #05 加权方法数 (WMC)**

---
