import http from "@/http/config";

export const getList = (params) => {
  return http({
    url: "/api/admin/sysUser/list",
    method: "get",
    params
  })
}

export const save = (data) => {
  return http({
    url: "/api/admin/sysUser/save",
    method: "post",
    data
  })
}

export const remove = (data) => {
  return http({
    url: "/api/admin/sysUser/remove",
    method: "post",
    data
  })
}

export const resetPass = (data) => {
  return http({
    url: "/api/admin/sysUser/resetPass",
    method: "post",
    data
  })
}