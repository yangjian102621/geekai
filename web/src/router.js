// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import { createRouter, createWebHistory } from 'vue-router'
import { getSystemInfo, getMenus } from '@/store/cache'

// PC 端主页面路由
const homeRoutes = {
  name: 'home',
  path: '/home',
  redirect: '/chat',
  component: () => import('@/views/Home.vue'),
  children: [
    {
      name: 'chat',
      path: '/chat',
      meta: { title: '创作中心' },
      component: () => import('@/views/ChatPlus.vue'),
    },
    {
      name: 'chat-id',
      path: '/chat/:id',
      meta: { title: '创作中心' },
      component: () => import('@/views/ChatPlus.vue'),
    },
    {
      name: 'image-mj',
      path: '/mj',
      meta: { title: 'MidJourney 绘画中心' },
      component: () => import('@/views/ImageMj.vue'),
    },
    {
      name: 'member',
      path: '/member',
      meta: { title: '会员充值中心' },
      component: () => import('@/views/Member.vue'),
    },
    {
      name: 'chat-app',
      path: '/apps',
      meta: { title: '应用中心' },
      component: () => import('@/views/ChatApps.vue'),
    },
    {
      name: 'images',
      path: '/gallery',
      meta: { title: '作品展示' },
      component: () => import('@/views/Gallery.vue'),
    },
    {
      name: 'user-invitation',
      path: '/invite',
      meta: { title: '推广计划' },
      component: () => import('@/views/Invitation.vue'),
    },
    {
      name: 'powerLog',
      path: '/powerLog',
      meta: { title: '消费日志' },
      component: () => import('@/views/PowerLog.vue'),
    },
    {
      name: 'xmind',
      path: '/xmind',
      meta: { title: '思维导图' },
      component: () => import('@/views/MarkMap.vue'),
    },
    {
      name: 'image',
      path: '/image',
      meta: { title: 'AI图像生成' },
      component: () => import('@/views/Image.vue'),
    },
    {
      name: 'ppt',
      path: '/ppt',
      meta: { title: 'PPT 生成' },
      component: () => import('@/views/PPTCreate.vue'),
    },
    {
      name: 'suno',
      path: '/suno',
      meta: { title: 'Suno音乐创作' },
      component: () => import('@/views/Suno.vue'),
    },
    {
      name: 'ExternalLink',
      path: '/external',
      component: () => import('@/views/ExternalPage.vue'),
    },
    {
      name: 'song',
      path: '/song/:id',
      meta: { title: 'Suno音乐播放' },
      component: () => import('@/views/Song.vue'),
    },
    {
      name: 'video',
      path: '/video',
      meta: { title: '视频创作中心' },
      component: () => import('@/views/Video.vue'),
    },
    {
      name: 'jimeng',
      path: '/jimeng',
      meta: { title: '即梦AI' },
      component: () => import('@/views/Jimeng.vue'),
    },
  ],
}

