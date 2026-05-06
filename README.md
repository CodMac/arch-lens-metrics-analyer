# Arch-Lens-Metrics-Analyer

**Arch-Lens-Metrics-Analyer** 是 Arch-Lens 系列工具中的度量与缺陷检测引擎。它通过分析 `arch-lens-dep-analyer` 生成的依赖图数据（JSONL 格式），计算各项源码度量指标，并根据预设的架构缺陷协议判定系统中的设计缺陷。

## 核心特性

- **原子指标计算**：支持 ATFD (Access to Foreign Data)、TCC (Tight Class Cohesion) 等核心架构度量指标。
- **缺陷自动判定**：内置上帝类 (God Class)、循环依赖 (Circular Dependency) 等常见架构缺陷的识别规则。
- **解耦架构设计**：度量指标计算与缺陷判定规则完全分离，易于扩展新的指标或自定义判定逻辑。
- **高性能图分析**：基于内存图结构和 Tarjan 算法，快速处理大规模项目的依赖关系。

## 项目架构

项目采用分层架构设计，确保逻辑清晰且易于维护：

- `core/`: 核心领域模型，负责维护依赖图（Graph）结构。
- `loader/`: 数据加载层，负责解析依赖分析器输出的 JSONL 文件。
- `metrics/`: 原子指标层，实现独立的度量算法（如 ATFD, TCC, LCOM 等）。
- `detector/`: 缺陷探测层，根据架构缺陷协议组合各项指标进行最终判定。
- `cmd/`: 命令行入口，负责流程调度与结果输出。

## 快速开始

### 1. 环境要求
- Go 1.25+
- 已安装并运行 `arch-lens-dep-analyer` 获取依赖数据

### 2. 编译
```bash
go build -o arch-metrics ./cmd/main.go
```

### 3. 使用
运行度量分析需要指定依赖分析器输出的 `element.jsonl` 和 `relation.jsonl` 文件：

```bash
./arch-metrics -elem path/to/element.jsonl -rel path/to/relation.jsonl
```

## 检测指标与阈值

目前已支持的缺陷检测规则（参考 `doc/` 协议）：

| 缺陷类型 | 核心指标 | 判定阈值 |
| :--- | :--- | :--- |
| **上帝类 (God Class)** | ATFD, TCC, WMC | ATFD > 5 && TCC < 0.33 |
| **循环依赖 (Circular Dependency)** | SCC (强连通分量) | 存在环路 (Class 级别) |

## 路线图

- [ ] 支持 WMC (加权方法复杂度) 计算
- [ ] 增加 LCOM (Lack of Cohesion in Methods) 指标
- [ ] 支持上帝文件 (God File) 和功能依恋 (Feature Envy) 检测
- [ ] 导出检测报告为 JSON 或 HTML 格式

## 许可证
MIT License
