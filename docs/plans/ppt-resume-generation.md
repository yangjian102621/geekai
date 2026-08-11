# PPT 断点续生成（简化方案 · 修订版）

## 设计原则（按产品确认）

- **不新增数据库字段**，只复用现有 `geekai_ppt_jobs.slides`（JSON）。
- 分镜结构已在 [`vo.PPTSlideData`](api/store/vo/ppt_job.go) 中体现：`theme/title/points/image_prompt` + `image_url`。
- **大模型生成分镜后立刻整表写入 `slides`**：每条记录带完整分镜字段，**`image_url` 为空**，表示「待生图」。
- **列表「继续生成」条件**：`status !== 'completed'` 即可展示（实现时需排除无意义的态：如从未写入分镜的异常数据；**`processing` 时按钮应禁用或显示加载**，避免与进行中的任务重复提交）。
- **继续生成**：读出 `slides`，**有 `image_url` 则跳过**，无则取当条分镜调生图 API；**每成功一张就更新一次数据库**。
- **并发更新 `slides`**：多 goroutine 同时写同一任务的 `slides` 字段时，必须 **加锁**（推荐在 [`PptService`](api/service/ppt/ppt_service.go) 内 **`sync.Mutex` 按 `task_id` 维度**保护「读 JSON → 改一条 slide → 写回」的整段逻辑，避免后写覆盖先写）。

## 后端改动要点

### 1. `RunTask` 流程调整（[`ppt_service.go`](api/service/ppt/ppt_service.go)）

1. `GenerateSlides` 成功后，将 `[]slidePlan` 转为 `[]SlideData` / `vo.PPTSlides`，**全部 `ImageURL` 置空**，一次性 `Save` / `Updates` 写入 `slides`，并 `SetTotalSlides(total)`。
2. **`completed_slides` 语义**：表示「已有配图」的页数，等于 `slides` 中 `image_url` 非空的条数（或单独在更新时维护，与 `len(slides)` 区分）。
3. 抽取 **`runSlideImageGeneration`**：对「需要生图」的条目并发生图；每条成功后调用新的 **`applySlideImage(taskID, slideIndex, url)`**（内部持锁、按 index 更新对应元素的 `image_url`，再扣算力）。

### 2. 替换「仅 append」的 `UpdateProgress`

- 现有 `UpdateProgress` 是 **append** 一条 slide，与「先写满占位再填图」冲突。
- 改为 **按 `slide_index` 原地更新** 指定项的 `image_url`（及必要时 `thumb`），并在持锁下 **整份序列化写回**。

### 3. `ResumeTask` / `POST .../resume`

- 校验用户、任务存在；若 `status === processing` 建议直接拒绝或返回「任务进行中」。
- 若 `status !== completed`：加载 `slides`，筛出 `image_url` 为空的项，**仅对这些项**走与 `RunTask` 相同的生图 + `applySlideImage`。
- 全部非空后：`UpdateStatus(completed)`，清空或保留 `err_msg` 按产品定。
- 若某次仍失败：可 `MarkAsFailed` 保留已生成页（与现网一致）。

### 4. 失败路径

- 任一页生图失败时，仍可 `MarkAsFailed`，但 **分镜与已生成图均已在 `slides` 中**，续跑只需补空 `image_url`。

## 前端（[`PPTCreate.vue`](web/src/views/PPTCreate.vue)）

- 列表与详情：当 **`status !== 'completed'`**（且建议 **`status !== 'processing'` 才可点**）显示「继续生成」，**`el-tooltip`** 文案：继续生成尚未完成的幻灯片，已有页面保持不变。
- 调用 `POST /api/v1/tasks/:id/resume`（路由注册方式同既有 export）。

## 与旧方案差异

- **不增加** `slide_plans` 等字段；**不重复存储**分镜，全部以 `slides` 中「无图占位」表达待生成状态。

## 验证

- 人为让第 N 页失败：列表出现继续生成，点击后仅补第 N 页及之后空图项。
- 并发压测：多页同时完成时 `slides` JSON 无丢失、无覆盖。

## 实施状态（已实现）

- 后端：`saveSlidesOutline` + `ApplySlideImage`（`sync.Map` 按 `task_id` 互斥）、`runSlideImageJobs`、`ResumeTask`、`POST /api/v1/tasks/:task_id/resume`。
- 前端：`PPTCreate.vue` 列表/详情「继续生成」+ `el-tooltip`，`processing` 时禁用或 loading。
- **旧任务**：若 `slides` 条数小于 `total_slides`（无完整占位），`ResumeTask` 返回不可续跑。
