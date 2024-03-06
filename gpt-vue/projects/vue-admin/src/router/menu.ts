import {
  IconUser,
  IconDashboard
} from "@arco-design/web-vue/es/icon";

const menu = [
  {
    path: '/dashboard',
    name: 'Dashboard',
    meta: {
      title: "仪表盘",
      icon: IconDashboard
    },
    component: () => import('@/views/DashboardView.vue')
  },
  {
    path: '/user',
    name: 'User',
    meta: {
      title: "用户管理",
      icon: IconUser,
    },
    component: () => import('@/views/User/UserContainer.vue')
  },
];

export default menu;
