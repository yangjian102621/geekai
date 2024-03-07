import http from "@/http/config";

export const userLogin = (data: {
  username: string;
  password: string;
}) => {
  return http({
    url: "/api/admin/login",
    method: "post",
    data,
  });
};

export const userLogout = () => {
  return http({
    url: "/api/admin/logout",
    method: "get",
  });
};

export const getSession = () => {
  return http({
    url: "/api/admin/session",
    method: "get",
  });
};


export const loginLog = (params?: Record<string, unknown>) => {
  return http({
    url: "/api/admin/user/loginLog",
    method: "get",
    params
  })
}