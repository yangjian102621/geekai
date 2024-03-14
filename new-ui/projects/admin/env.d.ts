/// <reference types="vite/client" />
declare module "*.vue" {
  import { DefineComponent } from "vue";
  // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/ban-types
  const component: DefineComponent<{}, {}, any>;
  export default component;
}
import 'vue-router'
declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    icon?: any
    permission?: boolean | string | string[]
  }
}

declare const __AUTH_KEY: string;