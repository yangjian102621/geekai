<template>
  <!-- 门户站点根布局：固定顶栏 + 内容 + 页脚；/space 页面自带顶部留白并隐藏页脚 -->
  <div class="site site-shell bg-canvas" :class="isSpace ? 'pt-16' : ''">
    <TopNav />
    <main class="flex-1">
      <router-view />
    </main>
    <Footer v-if="!isSpace" />
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import TopNav from './TopNav.vue'
import Footer from './Footer.vue'

const route = useRoute()
const isSpace = computed(() => route.path.startsWith('/space'))

// site.css 中的 body.site 规则依赖该 class 提供站点配色与字体
onMounted(() => document.body.classList.add('site'))
onBeforeUnmount(() => document.body.classList.remove('site'))
</script>
