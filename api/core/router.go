package core

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core/middleware"
	"geekai/core/types"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// adminPrefix 是管理后台路由的唯一根前缀。
// 全仓只有本文件可以出现它的字面量，其它地方一律通过 Router 的方法取得。
const adminPrefix = "/api/admin"

// Router 是全局唯一的路由注册入口。
//
// 设计约束（刻意的，写在这里以免后来者"顺手优化"回去）：
//  1. Router 不持有、不暴露 *gin.Engine。handler 层物理上无法拿到引擎。
//  2. 四个 group 在 NewRouter 里一次性建成，安全级别此后不可变；
//     Router 不提供任何 Use() 之类"事后追加中间件"的方法。
//  3. 唯一允许出现 /api/admin 前缀的是 AdminGroup 与 AdminPublicGroup；
//     PublicGroup / UserGroup 传入 /api/admin 前缀在启动期直接 panic。
//
// 为什么这样能消灭「gin 的 Use() 只对之后注册的路由生效」这个语义：
// gin 的 Group() 会复制父组的 handler 链，Use() 只往本组链上追加、只影响之后注册的路由。
// 于是"一个 group 同时承载公开与受保护两种级别、中间插一行 Use()"这种写法天然只能靠作者
// 的记忆保证正确。本类型把它拆成两个不可再组合的对象：安全级别由方法名决定，且在建组时
// 就已经烧进 handler 链，两次调用拿到的是两个不同对象 —— "忘了 Use"这个动作没有地方可以发生。
//
// 不给 *gin.Engine 提供任何 getter（Engine() / Gin() / GetEngine() 一律不加）。
// 测试需要引擎时用 NewServerWithEngine 拿自己传进去的那个变量，而不是从 AppServer 里取。
type Router struct {
	pub      *gin.RouterGroup // 无鉴权，挂在 "/"
	user     *gin.RouterGroup // 已挂 UserAuthMiddleware，挂在 "/"
	admin    *gin.RouterGroup // 已挂 AdminAuthMiddleware，挂在 "/api/admin"
	adminPub *gin.RouterGroup // 无鉴权，挂在 "/api/admin"，只服务于 login / logout
}

// NewRouter 纯函数构造：不改 gin 全局模式、不注册任何业务路由。
func NewRouter(engine *gin.Engine, cfg *types.AppConfig, rds *redis.Client) *Router {
	userSecret, adminSecret := "", ""
	if cfg != nil {
		userSecret = cfg.Session.SecretKey
		adminSecret = cfg.AdminSession.SecretKey
	}
	return &Router{
		pub:      engine.Group("/"),
		user:     engine.Group("/", middleware.UserAuthMiddleware(userSecret, rds)),
		admin:    engine.Group(adminPrefix, middleware.AdminAuthMiddleware(adminSecret, rds)),
		adminPub: engine.Group(adminPrefix),
	}
}

// PublicGroup 公开路由。额外中间件在"建组时"一次传入，不存在事后补挂。
//
// 传入 /api/admin 前缀会 panic：管理后台接口只可能落在带 AdminAuthMiddleware 的链上，
// 这条不变量是白名单式的 —— 将来真要加一个公开的 admin 接口，必须改本文件，改动会出现在
// core 包的 diff 里被 reviewer 看见，而不是散落在 22 个 admin handler 的某一行。
func (r *Router) PublicGroup(prefix string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	if isAdminPrefix(prefix) {
		panic("core.Router: /api/admin 前缀必须使用 AdminGroup 或 AdminPublicGroup，PublicGroup 收到: " + prefix)
	}
	return r.pub.Group(prefix, handlers...)
}

// UserGroup 需要用户登录的路由。
func (r *Router) UserGroup(prefix string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	if isAdminPrefix(prefix) {
		panic("core.Router: /api/admin 前缀必须使用 AdminGroup 或 AdminPublicGroup，UserGroup 收到: " + prefix)
	}
	return r.user.Group(prefix, handlers...)
}

// AdminGroup 需要管理员登录的路由。
//
// prefix 按"完整路径"传入（与改造前 h.App.Engine.Group(P) 的 P 逐字相同），
// 由本方法剥掉 /api/admin 基址后挂到已带 AdminAuthMiddleware 的 admin 组上。
// 因此 `AdminGroup("/api/admin/redeem/")` 的 BasePath 就是 `/api/admin/redeem/`，
// 与改造前 `h.App.Engine.Group("/api/admin/redeem/")` 完全一致 —— 搬迁时路径零变化。
func (r *Router) AdminGroup(prefix string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return r.admin.Group(strings.TrimPrefix(prefix, adminPrefix), handlers...)
}

// AdminPublicGroup 是 /api/admin 下唯一的安全出口：不带参数、只此一个实例。
// 只服务于管理后台的 login / logout —— 这两个接口必须在"取得管理员身份之前"可调用。
func (r *Router) AdminPublicGroup() *gin.RouterGroup { return r.adminPub }

// StaticFiles 静态目录（当前只有 /update），走公开链，继承全局中间件。
//
// 这是唯一的静态出口，它挂在 pub 组下 —— 静态目录不需要鉴权，但也不会因此获得
// 注册 API 路由的能力（Static 只注册一条 `GET /*filepath`）。
func (r *Router) StaticFiles(prefix, dir string, handlers ...gin.HandlerFunc) {
	r.pub.Group("", handlers...).Static(prefix, dir)
}

// RouteRegistrar 所有 handler 的注册契约。
// 注册表在 main.go 的 registerRoutes 里显式列出并有数量断言。
type RouteRegistrar interface {
	RegisterRoutes(r *Router)
}

// isAdminPrefix 判断一个前缀是否落在管理后台命名空间内。
func isAdminPrefix(p string) bool {
	return p == adminPrefix || strings.HasPrefix(p, adminPrefix+"/")
}
