<template>
  <h1>{{ title }}</h1>
</template>

<script setup>
import {ref} from "vue"
import {useRouter} from "vue-router";
import {checkSession} from "@/action/session";
import {isMobile} from "@/utils/libs";

const title = ref("HI, ChatGPT PLUS!");
const router = useRouter();
checkSession().then(() => {
  if (isMobile()) {
    router.push("/mobile")
  } else {
    router.push("/chat")
  }
}).catch(() => {
  router.push("login")
})

</script>
