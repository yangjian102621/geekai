<template>
  <div class="create-center">
    <van-nav-bar title="AI 创作中心" fixed :safe-area-inset-top="true">
      <template #left>
        <div class="nav-left">
          <i class="iconfont icon-mj"></i>
        </div>
      </template>
    </van-nav-bar>

    <div class="create-content">
      <van-tabs 
        v-model:active="activeTab" 
        animated 
        sticky 
        :offset-top="44"
        @change="onTabChange"
      >
        <van-tab title="MJ绘画" name="mj" v-if="activeMenu.mj">
          <div class="tab-content">
            <image-mj />
          </div>
        </van-tab>
        <van-tab title="SD绘画" name="sd" v-if="activeMenu.sd">
          <div class="tab-content">
            <image-sd />
          </div>
        </van-tab>
        <van-tab title="DALL·E" name="dalle" v-if="activeMenu.dall">
          <div class="tab-content">
            <image-dall />
          </div>
        </van-tab>
        <van-tab title="音乐创作" name="suno" v-if="activeMenu.suno">
          <div class="tab-content">
            <suno-create />
          </div>
        </van-tab>
        <van-tab title="视频生成" name="video" v-if="activeMenu.video">
          <div class="tab-content">
            <video-create />
          </div>
        </van-tab>
        <van-tab title="即梦AI" name="jimeng" v-if="activeMenu.jimeng">
          <div class="tab-content">
            <jimeng-create />
          </div>
        </van-tab>
      </van-tabs>
    </div>
  </div>
</template>

<script setup>
import { httpGet } from '@/utils/http'
import ImageDall from '@/views/mobile/pages/ImageDall.vue'
import ImageMj from '@/views/mobile/pages/ImageMj.vue'
import ImageSd from '@/views/mobile/pages/ImageSd.vue'
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// 临时组件，实际项目中需要创建对应的移动端组件
const SunoCreate = { template: '<div class="placeholder">Suno音乐创作功能开发中...</div>' }
const VideoCreate = { template: '<div class="placeholder">视频生成功能开发中...</div>' }
const JimengCreate = { template: '<div class="placeholder">即梦AI功能开发中...</div>' }

const route = useRoute()
const router = useRouter()
const activeTab = ref('mj')
const menus = ref([])
const activeMenu = ref({
  mj: false,
  sd: false,
  dall: false,
  suno: false,
  video: false,
  jimeng: false,
})

// 监听路由参数变化
watch(() => route.query.tab, (newTab) => {
  if (newTab && activeMenu.value[newTab]) {
    activeTab.value = newTab
  }
}, { immediate: true })

// Tab切换处理
const onTabChange = (name) => {
  router.replace({ 
    path: route.path, 
    query: { ...route.query, tab: name } 
  })
}

onMounted(() => {
  fetchMenus()
})

const fetchMenus = () => {
  httpGet('/api/menu/list').then((res) => {
    menus.value = res.data
    activeMenu.value = {
      mj: menus.value.some((item) => item.url === '/mj'),
      sd: menus.value.some((item) => item.url === '/sd'),
      dall: menus.value.some((item) => item.url === '/dalle'),
      suno: menus.value.some((item) => item.url === '/suno'),
      video: menus.value.some((item) => item.url === '/video'),
      jimeng: menus.value.some((item) => item.url === '/jimeng'),
    }
    
    // 如果没有指定tab，默认选择第一个可用的
    if (!route.query.tab) {
      const firstAvailable = Object.keys(activeMenu.value).find(key => activeMenu.value[key])
      if (firstAvailable) {
        activeTab.value = firstAvailable
      }
    }
  }).catch((e) => {
    console.error('获取菜单失败：', e.message)
  })
}
</script>

<style lang="scss" scoped>
.create-center {
  min-height: 100vh;
  background: var(--van-background);

  .nav-left {
    display: flex;
    align-items: center;
    
    .iconfont {
      font-size: 20px;
      color: var(--van-primary-color);
    }
  }

  .create-content {
    padding-top: 44px; // nav-bar height
    
    :deep(.van-tabs__nav) {
      background: var(--van-background);
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    }

    :deep(.van-tab) {
      font-weight: 500;
    }

    :deep(.van-tab--active) {
      font-weight: 600;
    }

    .tab-content {
      min-height: calc(100vh - 88px);
      
      .placeholder {
        display: flex;
        align-items: center;
        justify-content: center;
        height: 400px;
        color: var(--van-gray-6);
        font-size: 16px;
      }
    }
  }
}
</style>