const routes = [
  {
    name: 'Index',
    path: '/',
    meta: { title: '首页' },
    component: () => import('@/views/Index.vue'),
  },
  {
    name: 'chat-export',
    path: '/chat/export',
    meta: { title: '导出会话记录' },
    component: () => import('@/views/ChatExport.vue'),
  },

  {
    name: 'login',
    path: '/login',
    meta: { title: '用户登录' },
    component: () => import('@/views/Login.vue'),
  },
  {
    name: 'register',
    path: '/register',
    meta: { title: '用户注册' },
    component: () => import('@/views/Login.vue'),
  },
  {
    path: '/admin/login',
    name: 'admin-login',
    meta: { title: '控制台登录' },
    component: () => import('@/views/admin/Login.vue'),
  },

  {
    name: 'admin',
    path: '/admin',
    redirect: '/admin/dashboard',
    component: () => import('@/views/admin/Home.vue'),
    meta: { title: 'Geek-AI 控制台' },
    children: [
      {
        path: '/admin/dashboard',
        name: 'admin-dashboard',
        meta: { title: '仪表盘' },
        component: () => import('@/views/admin/Dashboard.vue'),
      },
      {
        path: '/admin/config/basic',
        name: 'admin-config-basic',
        meta: { title: '基础配置' },
        component: () => import('@/views/admin/settings/BasicConfig.vue'),
      },
      {
        path: '/admin/config/power',
        name: 'admin-config-power',
        meta: { title: '算力配置' },
        component: () => import('@/views/admin/settings/PowerConfig.vue'),
      },
      {
        path: '/admin/config/payment',
        name: 'admin-config-payment',
        meta: { title: '支付配置' },
        component: () => import('@/views/admin/settings/PaymentConfig.vue'),
      },
      {
        path: '/admin/config/storage',
        name: 'admin-config-storage',
        meta: { title: '存储配置' },
        component: () => import('@/views/admin/settings/StorageConfig.vue'),
      },
      {
        path: '/admin/config/sms',
        name: 'admin-config-sms',
        meta: { title: '短信配置' },
        component: () => import('@/views/admin/settings/SmsConfig.vue'),
      },
      {
        path: '/admin/config/smtp',
        name: 'admin-config-smtp',
        meta: { title: '邮件配置' },
        component: () => import('@/views/admin/settings/SmtpConfig.vue'),
      },
      {
        path: '/admin/config/plugin',
        name: 'admin-config-plugin',
        meta: { title: '插件配置' },
        component: () => import('@/views/admin/settings/PluginConfig.vue'),
      },
      {
        path: '/admin/moderation/config',
        name: 'admin-config-moderation',
        meta: { title: '文本审查配置' },
        component: () => import('@/views/admin/moderation/ModerationConfig.vue'),
      },
      {
        path: '/admin/moderation/list',
        name: 'admin-moderation-list',
        meta: { title: '文本审核记录' },
        component: () => import('@/views/admin/moderation/ModerationList.vue'),
      },
      {
        path: '/admin/config/markmap',
        name: 'admin-config-markmap',
        meta: { title: '思维导图配置' },
        component: () => import('@/views/admin/settings/MarkMapConfig.vue'),
      },
      {
        path: '/admin/config/notice',
        name: 'admin-config-notice',
        meta: { title: '公告配置' },
        component: () => import('@/views/admin/settings/NoticeConfig.vue'),
      },
      {
        path: '/admin/config/agreement',
        name: 'admin-config-agreement',
        meta: { title: '用户协议' },
        component: () => import('@/views/admin/settings/AgreementConfig.vue'),
      },
      {
        path: '/admin/config/privacy',
        name: 'admin-config-privacy',
        meta: { title: '隐私声明' },
        component: () => import('@/views/admin/settings/PrivacyConfig.vue'),
      },
      {
        path: '/admin/config/menu',
        name: 'admin-config-menu',
        meta: { title: '菜单配置' },
        component: () => import('@/views/admin/settings/MenuConfig.vue'),
      },
      {
        path: '/admin/user',
        name: 'admin-user',
        meta: { title: '用户管理' },
        component: () => import('@/views/admin/Users.vue'),
      },
      {
        path: '/admin/app',
        name: 'admin-app',
        meta: { title: '应用列表' },
        component: () => import('@/views/admin/Apps.vue'),
      },
      {
        path: '/admin/app/type',
        name: 'admin-app-type',
        meta: { title: '应用分类' },
        component: () => import('@/views/admin/AppType.vue'),
      },
      {
        path: '/admin/apikey',
        name: 'admin-apikey',
        meta: { title: 'API-KEY 管理' },
        component: () => import('@/views/admin/ApiKey.vue'),
      },
      {
        path: '/admin/chat/model',
        name: 'admin-chat-model',
        meta: { title: '语言模型' },
        component: () => import('@/views/admin/ChatModel.vue'),
      },
      {
        path: '/admin/product',
        name: 'admin-product',
        meta: { title: '充值产品' },
        component: () => import('@/views/admin/Product.vue'),
      },
      {
        path: '/admin/order',
        name: 'admin-order',
        meta: { title: '充值订单' },
        component: () => import('@/views/admin/Order.vue'),
      },
      {
        path: '/admin/redeem',
        name: 'admin-redeem',
        meta: { title: '兑换码管理' },
        component: () => import('@/views/admin/Redeem.vue'),
      },
      {
        path: '/admin/loginLog',
        name: 'admin-loginLog',
        meta: { title: '登录日志' },
        component: () => import('@/views/admin/LoginLog.vue'),
      },
      {
        path: '/admin/functions',
        name: 'admin-functions',
        meta: { title: '函数管理' },
        component: () => import('@/views/admin/Functions.vue'),
      },
      {
        path: '/admin/chats',
        name: 'admin-chats',
        meta: { title: '对话管理' },
        component: () => import('@/views/admin/records/ChatList.vue'),
      },
      {
        path: '/admin/images',
        name: 'admin-images',
        meta: { title: '绘图管理' },
        component: () => import('@/views/admin/records/ImageList.vue'),
      },
      {
        path: '/admin/records/suno',
        name: 'admin-suno',
        meta: { title: 'Suno音乐管理' },
        component: () => import('@/views/admin/records/Suno.vue'),
      },
      {
        path: '/admin/records/videos',
        name: 'admin-videos',
        meta: { title: '视频任务管理' },
        component: () => import('@/views/admin/records/Videos.vue'),
      },
      {
        path: '/admin/jimeng/jobs',
        name: 'admin-jimeng-jobs',
        meta: { title: '即梦AI任务' },
        component: () => import('@/views/admin/jimeng/JimengJobs.vue'),
      },
      {
        path: '/admin/jimeng/config',
        name: 'admin-jimeng-config',
        meta: { title: '即梦设置' },
        component: () => import('@/views/admin/jimeng/JimengConfig.vue'),
      },
      {
        path: '/admin/video/config',
        name: 'admin-video-config',
        meta: { title: '视频生成配置' },
        component: () => import('@/views/admin/video/VideoConfig.vue'),
      },
      {
        path: '/admin/ppt/config',
        name: 'admin-ppt-config',
        meta: { title: 'PPT 生成配置' },
        component: () => import('@/views/admin/settings/PPTConfig.vue'),
      },
      {
        path: '/admin/ppt/jobs',
        name: 'admin-ppt-jobs',
        meta: { title: 'PPT 任务列表' },
        component: () => import('@/views/admin/ppt/PPTJobs.vue'),
      },
      {
        path: '/admin/powerLog',
        name: 'admin-power-log',
        meta: { title: '算力日志' },
        component: () => import('@/views/admin/PowerLog.vue'),
      },
      {
        path: '/admin/manger',
        name: 'admin-manger',
        meta: { title: '管理员' },
        component: () => import('@/views/admin/Manager.vue'),
      },
      {
        path: '/admin/config/wechat',
        name: 'admin-wechat-config',
        meta: { title: '微信配置' },
        component: () => import('@/views/admin/settings/WechatConfig.vue'),
      },
    ],
  },

  {
    name: 'test',
    path: '/test',
    meta: { title: '测试页面' },
    component: () => import('@/views/test/Test.vue'),
  },

  {
    name: 'NotFound',
    path: '/:all(.*)',
    meta: { title: '页面没有找到' },
    component: () => import('@/views/404.vue'),
  },
]

