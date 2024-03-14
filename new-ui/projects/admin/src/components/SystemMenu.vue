<script lang="ts">
import { defineComponent, h } from "vue";
import type { Component, PropType } from "vue";
import type { RouteRecordRaw } from "vue-router";
import { SubMenu, MenuItem } from "@arco-design/web-vue";
const CustomMenuItem: Component = defineComponent({
  props: {
    tree: {
      type: Array as PropType<RouteRecordRaw[]>,
      default: () => [],
    },
  },
  setup: (props) => {
    return () =>
      props.tree?.map((item) => {
        const _icon = item.meta?.icon ? h(item.meta.icon) : undefined;
        const hasChildren = Array.isArray(item.children) && item.children.length;
        if (hasChildren) {
          return h(
            SubMenu,
            { title: item.meta.title, key: item.name },
            {
              default: () => h(CustomMenuItem, { tree: item.children }),
              icon: () => _icon,
            }
          );
        }
        return h(
          MenuItem,
          { key: item.name },
          { default: () => item.meta.title, icon: () => _icon }
        );
      });
  },
});
</script>

<script lang="ts" setup>
import { computed } from "vue";
import { useRoute } from "vue-router";
import router from "@/router";
import menu from "@/router/menu";
import { hasPermission } from "@/directives/permission";
defineProps({
  width: {
    type: [Number, String],
    default: 200,
  },
});
const route = useRoute();
const goto = (name: string) => router.push({ name });

const selectedKeys = computed(() => [route.name]);

const showMenu = computed(() => menu.filter((item: any) => hasPermission(item.meta?.permission)));
</script>

<template>
  <ALayoutSider :style="{ width, height: '100%' }">
    <AMenu :selectedKeys="selectedKeys" auto-open-selected @menu-item-click="goto">
      <CustomMenuItem :tree="showMenu" />
    </AMenu>
  </ALayoutSider>
</template>
