import { useAuthStore } from "@/stores/auth";

// 判断操作权限
export function hasPermission(permissionTag: string | string[] | boolean) {
  const authStore = useAuthStore();
  const { is_super_admin, permissions = [] } = authStore;
  if (is_super_admin) {
    return true;
  }
  if (Array.isArray(permissionTag)) {
    return permissionTag.every((tag) => permissions.includes(tag));
  }
  if (typeof permissionTag === "string") {
    return permissions.includes(permissionTag);
  }
  return permissionTag;
}

function checkPermission(el, binding) {
  if (!hasPermission(binding.value)) {
    el.parentNode && el.parentNode.removeChild(el);
  }
}

export const permission = {
  mounted(el, binding) {
    checkPermission(el, binding);
  },
  updated(el, binding) {
    checkPermission(el, binding);
  },
};
