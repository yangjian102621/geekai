import http from "@/http/config";

export const getList = (params?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/reward/list",
    method: "get",
    params
  })
}

export const remove = (data?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/reward/remove",
    method: "post",
    data
  })
}
