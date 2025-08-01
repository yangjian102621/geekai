<template>
  <div class="mobile-message-prompt">
    <div class="chat-item">
      <div ref="contentRef" :data-clipboard-text="content" class="content" v-html="content"></div>
      <div class="triangle"></div>
    </div>

    <div class="chat-icon">
      <van-image :src="icon" />
    </div>
  </div>
</template>

<script setup>
import Clipboard from 'clipboard'
import { showNotify } from 'vant'
import { onMounted, ref } from 'vue'

// eslint-disable-next-line no-unused-vars,no-undef
const props = defineProps({
  content: {
    type: Object,
    default: {
      text: '',
      files: [],
    },
  },
  icon: {
    type: String,
    default: '/images/user-icon.png',
  },
})
const contentRef = ref(null)
const content = computed(() => {
  return props.content.text
})
onMounted(() => {
  const clipboard = new Clipboard(contentRef.value)
  clipboard.on('success', () => {
    showNotify({ type: 'success', message: '复制成功', duration: 1000 })
  })
  clipboard.on('error', () => {
    showNotify({ type: 'danger', message: '复制失败', duration: 2000 })
  })
})
</script>

<style lang="scss">
.mobile-message-prompt {
  display: flex;
  justify-content: flex-end;

  .chat-icon {
    margin-left: 5px;

    .van-image {
      width: 32px;

      img {
        border-radius: 5px;
      }
    }
  }

  .chat-item {
    position: relative;
    padding: 0 5px 0 0;
    overflow: hidden;

    .triangle {
      width: 0;
      height: 0;
      border-top: 5px solid transparent;
      border-bottom: 5px solid transparent;
      border-left: 5px solid #98e165;
      position: absolute;
      right: 0;
      top: 10px;
    }

    .content {
      word-break: break-word;
      text-align: left;
      padding: 5px 10px;
      background-color: #98e165;
      color: #444444;
      font-size: 14px;
      border-radius: 5px;
      line-height: 1.5;
    }
  }
}

.van-theme-dark {
  .mobile-message-prompt {
    .chat-item {
      .triangle {
        border-left: 5px solid #223a34;
      }

      .content {
        background-color: #223a34;
        color: #c1c1c1;
      }
    }
  }
}
</style>
