<template>
  <div>
   支付回调
  </div>
</template>

<script setup>
import {useRouter} from "vue-router";
import {isMobile} from "@/utils/libs";

const router = useRouter()
console.log(router.currentRoute.value.query)
if (isMobile()) {
  router.push('/mobile/profile')
} else {
  window.close()
}
</script>
