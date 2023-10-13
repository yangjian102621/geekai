<template>
  <div class="page-apps">
    <div class="title">
      AI 助手应用中心
    </div>
    <div class="inner custom-scroll">
      <div class="app-list">
        <div class="list-item" v-for="item in list" :key="item.id">
          <div v-if="item.key !=='gpt'">
            <el-image :src="item.icon"/>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {ElMessage} from "element-plus";
import {httpGet} from "@/utils/http";

const list = ref([])
onMounted(() => {
  httpGet("/api/role/list?all=true").then((res) => {
    list.value = res.data
  }).catch(e => {
    ElMessage.error("获取应用失败：" + e.message)
  })
})
</script>

<style lang="stylus" scoped>
@import "@/assets/css/chat-app.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
