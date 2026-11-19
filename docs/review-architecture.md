# GeekAI 二开版：架构与工程治理评审

> 评审对象：`D:/studio/geekai`（基于开源 geekai 二开的智能体管理后台 + 营销站点）
> 评审方式：只读代码核验。所有结论标注 `文件:行号`；无法实证的单列「推断」。
> 分支：`feat-reconstruct` @ `568ff5cb`

---

## 一、本质判断（5 条）

**1. 这不是一个"AI 应用"，而是一条"代收代付 + 额度分销"的资金管道，外包了一层 UI。**
系统里真正持续发生的事只有三件：从上游（OpenAI/火山/可灵/即梦…）买算力、向用户收钱换成"积分"、把调用转出去。代码里没有任何 usage 计量（`ChatModel.Power` 是"每次调用固定积分"，token 数只被写进 `PowerLog.Remark` 的**中文字符串**，`api/handler/chat_handler.go:534`；`MaxContext` 字段全仓零引用）。所以**平台无法回答"这次调用花了多少上游成本"**——一个分销商不知道进货价，它的毛利率就是一个不可观测的随机量。

**2. 三大不变量（钱 / 额度 / 任务状态）全部靠"约定"活着，没有一个是"机制"。**
积分扣减本身写得不错（`api/service/user_service.go:67` 用 `WHERE power >= ?` 原子防负），但**每一个调用方都自己 `Begin()` 一个装饰性事务**：`tx := h.DB.Begin()` 之后调用 `IncreasePower`，而后者内部另开 `s.db.Transaction`。钱的那条腿永远不在调用方事务里。全仓 22 个增/扣积分调用点、5 段复制粘贴的异步退款循环、0 个幂等键（`PowerLog` 无唯一索引，`api/store/model/power_log.go:9`）。**"余额正确"这件事依赖 22 个 handler 作者都写对,而不是依赖数据库。**

**3. 权限不是"系统允许谁做什么"，而是"22 个路由注册文件各自的作者记不记得加一行 `.Use()`"。**
`api/handler/admin/` 下有 **22 处** `Engine.Group("/api/admin/...")`。gin 的 `Group()` 复制的是**引擎当时的** handler 链，`.Use()` 只作用于**当次**创建的 group。结果：15 个文件各自抄了一遍 `Use(AdminAuthMiddleware(...))`，**7 个文件漏了**（dashboard / jimeng / menu / order / powerLog / product / redeem）。"15 个文件重复写同一行"这件事本身就是漏防的自证。

**4. 仓库不是真理来源（source of truth），工作区才是。**
HEAD 的 `web/src` 是 **33 个 `.tsx/.ts`（React）**；工作区是 **97 个 `.vue/.js`（Vue）**，其中 **64 个从未 `git add`**。也就是说：`git clone && build` 得到的是一个**不存在于任何分支上的前端形态**；线上那套 Vue 前端在 Git 里不存在。而 `admin-web/` 已被整体删除合并进 `web/src/views/admin`。前端在近期经历了 **Vue → React → Vue 两次方向相反的完整重写**，最后一次还是未提交状态。

**5. 与上游的可合并性 = 0，且不是"冲突多"，是"没有基线可比"。**
`git remote -v` 只有自己的 fork（`zwbdzb/geekai`），**没有配置原作者 upstream remote**。`git diff origin/main...HEAD` = 545 文件 / +26,980 / −50,167，其中 `admin-web/` 178 文件全删、React 版 `web/src` 136 文件。结构性改动（删掉独立后台应用、把门户从 React 改回 Vue）会让任何后续"同步上游"变成手工重写——因为两侧已经不存在可三方合并的公共祖先语义。

---

## 二、分级问题清单

### P0（资损 / 数据泄露 / 不可维护，必须立即处理）

