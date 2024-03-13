<script lang="ts" setup>
import { computed } from "vue";
import { useRoute } from "vue-router";
import router from "@/router";
import showMenu from "@/router/menu";

defineProps({
  width: {
    type: [Number, String],
    default: 200,
  },
});
const route = useRoute();
const goto = (name: string) => router.push({ name });
const selectedKeys = computed(() => [route.name]);

</script>
<template>
  <ALayoutSider :style="{ width, height: '100%' }">
    <AMenu :selected-keys="selectedKeys" @menu-item-click="goto">
      <AMenuItem v-for="item in showMenu" :key="item.name">
        <template #icon>
          <component :is="item.meta?.icon" />
        </template>
        {{ item.meta.title }}
      </AMenuItem>
    </AMenu>
  </ALayoutSider>
</template>
