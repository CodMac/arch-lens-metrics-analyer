在“上帝类 (God Class)”的缺陷判定中，ATFD 是用来量化类对外部实体的侵入性（依恋情结 / Feature Envy）的核心指标。如果说 WMC 衡量的是自身的臃肿，TCC 衡量的是自身的散乱，那么 ATFD 衡量的就是对外部世界的“指手画脚”。

根据你 `arch-lens` 仓库中关于上帝类的架构审计协议，ATFD 的计算专注于**当前类中所有方法直接访问外部非标准库类的成员变量或 Getter/Setter 方法的去重总数**。

以下是严格按照统一格式为你输出的 **ATFD 指标规范文档**：

---

# 架构度量指标规范：访问外部数据数 (Access to Foreign Data - ATFD)

## 1. 指标定义

**访问外部数据数 (Access to Foreign Data, ATFD)** 是指在一个特定类（Class）的物理边界内，其内部包含的所有方法通过**直接字段访问或调用 Getter/Setter 方法的方式，所读取或改写的外部独立类（Foreign Classes）的去重属性（Attributes）总数**。

在 `arch-lens` 静态分析引擎对“上帝类 (God Class)”的判定算法中，ATFD 是衡量外部耦合度与封装破坏的核心分子。当一个类的 $ATFD > 5$（工业级高精度判定阈值）时，说明该类过度沉溺于操纵其他类的数据，严重违反了“迪米特法则（Law of Demeter）”与面向对象封装原则，属于典型的“依恋情结（Feature Envy）”坏味道。

---

## 2. 核心价值与局限性

* **价值**：ATFD 能够极其敏锐地捕捉到“数据管理器”或“贫血模型调度器”这类伪架构。它能直接揪出那些自身没有实质业务内聚、专门靠压榨和组装其他 POJO/DTO 类数据的巨型过程式类。
* **局限性**：ATFD 属于物理层面的外部属性去重计数。对于合法的“应用层门面组件（Facade）”或“业务流程编排器（Orchestrator）”，它们天然需要组装多方数据，其 ATFD 值往往不可避免地会超过阈值。因此，ATFD 必须排除基础设施并结合低内聚度（TCC）来联合确诊上帝类。

---

## 3. 计算方式与规则细节

在 `arch-lens` 引擎的语法分析与符号绑定阶段，算法遍历类中所有方法的代码体，识别所有的属性访问表达式（Property Access Expressions）和方法调用表达式（Method Call Expressions）。

### 3.1 严格计入 ATFD 的访问行为

1. **外部类的公开属性直接访问**：如 `otherObj.status`，计入外部数据访问。
2. **外部类的平凡属性读写器调用**：如 `otherObj.getStatus()` 或 `otherObj.setStatus(value)`。这类调用在语义上等同于直接操作数据，**计入 ATFD 算子池**。
3. **去重统计边界**：计数的最小单位是“外部类的去重属性（Field）”**。若方法 A 调用了 `User.getName()`，方法 B 也调用了 `User.getName()`，对当前类的 ATFD 贡献值**仅计为 1；若方法 A 同时调用了 `User.getName()` 和 `User.getAge()`，则贡献值**计为 2**。

### 3.2 严格排除在 ATFD 之外的元素

1. **基础设施与原生标准库**：对语言内置对象（如 Java 的 `java.lang.String`、`java.util.List`、Map 等）的属性或方法操作，**绝对不计入 ATFD**。
2. **外部类的行为方法调用（Behavioral Methods）**：若调用的外部类方法内部包含复杂的业务逻辑（非平凡的 Getter/Setter），则属于正常的控制流协同（Message Passing），不属于“访问外部静态数据”，不计入 ATFD。
3. **自身及继承成员**：当前类自身定义的成员字段，以及从父类继承下来的受保护/公开字段的访问，属于内部内聚，不计入 ATFD。

---

## 4. 典型误判场景与约束说明 (Edge Cases)

为了让静态分析引擎在面对现代复杂框架（如 Spring, MyBatis）时保持工业级的高精度，必须应用以下约束机制：

