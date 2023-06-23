import {createRouter, createWebHistory} from 'vue-router'
import {createApp} from 'vue'
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import App from './App.vue'
import ChatPlus from "@/views/ChatPlus.vue";
import NotFound from './views/404.vue'
import Home from "@/views/Home.vue";
import Login from "@/views/Login.vue"
import Register from "@/views/Register.vue";
import {createPinia} from "pinia";

const routes = [
    {name: 'home', path: '/', component: Home, meta: {title: 'ChatGPT-Plus'}},
    {name: 'login', path: '/login', component: Login, meta: {title: '用户登录'}},
    {name: 'register', path: '/register', component: Register, meta: {title: '用户注册'}},
    {name: 'plus', path: '/chat', component: ChatPlus, meta: {title: 'ChatGPT-智能助手V3'}},
    {
        path: '/admin/login',
        name: 'admin-login',
        meta: {title: 'Chat-Plus 控制台登录'},
        component: () => import('@/views/admin/Login.vue'),
    },
    {
        name: 'admin',
        path: '/admin',
        redirect: '/admin/welcome',
        component: () => import("@/views/admin/Home.vue"),
        meta: {title: 'ChatGPT-Plus 管理后台'},
        children: [
            {
                path: '/admin/welcome',
                name: 'admin-home',
                meta: {title: '系统首页'},
                component: () => import('@/views/admin/Welcome.vue'),
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
    {name: 'test', path: '/test', component: () => import('@/views/Test.vue'), meta: {title: '测试页面'}},
    {name: 'NotFound', path: '/:all(.*)', component: NotFound, meta: {title: '页面没有找到'}},
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

const app = createApp(App)
app.use(createPinia())
app.use(router).use(ElementPlus).mount('#app')


