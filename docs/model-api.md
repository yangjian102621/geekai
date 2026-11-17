# ThinkBuddy 模型 API

本期提供 `GET /v1/models` 和 `POST /v1/chat/completions`（JSON、SSE 流式、工具调用透传），使用本站积分。

## 启用

1. 部署本次后端，启动时自动创建 `thinkbuddy_model_tokens` 表；每个用户一条记录，只保存 SHA-256 摘要 `token_hash`，不保存令牌明文。
2. 后台沿用现有聊天模型与上游密钥配置。上游可填写 Higress 根地址、`/v1` 地址或完整 `/v1/chat/completions` 地址；上游密钥只放后台，不发给客户端。
3. 部署前端及 Nginx 的 `/v1/` 代理配置。如果测试环境使用独立 Ingress，也需将 `/v1/` 转发到后端，关闭 SSE 缓冲并设置 300 秒超时。否则仍可能返回静态站点的 405。
4. 登录网站，在「个人空间 → 账号信息 → 模型 API」生成令牌。令牌前缀为 `sk-`，只显示一次；重置后旧令牌立即失效，停用会删除摘要。请勿提交真实令牌到代码仓库。

## 客户端配置

在 ThinkBuddy 桌面端现有模型连接中新增自定义 OpenAI 兼容服务：

- 地址：`https://thinkbuddy-t.hanyoai.com/v1`（本地为 `http://localhost:8888/v1`）。
- 协议：**OpenAI Chat Completions**。
- API Key：网站生成的 `sk-...` 令牌。
- 模型：通过模型发现获取，或填写后台模型的 `value`。

OAuth 登录平台返回的令牌与本站模型令牌是独立凭据，不会自动切换已有客户端连接。本期不提供 Responses、Claude 或 Gemini 原生协议。

## 扣费规则

每个 HTTP 请求按照后台模型的「消耗积分」收取一次费用，`0` 为免费，不按 token 计价；桌面端一个任务可能发起多次模型请求。

积分先在数据库中原子扣除，同时写入积分明细。上游失败、响应无效、流式输出未产生回复时，退回本次积分并写退款记录。已有文本、思考或工具输出后中断，保留本次扣费。客户端提交的用户 ID 不参与扣费身份判断。

输出上限采用后台模型的 `max_tokens` 配置，客户端可指定更小值，每次请求只支持 `n=1`。请求体上限 16 MiB，单次调用最多 5 分钟。

极端情况下，进程在预扣后被强制终止，或退款时数据库不可用，可能需要根据积分记录和服务日志人工处理退款；本期不新增任务队列或独立结算系统。

## PowerShell 验证

```powershell
$thinkbuddyBase = 'https://thinkbuddy-t.hanyoai.com'
$thinkbuddyKey = Read-Host '输入本站模型令牌'
$thinkbuddyHeaders = @{ Authorization = "Bearer $thinkbuddyKey" }
Invoke-RestMethod "$thinkbuddyBase/v1/models" -Headers $thinkbuddyHeaders

$thinkbuddyBody = @{
  model = 'deepseek-ai/DeepSeek-V4.1-Flash-1017'
  messages = @(@{ role = 'user'; content = '你好' })
  stream = $false
} | ConvertTo-Json -Depth 10
Invoke-RestMethod "$thinkbuddyBase/v1/chat/completions" -Method Post `
  -Headers $thinkbuddyHeaders -ContentType 'application/json; charset=utf-8' `
  -Body ([System.Text.Encoding]::UTF8.GetBytes($thinkbuddyBody)) |
  ConvertTo-Json -Depth 20
```

验证后在网站积分明细确认扣费。流式测试将 `stream` 改为 `true`，使用支持 SSE 的客户端查看分段输出。

## 回滚

回退本功能的提交并重新部署前后端、恢复 Nginx 配置。新增令牌表可保留，不影响旧版运行；不需要改动现有用户、积分、模型或上游密钥表的数据。已发生的积分记录保留。
