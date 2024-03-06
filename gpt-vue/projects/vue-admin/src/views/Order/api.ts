import http from "@/http/config";

export const getList = (params?: Record<string, unknown>) => {
  return http({
    url: "/admin/order/list",
    methods: "get",
    params
  })
}