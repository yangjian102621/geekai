import http from "@/http/config";

export const getConfig = (params) => {
  return http({
    url: "/api/admin/config/get",
    method: "get",
    params
  })
}

export const save = (data) => {
  return http({
    url: "/api/admin/config/update",
    method: "post",
    data
  })
}

export const modelList = (params) => {
  return http({
    url: "/api/admin/model/list",
    method: "get",
    params
  })
}