**[P0] 7 组管理接口零鉴权，其中兑换码接口可匿名无限造钱**
— 证据：`api/handler/admin/redeem_handler.go:33-38`（`Group("/api/admin/redeem/")` 后无 `.Use(...)`，紧接 `GET list` / `POST create` / `POST set` / `GET remove` / `POST export`）；同构漏防：`dashboard_handler.go:33`、`jimeng_handler.go:41`、`menu_handler.go:32`、`order_handler.go:34`、`power_log_handler.go:33`、`product_handler.go:34`。对照已补的：`api_key_handler.go:36-38`。造钱链路：`POST /api/admin/redeem/create`（`redeem_handler.go:152-193`，`Power` 直接取请求体，无上限校验）→ 得到 enabled 兑换码 → 登录任一账号 → `POST /api/redeem/verify`（`api/handler/redeem_handler.go:45,71`）→ 积分到账。
— 为什么是本质问题：鉴权被实现成"每个路由文件的一个可选装饰"，于是权限边界 = 22 份代码的作者记忆之和。这是结构性的，不是遗漏。
— 修复方向：在 `api/main.go` 一次性建立 `adminGroup := s.Engine.Group("/api/admin")` 并 `adminGroup.Use(AdminGuard)`，把 group 实例**注入**给所有 admin handler（改 `RegisterRoutes(group *gin.RouterGroup)` 签名），从类型上消灭"新文件忘记加"的可能；`/api/admin/redeem/create` 的 `Power` 加上限与 RBAC。

**[P0] 匿名 SSRF：`GET /api/download?url=` 无校验、无超时、回吐响应体**
— 证据：注册在 `api/handler/net_handler.go:49`（`h.App.Engine.GET`，故意在鉴权 group 之外），实现在 `net_handler.go:156-186`：不校验 scheme/host、无 allowlist、`client := &http.Client{}` 无 `Timeout`、默认跟随 302、`io.Copy(c.Writer, r.Body)` 原样输出。
— 后果：读云元数据（`169.254.169.254`）、探内网 MySQL/Redis、以及**从内网回打上一条里那些匿名 admin GET 接口**；无超时 → 连接与 goroutine 耗尽。
— 修复方向：改为只允许"本站 OSS 域名的路径签名下载"（服务端存储 `file` 表已有 `file.go` model），彻底不再由客户端传 URL；若必须保留代理，需 host allowlist + 禁止私网段 + 5s Timeout + 不跟随跳转。

**[P0] 支付回调可双倍到账（钱进了，积分加两次）**
— 证据：`api/handler/payment_handler.go:405` `paySuccess`：`h.lock.Lock()` 是**进程内** mutex（多副本即失效）；`:431` 先 `IncreasePower`（内部独立事务，`api/service/user_service.go:65`）**提交**，`:447` 才 `h.DB.Updates(&order)` 置已支付。两步之间失败 → 返回 error → `AlipayNotify:479-484` 回 `"fail"` → 支付宝重试 → order 仍是 `OrderNotPaid`（`:411` 只读判断，不是条件更新）→ **再充一次**。无 DB 级幂等：`PowerLog` 无唯一索引（`api/store/model/power_log.go:9-20`），订单状态更新没有 `WHERE status = unpaid`。
— 反例自查：回调失败返回 "fail" 让网关重试，这一步是对的（不是"钱进了不发货"）。问题恰好在"重试是安全的"这个隐含假设上——重试不安全。
— 修复方向：`UPDATE geekai_orders SET status=paid, trade_no=? WHERE order_no=? AND status=unpaid`，`RowsAffected==1` 才执行加积分；`geekai_power_logs` 加 `UNIQUE KEY (ref_type, ref_id)`。

**[P0] "假事务"：`tx.Begin()` 与积分操作不在同一事务，导致可刷积分**
— 证据 1（注册刷邀请）：`api/handler/user_handler.go:384` `tx := h.DB.Begin()`；`:399` `reg_num` 用 **`h.DB`**（不是 tx）更新；`:403` `h.userService.IncreasePower(...)` 内部用 `s.db.Transaction`（`user_service.go:31`）独立提交；`:412` 只有 `InviteLog` 走 `tx`。于是 `tx.Rollback()`（`:414`）回滚不了积分、回滚不了 `reg_num`。构造"注册后段失败"即可**重复领取邀请奖励**。
— 证据 2（兑换码无限核销）：`api/handler/redeem_handler.go:69` `tx := h.DB.Begin()`；`:71` `IncreasePower` 独立提交；`:87` `tx.Updates(&item)` 写核销标记。标记写失败 → 积分已到账、码仍未使用 → 同一个码反复兑。
— 反例自查：`api/handler/model_api_chat.go:90-100` 用 `defer` 在"未产生回复"时退款，是正确写法；`user_service.go` 内部事务本身也正确。**问题不在 service，在调用方把 `tx` 当装饰品。**
— 修复方向：`IncreasePower/DecreasePower` 增加接受 `tx *gorm.DB` 的重载（或改为 `func(tx *gorm.DB, ...)` 强制传入），从签名上禁止"自己开事务"；删除所有 handler 里裸的 `h.DB.Begin()`。

