<template>
  <div class="mobile-message-reply">
    <div class="chat-icon">
      <van-image :src="icon"/>
    </div>

    <div class="chat-item">
      <div class="triangle"></div>
      <div class="content-box">
        <div ref="contentRef" :data-clipboard-text="orgContent" class="content" v-html="content"></div>
      </div>

    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"

import Clipboard from "clipboard";
import {showNotify} from "vant";

const props = defineProps({
  content: {
    type: String,
    default: '',
  },
  orgContent: {
    type: String,
    default: '',
  },
  icon: {
    type: String,
    default: '/images/gpt-icon.png',
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
.mobile-message-reply {
  display flex
  justify-content: flex-start;

  .chat-icon {
    margin-right 5px

    .van-image {
      width 32px

      img {
        border-radius 5px
      }
    }
  }

  .chat-item {
    display: inline-block;
    position: relative;
    padding: 0 0 0 5px;
    overflow: hidden;

    .triangle {
      width: 0;
      height: 0;
      border-top: 5px solid transparent;
      border-bottom: 5px solid transparent;
      border-right: 5px solid #fff;
      position: absolute;
      left: 0;
      top: 13px;
    }

    .content-box {

      display flex
      flex-direction row

      .content {
        text-align left
        width 100%
        overflow-x auto
        min-height 20px;
        word-break break-word;
        padding: 5px 10px;
        color #444444
        background-color: #ffffff;
        font-size: 16px
        border-radius: 5px;

        p:last-child {
          margin-bottom: 0
        }

        p:first-child {
          margin-top 0
        }

        p {
          code {
            color #2b2b2b
            background-color #c1c1c1
            padding 2px 5px
            border-radius 5px
          }

          img {
            max-width 100%
          }
        }
      }
    }

  }
}


.van-theme-dark {
  .mobile-message-reply {
    .chat-item {
      .triangle {
        border-right: 5px solid #404042;
      }

      .content-box {
        .content {
          color #c1c1c1
          background-color: #404042;

          p > code {
            color #c1c1c1
            background-color #2b2b2b
          }
        }
      }

    }
  }

}
</style>