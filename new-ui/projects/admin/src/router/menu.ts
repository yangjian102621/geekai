import {
  IconUser,
  IconDashboard,
  IconOrderedList,
  IconCalendar,
  IconHeartFill,
  IconCodeSquare,
  IconMessage,
  IconSettings,
  IconUserGroup,
  IconLock,
  IconCodepen,
  IconWechatpay,
  IconRobot,
} from "@arco-design/web-vue/es/icon";

const menu = [
  {
    path: "/dashboard",
    name: "Dashboard",
    meta: {
      title: "仪表盘",
      icon: IconDashboard,
    },
    component: () => import("@/views/DashboardView.vue"),
  },
  {
    path: "/user",
    name: "User",
    meta: {
      title: "用户管理",
      icon: IconUser,
    },
    component: () => import("@/views/User/UserContainer.vue"),
  },
  {
    path: "/role",
    name: "Role",
    meta: {
      title: "角色管理",
      icon: IconUserGroup,
    },
    component: () => import("@/views/Role/RoleContainer.vue"),
  },
  {
    path: "/chatModel",
    name: "ChatModel",
    meta: {
      title: "语言模型",
      icon: IconCodepen,
    },
    component: () => import("@/views/ChatModel/ChatModelContainer.vue"),
  },
  {
    path: "/product",
    name: "Product",
    meta: {
      title: "充值产品",
      icon: IconWechatpay,
    },
    component: () => import("@/views/Product/ProductContainer.vue"),
  },
  {
    path: "/apiKey",
    name: "ApiKey",
    meta: {
      title: "APIKEY",
      icon: IconLock,
    },
    component: () => import("@/views/ApiKey/ApiKeyContainer.vue"),
  },
  {
    path: "/order",
    name: "Order",
    meta: {
      title: "充值订单",
      icon: IconOrderedList,
    },
    component: () => import("@/views/Order/OrderContainer.vue"),
  },

  {
    path: "/reward",
    name: "Reward",
    meta: {
      title: "众筹管理",
      icon: IconHeartFill,
    },
    component: () => import("@/views/Reward/RewardContainer.vue"),
  },
  {
    path: "/functions",
    name: "Functions",
    meta: {
      title: "函数管理",
      icon: IconCodeSquare,
    },
    component: () => import("@/views/Functions/FunctionsContainer.vue"),
  },
  {
    path: "/chats",
    name: "Chats",
    meta: {
      title: "对话管理",
      icon: IconMessage,
    },
    component: () => import("@/views/Chats/ChatsContainer.vue"),
  },
  {
    path: "/system",
    name: "System",
    meta: {
      title: "系统设置",
      icon: IconSettings,
    },
    component: () => import("@/views/System/SystemContainer.vue"),
  },
  {
    path: "/loginLog",
    name: "LoginLog",
    meta: {
      title: "登录日志",
      icon: IconCalendar,
    },
    component: () => import("@/views/LoginLog.vue"),
  },
  {
    path: "/sysAdmin",
    name: "SysAdmin",
    meta: {
      title: "系统管理员",
      icon: IconRobot,
    },
    component: () => import("@/views/SysAdmin/SysAdminContainer.vue"),
  },
];

export default menu;