**[P0] 异步任务退款非原子，五份复制粘贴，多副本必双退**
— 证据（同一段逻辑抄了五遍）：`api/service/image/image_service.go:246-264`、`mj/service.go:268-281`、`suno/suno_service.go:404-418`、`video/video_service.go:315-329`、`jimeng/jimemg_service.go:525-538`。形态一致：`Find(status=failed AND power>0)` → `IncreasePower`（独立事务提交）→ `UpdateColumn("power", 0)`（另一条独立语句）。
— 后果 1：两步之间崩溃/DB 错误 → 下一轮循环再次退款，**永不收敛**。
— 后果 2：`IncreasePower` 出错时 `continue`（如用户已删除 → `RowsAffected==0` → 返回错误），job 永远停留在 `power>0`，循环每 5/10 秒重复扫描同一批 → 无限重试 + 全表条件扫描（`:247` 无 `limit`）。
— 后果 3：五个 goroutine 在 `api/main.go:183-219` **无条件启动、无分布式锁**，水平扩容即双退。
— 反例自查：video 的 adapter 抽象（`api/service/video/adapters/`）做得很好，**但退款逻辑没进 adapter**——说明抽象做在了"HTTP 形状"上，而真正共享的"任务生命周期 + 计费"没被抽象。
— 修复方向：单一 `TaskLedger.Settle(jobId, outcome)`：`UPDATE jobs SET power=0 WHERE id=? AND power>0` 抢占（`RowsAffected==1` 才退款）+ 退款与标记同事务 + `FOR UPDATE SKIP LOCKED` 取任务替代 `Find` 全扫。

**[P0] 管理员禁用 / 登出 / 重置密码三处全部失效，且无 RBAC**
— 禁用失效：`api/handler/admin/admin_handler.go:114` `redis.Set(..., "admin/{id}", token, 0)` —— **TTL=0 永久**；`api/core/middleware/auth.go:67-112` `AdminAuthMiddleware` 只校验 JWT 签名 + redis key 存在，**不查 `admin.status`、不查该 admin 是否还存在**；`Enable`（`admin_handler.go:279-300`）只改 DB，不删 redis key → **被禁用的管理员永久保有全部权限**。
— 登出失效：`Logout`（`admin_handler.go:135-141`）用 `h.GetUserKey(c)`，而 `GetUserKey`（`api/handler/base_handler.go:57-63`）读的是 `LoginUserID`，admin 路由只设 `AdminUserID`（`api/core/types/session.go:10-11`）；且该路由注册在 `group.Use()` **之前**（`admin_handler.go:57` vs `:60`）。两者叠加 → key 为空字符串 → **`redis.Del("")`，什么都没删**。
— 无 RBAC：`api/store/model/admin_user.go:7-17` **没有 role 字段**；`Save`（`admin_handler.go:195-223`）任何管理员都能新建管理员（平行越权/持久化）；而 `ResetPass`（`:264-300`）却做了超管判断 → 授权检查是**逐函数凭手感加的**。
— 顺带一个死锁 bug：`ResetPass:290` `utils.GenPassword(data.Password, user.Salt)` 用的是**操作者自己**的 salt，`:293` 却写到 `data.Id` 那一行 → **重置任何其他人的密码都会让该账号永久无法登录**。
— 修复方向：admin 会话改为每次从 DB 校验 `status` + `role`；`AdminGuard(perm string)` 按权限码声明式挂载；`Enable=false`/`Remove` 时同步 `redis.Del("admin/{id}")`。

