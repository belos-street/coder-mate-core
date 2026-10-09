# code-mate-core 项目指南

> 本项目是将 TypeScript 版本的语法高亮引擎（coder-mate-core）复刻为 Go + WASM 实现。
> 让 LLM/Agent 快速理解项目结构和规范，高效协作开发。

---

## 项目概述

**目标**：将原 TypeScript 版本的语法高亮引擎复刻为 Go 实现，编译为 WASM 供前端调用。

**核心约束**：
- 输出确定性：相同输入代码必须产生相同的 HTML/Token 流
- Go 1.23+（见 go.mod），零外部依赖，仅用标准库
- 目录结构遵循 Go 标准项目布局

**当前状态**：Go 实现已完成（17 语言 / 10 主题 / WASM / 并发安全 / CI 门禁）。

> 注：TypeScript 参考实现（`ts/` 目录）已随提交 `896d07e` 移除，文中残留的
> `ts/lib/` 路径仅作历史参考；行为验证以仓库内测试与 golden 基准为准。

---

## 目录结构说明

```
code-mate-core/
├── .agents/                  # Agent 工作区（LLM 协作所需的规范和文档）
│   ├── skills/               # 编码规范和领域知识（Agent 开发时参考）
│   │   ├── belos-street/     # 个人编码风格和最佳实践
│   │   ├── brainstorming/    # 需求分析和设计讨论流程
│   │   ├── golang-best-practices/  # Go 开发规范（WASM/concurrency/testing）
│   │   ├── vibe-flow/        # Web/SaaS 全流程开发
│   │   └── writing-plans/    # 实现计划编写规范
│   └── docs/                 # 项目文档（需求、设计、计划）
│       ├── go-migration-requirements.md  # 需求索引文件（入口）
│       ├── specs/            # 详细规格文档（01-10，按 Phase 组织）
│       └── plans/            # 实现阶段计划
├── core/ language/ theme/ api.go  # Go 实现源码
├── cmd/wasm/                 # WASM 入口
└── agents.md                 # 本文件 - 项目入口指南
```

## 目录角色

### `.agents/skills/` — Agent 技能与规范

当 Agent 执行编码任务时，应参考此目录下的规范：

| 技能 | 用途 | 触发条件 |
|------|------|----------|
| `golang-best-practices` | Go 开发全规范 | 编写任何 `.go` 文件前必须加载 |
| `belos-street` | 个人编码风格 | 需要遵循特定代码风格时加载 |
| `brainstorming` | 需求分析与设计 | 讨论新功能/设计决策时使用 |
| `writing-plans` | 实现计划编写 | 编写多步骤实现计划时使用 |

**特别说明**：`golang-best-practices/references/` 下包含 37 个 Go 专题参考（WASM、并发、测试、内存管理等），开发时按需查阅。

### `.agents/docs/` — 项目文档

存放需求文档、设计文档、实现计划等。当前阶段的核心文档：

| 文档 | 说明 |
|------|------|
| [go-migration-requirements.md](.agents/docs/go-migration-requirements.md) | 需求索引文件（入口） |
| [specs/](.agents/docs/specs/) | 详细规格文档（01-10），按 Phase 组织 |
| [plans/](.agents/docs/plans/) | 实现阶段计划 |

### `ts/` — TypeScript 参考实现（已移除）

TypeScript 参考实现已随提交 `896d07e` 从仓库移除（原 [coder-mate-core-js](https://github.com/belos-street/coder-mate-core-js) 仍可在线访问）。行为一致性验证以仓库内测试（内联期望 token 流 / 确定性快照 / golden 基准）为准。

---

## 开发工作流

1. **理解需求** → 阅读 `.agents/docs/go-migration-requirements.md`
2. **撰写计划** → 使用 `writing-plans` skill 创建详细的实现计划
3. **编码实现** → 参考 `golang-best-practices` skill 与既有 Go 源码
4. **验证行为** → 内联期望 token 断言 / 确定性快照测试 / CI 四道门禁（gofmt/vet/test/race）

---

## 快速上手

**首次进入项目的 Agent 应执行**：
1. 阅读本文件（agents.md）了解项目全貌
2. 阅读 [需求索引](.agents/docs/go-migration-requirements.md)，按当前 Phase 加载对应 spec
3. 浏览 `core/` 与 `language/` 理解 FSM 分词与适配器架构
4. 加载 `golang-best-practices` skill 获取 Go 开发规范

---

*本文件是 LLM/Agent 协作的入口指南，保持简洁准确。*
