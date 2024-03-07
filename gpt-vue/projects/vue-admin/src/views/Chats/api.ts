import http from "@/http/config";

export const getList = (data) => {
  return http({
    url: "/api/admin/chat/list",
    method: "post",
    data
  })
}

export const message = (data) => {
  return http({
    url: "/api/admin/chat/message",
    method: "post",
    data
  })
}

export const history = (params) => {
  return http({
    url: "/api/admin/chat/history",
    method: "get",
    params
  })
}

export const remove = (params) => {
  return http({
    url: "/api/admin/chat/remove",
    method: "get",
    params
  })
}