**[P0] 管理员可无痕伪造任意用户身份；现有审计日志字段恒为 0**
— 证据：`api/handler/admin/user_handler.go:539-569` `GenLoginLink` —— 为任意 `user_id` 签发用户 JWT 并 `redis.Set("users/{id}", token, 0)`（永久、无 TTL），**不写任何审计记录**，任何管理员可调（无 RBAC），前端入口在 `web/src/views/admin/Users.vue:437`。
— 连带：积分强改的审计字段 `Remark: "...管理员ID：%d"` 用的是 `h.GetLoginUserId(c)`（`admin/user_handler.go:371`），在 admin 路由下恒为 **0** → **现有审计日志从上线第一天起就是废的**。
— 为什么是本质问题：`LoginUserID` 与 `AdminUserID` 是两套身份命名空间，但 handler 层没有把它们隔开，于是一切"谁干的"记录都静默退化成 0，而系统没有任何编译期/运行期机制能发现这件事。
— 修复方向：`GenLoginLink` 改为"短期一次性 ticket + 强制记录 impersonation 审计"；admin 侧统一 `h.AdminId(c)`，并在 CI 里禁止 admin 包内出现 `GetLoginUserId`。

**[P0] Stored XSS → 管理员 token 被窃 → 全平台接管（完整可走通的链路）**
— 环 1：`web/src/views/admin/records/ChatList.vue:179` `<div class="chat-line" v-html="dialogContent">` —— **管理后台直接渲染用户提交的对话内容**；同构 `Suno.vue:132`。全仓无 DOMPurify/sanitize（`web/package.json` 无此依赖）。
— 环 2：管理员 token 存 localStorage（`web/src/store/session.js:28-34`），且被挂到**每一个**请求上（`web/src/utils/http.js:47`）。
— 环 3：拿到 token 后 → 无鉴权的 `/api/admin/*`（P0-1）+ 永久不过期的会话（P0-6）+ `GenLoginLink`（P0-7）。
— 攻击路径：任意注册用户提交含 `<img src=x onerror=...>` 的对话 → 管理员打开"对话管理"→ token 外带 → 接管。
— 修复方向：admin 侧 v-html 一律走 sanitize（或干脆用 `<pre>{{}}</pre>`）；token 移出 localStorage（httpOnly Cookie 或独立 admin 子域）；管理后台与门户拆 origin。

**[P0] 仓库不可构建 / 上游不可合并（工程治理第一性问题）**
— 证据：`git ls-tree -r HEAD --name-only web/src` = 33 个文件，全是 `.tsx/.ts`（`api/client.ts`、`app/router.tsx`、`AppLayout.tsx`）；工作区 `web/src` = 97 个 `.vue/.js`，`git status --porcelain web/src` 有 59 条删除 + 大量未跟踪。即 **当前可运行的前端从未进入版本库**。
— 无上游基线：`git remote -v` 只有 `origin → github.com/zwbdzb/geekai.git`，**未配置原作者 remote**；`git diff origin/main...HEAD` = 545 files / +26,980 / −50,167，含 `admin-web/` 178 文件全删。
— 为什么是本质问题：fork 的可持续性 = "能算出与上游的差"。现在算不出来（没基线、没 vendor、目标代码未入库），所以**任何"同步上游"都是手工重写，且无法验证重写对不对**。叠加"前端 3 个月 Vue→React→Vue 两次逆向重写"，说明团队没有把"可回滚"当成不变量。
— 修复方向：见第三节第 3 件。

### P1（显著技术债，1-2 个迭代内）

**[P1] 定价与计量彻底脱钩：按"次"卖，按"token"买**
— 证据：`api/store/model/chat_model.go:16` `Power int smallint // 消耗积分点数`，`api/handler/chat_handler.go:528-541` `subUserPower` 固定扣 `ChatModel.Power`，token 数被格式化进 `Remark` 字符串（`:536`）。上游 `usage.prompt_tokens/completion_tokens` 其实**已经解析到了**（`chat_openai_handler.go:31-33`、`chat_handler.go:585,617`），只是被丢弃去写字符串。
— 为什么是本质问题：分销业务的利润 = 收入 − 上游成本；当收入与上游计费口径不同源，毛利不可测、套利不可检测。
— 修复方向：`geekai_chat_messages` 已有 `total_tokens` 列（`chat_message.go:17`）但只用于展示；`PowerLog` 增加 `prompt_tokens/completion_tokens/upstream_cost` 三列，定价模型改为 `power = ceil(in*t1 + out*t2)`。

