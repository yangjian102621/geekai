<template>
  <!-- 门户站点根布局：固定顶栏 + 内容 + 页脚；/space 页面自带顶部留白并隐藏页脚 -->
  <div
    class="site site-shell bg-canvas"
    :class="[isSpace ? 'pt-16' : '', isThemed ? 'home-shell' : '', isSpace ? 'space-shell' : '']"
  >
    <TopNav />
    <main class="flex-1">
      <router-view />
    </main>
    <Footer v-if="!isSpace" />
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import TopNav from './TopNav.vue'
import Footer from './Footer.vue'

const route = useRoute()
const isSpace = computed(() => route.path.startsWith('/space'))
// 首页与邀请页共用同一套皮肤
const isThemed = computed(() => route.path === '/' || route.path.startsWith('/invite'))

// site.css 中的 body.site 规则依赖该 class 提供站点配色与字体；主题页额外挂对应皮肤
onMounted(() => {
  document.body.classList.add('site')
  document.body.classList.toggle('home-shell', isThemed.value)
  document.body.classList.toggle('space-shell', isSpace.value)
})
watch(isThemed, (v) => document.body.classList.toggle('home-shell', v))
watch(isSpace, (v) => document.body.classList.toggle('space-shell', v))
onBeforeUnmount(() => {
  document.body.classList.remove('site')
  document.body.classList.remove('home-shell')
  document.body.classList.remove('space-shell')
})
</script>
