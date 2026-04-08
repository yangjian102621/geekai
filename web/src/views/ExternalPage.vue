<template>
  <div
    class="page-iframe"
    v-loading="loading"
    style="--el-color-primary: #47fff1"
    element-loading-text="页面正在加载中..."
    element-loading-background="rgba(122, 122, 122, 0.8)"
  >
    <iframe :src="externalUrl" class="iframe" @load="onIframeLoad"></iframe>
  </div>
</template>
<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const loading = ref(true)
const router = useRouter()
const externalUrl = computed(() => {
  loading.value = true
  return router.currentRoute.value.query.url || 'about:blank'
})

// 设置标题
document.title = router.currentRoute.value.query.title

const onIframeLoad = () => {
  loading.value = false
}
</script>

<style scoped lang="scss">
.page-iframe {
  width: 100%;
  height: 100vh;
  overflow: hidden;

  .iframe {
    width: 100%;
    height: 100%;
    border: none;
  }
}
</style>
