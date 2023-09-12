<template>
  <div class="chat-line chat-line-prompt">
    <div class="chat-line-inner">
      <div class="chat-icon">
        <img :src="icon" alt="User"/>
      </div>

      <div class="chat-item">
        <div class="content" v-html="content"></div>
        <div class="bar" v-if="createdAt !== ''">
          <span class="bar-item"><el-icon><Clock/></el-icon> {{ createdAt }}</span>
          <span class="bar-item">tokens: {{ finalTokens }}</span>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import {defineComponent} from "vue"
import {Clock} from "@element-plus/icons-vue";
import {httpPost} from "@/utils/http";

export default defineComponent({
  name: 'ChatPrompt',
  components: {Clock},
  methods: {},
  props: {
    content: {
      type: String,
      default: '',
    },
    icon: {
      type: String,
      default: 'images/user-icon.png',
    },
    createdAt: {
      type: String,
      default: '',
    },
    tokens: {
      type: Number,
      default: 0,
    },
    model: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      finalTokens: this.tokens
    }
  },
  mounted() {
    if (!this.finalTokens) {
      httpPost("/api/chat/tokens", {text: this.content, model: this.model}).then(res => {
        this.finalTokens = res.data;
      })
    }
  }
})
</script>

<style lang="stylus">
.chat-line-prompt {
  background-color #ffffff;
  justify-content: center;
  width 100%
  padding-bottom: 1.5rem;
  padding-top: 1.5rem;
  border-bottom: 1px solid #d9d9e3;

  .chat-line-inner {
    display flex;
    width 100%;
    max-width 900px;
    padding-left 10px;

    .chat-icon {
      margin-right 20px;

      img {
        width: 36px;
        height: 36px;
        border-radius: 10px;
        padding: 1px;
      }
    }

    .chat-item {
      position: relative;
      padding: 0 5px 0 0;
      overflow: hidden;

      .content {
        word-break break-word;
        padding: 6px 10px;
        color #374151;
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

      .bar {
        padding 10px;

        .bar-item {
          background-color #f7f7f8;
          color #888
          padding 3px 5px;
          margin-right 10px;
          border-radius 5px;

          .el-icon {
            position relative
            top 2px;
          }
        }
      }
    }
  }


}
</style>