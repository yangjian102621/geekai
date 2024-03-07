import { h } from "vue";
import type { Component, ComponentInternalInstance } from "vue";
import { Modal, Drawer } from "@arco-design/web-vue";
import type { ModalConfig, DrawerConfig } from "@arco-design/web-vue";
import app from "@/main";

interface Config {
  nodeProps?: (...arg: any) => Record<string, any>;
  popupProps?: (
    arg: any[],
    exposed: () => ComponentInternalInstance["exposed"]
  ) => Omit<ModalConfig | DrawerConfig, "content"> & {
    [key: string]: any;
  };
  type?: "drawer" | "modal";
}

const component = {
  modal: Modal,
  drawer: Drawer,
};
function usePopup(node: Component, config: Config) {
  const { nodeProps, popupProps, type = "modal" } = config;

  return (...arg: any[]) => {
    const content = h(node, nodeProps ? nodeProps(arg) : {});
    const popupNode = component[type];
    // 获取全局组件的上下文
    popupNode._context = app._context;
    popupNode.open({
      content: () => content,
      ...popupProps?.(arg, () => content?.component?.exposed as any),
    });
  };
}

export default usePopup;
