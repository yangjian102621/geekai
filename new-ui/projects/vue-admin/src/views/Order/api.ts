import http from "@/http/config";

export const getList = (data?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/order/list",
    method: "post",
    data
  })
}