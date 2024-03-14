import type { RouteRecordRaw } from "vue-router";
const system: RouteRecordRaw[] = [
  {
    path: "admin",
    name: "SysAdmin",
    meta: {
      title: "系统管理员",
    },
    component: () => import("@/views/SysAdmin/SysAdminContainer.vue"),
  },
  {
    path: "permission",
    name: "SysPermission",
    meta: {
      title: "权限配置",
    },
    component: () => import("@/views/SysPermission/SysPermissionContainer.vue"),
  },
  {
    path: "role",
    name: "SysRole",
    meta: {
      title: "角色管理",
    },
    component: () => import("@/views/SysRole/SysRoleContainer.vue"),
  },
  {
    path: "loginLog",
    name: "LoginLog",
    meta: {
      title: "登录日志",
    },
    component: () => import("@/views/LoginLog.vue"),
  },
]

export default system