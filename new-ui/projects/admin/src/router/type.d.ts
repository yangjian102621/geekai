import 'vue-router'
declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    icon?: any
    permission?: boolean | string | string[]
  }
}