import http from "@/http/config";

export const getList = (params) => {
  return http({
    url: "/api/admin/sysRole/list",
    method: "get",
    params
  })
}

export const save = (data) => {
  return http({
    url: "/api/admin/sysRole/save",
    method: "post",
    data
  })
}

export const remove = (data) => {
  return http({
    url: "/api/admin/sysRole/remove",
    method: "post",
    data
  })
}