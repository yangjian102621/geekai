# ChatGPT-Plus

**ChatGPT-PLUS** 基于 AI 大语言模型 API 实现的 AI 助手全套开源解决方案，自带运营管理后台，开箱即用。集成了 OpenAI, Azure,
ChatGLM,讯飞星火，文心一言等多个平台的大语言模型。集成了 MidJourney 和 Stable Diffusion AI绘画功能。主要有如下特性：

* 完整的开源系统，前端应用和后台管理系统皆可开箱即用。
* 基于 Websocket 实现，完美的打字机体验。
* 内置了各种预训练好的角色应用，比如小红书写手，英语翻译大师，苏格拉底，孔子，乔布斯，周报助手等。轻松满足你的各种聊天和应用需求。
* 支持 OPenAI，Azure，文心一言，讯飞星火，清华 ChatGLM等多个大语言模型。
* 支持 MidJourney / Stable Diffusion AI 绘画集成，开箱即用。
* 支持使用个人微信二维码作为充值收费的支付渠道，无需企业支付通道。
* 已集成支付宝支付功能，支持多种会员套餐和点卡购买功能。
* 集成插件 API 功能，可结合大语言模型的 function 功能开发各种强大的插件，已内置实现了微博热搜，今日头条，今日早报和 AI
  绘画函数插件。

## 功能截图

### PC 端聊天界面

![ChatGPT Chat Page](/docs/imgs/gpt.gif)

### AI 对话界面

![ChatGPT new Chat Page](/docs/imgs/chat-new.png)

### MidJourney 专业绘画界面

![mid-journey](/docs/imgs/mj_image.jpg)

### Stable-Diffusion 专业绘画页面

![Stable-Diffusion](/docs/imgs/sd_image.jpg)
![Stable-Diffusion](/docs/imgs/sd_image_detail.jpg)

### 绘图作品展

![ChatGPT image_list](/docs/imgs/image-list.png)

### AI应用列表

![ChatGPT-app-list](/docs/imgs/app-list.jpg)

### 会员充值

![会员充值](/docs/imgs/member.png)

### 自动调用函数插件

![ChatGPT function plugin](/docs/imgs/plugin.png)
![ChatGPT function plugin](/docs/imgs/mj.jpg)

### 管理后台

![ChatGPT admin](/docs/imgs/admin_dashboard.png)
![ChatGPT admin](/docs/imgs/admin_config.jpg)
![ChatGPT admin](/docs/imgs/admin_models.jpg)
![ChatGPT admin](/docs/imgs/admin_user.png)

### 移动端 Web 页面

![Mobile chat list](/docs/imgs/mobile_chat_list.png)
![Mobile chat session](/docs/imgs/mobile_chat_session.png)
![Mobile chat setting](/docs/imgs/mobile_user_profile.png)
![Mobile chat setting](/docs/imgs/mobile_pay.png)

### 7. 体验地址

