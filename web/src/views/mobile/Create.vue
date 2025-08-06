<template>
  <div class="create-center">
    <div class="create-content p-3">
      <CustomTabs
        :model-value="activeTab"
        @update:model-value="activeTab = $event"
        @tab-click="onTabChange"
      >
        <CustomTabPane name="mj" label="MJ绘画">
          <div class="tab-content">
            <image-mj />
          </div>
        </CustomTabPane>
        <CustomTabPane name="sd" label="SD绘画">
          <div class="tab-content">
            <image-sd />
          </div>
        </CustomTabPane>
        <CustomTabPane name="dalle" label="DALL·E">
          <div class="tab-content">
            <image-dall />
          </div>
        </CustomTabPane>
        <CustomTabPane name="suno" label="音乐创作">
          <div class="tab-content">
            <suno-create />
          </div>
        </CustomTabPane>
        <CustomTabPane name="video" label="视频生成">
          <div class="tab-content">
            <video-create />
          </div>
        </CustomTabPane>
        <CustomTabPane name="jimeng" label="即梦AI">
          <div class="tab-content">
            <jimeng-create />
          </div>
        </CustomTabPane>
      </CustomTabs>
    </div>
  </div>
</template>

<script setup>
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { httpGet } from '@/utils/http'
import ImageDall from '@/views/mobile/pages/ImageDall.vue'
import ImageMj from '@/views/mobile/pages/ImageMj.vue'
import ImageSd from '@/views/mobile/pages/ImageSd.vue'
import { Button, Field, Image, showNotify } from 'vant'
import { h, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// 删除 SunoCreate、VideoCreate、JimengCreate 相关 setup 代码和渲染逻辑

const route = useRoute()
const router = useRouter()
const activeTab = ref(route.query.tab || 'mj')
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
watch(
  () => route.query.tab,
  (newTab) => {
    if (newTab && activeMenu.value[newTab]) {
      activeTab.value = newTab
    }
  },
  { immediate: true }
)

// Tab切换处理
const onTabChange = (name) => {
  router.replace({
    path: route.path,
    query: { ...route.query, tab: name },
  })
}

onMounted(() => {
  fetchMenus()
})

const fetchMenus = () => {
  httpGet('/api/menu/list')
    .then((res) => {
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
        const firstAvailable = Object.keys(activeMenu.value).find((key) => activeMenu.value[key])
        if (firstAvailable) {
          activeTab.value = firstAvailable
        }
      } else {
        // 如果当前选中的tab不可用，选择第一个可用的
        if (!activeMenu.value[route.query.tab]) {
          const firstAvailable = Object.keys(activeMenu.value).find((key) => activeMenu.value[key])
          if (firstAvailable) {
            activeTab.value = firstAvailable
            router.replace({
              path: route.path,
              query: { ...route.query, tab: firstAvailable },
            })
          }
        }
      }
    })
    .catch((e) => {
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

      // 新增组件样式
      .suno-create,
      .video-create,
      .jimeng-create {
        padding: 20px;

        .create-header {
          text-align: center;
          margin-bottom: 24px;

          h3 {
            font-size: 20px;
            font-weight: 600;
            color: var(--van-text-color);
            margin: 0 0 8px 0;
          }

          p {
            font-size: 14px;
            color: var(--van-gray-6);
            margin: 0;
          }
        }

        .create-form {
          margin-bottom: 24px;

          .van-field {
            margin-bottom: 16px;
          }

          .van-button {
            margin-top: 16px;
          }
        }

        .result-area {
          background: var(--van-cell-background);
          border-radius: 12px;
          padding: 16px;
          text-align: center;

          h4 {
            font-size: 16px;
            font-weight: 600;
            color: var(--van-text-color);
            margin: 0 0 12px 0;
          }

          audio,
          video,
          .van-image {
            width: 100%;
            max-width: 300px;
            margin-bottom: 12px;
            border-radius: 8px;
          }

          .van-button {
            margin-top: 8px;
          }
        }
      }
    }
  }
}
</style>
