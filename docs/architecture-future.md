# Future Architecture Blueprint

## 1. 架构愿景

从“单体可用系统”演进到“可解释推荐 + 事件驱动 + 多租户治理”的 Venture Infra。

目标能力：
- 支持多城市、多机构、多策略并行运行
- 推荐可解释且可审计
- 实验可控，策略可灰度，回滚可分钟级

## 2. 分层架构（目标态）

### L1 接入层
- API Gateway（鉴权、限流、租户路由、审计）
- BFF（Founder 端 / Investor 端）

### L2 业务域层
- Identity Domain（用户与权限）
- Matching Domain（召回、排序、解释）
- Pipeline Domain（线索推进与状态机）
- Content/Policy Domain（政策与资讯）
- Org Domain（机构协作与团队视图）

### L3 智能层
- Feature Store（在线/离线特征）
- Retrieval Service（规则召回 + 向量召回）
- Rank Service（可解释重排）
- Experiment Service（A/B 与策略版本控制）

### L4 数据层
- OLTP：PostgreSQL
- Vector：Qdrant
- Event Store：Scylla / Kafka（按阶段演进）
- Lakehouse（后续）：用于离线训练与归因分析

## 3. 事件驱动数据流

### 关键事件
- `profile_completed`
- `match_exposed`
- `match_clicked`
- `contact_initiated`
- `contact_replied`
- `meeting_scheduled`
- `feedback_submitted`

### 流程
1. 事件进入总线
2. 实时聚合生成在线特征
3. 批处理生成离线特征与训练样本
4. 新策略上线并灰度
5. 观察指标，保留或回滚

## 4. 推荐系统可解释架构

### 4.1 三段式决策
- Recall：大范围召回
- Rank：多目标重排
- Reasoner：解释生成与证据绑定

### 4.2 解释输出标准
每条推荐至少包含：
- 匹配分
- 关键匹配因子 Top3
- 潜在风险因子 Top2
- 建议动作（例如：先补充财务里程碑再联系）

## 5. 多租户与隔离

- 逻辑隔离：租户级数据权限、配额、策略空间
- 计算隔离：关键队列与任务按租户配权
- 观测隔离：租户级 SLA、错误率、转化漏斗

## 6. 安全与合规架构

- 全链路 `request_id` + 审计日志
- 敏感字段分级访问（默认最小权限）
- 模型策略版本留痕（满足复盘与合规审查）
- 明确金融边界：平台提供撮合与信息服务，不提供未许可的投资建议承诺

## 7. 演进策略

### Stage A（0-3 个月）
- 夯实当前单体，补齐指标与反馈
- 建立基础实验框架

### Stage B（3-9 个月）
- 拆分 Matching/Feature/Experiment 关键能力
- 引入事件总线，形成在线/离线闭环

### Stage C（9-18 个月）
- 多租户平台化
- 跨区域部署与策略自治
- 智能解释系统标准化
