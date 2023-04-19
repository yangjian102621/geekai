<template>
  <div class="chat-line chat-line-left">
    <div class="chat-icon">
      <img :src="icon" alt="ChatGPT">
    </div>

    <div class="chat-item">
      <div class="triangle"></div>
      <div class="content-box">
        <div class="content" v-html="content"></div>
        <div class="tool-box">
          <el-tooltip
              class="box-item"
              effect="dark"
              content="复制回答"
              placement="top"
          >
            <el-button type="info" class="copy-reply" :data-clipboard-text="orgContent" plain>
              <el-icon>
                <DocumentCopy/>
              </el-icon>
            </el-button>
          </el-tooltip>
        </div>
      </div>

    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue"
import {randString} from "@/utils/libs";
import {DocumentCopy} from "@element-plus/icons-vue";

export default defineComponent({
  name: 'ChatReply',
  components: {DocumentCopy},
  props: {
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
      default: 'images/gpt-icon.png',
    }
  },
  data() {
    return {
      id: randString(32),
      clipboard: null,
    }
  },

})
</script>

<style lang="stylus">
.chat-line-left {
  justify-content: flex-start;

  .chat-icon {
    margin-right 5px;

    img {
      border-radius 5px;
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
        min-height 20px;
        word-break break-word;
        padding: 8px 10px;
        color var(--content-color)
        background-color: #fff;
        font-size: var(--content-font-size);
        border-radius: 5px;

        p:last-child {
          margin-bottom: 0
        }

        p:first-child {
          margin-top 0
        }

        p > code {
          color #cc0000
          background-color #f1f1f1
        }
      }

      .tool-box {
        padding-left 10px;
        font-size 16px;

        .el-button {
          height 20px
          padding 5px 2px;
        }
      }
    }

  }
}

</style>
