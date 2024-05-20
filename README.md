# GeekAI
### 本项目已经正式更名为 GeekAI，请大家及时更新代码克隆地址。

**GeekAI** 基于 AI 大语言模型 API 实现的 AI 助手全套开源解决方案，自带运营管理后台，开箱即用。集成了 OpenAI, Azure,
ChatGLM,讯飞星火，文心一言等多个平台的大语言模型。集成了 MidJourney 和 Stable Diffusion AI绘画功能。

主要特性：

- 完整的开源系统，前端应用和后台管理系统皆可开箱即用。
- 基于 Websocket 实现，完美的打字机体验。
- 内置了各种预训练好的角色应用，比如小红书写手，英语翻译大师，苏格拉底，孔子，乔布斯，周报助手等。轻松满足你的各种聊天和应用需求。
- 支持 OPenAI，Azure，文心一言，讯飞星火，清华 ChatGLM等多个大语言模型。
- 支持 Suno 文生音乐
- 支持 MidJourney / Stable Diffusion AI 绘画集成，文生图，图生图，换脸，融图。开箱即用。
- 支持使用个人微信二维码作为充值收费的支付渠道，无需企业支付通道。
- 已集成支付宝支付功能，微信支付，支持多种会员套餐和点卡购买功能。
- 集成插件 API 功能，可结合大语言模型的 function 功能开发各种强大的插件，已内置实现了微博热搜，今日头条，今日早报和 AI
  绘画函数插件。

### 🚀 更多功能请查看 [GeekAI-PLUS](https://github.com/yangjian102621/geekai-plus)

- [x] 更友好的 UI 界面
- [x] 支持 Dall-E 文生图功能
- [x] 支持文生思维导图
- [x] 支持为模型绑定指定的 API KEY，支持为角色绑定指定的模型等功能
- [x] 支持网站 Logo 版权等信息的修改

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

### 体验地址

> 免费体验地址：[https://ai.r9it.com/chat](https://ai.r9it.com/chat) <br/>
> **注意：请合法使用，禁止输出任何敏感、不友好或违规的内容！！！**

## 快速部署

请参考文档 [**GeekAI 快速部署**](https://ai.r9it.com/docs/install/)。

## 使用须知

1. 本项目基于 Apache2.0 协议，免费开放全部源代码，可以作为个人学习使用或者商用。
2. 如需商用必须保留版权信息，请自觉遵守。确保合法合规使用，在运营过程中产生的一切任何后果自负，与作者无关。

## 项目地址

* Github 地址：https://github.com/yangjian102621/geekai
* 码云地址：https://gitee.com/blackfox/geekai

## 客户端下载

目前已经支持 Win/Linux/Mac/Android 客户端，下载地址为：https://github.com/yangjian102621/geekai/releases/tag/v3.1.2

## TODOLIST

* [ ] 支持基于知识库的 AI 问答
* [ ] 会员邀请注册推广功能
* [ ] 微信支付功能

## 项目文档

最新的部署视频教程：[https://www.bilibili.com/video/BV1Cc411t7CX/](https://www.bilibili.com/video/BV1Cc411t7CX/)

详细的部署和开发文档请参考 [**GeekAI 文档**](https://ai.r9it.com/docs/)。

加微信进入微信讨论群可获取 **一键部署脚本（添加好友时请注明来自Github!!!）。**

![微信名片](docs/imgs/wx.png)

## 参与贡献

个人的力量始终有限，任何形式的贡献都是欢迎的，包括但不限于贡献代码，优化文档，提交 issue 和 PR 等。

#### 特此声明：由于个人时间有限，不接受在微信或者微信群给开发者提 Bug，有问题或者优化建议请提交 Issue 和 PR。非常感谢您的配合！

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

![Star History Chart](https://api.star-history.com/svg?repos=yangjian102621/geekai&type=Date)