const mobileRoutes = {
  name: 'mobile',
  path: '/mobile',
  meta: { title: '首页' },
  component: () => import('@/views/mobile/Home.vue'),
  redirect: '/mobile/index',
  children: [
    {
      path: '/mobile/index',
      name: 'mobile-index',
      component: () => import('@/views/mobile/Index.vue'),
    },
    {
      meta: { title: 'AI对话' },
      path: '/mobile/chat',
      name: 'mobile-chat',
      component: () => import('@/views/mobile/ChatList.vue'),
    },
    {
      meta: { title: '创作中心' },
      path: '/mobile/create',
      name: 'mobile-create',
      component: () => import('@/views/mobile/Create.vue'),
    },
    {
      meta: { title: '发现' },
      path: '/mobile/discover',
      name: 'mobile-discover',
      component: () => import('@/views/mobile/Discover.vue'),
    },
    {
      meta: { title: '个人中心' },
      path: '/mobile/profile',
      name: 'mobile-profile',
      component: () => import('@/views/mobile/Profile.vue'),
    },
    {
      meta: { title: '会员充值' },
      path: '/mobile/member',
      name: 'mobile-member',
      component: () => import('@/views/mobile/Member.vue'),
    },
    {
      meta: { title: '作品展示' },
      path: '/mobile/gallery',
      name: 'mobile-gallery',
      component: () => import('@/views/mobile/pages/Gallery.vue'),
    },
    {
      path: '/mobile/chat/session',
      name: 'mobile-chat-session',
      component: () => import('@/views/mobile/ChatSession.vue'),
    },

    {
      meta: { title: '应用中心' },
      path: '/mobile/apps',
      name: 'mobile-apps',
      component: () => import('@/views/mobile/Apps.vue'),
    },
    // 新增的功能页面路由
    {
      meta: { title: '消费日志' },
      path: '/mobile/power-log',
      name: 'mobile-power-log',
      component: () => import('@/views/mobile/PowerLog.vue'),
    },
    {
      meta: { title: '推广计划' },
      path: '/mobile/invite',
      name: 'mobile-invite',
      component: () => import('@/views/mobile/Invite.vue'),
    },
    {
      meta: { title: '设置' },
      path: '/mobile/settings',
      name: 'mobile-settings',
      component: () => import('@/views/mobile/Settings.vue'),
    },
    {
      meta: { title: 'Suno音乐创作' },
      path: '/mobile/suno',
      name: 'mobile-suno',
      component: () => import('@/views/mobile/SunoCreate.vue'),
    },
    {
      meta: { title: '视频生成' },
      path: '/mobile/video',
      name: 'mobile-video',
      component: () => import('@/views/mobile/VideoCreate.vue'),
    },
    {
      meta: { title: '即梦AI' },
      path: '/mobile/jimeng',
      name: 'mobile-jimeng',
      component: () => import('@/views/mobile/JimengCreate.vue'),
    },
  ],
}

// console.log(MY_VARIABLE)
const router = createRouter({
  history: createWebHistory(),
  routes: routes,
})

let prevRoute = null
// dynamic change the title when router change
router.beforeEach((to, from, next) => {
  document.title = to.meta.title
  prevRoute = from
  next()
})

// 检测是否启用了移动端（同步：顶层 await）
const res = await getSystemInfo()
const data = res.data
if (data && data.enable_mobile_site) {
  router.addRoute(mobileRoutes)
}

// 获取所有的菜单列表，然后把禁用的菜单从router中移除
const menus = await getMenus()
homeRoutes.children = homeRoutes.children.filter((route) => {
  return !menus[route.path] || menus[route.path]?.enabled
})
// 添加主页面路由
router.addRoute(homeRoutes)

export { prevRoute, router, mobileRoutes }
