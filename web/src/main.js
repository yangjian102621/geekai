import {createRouter, createWebHistory} from 'vue-router'
import {createApp} from 'vue'
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import App from './App.vue'
import ChatPlus from "@/views/ChatPlus.vue";
import Chat from "@/views/Chat.vue";
import NotFound from './views/404.vue'
import TestPage from './views/Test.vue'
import Home from "@/views/Home.vue";
import './utils/prototype'

const routes = [
    {
        name: 'home', path: '/', component: Home, meta: {
            title: 'ChatGPT-Plus'
        }
    },
    {
        name: 'plus', path: '/plus', component: ChatPlus, meta: {
            title: 'ChatGPT-Plus for PC'
        }
    },
    {
        name: 'mobile', path: '/mobile', component: Chat, meta: {
            title: 'ChatGPT-Plus for Mobile'
        }
    },
    {
        name: 'test', path: '/test', component: TestPage, meta: {
            title: '测试页面'
        }
    },
    {
        name: 'NotFound', path: '/:all(.*)', component: NotFound, meta: {
            title: '页面没有找到'
        }
    },
]

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


