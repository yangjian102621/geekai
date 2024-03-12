import http from "@/http/config";

export const getList = (params?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/role/list",
    method: "get",
    params,
  });
};
export const save = (data?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/role/save",
    method: "post",
    data,
  });
};
export const deleting = (id: string | number) => {
  return http({
    url: `/api/admin/role/remove?id=${id}`,
    method: "get",
  });
};
export const setStatus = (data) => {
  return http({
    url: `/api/admin/role/set`,
    method: "post",
    data,
  });
};