**[P1] 入参长度无上限 → 对转售型业务的直接套利口**
— 证据：`MaxContext` 在**全仓没有任何引用**（`grep -rn MaxContext api/` 为空）；`api/handler/model_api_chat.go:21` `maxModelRequestBytes = 16 << 20`（16MB）；`:141-167` `limitChatTokens` 只约束 `max_tokens`（出参），不约束 messages 总长。→ 单次请求可塞数百万 token 入参，却按 1 次调用计费。
— 修复方向：接入 tokenizer 对 `messages` 预计算，超 `ChatModel.MaxContext` 直接 400；定价按上一条改计量。

**[P1] 启动期执行破坏性 DDL，且守卫常量是死代码**
— 证据：`api/service/migration_service.go:24-26` 定义了 `MigrationStatusKey`/`MigrationCompleted`，**全仓再未出现**（grep 仅命中定义行）→ 无任何守卫；`api/main.go:159-161` 每次启动 `StartMigrate()`；`:452-457` `TableMigration()` 顺序执行 `fullTableMigration` → `fixTableConstraints` → `incrementalTableMigration`，其中 `:324` `ALTER TABLE geekai_chat_roles DROP COLUMN marker`、`:395` `ALTER TABLE geekai_users DROP COLUMN chat_roles_json` 是**不可逆删列**；`api/store/mysql.go:84-95` 启动期还做 `ALTER TABLE ... RENAME`。
— 为什么是本质问题：schema 状态被记在 Redis（一个缓存），而不是 schema 自己（版本表）。Redis 冲掉 / 换实例 / 双副本同时启动 = 重跑破坏性 DDL。
— 修复方向：`geekai_schema_migrations(version, applied_at)` 表 + 每个迁移一个带版本号的可重入函数 + 启动时 `GET_LOCK('geekai_migrate')`；DDL 从进程启动路径移出，做成显式命令。

**[P1] 签到防重依赖本地 LevelDB，多副本/容器重建即失效**
— 证据：`api/handler/user_handler.go:772-783` `h.levelDB.Get(key)` → `h.levelDB.Put(key, true)`（`key = signin/{userId}/{date}`），`IncreasePower` 在 `:781`。LevelDB 目录固定在 `./data`（`api/store/leveldb.go:22` `leveldb.OpenFile("data", nil)`）—— 本地文件、单进程独占、容器重启即丢；`Get`+`Put` 非原子，双击即双领。
— 修复方向：Redis `SETNX signin/{uid}/{date} 1 EX 90000` 占据，成功才发放。

**[P1] "三层架构"不存在：handler 直连 GORM，service 是上游 HTTP 客户端，store 不是层**
— 证据：`api/store/` 只有 `mysql.go(117) / redis.go(31) / redis_queue.go(57) / leveldb.go(110)` 共 **315 行连接代码 + model + vo**，无 DAO/repository；`user_handler.go` 单文件 21 处 `h.DB.` 直接查询；`api/service/` 下除 `user_service.go(93 行)` 外，image/mj/suno/jimeng/video/sora 全是"调上游 + 轮询"的客户端。
— 反例自查（给分项）：依赖方向**确实是单向、无环**——`go list -deps` 验证 `service` 不 import `handler`，`store` 不 import `service`/`handler`。所以不是"分层被破坏"，而是"层根本没建起来，但至少没腐化成循环依赖"。
— 修复方向：不追求补 DAO。真正该抽的是"任务 + 账本"（见第二节 P0-5 与第三节第 1 件）。

**[P1] 上游供应商抽象做在了错误的层**
— 证据：`api/service/video/adapters/` 有 7 个 adapter + `map[string]VideoAdapter` 注册表（`video_service.go:34,56-84`），`VideoAdapter` 接口声明**在 `veo_adapter.go:23`**（不在 `adapter.go`，后者只有两个 struct）——抽象是在加第 2 个供应商时被动长出来的。但 mj / suno / jimeng / image / library **各自一套**：5 张 job 表（`image_job/mj_job/suno_job/video_job/jimeng_job`）、5 个轮询 goroutine、5 段退款复制品。
— "新增一个供应商"的真实成本：**video 类** = 1 个 adapter 文件（150-400 行）+ config 条目，最低；**非 video 类** = 新表 + model + user handler(400-600 行) + service(含轮询+退款) + admin handler + admin Vue 页 + `main.go` 两处 fx 注册 + 路由。二者相差一个数量级。
— 为什么是本质问题：`VideoAdapter` 抽象的是"HTTP 形状"（`CreateTask`/`QueryTask`），而真正跨供应商共享的不变量是"任务生命周期 + 计费 + 轮询"——恰恰没被抽象。
— 修复方向：把 `VideoAdapter` 降级为 `ProviderAdapter`（只管 HTTP），在其上加 `TaskRunner`（统一 job 表 + 统一轮询 + 统一 settle）。

