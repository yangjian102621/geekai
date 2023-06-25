<template>
  <div class="message-reply">
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

<style lang="stylus">
.message-reply {
  justify-content: flex-end;

  .chat-icon {
    margin-left 5px;

    img {
      border-radius 5px;
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
      padding: 6px 10px;
      background-color: #98E165;
      color var(--content-color);
      font-size: var(--content-font-size);
      border-radius: 5px;
      line-height 1.5
    }
  }
}
</style>