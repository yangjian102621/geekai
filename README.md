# ChatGPT-Plus

**ChatGPT-PLUS** 是基于 OpenAI API 实现的 ChatGPT 聊天系统。主要有如下特性：

* 完整的开源系统，前端应用和后台管理系统皆可开箱即用。
* 聊天体验跟 ChatGPT 官方版本完全一致。
* 内置了各种预训练好的角色，比如小红书写手，英语翻译大师，苏格拉底，孔子，乔布斯，周报助手等。轻松满足你的各种聊天和应用需求。

## 功能截图

### PC 端聊天界面

![ChatGPT Chat Page](docs/imgs/gpt.gif)

### 新版聊天界面

![ChatGPT new Chat Page](docs/imgs/chat-new.png)

### 自动调用函数插件

![ChatGPT function plugin](docs/imgs/plugin.png)

### 用户设置

![ChatGPT user profle](docs/imgs/user_profile.png)

### 登录页面

![ChatGPT Login](docs/imgs/login.png)

### 管理后台

![ChatGPT admin](docs/imgs/admin_dashboard.png)

![ChatGPT admin](docs/imgs/admin_config.png)

![ChatGPT admin](docs/imgs/admin_user.png)

### 移动端 Web 页面

![Mobile chat list](/docs/imgs/mobile_chat_list.png)
![Mobile chat session](/docs/imgs/mobile_chat_session.png)
![Mobile chat setting](/docs/imgs/mobile_chat_setting.png)

### 7. 体验地址

> 免费体验地址：[https://www.chat-plus.net/chat](https://www.chat-plus.net/chat) <br/>
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

## TODOLIST

* [ ] 整合 Midjourney AI 绘画 API
* [ ] 开发移动端聊天页面
* [x] 接入微信支付功能
* [x] 支持 ChatGPT 函数功能，通过函数实现插件
* [ ] 接入语音和 TTS API，支持语音聊天
* [ ] 开发手机 App 客户端

## Docker 快速部署

>
鉴于最新不少网友反馈在部署的时候遇到一些问题，大部分问题都是相同的，所以我这边做了一个视频教程 [五分钟部署自己的 ChatGPT 服务](https://www.bilibili.com/video/BV1H14y1B7Qw/)。
> 习惯看视频教程的朋友可以去看视频教程，视频的语速比较慢，建议 2 倍速观看。

V3.0.0 版本以后已经支持使用容器部署了，跳过所有的繁琐的环境准备，一条命令就可以轻松部署上线。

### 1. 导入数据库

首先我们需要创建一个 MySQL 容器，并导入初始数据库。

```shell
cd docker/mysql
# 创建 mysql 容器
docker-compose up -d
# 导入数据库
docker exec -i chatgpt-plus-mysql sh -c 'exec mysql -uroot -p12345678' < ../../database/chatgpt_plus.sql
```

如果你本地已经安装了 MySQL 服务，那么你只需手动导入数据库即可。

```shell
# 连接数据库
mysql -u username -p password
# 导入数据库
source database/chatgpt_plus.sql
```

### 2. 修改配置文档

修改配置文档 `docker/conf/config.toml` 配置文档，修改代理地址和管理员密码：

```toml
Listen = "0.0.0.0:5678"
ProxyURL = ["YOUR_PROXY_URL"] # 替换成你本地代理，如：http://127.0.0.1:7777
#ProxyURL = "" 如果你的服务器本身就在墙外，那么你直接留空就好了
MysqlDns = "root:12345678@tcp(172.22.11.200:3307)/chatgpt_plus?charset=utf8&parseTime=True&loc=Local"
StartWechatBot = false # 是否启动微信机器人，默认关闭，如果设置为 TRUE 则启动服务的时候需要微信扫码登录
EnabledMsgService = false # 注册时是否开启短信验证功能，该功能需要配合短信服务一起使用

[Session]
  SecretKey = "azyehq3ivunjhbntz78isj00i4hz2mt9xtddysfucxakadq4qbfrt0b7q3lnvg80"
  Name = "CHAT_SESSION_ID"
  Path = "/"
  Domain = ""
  MaxAge = 86400
  Secure = false
  HttpOnly = false
  SameSite = 2

[Manager]
  Username = "admin"
  Password = "admin123" # 如果是生产环境的话，这里管理员的密码记得修改
  
[ApiConfig] # 插件 API 服务配置，此为第三方插件服务，如需使用请联系作者开通
  ApiURL = "{URL}"
  AppId = "{APP_ID}"
  Token = "{TOKEN}"

[SmsConfig] # 阿里云短信服务配置
  AccessKey = "{YOUR_ACCESS_KEY}"
  AccessSecret = "{YOUR_SECRET_KEY}"
  Product = "Dysmsapi"
  Domain = "dysmsapi.aliyuncs.com"

```

修改 nginx 配置文档 `docker/conf/nginx/conf.d/chatgpt-plus.conf`，把后端转发的地址改成当前主机的内网 IP 地址。

```shell
 # 这里配置后端 API 的转发
location /api/ {
       proxy_http_version 1.1;
       proxy_connect_timeout 300s;
       proxy_read_timeout 300s;
       proxy_send_timeout 12s;
       proxy_set_header Host $host;
       proxy_set_header X-Real-IP $remote_addr;
       proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
       proxy_set_header Upgrade $http_upgrade;
       proxy_set_header Connection $connection_upgrade;
       proxy_pass http://172.28.173.76:6789; # 这里改成后端服务的内网 IP 地址
}
```

### 3. 启动应用

```shell
cd docker
docker-compose up -d
```

* 前端访问地址：http://localhost:8080/chat
* 后台管理地址：http://localhost:8080/admin
* 移动端地址：http://localhost:8080/mobile

> 注意：你得访问后台管理系统 http://localhost:8080/admin
> 输入你前面配置文档中设置的管理员用户名和密码登录。
> 然后进入 `API KEY 管理` 菜单，添加一个 OpenAI 的 API KEY 才可以正常开启 AI 对话。

![add API Key](docs/imgs/apikey_add.png)

最后登录前端聊天页面 [http://localhost:8080/chat](http://localhost:8080/chat)
你可以注册新用户，也可以使用系统默认有个账号：`geekmaster/12345678` 登录聊天。

祝你使用愉快！！！

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
**尤其是新版本的开发计划比较大，包括各种语言的后端 API 实现，本人精力有限，希望借助社区的力量来完成这些 API 的开发。**

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

![微信打赏](docs/imgs/wechat-pay.png)


