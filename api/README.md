# chatgpt-plus-go

chatgpt-plus 后端 API Go 语言实现。技术选型采用 Gin + Mysql 架构，依赖注入使用的是 fx 框架，ORM 采用的是 GORM 框架。

## 支付说明

- 支付配置统一在 `system.payment` 中管理，包含 `alipay`、`wxpay`、`epay` 和 `stripe`。
- Stripe 走 Checkout 一次性支付，不支持订阅。
- 商品表新增 `stripe_price` 字段，作为 Stripe 专用美元价格；旧数据会在迁移时自动用 `price` 回填。
- Stripe 支付回调以 webhook 为准，成功后回到会员页并由现有订单轮询刷新状态。

