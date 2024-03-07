<script lang="ts" setup>
import { ref } from "vue";
import SystemBaseConfig from "./SystemBaseConfig.vue";
import SystemChatConfig from "./SystemChatConfig.vue";
import SystemNoticeConfig from "./SystemNoticeConfig.vue";
const tabsList = [
  { key: "1", title: "基本设置", components: SystemBaseConfig },
  { key: "2", title: "模型设置", components: SystemChatConfig },
  { key: "3", title: "公告设置", components: SystemNoticeConfig },
];

const activeKey = ref(tabsList[0].key);
</script>
<template>
  <a-tabs v-model:active-key="activeKey" lazy-load>
    <a-tab-pane v-for="item in tabsList" :key="item.key" :title="item.title">
      <component :is="item.components" />
    </a-tab-pane>
  </a-tabs>
</template>
