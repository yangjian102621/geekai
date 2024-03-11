import http from "@/http/config";

export const getList = (params) => {
  return http({
    url: "/api/admin/function/list",
    method: "get",
    params
  })
}

export const save = (data) => {
  return http({
    url: "/api/admin/function/save",
    method: "post",
    data
  })
}

export const remove = (params) => {
  return http({
    url: "/api/admin/function/remove",
    method: "get",
    params
  })
}

export const setStatus = (data) => {
  return http({
    url: "/api/admin/function/set",
    method: "post",
    data
  })
}

export const token = (params) => {
  return http({
    url: "/api/admin/function/token",
    method: "get",
    params
  })
}
