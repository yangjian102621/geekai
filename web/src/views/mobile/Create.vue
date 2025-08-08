<template>
  <div class="create-center">
    <div class="create-content p-3">
      <CustomTabs
        :model-value="activeTab"
        @update:model-value="activeTab = $event"
        @tab-click="onTabChange"
      >
        <CustomTabPane name="mj">
          <template #label>
            <i class="iconfont icon-mj mr-1"></i>
            <span>Midjourney</span>
          </template>
          <div class="tab-content">
            <image-mj />
          </div>
        </CustomTabPane>
        <CustomTabPane name="sd">
          <template #label>
            <i class="iconfont icon-sd mr-1"></i>
            <span>StableDiffusion</span>
          </template>
          <div class="tab-content">
            <image-sd />
          </div>
        </CustomTabPane>
        <CustomTabPane name="dalle">
          <template #label>
            <i class="iconfont icon-dalle mr-1"></i>
            <span>Dalle</span>
          </template>
          <div class="tab-content">
            <image-dall />
          </div>
        </CustomTabPane>
      </CustomTabs>
    </div>
  </div>
</template>

<script setup>
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import ImageDall from '@/views/mobile/pages/ImageDall.vue'
import ImageMj from '@/views/mobile/pages/ImageMj.vue'
import ImageSd from '@/views/mobile/pages/ImageSd.vue'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const activeTab = ref(route.query.tab || 'mj')
const menus = ref([])
const activeMenu = ref({
  mj: false,
  sd: false,
  dall: false,
})

// Tab切换处理
const onTabChange = (name) => {
  router.replace({
    path: route.path,
    query: { ...route.query, tab: name },
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
