import {
  IconUser,
  IconDashboard,
  IconOrderedList,
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
  {
    path: '/order',
    name: 'Order',
    meta: {
      title: "充值订单",
      icon: IconOrderedList,
    },
    component: () => import('@/views/Order/OrderContainer.vue')
  },
];

export default menu;
