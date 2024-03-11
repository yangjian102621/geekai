import usePopup, { type Config } from "./usePopup";
import { Message } from "@arco-design/web-vue";
import type { Component } from "vue";
import type { BaseResponse } from "@gpt-vue/packages/type";
interface Arg {
  reload?: () => void;
  record?: Record<string, any>;
}

export default function (
  node: Component,
  api: (params?: any) => Promise<BaseResponse<any>>,
  config?: Config
): (arg: Arg) => void {
  const nodeProps = (arg: Arg[]) => {
    return {
      data: arg[0].record || {},
      ...config.nodeProps?.(arg),
    };
  };

  const popupProps = (arg: Arg[], getExposed) => {
    return {
      width: 750,
      maskClosable: false,
      onBeforeOk: async () => {
        const exposed = getExposed();
        const validateRes = await exposed?.formRef.value.validate();
        if (validateRes) {
          return false;
        }
        const { code } = await api(exposed?.form.value);
        if (code === 0) {
          Message.success("操作成功");
        }
        arg[0]?.reload?.();
        return code === 0;
      },
      ...config.popupProps?.(arg, getExposed),
    };
  };

  return usePopup(node, { nodeProps, popupProps });
}
