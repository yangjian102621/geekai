import http from "@/http/config";

export const getList = (params) => {
  return http({
    url: "/api/admin/sysPermission/list",
    method: "get",
    params
  })
}

export const save = (data) => {
  return http({
    url: "/api/admin/sysPermission/save",
    method: "post",
    data
  })
}

export const remove = (data) => {
  return http({
    url: "/api/admin/sysPermission/remove",
    method: "post",
    data
  })
}