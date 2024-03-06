import { createRouter, createWebHistory } from 'vue-router'
import CustomLayout from '@/components/CustomLayout.vue'
import menu from './menu'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: CustomLayout,
      redirect: () => menu[0].path,
      children: menu
    },
    {
      path: "/:pathMatch(.*)*",
      name: "404",
      component: () => import("@/views/NotFound.vue"),
    },
  ]
})

export default router
