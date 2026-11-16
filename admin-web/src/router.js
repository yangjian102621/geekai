import { createRouter, createWebHistory } from 'vue-router'

const routes = [
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
        meta: { title: '积分配置' },
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
        path: '/admin/library',
        name: 'admin-library',
        meta: { title: '公共资料库' },
        component: () => import('@/views/admin/Library.vue'),
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
        path: '/admin/powerLog',
        name: 'admin-power-log',
        meta: { title: '积分日志' },
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
    path: '/login',
    meta: { title: '用户登录' },
    component: () => import('@/views/Login.vue'),
  },
  {
    path: '/register',
    meta: { title: '用户注册' },
    component: () => import('@/views/Login.vue'),
  },
  {
    path: '/:all(.*)',
    meta: { title: '页面没有找到' },
    component: () => import('@/views/404.vue'),
  },
]

const router = createRouter({ history: createWebHistory(), routes })
let prevRoute = null

router.beforeEach((to, from) => {
  if (!to.path.startsWith('/admin') && !['/login', '/register'].includes(to.path)) {
    window.location.assign(to.fullPath)
    return false
  }
  document.title = to.meta.title || 'GeekAI 管理后台'
  prevRoute = from
})

export { prevRoute, router }
