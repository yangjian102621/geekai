<template>
  <div class="chat-line chat-line-left">
    <div class="chat-icon">
      <img :src="icon" alt="ChatGPT">
    </div>

    <div class="chat-item">
      <div class="triangle"></div>
      <span class="content">{{ content }}</span>
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
.chat-line-left {
  justify-content: start;

  .chat-icon {
    margin-right 5px;
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

    .content {
      float: left;
      padding: 8px 10px;
      color var(--content-color)
      background-color: #fff;
      font-size: var(--content-font-size);
      border-radius: 5px;
    }
  }
}

</style>
