<template>
  <div>{{ title }}</div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {isMobile} from "@/utils/libs";
import {useRouter} from "vue-router"; // 导入useRouter函数

const title = ref("Loading page...");
const router = useRouter();
onMounted(() => {
  if (isMobile()) {
    router.push("mobile");
  } else {
    router.push("chat");
  }
})


</script>
