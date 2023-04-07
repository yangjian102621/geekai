<template>
  <div class="chat-line chat-line-left">
    <div class="chat-icon">
      <img :src="icon" alt="ChatGPT">
    </div>

    <div class="chat-item">
      <div class="triangle"></div>
      <div class="content" v-html="content"></div>
    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue"

export default defineComponent({
  name: 'ChatReply',
  props: {
    content: {
      type: String,
      default: '',
    },
    icon: {
      type: String,
      default: 'images/gpt-icon.png',
    }
  },
  data() {
    return {}
  },
})
</script>

<style lang="stylus">
.body-plus {
  .chat-line-left {
    justify-content: flex-start;

    .chat-icon {
      margin-right 5px;

      img {
        border-radius 50%;
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
        border-top: 10px solid transparent;
        border-bottom: 10px solid transparent;
        border-right: 10px solid #404042;
        position: absolute;
        left: 0;
        top: 13px;
      }

      .content {
        min-height 20px;
        word-break break-word;
        padding: 12px 15px;
        color var(--content-color)
        background-color: #404042;
        font-size: var(--content-font-size);
        border-radius: 5px;

        p {
          line-height 1.5

          code {
            color #f1f1f1
            background-color #202121
            padding 0 3px;
            border-radius 5px;
          }
        }

        p:last-child {
          margin-bottom: 0
        }

        p:first-child {
          margin-top 0
        }

      }
    }
  }
}

</style>