**[P1] 硬依赖上游作者的商业服务，自营交付时是可用性 + 合规双风险**
— 证据：`api/core/types/geekai.go:13` `GeekAPIURL = "https://sapi.geekai.me"`，可被 `GEEK_API_URL` 覆盖（`:15-19`，说明作者留了口子，但没有自建服务端）。走它的功能：行为验证码（`captcha_service.go:39,56,74,91`）、微信扫码登录（`wxlogin_service.go:67`）、微博/早报（`function_handler.go:118,151`）；另外 `api/core/app_server.go:67-80` 每次启动把 `host.Info().HostID` 上报到该域名。
— 修复方向：验证码换自建（`user_mobile_login.go:116` 已有 Redis 限流基础），微信扫码走自己的公众号/开放平台，删除 install 上报。

**[P1] 前端单 axios 单例：admin token 被挂到每个请求**
— 证据：`web/src/utils/http.js:43-53` 请求拦截器无条件写入 `config.headers['Authorization'] = getUserToken()` **和** `config.headers['Admin-Authorization'] = getAdminToken()`；`:39` `axios.defaults.withCredentials = true` 全局生效。→ 门户站点的任何请求（含跨源）都携带 admin 凭证；admin 与 site 共享同一实例，无法按域收紧。
— 修复方向：两个 axios 实例（`userClient` / `adminClient`），admin 实例只在 `/api/admin/*` 前缀下注入 header，并迁移到 httpOnly cookie。

**[P1] 前端零路由守卫**
— 证据：`web/src/router.js:321-325` `router.beforeEach` **只设置 `document.title`**，无 token 检查、无 `/admin/login` 重定向。40+ 个 `/admin/*` 路由任何人可直接访问；后台安全完全押在"接口会返回 401"上——而接口本身没鉴权（P0-1）。
— 修复方向：`beforeEach` 校验 admin session（调 `/api/admin/session`）+ 权限码路由 meta，做双保险。

**[P1] 依赖与构建产物卫生**
— 证据：`marked@15` 与 `markdown-it@13` 并存（实际使用：`markdown-it` 4 文件 / `marked` 1 文件）；`vant@4.5` 仍被 2 个文件 import（mobile 视图已随 diff 删除）；`axios@0.27.2`（0.x 已 EOL）；`package-lock.json` + `pnpm-lock.yaml` + `pnpm-workspace.yaml` 三套包管理元数据；工作区存在 `web/dist`（33MB）、`web/node_modules`、`web/vite.config.js.timestamp-1794989911539-f59962d7091e48.mjs`（Vite 配置编译进程异常退出的残留）。
— 修复方向：定 pnpm 单一锁文件；删 vant / 二选一 markdown 渲染器；`dist`/`node_modules`/`*.timestamp-*.mjs` 入 `.gitignore` 并清理。

### P2（优化项）

- `TestHandler` 全空实现却注册了 `/api/test/sse`（`api/handler/test_handler.go:26-33,36-53`）。
- 状态变更走 GET：`admin/logout`（`admin_handler.go:57`）、`admin/list/remove`（`:63`）、`upload/remove`（`net_handler.go:47`）、`admin/redeem/remove`（`admin/redeem_handler.go:37`）。
- `DownloadHandler` 把 `updateDir` 整个目录公开静态服务（`api/handler/download_handler.go:63-64`）。
- 超大文件：`user_handler.go` 866、`chat_handler.go` 797、`migration_service.go` 834、`jimemg_service.go` 897、`web/src/store/data/video_params.js` 1251、`LoginDialog.vue` 1096。
- suno 的 `SyncTaskProgress` 在 `main.go:204-205` 被注释掉，但 `suno_service.go:325-420` 的循环仍在（若被别处再 invoke 会重复跑）。
- `api/service/video/adapters/adapter.go` 只有 2 个 struct，没有接口 —— 文件名与内容不符，建议把 `VideoAdapter` 移回此处。
- `genLoginLink` 类 GET 接口返回可登录凭证（`web/src/views/admin/Users.vue:437` httpGet）→ 会进浏览器历史/代理日志。

