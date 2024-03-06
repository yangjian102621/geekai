import { Notification } from "@arco-design/web-vue";
import createInstance from "@gpt-vue/packages/request"
import type { BaseResponse } from "@gpt-vue/packages/type";

export const uploadUrl = import.meta.env.VITE_PROXY_BASE_URL + "/common/upload/minio";

export const instance = createInstance()

instance.interceptors.request.use((config) => {
  config.headers[__AUTH_KEY] = localStorage.getItem(__AUTH_KEY);
  config.headers["Authorization"] = localStorage.getItem(__AUTH_KEY);
  return config;
});

instance.interceptors.response.use(
  (response) => {
    const { data }: { data: BaseResponse<unknown> } = response
    if (data && typeof data === "object" && data.code !== 0) {
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
