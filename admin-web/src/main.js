import '@/assets/css/tailwind.css'
import '@/assets/css/custom-scroll.scss'
import '@/assets/iconfont/iconfont.css'
import '@/assets/css/common.scss'
import '@/assets/css/theme-dark.scss'
import '@/assets/css/theme-light.scss'
import 'element-plus/dist/index.css'
import 'vant/lib/index.css'
import ElementPlus from 'element-plus'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import { router } from '@/router'
import { useSharedStore } from '@/store/sharedata'
import App from './App.vue'

const pinia = createPinia()
const store = useSharedStore(pinia)
store.setTheme(store.theme)
createApp(App).use(pinia).use(router).use(ElementPlus).mount('#app')