---

## 三、我验证过的 / 我推断但未验证的

### 已验证（有 file:line 实证，可复现）
1. 22 个 admin 路由组，7 个无鉴权中间件（逐个 `grep -n "\.Use(" api/handler/admin/*.go` 对照 `grep -n 'Group("/api/admin'`）。
2. `/api/admin/redeem/create` 无鉴权且 `Power` 取自请求体（`admin/redeem_handler.go:152-176`）。
3. `/api/download` 无鉴权、无 URL 校验、无超时（`net_handler.go:49,156-186`）。
4. 支付回调"先加积分、后改订单状态、失败则让网关重试"（`payment_handler.go:431,447,479-484`）。
5. 3 处"假事务"（`user_handler.go:384-414`、`redeem_handler.go:69-95`、`admin/jimeng_handler.go:131-159` 的 `tx` 同样未被 `IncreasePower` 使用）。
6. 5 段退款循环逐行同构（行号见 P0-5），且都无分布式锁、都在 `main.go:183-219` 无条件启动。
7. `admin/*` 中 `Enable` 不删 redis、`Logout` 读错上下文键、`GenLoginLink` 无审计、`ResetPass` 用错 salt（行号见 P0-6/P0-7）。
8. `MaxContext` 全仓零引用；`ChatModel.Power` 为固定值；token 数只进 `Remark` 字符串。
9. `MigrationStatusKey` 常量定义后从未使用（grep 只命中定义行）。
10. `store/` 无 DAO（仅 315 行连接代码）；依赖方向单向无环（`go list -deps` 验证）。
11. HEAD 的 `web/src` 是 `.tsx`、工作区是 `.vue`，64 个前端源文件未跟踪；无 upstream remote；545 文件 diff。
12. `v-html` 渲染用户内容且无 sanitize（`admin/records/ChatList.vue:179`、`Suno.vue:132` 等 6 处）。
13. `router.js:321-325` 无守卫；`http.js:46-47` 双 token 注入。
14. 桌面端**没有**污染主链路：后端只有 `download_handler.go:112-184` 的安装包分发，属于官网功能；未发现 `isDesktop` 之类分支。

### 推断但未验证（需要线上/环境确认后才能定性）
1. **多副本是否已在跑**：所有并发类结论（双退、双充、签到、`h.lock`）在**单副本**下大部分降级为"崩溃窗口"而非"必现"。代码里没有任何"只允许单实例"的约束（无 leader election、无文件锁），也没有文档说明，故按"随时可能扩容"评审。
2. **`smallint` 溢出行为**：`PowerLog.Amount`（`power_log.go:14`）与 `ChatModel.Power`（`chat_model.go:16`）是 `smallint`（上限 32767），与 `User.Power` 的 `int` 不一致。若 MySQL 非 strict mode 会静默截断 → 账本与余额不一致；strict mode 则写入报错。取决于部署 `sql_mode`，未实测。
3. **CORS 缺失的实际影响**：全仓无 CORS 中间件（`grep -rn "CORS" api/` 为空），但 `AGENTS.md` 声称"中间件涵盖 CORS"。生产走同源反代时不成问题；Electron（`file://`/`app://` 源）与小程序跨源直连时会失败。`desktop/index.js` 未发现 API 基址配置，未验证其实际调用方式。
4. **小程序形态**：`miniprogram/` 目录为空且被 `.gitignore` 忽略（`miniprogram`），第四种交付形态的代码不在本仓库。若在别处维护，则"同一后端服务四形态"的耦合面无法从本仓评估。
5. **`admin/records/ChatList.vue` 的 `dialogContent`**：我确认了模板级 `v-html`，但未逐行读它的 computed 是否做过转义。按"未发现 sanitize 调用"记为 P0-8 环 1，建议工程侧复核。
6. **`GetLoginUser` 用 `h.DB.Where("id", userId)` 传 `interface{}`（float64）**（`base_handler.go:97`，值来自 JWT claims）：是否影响索引使用，未实测。
7. **`libapi/library_integration_test.go:164` 调用外部 pytest**（`tests/multi_user/test_platform_library_integration.py`）：该测试目录不在仓库内，`go test ./...` 是否会因此不稳定，未实测。

