import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from "@/stores/auth";
import CustomLayout from '@/components/CustomLayout.vue'
import { hasPermission } from "@/directives/permission";
import menu from './menu'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    permission?: string
  }
}

const whiteListRoutes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/LoginView.vue"),
  },
  {
    path: "/403",
    name: "403",
    component: () => import("@/views/NoPermission.vue"),
  },
  {
    path: "/:pathMatch(.*)*",
    name: "404",
    component: () => import("@/views/NotFound.vue"),
  },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: CustomLayout,
      redirect: () => menu[0].path,
      children: menu
    },
    ...whiteListRoutes
  ]
})

const whiteList = whiteListRoutes.map((i) => i.name);

router.beforeEach((to, _, next) => {
  const authStore = useAuthStore();
  authStore.init()
  if (typeof to.name === "string" && whiteList.includes(to.name)) {
    if (authStore.token && to.name === "Login") {
      next({ path: menu[0].path });
      return;
    }
    next();
    return;
  }
  if (!authStore.token) {
    authStore.$reset();
    next({ name: "Login" });
    return;
  }
  if (to.meta.permission) {
    next(!hasPermission(to.meta.permission) ? { name: "403" } : undefined);
  }
  next();
});

export default router