> 免费体验地址：[https://ai.r9it.com/chat](https://ai.r9it.com/chat) <br/>
> **注意：请合法使用，禁止输出任何敏感、不友好或违规的内容！！！**

## 使用须知

1. 本项目基于 MIT 协议，免费开放全部源代码，可以作为个人学习使用或者商用。
2. 如需商用必须保留版权信息，请自觉遵守。确保合法合规使用，在运营过程中产生的一切任何后果自负，与作者无关。

## 项目介绍

这一套完整的系统，包括前端聊天应用和一个后台管理系统。系统有用户鉴权，你可以自己使用，也可以部署直接给 C 端用户提供
ChatGPT 的服务。

### 项目的技术架构

新版的系统前后端都进行大改动的重构，后端还是用的 Gin Web 框架，但是作者整合了 fx 自动注入框架，整个后端应用结构非常简洁，特别适合二次开发。
另外，数据存储用 MySQL 替换了 leveldb, 因为要对 C 端，后期会涉及到很多业务数据查询统计，leveldb 已经完全不够用了。

> Gin + fx + MySQL

3.0 版本之后会陆续添加其他语言的 API 实现，比如 PHP，Java 等。考虑到作者精力有限，api 目录已经添加了，有兴趣的同学自主去认领各自擅长的语言去实现。

前端的框架还是:

> Vue3 + Element-Plus

前后台的页面风格已经全部变了，几乎所有页面样式代码都重写了。逻辑代码还是沿用之前的，毕竟功能没有太大的变化。

此次重构改版主要是为了后面功能的扩展准备了。

新版本已经实现的功能如下：

1. 引入用户体系，新增用户注册和登录功能。
2. 聊天页面改版，实现了跟 ChatGPT 官方版本一致的聊天体验。
3. 创建会话的时候可以选择聊天角色和模型。
4. 新增聊天设置功能，用户可以导入自己的 API KEY
5. 保存聊天记录，支持聊天上下文。
6. 重构后台管理模块，更友好，扩展性更好的后台管理系统。
7. 引入 ip2region 组件，记录用户的登录IP和地址。
8. 支持会话搜索过滤。
9. 支持微信支付充值

## 项目地址

* Github 地址：https://github.com/yangjian102621/chatgpt-plus
* 码云地址：https://gitee.com/blackfox/chatgpt-plus

## 客户端下载

目前已经支持 Win/Linux/Mac/Android 客户端，下载地址为：https://github.com/yangjian102621/chatgpt-plus/releases/tag/v3.1.2

## TODOLIST

* [x] 整合 Midjourney AI 绘画 API
* [x] 开发移动端聊天页面
* [x] 接入微信收款功能
* [x] 支持 ChatGPT 函数功能，通过函数实现插件
* [x] 开发桌面版应用
* [x] 开发手机 App 客户端
* [x] 支付宝支付功能
* [ ] 支持基于知识库的 AI 问答
* [ ] 会员推广功能
* [ ] 微信支付功能

## 快速部署

请参考 [ChatGPT-Plus 部署文档](https://ai.r9it.com/docs/)。

## 本地开发调试

本地开发同样要分别运行前端和后端程序。

### 运行后端程序

1. 同样你首先要 [导入数据库](#1-导入数据库)
2. 然后 [修改配置文档](#2-修改配置文档)
3. 运行后端程序：

    ```shell
    cd api 
    # 1. 先下载依赖
    go mod tidy
    # 2. 运行程序
    go run main.go
    # 如果你安装了 fresh 可以使用 fresh 实现热启动
    fresh -c fresh.conf
    ```

### 运行前端程序

同样先拷贝配置文档：

```shell
cd web
cp .env.production .env.development
```

编辑 `.env.development` 文件，修改后端 API 的访问路径：

```ini
VUE_APP_API_HOST=http://localhost:5678
VUE_APP_WS_HOST=ws://localhost:5678
```

配置好了之后就可以运行前端应用了：

```
# 安装依赖
npm install
# 运行
npm run dev
```

* 前端页面：http://localhost:8888/chat
* 后台管理页面：http://localhost:8888/admin

## 项目打包

由于本项目是采用异构开发的方式，所项目打包分成两步：首先编译后端程序，然后再打包前端应用。

### 打包前端

```shell
cd web
npm run build
```

### 打包后端

你可以根据个人需求将项目打包成 windows/linux/darwin 平台项目。

```shell
cd api
# for all platforms
make clean all
# for linux only
make clean linux
```

打包后的可执行文件在 `bin` 目录下。

## 参与贡献

个人的力量始终有限，任何形式的贡献都是欢迎的，包括但不限于贡献代码，优化文档，提交 issue 和 PR 等。

如果有兴趣的话，也可以加微信进入微信讨论群（**添加好友时请注明来自Github!!!**）。

![微信名片](docs/imgs/wx.png)

#### 特此声明：不接受在微信或者微信群给开发者提 Bug，有问题或者优化建议请提交 Issue 和 PR。非常感谢您的配合！

### Commit 类型

* feat: 新特性或功能
* fix: 缺陷修复
* docs: 文档更新
* style: 代码风格或者组件样式更新
* refactor: 代码重构，不引入新功能和缺陷修复
* opt: 性能优化
* chore: 一些不涉及到功能变动的小提交，比如修改文字表述，修改注释等

## 打赏

如果你觉得这个项目对你有帮助，并且情况允许的话，可以请作者喝杯咖啡，非常感谢你的支持～

![打赏](docs/imgs/donate.png)

![Star History Chart](https://api.star-history.com/svg?repos=yangjian102621/chatgpt-plus&type=Date)



