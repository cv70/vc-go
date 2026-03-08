# API Contract

统一契约真相源：
- `openapi/openapi.yaml`

## 1. Source of Truth

- 主契约文件：`openapi/openapi.yaml`
- 如有冲突，以 `openapi/openapi.yaml` 为准

## 2. 契约治理规则

- 版本：OpenAPI `3.1.0`
- 路径前缀：`/api/v1`
- 鉴权：`Bearer JWT`
- 响应信封：`code/message/request_id/data`
- 时间格式：ISO-8601 UTC
- 分页：Cursor 优先

## 3. 兼容性策略

- 兼容增强：新增字段，不删除旧字段
- 非兼容变更：升级主版本路径（如 `/api/v2`）
- 错误码语义稳定，不复用

## 4. 变更流程

1. 先改 `openapi/openapi.yaml`
2. 本地执行 lint
3. 再更新后端实现与测试
4. 最后更新 SDK 与调用方

## 5. 本地检查命令

```bash
npx --yes @stoplight/spectral-cli lint openapi/openapi.yaml -r .spectral.yaml
```

## 6. 前端类型与 Client 生成

在 `frontend/` 目录执行：

```bash
npm run api:gen
```

对应脚本：
- `api:types`：生成 `src/api/types.ts`
- `api:client`：生成 `src/api/client/`

## 7. CI 检查

仓库已配置 GitHub Actions：
- `.github/workflows/openapi-lint.yml`

触发条件：
- PR 修改 `openapi/**` 或 `.spectral.yaml`
- push 到 `main/master` 且涉及 OpenAPI 文件
