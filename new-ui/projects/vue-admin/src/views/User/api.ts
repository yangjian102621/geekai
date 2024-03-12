import http from "@/http/config";

export const getList = (params?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/user/list",
    method: "get",
    params,
  });
};

export const save = (data?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/user/save",
    method: "post",
    data,
  });
};

export const deletApi = (id: string | number) => {
  return http({
    url: `/api/admin/user/remove?id=${id}`,
    method: "get",
  });
};

export const getRole = () => {
  return http({
    url: `/api/admin/role/list`,
    method: "get",
  });
};

export const getModel = () => {
  return http({
    url: `/api/admin/model/list`,
    method: "get",
  });
};

export const resetPassword = (data) => {
  return http({
    url: `/api/admin/user/resetPass`,
    method: "post",
    data,
  });
};
