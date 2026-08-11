# Repository Guidelines

## 项目结构与模块组织
- `api/`：Go + Gin 后端，包含 `core/` 业务、`handler/` 控制器、`service/` 接口调用以及 `store/` 数据访问，`config.toml` 存放默认配置，`Makefile` 用于多架构交叉构建。
- `web/`：Vue3 + Vite 前端，源码集中在 `src/`，`public/` 为静态资源，`dist/` 存放构建结果并可被 `api/static/` 或 `desktop/` 引用。
- `desktop/`：Electron 客户端入口为 `index.js`，配合 `electron-builder` 可打包 AppImage/DMG/NSIS。
- `miniprogram/`、`docs/`、`database/` 与 `config/` 分别承载小程序壳、部署文档、SQL 脚本及全局 YAML 配置；`build/` 包含 Dockerfile、安装脚本。

## 构建、测试与开发命令
- `docker-compose up -d`：根目录拉起全部容器，需提前准备好 MySQL、Redis 与模型密钥。
- `cd api && go run main.go`：本地热调试；`make amd64` / `make arm64` 生成无 CGO 二进制至 `api/bin/` 便于镜像打包。
- `cd web && pnpm install && pnpm dev --host`：Vite 开发模式；`pnpm build` 产出静态文件；`pnpm lint` 运行 ESLint 自动修复。
- `cd desktop && npm install && npm run start`：调试 Electron；`npm run package` 通过 electron-builder 生成多平台安装包。

## 编码风格与命名规范
- Go 代码必须经过 `gofmt`/`goimports`，保持 tab 缩进与驼峰命名；HTTP 路由遵循 `/api/v1/resources` 模式，与 handler 函数命名 (`ResourceHandler`) 对应。
- Vue 组件文件使用 PascalCase（如 `ChatPanel.vue`），Pinia store 与工具采用 kebab-case 文件名（如 `chat-session.ts`）；统一通过 ESLint、Tailwind 与 `postcss.config.js` 约束样式。

## 测试指南
- `cd api && go test ./... -race` 是最低要求，新增 service/handler 需补 `_test.go` 并用 mock 隔离第三方 API；涉及时序逻辑可新增 `Test*Integration` 验证。
- 前端暂未启用单测框架，至少运行 `pnpm lint` 并在 PR 中附关键页面截图或录屏证明交互可用；桌面端如修改构建脚本，需在 macOS/Linux/Windows 中至少验证一个安装包。

## 提交与 Pull Request 规范
- 参考历史记录（如“支持腾讯云短信服务”），提交信息使用中文动词开头、聚焦单一变更，并可加子系统前缀：`web: 优化聊天动画`。
- PR 描述需包含变更背景、实现概述、验证方式（命令、截图或日志）与关联 issue/任务号；涉及配置或部署脚本，还要说明回滚流程并 @ 相关 reviewer。

## 安全与配置提示
- 禁止提交真实密钥，请复制 `config.sample.toml` 或 `config/config.yaml` 生成私有文件，并用 `git update-index --skip-worktree` 忽略。
- 对象存储、短信、支付等凭证统一放入 Vault 或 CI Secret，代码中仅引用占位常量；`docs/` 中同步记录新增敏感字段与启用步骤。
