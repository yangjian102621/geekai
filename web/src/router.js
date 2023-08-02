import {createRouter, createWebHistory} from "vue-router";

const routes = [
    {
        name: 'home',
        path: '/',
        meta: {title: 'ChatGPT-Plus'},
        component: () => import('@/views/Home.vue'),
    },
    {
        name: 'login',
        path: '/login',
        meta: {title: '用户登录'},
        component: () => import('@/views/Login.vue'),
    },
    {
        name: 'register',
        path: '/register',

        meta: {title: '用户注册'},
        component: () => import('@/views/Register.vue'),
    },
    {
        name: 'plus',
        path: '/chat',
        meta: {title: 'ChatGPT-智能助手V3'},
        component: () => import('@/views/ChatPlus.vue'),
    },
    {
        path: '/admin/login',
        name: 'admin-login',
        meta: {title: 'Chat-Plus 控制台登录'},
        component: () => import('@/views/admin/Login.vue'),
    },
    {
        name: 'admin',
        path: '/admin',
        redirect: '/admin/dashboard',
        component: () => import("@/views/admin/Home.vue"),
        meta: {title: 'ChatGPT-Plus 管理后台'},
        children: [
            {
                path: '/admin/dashboard',
                name: 'admin-dashboard',
                meta: {title: '仪表盘'},
                component: () => import('@/views/admin/Dashboard.vue'),
            },
            {
                path: '/admin/system',
                name: 'admin-system',
                meta: {title: '系统设置'},
                component: () => import('@/views/admin/SysConfig.vue'),
            },
            {
                path: '/admin/user',
                name: 'admin-user',
                meta: {title: '用户管理'},
                component: () => import('@/views/admin/UserList.vue'),
            },
            {
                path: '/admin/role',
                name: 'admin-role',
                meta: {title: '角色管理'},
                component: () => import('@/views/admin/RoleList.vue'),
            },
            {
                path: '/admin/apikey',
                name: 'admin-apikey',
                meta: {title: 'API-KEY 管理'},
                component: () => import('@/views/admin/ApiKey.vue'),
            },
            {
                path: '/admin/reward',
                name: 'admin-reward',
                meta: {title: '众筹管理'},
                component: () => import('@/views/admin/RewardList.vue'),
            },
            {
                path: '/admin/loginLog',
                name: 'admin-loginLog',
                meta: {title: '登录日志'},
                component: () => import('@/views/admin/LoginLog.vue'),
            },
            {
                path: '/admin/demo/form',
                name: 'admin-form',
                meta: {title: '表单页面'},
                component: () => import('@/views/admin/demo/Form.vue'),
            },
            {
                path: '/admin/demo/table',
                name: 'admin-table',
                meta: {title: '数据列表'},
                component: () => import('@/views/admin/demo/Table.vue'),
            },
            {
                path: '/admin/demo/import',
                name: 'admin-import',
                meta: {title: '导入数据'},
                component: () => import('@/views/admin/demo/Import.vue'),
            },
            {
                path: '/admin/demo/editor',
                name: 'admin-editor',
                meta: {title: '富文本编辑器'},
                component: () => import('@/views/admin/demo/Editor.vue'),
            },
        ]
    },

    {
        path: '/mobile/chat/session',
        name: 'mobile-chat-session',
        component: () => import('@/views/mobile/ChatSession.vue'),
    },
    {
        name: 'mobile',
        path: '/mobile',
        meta: {title: 'ChatGPT-智能助手V3'},
        component: () => import('@/views/mobile/Home.vue'),
        redirect: '/mobile/chat/list',
        children: [
            {
                path: '/mobile/chat/list',
                name: 'mobile-chat-list',
                component: () => import('@/views/mobile/ChatList.vue'),
            },
            {
                path: '/mobile/setting',
                name: 'mobile-setting',
                component: () => import('@/views/mobile/Setting.vue'),
            },
            {
                path: '/mobile/profile',
                name: 'mobile-profile',
                component: () => import('@/views/mobile/Profile.vue'),
            },
        ]
    },
    {
        name: 'test',
        path: '/test',
        meta: {title: '测试页面'},
        component: () => import('@/views/Test.vue'),
    },
    {
        name: 'NotFound',
        path: '/:all(.*)',
        meta: {title: '页面没有找到'},
        component: () => import('@/views/404.vue'),
    },
]

// console.log(MY_VARIABLE)
const router = createRouter({
    history: createWebHistory(),
    routes: routes,
})

// dynamic change the title when router change
router.beforeEach((to, from, next) => {
    if (to.meta.title) {
        document.title = `${to.meta.title} | ChatGPT-PLUS`
    }
    next()
})

export default router;