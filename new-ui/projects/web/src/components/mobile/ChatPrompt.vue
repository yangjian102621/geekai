<template>
  <div class="mobile-message-prompt">
    <div class="chat-item">
      <div ref="contentRef" :data-clipboard-text="content" class="content" v-html="content"></div>
      <div class="triangle"></div>
    </div>

    <div class="chat-icon">
      <van-image :src="icon"/>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import Clipboard from "clipboard";
import {showNotify} from "vant";

const props = defineProps({
  content: {
    type: String,
    default: '',
  },
  icon: {
    type: String,
    default: '/images/user-icon.png',
  }
});
const contentRef = ref(null)
onMounted(() => {
  const clipboard = new Clipboard(contentRef.value);
  clipboard.on('success', () => {
    showNotify({type: 'success', message: '复制成功', duration: 1000})
  })
  clipboard.on('error', () => {
    showNotify({type: 'danger', message: '复制失败', duration: 2000})
  })
})
</script>

<style lang="stylus">
.mobile-message-prompt {
  display flex
  justify-content: flex-end

  .chat-icon {
    margin-left 5px

    .van-image {
      width 32px

      img {
        border-radius 5px
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
      border-left: 5px solid #98E165;
      position: absolute;
      right: 0;
      top: 10px;
    }

    .content {
      word-break break-word;
      text-align left
      padding: 5px 10px;
      background-color: #98E165;
      color #444444
      font-size: 14px
      border-radius: 5px
      line-height 1.5
    }
  }
}

.van-theme-dark {
  .mobile-message-prompt {
    .chat-item {

      .triangle {
        border-left: 5px solid #223A34
      }

      .content {
        background-color: #223A34
        color #c1c1c1
      }
    }
  }
}
</style>