* **约束 A（链式调用的间接穿透识别）**：在形如 `order.getUser().getAddress().getCity()` 的非赋值链式调用中，引擎通过符号表和类型推导，**必须将 `User.address`、`Address.city` 这两个暴露出来的外部物理属性同时计入 ATFD 池中**。
* **约束 B（三方 DTO 与 自动生成 POJO 隔离）**：如果项目大量使用自动生成工具（如 MapStruct 转换器，或者纯 DTO 互转类），这些类中充满了大量的属性搬运（`dto.setX(entity.getX())`）。这类纯结构体映射类应当在引擎预处理阶段通过包名白名单（如 `*.dto.*`, `*.vo.*`）进行过滤，不转入上帝类判定流。
* **约束 C（流式 API / Builder 模式防御）**：在形如 `QueryWrapper.builder().eq("id", 1).like("name", "a").build()` 的流式调用中，虽然调用了大量的方法，但其本质是构建器模式的配置行为，并非蚕食外部类的数据资产。引擎应识别其返回类型是否为当前链的同一个 Builder 实例，避免误判。

---

## 5. 详细代码计算实例

以下是一段经典的 Java 业务代码（`OrderManager`），我们来对其内部所有方法操作外部数据的行为进行静态扫描，并精密推导其 **ATFD** 值：

```java
package com.archlens.demo;

import com.archlens.demo.model.Order; // 外部类 A
import com.archlens.demo.model.User;  // 外部类 B

public class OrderManager {
    private String managerId; // 自身属性

    // 核心业务方法：评估订单
    public void evaluateOrder(Order order, User user) {
        // 1. 访问当前类自身字段，排除
        System.out.println("Manager: " + this.managerId); 

        // 2. 访问外部类 User 的属性（通过标准 Getter） -> 命中属性 1: User.vipLevel (ATFD +1)
        if (user.getVipLevel() > 3) { 
            
            // 3. 访问外部类 Order 的属性 -> 命中属性 2: Order.price (ATFD +1)
            double discountPrice = order.getPrice() * 0.8;
            
            // 4. 改写外部类 Order 的属性 -> 命中属性 3: Order.finalPrice (ATFD +1)
            order.setFinalPrice(discountPrice); 
        }

        // 5. 重复访问外部类 User 的属性 -> 属性 1 已存在，去重，排除
        System.out.println("User log: " + user.getVipLevel());

        // 6. 调用外部类的复杂行为方法（非平凡 Getter） -> 属于业务协同，排除
        user.renovateProfileStatus(); 
    }
}

```

### 计算结果对撞表：

| 行号 | 调用的目标表达式/方法符号 | 符号归属的物理类 | 提取的底层属性资产 | 是否计入 ATFD 池 | 判定原理解析 |
| --- | --- | --- | --- | --- | --- |
| **11** | `this.managerId` | `OrderManager` | `managerId` | **No** | 访问的是当前类自身字段，属于内聚行为 |
| **14** | `user.getVipLevel()` | `User` | `vipLevel` | **Yes (1)** | 首次访问外部类 `User` 的属性，入池 |
| **17** | `order.getPrice()` | `Order` | `price` | **Yes (2)** | 首次读取外部类 `Order` 的属性，入池 |
| **20** | `order.setFinalPrice(...)` | `Order` | `finalPrice` | **Yes (3)** | 首次修改外部类 `Order` 的属性，入池 |
| **24** | `user.getVipLevel()` | `User` | `vipLevel` | **No** | 属性 `User.vipLevel` 之前已在池中，**严格去重** |
| **27** | `user.renovateProfileStatus()` | `User` | 无 | **No** | 该方法内部包含复杂业务逻辑，属于行为调用而非数据访问 |

### ATFD 最终公式落地

根据基本定义，将当前类中提取出的外部有效属性集合进行收拢：


$$\text{ForeignDataPool}(OrderManager) = \{\text{"User.vipLevel"}, \text{"Order.price"}, \text{"Order.finalPrice"}\}$$

* **最终涉及的外部物理类总数**：2 个 (`User`, `Order`)
* **最终去重后的外部属性访问总数 (ATFD)** = **3**

---

**ArchLens 度量指标规范 - #06 访问外部数据数 (ATFD)**

---
