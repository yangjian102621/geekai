<template>
  <h1>{{ title }}</h1>
</template>

<script setup>
import {ref} from "vue"
import {useRouter} from "vue-router";
import {checkSession} from "@/action/session";

const title = ref("HI, ChatGPT PLUS!");
const router = useRouter();
checkSession().then(() => {
  router.push("chat")
}).catch(() => {
  router.push("login")
})

</script>
