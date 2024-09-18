<template>
  <div>
   支付回调
  </div>
</template>

<script setup>
import {useRouter} from "vue-router";

const router = useRouter()
console.log(router.currentRoute.value.query)
window.close()
</script>
