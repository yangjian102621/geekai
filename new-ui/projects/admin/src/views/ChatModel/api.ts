import http from "@/http/config";

export const getList = (params?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/model/list",
    method: "get",
    params,
  });
};
export const save = (data?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/model/save",
    method: "post",
    data,
  });
};
export const deleting = (id: string | number) => {
  return http({
    url: `/api/admin/model/remove`,
    method: "post",
    data: { id },
  });
};
export const setStatus = (data) => {
  return http({
    url: `/api/admin/model/set`,
    method: "post",
    data,
  });
};
