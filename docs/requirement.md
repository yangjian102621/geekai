# 📄 PRD: 智能演示文稿 (PPT) 生成后端 API (Agentic Workflow)

## 1. 项目概述 (Project Overview)

本项目旨在开发一个基于 Go (Gin 框架) 的后端 API 服务。该服务接收用户提供的「文本大纲」和「设计提示词（可选）」，通过大语言模型 (LLM) 进行结构化内容拆解与提示词工程，并利用并发调度图像生成模型 (如 Nano Banana 2) 生成幻灯片配图。

为保证系统稳定性和极佳的用户体验，系统采用**异步任务+前端定时轮询（3秒/次）**的架构。后端负责核心的并发控制、速率限制（Rate Limiting）和状态管理。

## 2. 技术栈约束 (Tech Stack)

- **语言/框架**: Go 1.21+, Gin Web Framework
- **并发控制**: `golang.org/x/sync/errgroup`
- **速率限制**: `golang.org/x/time/rate`
- **状态管理**: `sync.Map` (内存级，后续可扩展至 Redis)
- **数据交互**: JSON 格式，RESTful 风格

## 3. 核心 API 规范 (API Specification)

### 3.1 创建生成任务 (Create Task)

接收用户的原始文本和要求，立即返回任务 ID，不阻塞等待生成结果。

- **Endpoint**: `POST /api/v1/tasks/generate-slides`
- **Request Body**:

```json
{
  "content": "用户的 Markdown 笔记或大纲内容 (必填)",
  "prompt": "用户的附加设计要求，如：极简商务风、卡通手绘风 (选填)",
  "language": "用户选择的语言，如：中文、英文,必填",
  "pages": "用户选择的页数，如：5-20,必填"
}
```

- **Response Body**:

```json
{
  "code": 200,
  "message": "Task created successfully",
  "data": {
    "task_id": "uuid-v4-string",
    "status": "pending"
  }
}
```

### 3.2 查询任务进度 (Query Task Progress)

前端每隔 3 秒调用一次此接口，获取最新进度和已生成的幻灯片数据。

- **Endpoint**: `GET /api/v1/tasks/:task_id`
- **Response Body**:

```json
{
  "code": 200,
  "data": {
    "task_id": "uuid-v4-string",
    "status": "processing", // 状态枚举: pending, processing, completed, failed
    "progress": {
      "total_slides": 10,
      "completed_slides": 3,
      "percentage": 30
    },
    "slides": [
      {
        "slide_index": 1,
        "theme": "早晨环节",
        "title": "第一关：起床咕噜咕噜喝温水！",
        "points": ["水量：300-500ml", "唤醒身体小怪兽"],
        "image_url": "https://example.com/generated-image-url-1.png"
      }
    ],
    "error_message": "" // 仅在 status 为 failed 时返回具体错误
  }
}
```

## 4. 核心工作流与状态机 (State Machine & Workflow)

### 4.1 任务初始化阶段

1. 接收 `POST` 请求，生成全局唯一的 `task_id`。
2. 在 `TaskManager` (基于 `sync.Map`) 中初始化任务实例，初始化状态 `status = pending`，并加上并发安全锁 (`sync.Mutex`)。
3. 启动独立的后台 Goroutine 接管耗时业务，主线程立即向客户端返回 `task_id`。

### 4.2 后台异步执行阶段

1. **LLM 解析 (状态变更: processing)**
   调用大语言模型解析传入的 `content` 和 `prompt`。必须注入以下 System Prompt 常量约束大模型输出：
   > # Role

> 你是一位顶级的专业演示文稿（PPT）策划专家和 AI 图像提示词（Prompt）工程师。任务是根据用户提供的「内容大纲」或「设计要求」，生成一套逻辑清晰、视觉风格高度统一的幻灯片分镜数据。

> # Rules

> 1. 全局风格锚定：根据大纲推断或遵循用户要求的全局视觉风格。所有配图必须严格遵循此风格。
> 2. 结构化拆解：合理拆分为多张幻灯片，单页最多 3-4 个简短要点。
> 3. 视觉转译：为每页构思具体的画面描述 (image_prompt)。必须包含前缀 `[全局风格描述]`，必须包含后缀 `[画面左侧或右侧留出干净的纯色或虚化空间，用于排版文字。绝对不要在图片中生成任何英文字母、汉字或乱码。]`
> 4. 严格输出合法的纯 JSON 数组 `[{"slide_index": 1, "theme": "...", "title": "...", "points": ["..."], "image_prompt": "..."}]`，禁止使用 Markdown 标记。

2. **获取总数**
   解析成功后，更新该任务的 `total_slides` 字段，`completed_slides` 保持为 0。
3. **并发绘图与限流调度 (核心挑战)**

- 启动 `errgroup.WithContext` 进行并发调度，设置 `g.SetLimit(3)` 控制最大并发线程。
- 引入令牌桶限流器 `rate.NewLimiter(rate.Every(1*time.Second), 1)` 控制 API 外部请求频率（QPS = 1）。
- 必须实现**指数退避重试 (Exponential Backoff)**：遇到 HTTP 429 错误时，等待 2s, 4s, 8s 后重试，最高重试 3 次。
- **进度更新**：每成功生成一张图片，必须获取 `Task` 实例的写锁，将 `completed_slides` 加 1，并将完整的 `SlideData` 追加到任务的 `slides` 数组中，以便前端下一秒轮询时能拿到最新切片。

### 4.3 任务收尾阶段

- 正常结束：`errgroup.Wait()` 返回 nil，状态变更为 `completed`。
- 异常中断：发生超过最大重试次数的错误或解析失败，状态变更为 `failed`，并记录错误原因至 `error_message`。

请按照以下顺序逐步实现 Go 后端代码，每完成一步请与我确认：

1. **核心模型与接口定义**：定义 `Task` 结构体（包含 sync.Mutex）、DTO 结构体以及 Gin 的路由组搭建。
2. **状态管理器实现**：实现一个基于 `sync.Map` 的 `TaskManager`，提供安全的 `CreateTask`、`GetTask`、`UpdateProgress` 和 `MarkAsFailed` 方法。
3. **限流与重试 HTTP 客户端**：实现 `callImageAPIWithRetry` 函数，内部封装 `rate.Limiter` 的等待逻辑和针对 429 状态码的指数退避重试逻辑（可先用 Mock 数据代替真实 HTTP 请求）。
4. **组装核心 Handler**：实现 `POST` 和 `GET` 控制器。在 `POST` 的后台 Goroutine 中，串联 LLM Mock 解析和 `errgroup` 并发画图流程，确保在循环中安全地调用 `TaskManager.UpdateProgress` 更新进度和结果数组。
