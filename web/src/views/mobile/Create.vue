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

// 创建缺失的移动端组件
const SunoCreate = {
  name: 'SunoCreate',
  setup() {
    const prompt = ref('')
    const duration = ref(30)
    const loading = ref(false)
    const result = ref('')

    const generateMusic = () => {
      if (!prompt.value.trim()) {
        showNotify({ type: 'warning', message: '请输入音乐描述' })
        return
      }
      loading.value = true
      // TODO: 调用Suno API
      setTimeout(() => {
        loading.value = false
        showNotify({ type: 'success', message: '音乐生成功能开发中' })
      }, 2000)
    }

    const downloadMusic = () => {
      // TODO: 实现下载功能
      showNotify({ type: 'primary', message: '下载功能开发中' })
    }

    return () =>
      h('div', { class: 'suno-create' }, [
        h('div', { class: 'create-header' }, [h('h3', '音乐创作'), h('p', 'AI驱动的音乐生成工具')]),
        h('div', { class: 'create-form' }, [
          h(Field, {
            value: prompt.value,
            onInput: (val) => {
              prompt.value = val
            },
            type: 'textarea',
            placeholder: '描述您想要的音乐风格、情感或主题...',
            rows: 4,
            maxlength: 500,
            'show-word-limit': true,
          }),
          h(Field, {
            value: duration.value,
            onInput: (val) => {
              duration.value = val
            },
            label: '时长',
            type: 'number',
            placeholder: '音乐时长（秒）',
          }),
          h(
            Button,
            {
              type: 'primary',
              size: 'large',
              block: true,
              loading: loading.value,
              onClick: generateMusic,
            },
            '生成音乐'
          ),
        ]),
        result.value
          ? h('div', { class: 'result-area' }, [
              h('h4', '生成结果'),
              h('audio', { src: result.value, controls: true }),
              h(Button, { size: 'small', onClick: downloadMusic }, '下载'),
            ])
          : null,
      ])
  },
}

const VideoCreate = {
  name: 'VideoCreate',
  setup() {
    const prompt = ref('')
    const duration = ref(10)
    const loading = ref(false)
    const result = ref('')

    const generateVideo = () => {
      if (!prompt.value.trim()) {
        showNotify({ type: 'warning', message: '请输入视频描述' })
        return
      }
      loading.value = true
      // TODO: 调用视频生成API
      setTimeout(() => {
        loading.value = false
        showNotify({ type: 'success', message: '视频生成功能开发中' })
      }, 2000)
    }

    const downloadVideo = () => {
      // TODO: 实现下载功能
      showNotify({ type: 'primary', message: '下载功能开发中' })
    }

    return () =>
      h('div', { class: 'video-create' }, [
        h('div', { class: 'create-header' }, [h('h3', '视频生成'), h('p', 'AI驱动的视频创作工具')]),
        h('div', { class: 'create-form' }, [
          h(Field, {
            value: prompt.value,
            onInput: (val) => {
              prompt.value = val
            },
            type: 'textarea',
            placeholder: '描述您想要的视频内容、风格或场景...',
            rows: 4,
            maxlength: 500,
            'show-word-limit': true,
          }),
          h(Field, {
            value: duration.value,
            onInput: (val) => {
              duration.value = val
            },
            label: '时长',
            type: 'number',
            placeholder: '视频时长（秒）',
          }),
          h(
            Button,
            {
              type: 'primary',
              size: 'large',
              block: true,
              loading: loading.value,
              onClick: generateVideo,
            },
            '生成视频'
          ),
        ]),
        result.value
          ? h('div', { class: 'result-area' }, [
              h('h4', '生成结果'),
              h('video', { src: result.value, controls: true }),
              h(Button, { size: 'small', onClick: downloadVideo }, '下载'),
            ])
          : null,
      ])
  },
}

const JimengCreate = {
  name: 'JimengCreate',
  setup() {
    const prompt = ref('')
    const negativePrompt = ref('')
    const steps = ref(20)
    const loading = ref(false)
    const result = ref('')

    const generateImage = () => {
      if (!prompt.value.trim()) {
        showNotify({ type: 'warning', message: '请输入图像描述' })
        return
      }
      loading.value = true
      // TODO: 调用即梦AI API
      setTimeout(() => {
        loading.value = false
        showNotify({ type: 'success', message: '即梦AI功能开发中' })
      }, 2000)
    }

    const downloadImage = () => {
      // TODO: 实现下载功能
      showNotify({ type: 'primary', message: '下载功能开发中' })
    }

    return () =>
      h('div', { class: 'jimeng-create' }, [
        h('div', { class: 'create-header' }, [h('h3', '即梦AI'), h('p', '专业的AI图像生成工具')]),
        h('div', { class: 'create-form' }, [
          h(Field, {
            value: prompt.value,
            onInput: (val) => {
              prompt.value = val
            },
            type: 'textarea',
            placeholder: '描述您想要的图像内容...',
            rows: 4,
            maxlength: 500,
            'show-word-limit': true,
          }),
          h(Field, {
            value: negativePrompt.value,
            onInput: (val) => {
              negativePrompt.value = val
            },
            type: 'textarea',
            placeholder: '负面提示词（可选）',
            rows: 2,
            maxlength: 200,
          }),
          h(Field, {
            value: steps.value,
            onInput: (val) => {
              steps.value = val
            },
            label: '步数',
            type: 'number',
            placeholder: '生成步数',
          }),
          h(
            Button,
            {
              type: 'primary',
              size: 'large',
              block: true,
              loading: loading.value,
              onClick: generateImage,
            },
            '生成图像'
          ),
        ]),
        result.value
          ? h('div', { class: 'result-area' }, [
              h('h4', '生成结果'),
              h(Image, { src: result.value, fit: 'cover' }),
              h(Button, { size: 'small', onClick: downloadImage }, '下载'),
            ])
          : null,
      ])
  },
}

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
      console.log(res)

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
