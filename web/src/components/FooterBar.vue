<template>
  <div class="foot-container">
    <div class="footer text-base">
      <div>
        <span>
          {{ title }} -
          {{ version }}
        </span>
      </div>
      <div class="flex justify-center text-sm">
        <span class="mr-2">{{ copyRight }}</span>
      </div>
      <div class="flex justify-center text-sm">
        <template v-if="icp">
          <a href="https://beian.miit.gov.cn" target="_blank">ICP备案：{{ icp }}</a>
        </template>
        <template v-if="gaBeian">
          <span>|</span>
          <img :src="gaBeianImg" class="w-4 h-4 mx-1" alt="beian" />
          <a
            :href="`http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=${getCodeNum(
              gaBeian
            )}`"
            target="_blank"
            >{{ gaBeian }}</a
          >
        </template>
      </div>
    </div>
  </div>
</template>
<script setup>
import { getSystemInfo } from '@/store/cache'
import { showMessageError } from '@/utils/dialog'
import { ref } from 'vue'

const title = ref('')
const version = ref(import.meta.env.VITE_VERSION)
const gitURL = ref(import.meta.env.VITE_GITHUB_URL)
const copyRight = ref('')
const icp = ref('')
const gaBeian = ref('')
const props = defineProps({
  textColor: {
    type: String,
    default: '#ffffff',
  },
})
const gaBeianImg = ref('/images/ga_beian.png')

// 获取系统配置
getSystemInfo()
  .then((res) => {
    title.value = res.data.title ?? import.meta.env.VITE_TITLE
    copyRight.value =
      (res.data.copyright ? res.data.copyright : '极客学长') +
      ' © 2023 - ' +
      new Date().getFullYear() +
      ' All rights reserved'
    icp.value = res.data.icp
    gaBeian.value = res.data.ga_beian
  })
  .catch((e) => {
    showMessageError('获取系统配置失败：' + e.message)
  })

// 获取公安备案号
const getCodeNum = (code) => {
  // 提取数字
  try {
    const num = code.match(/\d+/)
    if (num) {
      return num[0]
    }
  } catch (e) {
    return ''
  }
  return ''
}
</script>

<style scoped lang="scss">
.foot-container {
  // 仅在 PC 端 fixed，移动端正常流式布局
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100%;
  display: flex;
  justify-content: center;
  // background: var(--theme-bg);
  margin-top: -4px;

  @media (max-width: 768px) {
    margin-top: 2rem !important;
    position: static;
    margin-top: 0;
  }

  .footer {
    // max-width: 400px;
    text-align: center;
    font-size: 14px;
    padding: 20px;
    width: 100%;

    a {
      color: var(--text-color);

      &:hover {
        text-decoration: underline;
      }
    }
    span {
      color: var(--text-color);
    }
  }
}
</style>
