# 即梦 AI 配置功能说明

## 功能概述

即梦 AI 配置功能允许管理员通过 Web 界面配置即梦 AI 的 API 密钥和算力消耗设置，支持动态配置更新，无需重启服务。

## 功能特性

### 1. 秘钥配置

- AccessKey 和 SecretKey 配置
- 支持密码显示/隐藏
- 连接测试功能

### 2. 算力配置

- 文生图算力消耗
- 图生图算力消耗
- 图片编辑算力消耗
- 图片特效算力消耗
- 文生视频算力消耗
- 图生视频算力消耗

### 3. 动态配置

- 配置实时生效
- 无需重启服务
- 支持配置验证

## API 接口

### 获取配置

```
GET /api/admin/jimeng/config
```

### 更新配置

```
POST /api/admin/jimeng/config
Content-Type: application/json

{
  "config": {
    "access_key": "your_access_key",
    "secret_key": "your_secret_key",
    "power": {
      "text_to_image": 10,
      "image_to_image": 15,
      "image_edit": 20,
      "image_effects": 25,
      "text_to_video": 30,
      "image_to_video": 35
    }
  }
}
```

### 测试连接

```
POST /api/admin/jimeng/config/test
Content-Type: application/json

{
  "config": {
    "access_key": "your_access_key",
    "secret_key": "your_secret_key"
  }
}
```

## 前端页面

### 访问路径

管理后台 -> 即梦 AI -> 配置设置

### 页面功能

1. **秘钥配置标签页**

   - AccessKey 输入框（密码模式）
   - SecretKey 输入框（密码模式）
   - 测试连接按钮

2. **算力配置标签页**

   - 各种任务类型的算力消耗配置
   - 数字输入框，支持 1-100 范围
   - 提示信息说明

3. **操作按钮**
   - 保存配置
   - 重置配置

## 配置存储

配置存储在数据库的`config`表中：

- 配置键：`jimeng`
- 配置值：JSON 格式的即梦 AI 配置

## 默认配置

如果配置不存在，系统会使用以下默认值：

```json
{
  "access_key": "",
  "secret_key": "",
  "power": {
    "text_to_image": 10,
    "image_to_image": 15,
    "image_edit": 20,
    "image_effects": 25,
    "text_to_video": 30,
    "image_to_video": 35
  }
}
```

## 使用流程

1. **初始配置**

   - 访问管理后台即梦 AI 配置页面
   - 填写 AccessKey 和 SecretKey
   - 点击"测试连接"验证配置
   - 调整各功能算力消耗
   - 保存配置

2. **配置更新**

   - 修改需要更新的配置项
   - 保存配置
   - 配置立即生效

3. **故障排查**
   - 使用"测试连接"功能验证 API 密钥
   - 检查配置是否正确保存
   - 查看服务日志

## 注意事项

1. **权限要求**

   - 只有管理员可以访问配置页面
   - 需要有效的管理员登录会话

2. **配置验证**

   - AccessKey 和 SecretKey 不能为空
   - 算力消耗必须大于 0
   - 建议先测试连接再保存配置

3. **服务影响**
   - 配置更新不会影响正在进行的任务
   - 新任务会使用更新后的配置
   - 客户端配置会在下次请求时更新

## 错误处理

1. **配置加载失败**

   - 使用默认配置
   - 记录错误日志

2. **连接测试失败**

   - 显示具体错误信息
   - 建议检查 API 密钥

3. **配置保存失败**
   - 显示错误信息
   - 保留原有配置

## 开发说明

### 后端文件

- `api/handler/admin/jimeng_handler.go` - 配置管理 API
- `api/service/jimeng/service.go` - 配置服务逻辑
- `api/core/types/jimeng.go` - 配置类型定义

### 前端文件

- `web/src/views/admin/jimeng/JimengSetting.vue` - 配置页面

### 数据库

- `config`表存储配置信息
- 配置键：`jimeng`
- 配置值：JSON 格式
