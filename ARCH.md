# ARCH.md - Current Architecture Entry

本文件是“当前可运行架构”的入口。它回答三个问题：
- 系统现在怎么跑
- 关键模块如何分工
- 如何平滑演进到未来架构

## 1. 当前系统快照

- 后端：Rust (Edition 2024) + Axum + Tokio
- 数据层：PostgreSQL（结构化主数据）
- 向量层：Qdrant（语义检索与匹配召回）
- 高吞吐记录：ScyllaDB（事件日志/行为数据）
- 当前 API 前缀：`/api/v1`

## 2. 当前模块结构

- `domain/`：业务规则（founder/investor/financing/policy/news）
- `datasource/`：DAO 与数据访问封装
- `infra/`：DB、Qdrant、Scylla 连接与初始化
- `main.go`：路由、注入、启动

## 3. 当前业务闭环

1. 用户建立资料（Founder 或 Investor）
2. 系统进行检索/匹配召回
3. 返回候选对象与推荐结果
4. 用户行为沉淀为后续优化数据

## 4. 当前架构优势与限制

### 优势
- 单体分层清晰，适合快速迭代
- 领域边界明确，便于后续服务拆分
- 向量检索已具备，能支撑语义匹配

### 限制
- 匹配解释层与反馈学习仍需体系化
- 指标采集与实验框架尚未完全产品化
- 多租户与跨区域部署能力待完善

## 5. 演进入口

未来架构的目标状态、分层治理、事件流与 AI 能力栈，见：
- `docs/architecture-future.md`

产品系统与业务机制，见：
- `docs/product-os.md`
