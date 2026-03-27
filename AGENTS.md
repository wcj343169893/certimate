# Agents 使用指南（为未来 AI 助手准备）

本文档为本仓库提供面向 AI/agent 的快速导览、建议任务和示例提示（prompts），便于将来把自动化 agent 安全、高效地接入项目协作。

## 项目概览

- 语言/平台：Go 后端（模块化在 `internal/`、`repository/`、`rest/` 等），前端在 `ui/`（Vite + React/TypeScript）。
- 主要职责：域名/证书申请与部署、自動化部署接入多云厂商（见 `internal/deployer`、`internal/applicant`）。
- 关键文件：`main.go`（程序入口）、`go.mod`、`Dockerfile`、`Makefile`、`ui/`（前端源码与构建配置）。

## 关键位置与职责

- `main.go`：启动流程，入口点，检查路由与依赖注入点。
- `internal/`：核心业务逻辑，含 `applicant`（申请器）、`deployer`（部署器）、`domain`（域名与 acme 相关）。
- `repository/`：持久化抽象与实现（数据库交互）。
- `routes/`、`rest/`：HTTP 接口与路由实现。
- `ui/`：前端运维界面，包含构建脚本与依赖（package.json、vite.config.ts）。

## 环境与运行（供 agent 使用的命令）

- 安装依赖（后端）：`go mod download`
- 运行后端：`go run .` 或 `make run`（参见 `Makefile`）
- 构建镜像：`docker build -t certimate .`（参见 `Dockerfile`）
- 前端：进入 `ui/`，`npm install`，`npm run dev` 或 `npm run build`

（提示：agent 在执行这些命令前应询问并确认是否可以在当前环境运行外部命令或修改文件）

## 推荐 agent 职责与权限边界

- 阅读与代码导航：只读访问仓库代码，生成代码摘要、依赖图、模块边界。
- 代码修改：应在创建分支与 PR 的流程下进行，遵循变更说明模板与单元测试要求。
- 测试与构建：在隔离环境或 CI 中运行 `go test`、前端构建，避免在用户机器上直接破坏状态。
- 敏感信息：绝不可将 secrets、API keys 或数据库凭证写入提交或日志中。

## 示例任务与 Prompt 模板

1. 快速代码概览（summary）

Prompt 模板：
"请读取仓库并给出 3-5 行的高层次摘要：项目目标、主要模块、运行方式、潜在技术债务。优先关注 `main.go`、`internal/`、`ui/`。"

2. 添加/修改功能（变更草案）

Prompt 模板：
"我要在 `internal/deployer` 中添加一个新厂商适配器，生成变更计划（文件/函数/测试），并给出一个简短的 PR 描述与需要的单元测试用例。"

3. 代码审查（PR review）

Prompt 模板：
"审查下列 diff（或 PR 链接）：指出逻辑错误、安全隐患、未覆盖的边界条件，并建议 3 个改进点与必要的测试。"

4. 重构建议

Prompt 模板：
"本模块 `internal/applicant` 有重复代码，用于多个厂商的 HTTP 请求包装。给出 1-2 个重构方案，包含接口定义、影响范围、以及回退策略。"

## 自动化工作流建议（agent 模式）

- `Explorer`（只读）: 自动读取文件结构、生成模块依赖图、提供变更影响评估。
- `LocalDev`（协助开发）: 帮助运行本地命令、生成代码片段、修改文件并创建分支/PR（需要明确授权）。
- `CI Assistant`（CI 环境）: 在 CI 上运行测试、生成覆盖率报告、提交修复建议（只在 CI 环境有权限时启用）。

每个模式需明确定义允许执行的 shell 命令、文件修改范围、以及是否能创建/合并 PR。

## 提交 / PR Checklist（agent 应自动检查）

- 是否包含相关单元测试（`go test`）？
- 是否通过 `go vet`、`golangci-lint`（如项目使用）？
- 前端改动是否包含 `ui/` 的构建验证？
- 文档或 README 是否同步更新？

## 常见问题与注意事项

- 本项目对云厂商的密钥和账号信息非常敏感。任何 agent 变更流程需确保 secrets 不入 repo。
- 在 Windows 环境有时会遇到路径或权限差异（例如 `TESSDATA_PREFIX` 之类的环境变量问题），agent 在运行环境特定命令前应检测 OS 类型。

## 快速上手流程（给 AI 的步骤）

1. 读取 `main.go`、`internal/`、`ui/` 三个位置的文件列表与 README。
2. 生成模块依赖与启动路径（3-5 行摘要）。
3. 等待用户确认是否运行测试或进行代码修改。

---

如果你希望我为常见任务生成具体的 prompt 库（例如 PR 审查、重构、测试修复），我可以把这些模板追加到 `agents.md` 或单独生成 `prompts/` 文件夹中的模板文件。
