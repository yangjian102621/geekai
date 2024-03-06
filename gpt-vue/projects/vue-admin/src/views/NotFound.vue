<script setup lang="ts">
import { useRouter } from "vue-router";
const router = useRouter();
const goBack = () => {
  if (window.history.length) {
    router.back();
    return;
  }
  router.replace("/");
};
</script>

<template>
  <AResult status="404" subtitle="找不到页面">
    <template #extra>
      <AButton type="primary" @click="goBack">返回</AButton>
    </template>
  </AResult>
</template>
