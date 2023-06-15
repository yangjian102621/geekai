<template>
  <div class="chat-line chat-line-right">
    <div class="chat-item">
      <div class="content" v-html="content"></div>
      <div class="triangle"></div>
    </div>

    <div class="chat-icon">
      <img :src="icon" alt="User"/>
    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue"

export default defineComponent({
  name: 'ChatPrompt',
  props: {
    content: {
      type: String,
      default: '',
    },
    icon: {
      type: String,
      default: 'images/user-icon.png',
    }
  },
  data() {
    return {}
  },
})
</script>

<style lang="stylus" scoped>
.chat-line-right {
  justify-content: flex-end;

  .chat-icon {
    margin-left 5px;

    img {
      border-radius 50%;
    }
  }

  .chat-item {
    position: relative;
    padding: 0 5px 0 0;
    overflow: hidden;

    .triangle {
      width: 0;
      height: 0;
      border-top: 6px solid transparent;
      border-bottom: 6px solid transparent;
      border-left: 6px solid #223A34;
      position: absolute;
      right: 0;
      top: 10px;
    }

    .content {
      word-break break-word;
      padding: 6px 10px;
      background-color: #223A34;
      color var(--content-color);
      font-size: var(--content-font-size);
      border-radius: 5px;
      overflow: auto;

      p {
        line-height 1.5
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
</style>