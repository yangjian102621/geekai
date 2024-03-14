import http from "@/http/config";

export const getList = (data?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/order/list",
    method: "post",
    data
  })
}

export const remove = (data) => {
  return http({
    url: "/api/admin/order/remove",
    method: "post",
    data
  })
}