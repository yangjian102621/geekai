import { createApp } from "vue";
import { createPinia } from "pinia";
import ArcoVue from "@arco-design/web-vue";
import ArcoVueIcon from "@arco-design/web-vue/es/icon";
import "@arco-design/web-vue/dist/arco.css";

import App from "./App.vue";
import router from "./router";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(ArcoVue);
app.use(ArcoVueIcon);

app.mount("#app");
app.config.warnHandler = (msg, vm, trace) => {
  if (msg.includes('Invalid prop name: "key" is a reserved property.')) {
    // 如果警告信息包含我们要屏蔽的内容，则不执行任何操作
    return;
  }
  console.warn(`[Vue warn]: ${msg}${trace}`);
};

export default app;
