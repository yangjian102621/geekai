<template>
  <div class="page-iframe" v-loading="loading"
       style="--el-color-primary:#47fff1"
       element-loading-text="页面正在加载中..."
       element-loading-background="rgba(122, 122, 122, 0.8)">
    <iframe :src="externalUrl" class="iframe" @load="onIframeLoad"></iframe>

  </div>
</template>
<script setup>
import {useRouter} from "vue-router";
import {computed, ref} from "vue";

const loading = ref(true)
const router = useRouter()
const externalUrl = computed(() => {
  loading.value = true
  return router.currentRoute.value.query.url || 'about:blank'
})

const onIframeLoad = () => {
  loading.value = false
}
</script>

<style scoped lang="stylus">
.page-iframe {
  width 100%
  height 100%
  overflow hidden

  .iframe {
    width 100%
    height 100%
    border none
  }
}
</style>
