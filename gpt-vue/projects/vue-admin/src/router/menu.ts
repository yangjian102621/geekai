import {
  IconUser,
  IconDashboard,
  IconOrderedList,
  IconCalendar,
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
  {
    path: '/reward',
    name: 'Reward',
    meta: {
      title: "众筹管理",
      icon: IconCalendar,
    },
    component: () => import('@/views/Reward/RewardContainer.vue')
  },
  {
    path: '/functions',
    name: 'Functions',
    meta: {
      title: "函数管理",
      icon: IconCalendar,
    },
    component: () => import('@/views/Functions/FunctionsContainer.vue')
  },
  {
    path: '/loginLog',
    name: 'LoginLog',
    meta: {
      title: "登录日志",
      icon: IconCalendar,
    },
    component: () => import('@/views/LoginLog.vue')
  },
];

export default menu;