---

## 四、如果只允许做 3 件事

### 第 1 件：把"钱"和"权限"从约定变成机制（1 个账本 + 1 个守卫）
证据锚点：P0-3/4/5（5 段退款复制品、3 处假事务、支付重试双充）+ P0-1/6/7（22 处路由各自鉴权、admin 禁用无效、无痕伪造）。
- **账本**：`LedgerService.Settle(tx, userID, delta, refType, refID)` —— 必须传 `tx`（签名强制，杜绝"装饰性事务"），`geekai_power_logs` 加 `UNIQUE(ref_type, ref_id)` 做幂等，余额变更与业务状态变更同事务。所有 22 个调用点改为传 ref。
- **守卫**：`main.go` 里 `adminGroup := Engine.Group("/api/admin")` 单点挂载 `AdminGuard`，`AdminGuard` 每次从 DB 校验 `status`+`role`；`RegisterRoutes(group *gin.RouterGroup)` 签名化，让"新文件忘记挂中间件"从语法上不可能。
- **追溯第一性原理**：平台本质是"代收代付 + 额度分销"。资金流的第一性不变量只有两条——**余额不能凭空出现/消失**、**只有被授权的人能改价改额**。这两条是物理约束，必须由数据库与统一入口保证；一旦交给 22 个 handler 各自保证，它就不再是约束而是祈祷。现在的证据是：同一段退款逻辑被抄了 5 遍且每份都非原子——这已经不是"某处写错"，是"没有机制"。

### 第 2 件：切断"用户内容 → 管理员权限"的可达路径（拆信任域）
证据锚点：P0-8（admin 页 v-html 渲染用户对话 + token 在 localStorage）+ P1（单 axios 双 token、`router.js` 无守卫、同 origin）。
- admin 与 site **拆成两个 origin**（哪怕同一份构建产物，不同域名/端口），cookie/token 各自 `SameSite` 且 admin 侧 httpOnly；admin 侧所有 `v-html` 过 sanitize。
- 附带必然结果：`AdminAuthMiddleware` 可以真正做"按 origin 拒绝"，把 P0-1 的漏防从"5 个文件忘了"降级为"跨域直接拒绝"。
- **追溯第一性原理**：这个系统把"任意陌生人的输入"和"能改全站价格的会话"放在同一个 origin 的同一个 localStorage 里。**信任域混同**是一个物理性的结构缺陷——只要两者同源，任何一次渲染疏忽（6 处 `v-html`）都等价于一次管理员权限泄露。不需要新增功能，只需要把两件不同信任级别的东西物理隔开。

### 第 3 件：让 Git 重新成为真理来源，并恢复上游基线
证据锚点：P0-9（HEAD 是 React、工作区是 Vue、64 文件未跟踪、无 upstream remote、545 文件 diff）。
- 提交当前 Vue 前端（`web/src` 全量）；加 `upstream` remote 指向原作者；`admin-web` 的删除单独一个 commit 并打 tag（保留可回退点）；`dist`/`node_modules`/`*.timestamp-*.mjs` 清理并 gitignore；CI 加 `git status --porcelain` 必须为空的门禁。
- **追溯第一性原理**：fork 二开的可持续性 = **能算出与上游的差**。这是可合并性的唯一物理前提——没有公共基线，就没有三方合并，只有"看着别人代码手工重写"。当前 `git diff origin/main...HEAD` 是 545 文件 / −50,167 行，而真正在跑的前端甚至不在其中；这意味着**任何"下次同步上游"都无法验证正确性，也无法回滚**。在钱和用户数据都在这个仓库里的前提下，"不可回滚"本身就是最高级别的架构风险。

---

### 附：不建议现在做的事
- 不建议补 DAO 层 / 补 service 层。当前依赖方向已单向无环（`go list -deps` 验证），强行分层只会增加 22 个 handler 的改动面，而 P0 里没有一条是由"缺少 DAO"直接导致的。
- 不建议"全面加测试"。测试应该长在第 1 件做完之后：账本和守卫有了唯一入口，才存在值得测的不变量；现在测 5 份复制粘贴的退款循环只会把 bug 固化。
