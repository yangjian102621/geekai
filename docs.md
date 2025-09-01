## 开发 3D 图片生成功能

对接 3D 图片生成接口，为当前系统添加 3D 模型生成功能，默认支持腾讯云和 Gitee 的图生 3D API 接口。

## 要求

1. 完成数据库设计，后端 API 设计，前端页面设计。
2. 完成前端功能页面以及后台管理页面，具体设计结构可以参考即梦 AI，在对应的模块建立独立的模块 。
3. 页面设计要精美，但是整体风格要跟整站风格一致。
4. 支持前端 3D 模型预览，支持 3D 模型下载。

## 腾讯云图生 3D API 接口文档

1. 提交任务： https://cloud.tencent.com/document/product/1804/120826
2. 查询任务： https://cloud.tencent.com/document/product/1804/120827
3. Golang SDK： https://gitee.com/TencentCloud/tencentcloud-sdk-go/blob/master/tencentcloud/ai3d/v20250513/client.go 依赖已经安装到本地了

## Gitee 图生 3D API 接口文档

1. 提交任务： https://ai.gitee.com/docs/openapi/v1#tag/3d-%E7%94%9F%E6%88%90/post/async/image-to-3d
2. 查询任务：https://ai.gitee.com/docs/openapi/v1#tag/%E5%BC%82%E6%AD%A5%E4%BB%BB%E5%8A%A1/get/task/{task_id}/get

首先，你需要认真阅读上述接口文档，然后按照接口文档的示例代码实现腾讯云和 Gitee 的图生 3D API 接口，并且将接口集成到现有的系统中。

📋 功能概述

     为现有的GeekAI-Plus系统添加3D图片生成功能，集成腾讯云和Gitee的图生3D API接口，包含完整的前后端功能和管理界面。

     🗄️ 数据库设计

     新增数据表：geekai_3d_jobs
     - id (uint): 主键
     - type (string): API类型 (tencent/gitee)
     - user_id (uint): 用户ID
     - power (int): 消耗算力
     - task_id (string): 第三方任务ID
     - img_url (string): 生成的3D模型文件地址
     - model (string): 使用的3D模型类型
     - status (string): 任务状态
     - err_msg (string): 错误信息
     - params (JSON): 任务参数(包含输入图片、提示词等所有参数)
     - created_at (int64): 创建时间

     🔧 后端API实现

     路由结构：/api/3d/*
     - POST /api/3d/generate - 创建3D生成任务
     - GET /api/3d/jobs - 获取任务列表
     - GET /api/3d/job/{id} - 获取任务详情
     - GET /api/3d/download/{id} - 下载3D模型
     - DELETE /api/3d/job/{id} - 删除任务

     核心服务：
     - service/3d/tencent_client.go - 腾讯云3D API客户端
     - service/3d/gitee_client.go - Gitee 3D API客户端
     - service/3d/service.go - 3D生成服务统一接口
     - handler/3d_handler.go - HTTP处理器
     - store/vo/3d_job.go - 数据模型

     🎨 前端界面设计

     用户端页面：/3d - 3D生成主页面
     - 参考JiMeng.vue的设计风格和布局
     - 使用CustomTab组件分离平台参数：
       - Tab 1: "魔力方舟" (Gitee平台参数)
       - Tab 2: "腾讯混元" (腾讯云平台参数)
     - 每个Tab内包含：
       - 图片上传区域
       - 模型选择下拉框
       - 算力消耗实时显示
       - 平台特定的参数配置
       - 生成按钮
     - 任务列表和状态显示
     - 集成3D模型预览器 (three.js)
     - 模型下载功能

     移动端适配：
     - mobile/3dCreate.vue - 移动端3D生成页面
     - 保持Tab切换功能
     - 响应式设计，触控优化

     🛠️ 管理后台

     管理功能：
     - admin/3d/3dJobs.vue - 任务管理列表
     - admin/3d/3dSetting.vue - API配置页面
     - 模型配置管理：
       - 分平台配置模型列表
       - 设置每个模型的算力消耗值
     - API密钥和端点配置

     🔌 API集成方案

     腾讯云集成：
     - 使用官方Golang SDK
     - 支持异步任务提交和状态查询

     Gitee集成：
     - HTTP客户端实现
     - 标准化响应处理

     🎯 核心功能特性

     - 平台切换：通过CustomTab在魔力方舟和腾讯混元间切换
     - 模型选择：每个平台支持不同的3D模型
     - 动态算力：切换模型时实时更新算力消耗显示
     - 参数隔离：不同平台的参数配置完全分离
     - 3D预览：集成Three.js实现模型预览
     - 统一体验：保持与JiMeng.vue相似的交互风格

     📱 用户体验

     - JiMeng.vue风格的简洁界面
     - Tab切换流畅的平台选择
     - 模型选择时算力消耗实时更新
     - 支持拖拽上传图片
     - 实时任务状态显示
     - 3D模型交互式预览

     这个设计将创建一个与现有JiMeng功能风格一致的3D生成模块，通过Tab分离实现平台参数的清晰管理。

整个实现严格按照现有系统的代码规范和架构模式，与 JiMeng 等模块保持一致的用户体验！
