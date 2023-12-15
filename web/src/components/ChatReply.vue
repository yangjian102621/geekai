<template>
  <div class="chat-line chat-line-reply">
    <div class="chat-line-inner">
      <div class="chat-icon">
        <img :src="icon" alt="ChatGPT">
      </div>

      <div class="chat-item">
        <div class="content" v-html="content"></div>
        <div class="bar" v-if="createdAt !== ''">
          <span class="bar-item"><el-icon><Clock/></el-icon> {{ createdAt }}</span>
          <span class="bar-item">tokens: {{ tokens }}</span>
          <el-tooltip
              class="box-item"
              effect="light"
              content="复制回答"
              placement="top"
          >
            <el-button type="info" class="copy-reply" :data-clipboard-text="orgContent">
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
import {Clock, DocumentCopy, Position} from "@element-plus/icons-vue";

export default defineComponent({
  name: 'ChatReply',
  components: {Position, Clock, DocumentCopy},
  props: {
    content: {
      type: String,
      default: '',
    },
    orgContent: {
      type: String,
      default: '',
    },
    createdAt: {
      type: String,
      default: '',
    },
    tokens: {
      type: Number,
      default: 0,
    },
    icon: {
      type: String,
      default: 'images/gpt-icon.png',
    }
  },
  data() {
    return {
      finalTokens: this.tokens
    }
  }
})
</script>

<style lang="stylus">
.common-layout {
  .chat-line-reply {
    justify-content: center;
    background-color: rgba(247, 247, 248, 1);
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
        padding: 0 0 0 5px;
        overflow: hidden;

        .content {
          min-height 20px;
          word-break break-word;
          padding: 6px 10px;
          color #374151;
          font-size: var(--content-font-size);
          border-radius: 5px;
          overflow auto;

          // control the image size in content

          img {
            max-width: 600px;
            border-radius: 10px;
          }

          p {
            line-height 1.5

            code {
              color #374151
              background-color #e7e7e8
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

          .code-container {
            position relative

            .hljs {
              border-radius 10px
              line-height 1.5
            }

            .copy-code-btn {
              position: absolute;
              right 10px
              top 10px
              cursor pointer
              font-size 12px
              color #c1c1c1

              &:hover {
                color #20a0ff
              }
            }

          }

          .lang-name {
            position absolute;
            right 10px
            bottom 50px
            padding 2px 6px 4px 6px
            background-color #444444
            border-radius 10px
            color #00e0e0
          }

          // 设置表格边框

          table {
            width 100%
            margin-bottom 1rem
            color #212529
            border-collapse collapse;
            border 1px solid #dee2e6;
            background-color #ffffff

            thead {
              th {
                border 1px solid #dee2e6
                vertical-align: bottom
                border-bottom: 2px solid #dee2e6
                padding 10px
              }
            }

            td {
              border 1px solid #dee2e6
              padding 10px
            }
          }
        }


        .bar {
          padding 10px;

          .bar-item {
            background-color #e7e7e8;
            color #888
            padding 3px 5px;
            margin-right 10px;
            border-radius 5px;

            .el-icon {
              position relative
              top 2px;
            }
          }

          .el-button {
            height 20px
            padding 5px 2px;
          }
        }

      }

      .tool-box {
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
