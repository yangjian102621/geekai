import {createRouter, createWebHistory} from 'vue-router'
import {createApp} from 'vue'
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import App from './App.vue'
import ChatPlus from "@/views/ChatPlus.vue";
import NotFound from './views/404.vue'
import TestPage from './views/Test.vue'
import Home from "@/views/Home.vue";
import Admin from "@/views/admin/Admin.vue";
import Login from "@/views/Login.vue"
import Register from "@/views/Register.vue";
import AdminLogin from "@/views/admin/Login.vue"

const routes = [
    {name: 'home', path: '/', component: Home, meta: {title: 'ChatGPT-Plus'}},
    {name: 'login', path: '/login', component: Login, meta: {title: '用户登录'}},
    {name: 'register', path: '/register', component: Register, meta: {title: '用户注册'}},
    {name: 'plus', path: '/chat', component: ChatPlus, meta: {title: 'ChatGPT-智能助手V3'}},
    {name: 'admin', path: '/admin', component: Admin, meta: {title: 'Chat-Plus 控制台'}},
    {name: 'admin-login', path: '/admin/login', component: AdminLogin, meta: {title: 'Chat-Plus 控制台登录'}},
    {name: 'test', path: '/test', component: TestPage, meta: {title: '测试页面'}},
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
        document.title = to.meta.title
    }
    next()
})

const app = createApp(App)
app.use(router).use(ElementPlus).mount('#app')


