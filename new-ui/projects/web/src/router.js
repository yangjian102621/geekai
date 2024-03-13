import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    name: "home",
    path: "/",
    redirect: "/chat",
    meta: { title: "首页" },
    component: () => import("@/views/Home.vue"),
    children: [
      {
        name: "chat",
        path: "/chat",
        meta: { title: "创作中心" },
        component: () => import("@/views/ChatPlus.vue"),
      },
      {
        name: "image-mj",
        path: "/mj",
        meta: { title: "MidJourney 绘画中心" },
        component: () => import("@/views/ImageMj.vue"),
      },
      {
        name: "image-sd",
        path: "/sd",
        meta: { title: "stable diffusion 绘画中心" },
        component: () => import("@/views/ImageSd.vue"),
      },
      {
        name: "member",
        path: "/member",
        meta: { title: "会员充值中心" },
        component: () => import("@/views/Member.vue"),
      },
      {
        name: "chat-role",
        path: "/apps",
        meta: { title: "应用中心" },
        component: () => import("@/views/ChatApps.vue"),
      },
      {
        name: "images",
        path: "/images-wall",
        meta: { title: "作品展示" },
        component: () => import("@/views/ImagesWall.vue"),
      },
      {
        name: "user-invitation",
        path: "/invite",
        meta: { title: "推广计划" },
        component: () => import("@/views/Invitation.vue"),
      },
      {
        name: "knowledge",
        path: "/knowledge",
        meta: { title: "我的知识库" },
        component: () => import("@/views/Knowledge.vue"),
      },
    ],
  },
  {
    name: "chat-export",
    path: "/chat/export",
    meta: { title: "导出会话记录" },
    component: () => import("@/views/ChatExport.vue"),
  },
  {
    name: "login",
    path: "/login",
    meta: { title: "用户登录" },
    component: () => import("@/views/Login.vue"),
  },
  {
    name: "register",
    path: "/register",

    meta: { title: "用户注册" },
    component: () => import("@/views/Register.vue"),
  },

  {
    name: "mobile",
    path: "/mobile",
    meta: { title: "ChatPlus-智能助手V3" },
    component: () => import("@/views/mobile/Home.vue"),
    redirect: "/mobile/chat",
    children: [
      {
        path: "/mobile/chat",
        name: "mobile-chat",
        component: () => import("@/views/mobile/ChatList.vue"),
      },
      {
        path: "/mobile/mj",
        name: "mobile-mj",
        component: () => import("@/views/mobile/ImageMj.vue"),
      },
      {
        path: "/mobile/profile",
        name: "mobile-profile",
        component: () => import("@/views/mobile/Profile.vue"),
      },
      {
        path: "/mobile/img-wall",
        name: "mobile-img-wall",
        component: () => import("@/views/mobile/ImgWall.vue"),
      },
    ],
  },
  {
    path: "/mobile/chat/session",
    name: "mobile-chat-session",
    component: () => import("@/views/mobile/ChatSession.vue"),
  },
  {
    path: "/mobile/chat/export",
    name: "mobile-chat-export",
    component: () => import("@/views/mobile/ChatExport.vue"),
  },

  {
    name: "test",
    path: "/test",
    meta: { title: "测试页面" },
    component: () => import("@/views/Test.vue"),
  },
  {
    name: "NotFound",
    path: "/:all(.*)",
    meta: { title: "页面没有找到" },
    component: () => import("@/views/404.vue"),
  },
];

// console.log(MY_VARIABLE)
const router = createRouter({
  history: createWebHistory(),
  routes: routes,
});

let prevRoute = null;
// dynamic change the title when router change
router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} | ${process.env.VUE_APP_TITLE}`;
  }
  prevRoute = from;
  next();
});

export { router, prevRoute };
