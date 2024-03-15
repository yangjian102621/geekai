import router from "@/router";
import { Notification } from "@arco-design/web-vue";
import createInstance from "@chatgpt-plus/packages/request"
import type { BaseResponse } from "@chatgpt-plus/packages/type";

export const uploadUrl = import.meta.env.VITE_PROXY_BASE_URL + "/api/admin/upload";

export const instance = createInstance(import.meta.env.VITE_PROXY_BASE_URL)

instance.interceptors.request.use((config) => {
  const TOKEN = JSON.parse(localStorage.getItem(__AUTH_KEY))?.token
  config.headers[__AUTH_KEY] = TOKEN;
  config.headers["Authorization"] = TOKEN;
  return config;
});

instance.interceptors.response.use(
  (response) => {
    const { data }: { data: BaseResponse<unknown> } = response
    if (data && typeof data === "object" && data.code > 0) {
      switch (data.code) {
        case 400: {
          localStorage.removeItem(__AUTH_KEY);
          router.push({ name: "Login" })
          break;
        }
        case 403: {
          router.replace({ name: "403" })
          break;
        }
      }
      Notification.error(data.message ?? '未知错误')
    }
    return { data, response } as any;
  },
  (error) => {
    const STATUS_CODE: any = {
      401: {
        msg: error.response.data || "没有操作权限！",
        event: null,
      },
      500: {
        msg: error.response.data || "系统正在部署升级中，请稍后再试！",
        event: null,
      },
    };

    const statusCodeEvent = STATUS_CODE?.[error.response.status];

    if (statusCodeEvent) {
      Notification.error(statusCodeEvent.msg);
      statusCodeEvent.event?.();
    }

    if (error.message.indexOf("timeout") !== -1) {
      Notification.error("连接超时");
    }
    return Promise.reject(error);
  }
);

function http<T = any>(config: any): Promise<BaseResponse<T>> {
  return instance(config).then((res) => {
    return res.data;
  }) as unknown as Promise<BaseResponse<T>>;
}

export function originHttp<T = any>(config: any) {
  return instance<T>(config as any).then((res) => res);
}

export default http;
