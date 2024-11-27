<template>
  <div class="chat-line chat-line-reply-list" v-if="listStyle === 'list'">
    <div class="chat-line-inner">
      <div class="chat-icon">
        <img :src="data.icon" alt="ChatGPT">
      </div>

      <div class="chat-item">
        <div class="content" v-html="md.render(processContent(data.content))"></div>
        <div class="bar" v-if="data.created_at">
          <span class="bar-item"><el-icon><Clock/></el-icon> {{ dateFormat(data.created_at) }}</span>
                    <span class="bar-item">tokens: {{ data.tokens }}</span>
          <span class="bar-item">
              <el-tooltip
                  class="box-item"
                  effect="dark"
                  content="复制回答"
                  placement="bottom"
              >
                <el-icon class="copy-reply" :data-clipboard-text="data.content">
                  <DocumentCopy/>
                </el-icon>
              </el-tooltip>
            </span>
          <span v-if="!readOnly">
            <span class="bar-item" @click="reGenerate(data.prompt)">
            <el-tooltip
                class="box-item"
                effect="dark"
                content="重新生成"
                placement="bottom"
            >
              <el-icon><Refresh/></el-icon>
            </el-tooltip>
          </span>

          <span class="bar-item" @click="synthesis(data.content)">
            <el-tooltip
                class="box-item"
                effect="dark"
                content="生成语音朗读"
                placement="bottom"
            >
              <i class="iconfont icon-speaker"></i>
            </el-tooltip>
          </span>
          </span>
          <!--          <span class="bar-item">-->
          <!--            <el-dropdown trigger="click">-->
          <!--              <span class="el-dropdown-link">-->
          <!--                <el-icon><More/></el-icon>-->
          <!--              </span>-->
          <!--              <template #dropdown>-->
          <!--                <el-dropdown-menu>-->
          <!--                  <el-dropdown-item :icon="Headset" @click="synthesis(orgContent)">生成语音</el-dropdown-item>-->
          <!--                </el-dropdown-menu>-->
          <!--              </template>-->
          <!--            </el-dropdown>-->
          <!--          </span>-->
        </div>
      </div>
    </div>
  </div>

  <div class="chat-line chat-line-reply-chat" v-else>
    <div class="chat-line-inner">
      <div class="chat-icon">
        <img :src="data.icon" alt="ChatGPT">
      </div>
      <div class="chat-item">
        <div class="content-wrapper">
          <div class="content" v-html="md.render(processContent(data.content))"></div>
        </div>
        <div class="bar" v-if="data.created_at">
          <span class="bar-item"><el-icon><Clock/></el-icon> {{ dateFormat(data.created_at) }}</span>
<!--          <span class="bar-item">tokens: {{ data.tokens }}</span>-->
          <span class="bar-item bg">
              <el-tooltip
                  class="box-item"
                  effect="dark"
                  content="复制回答"
                  placement="bottom"
              >
                <el-icon class="copy-reply" :data-clipboard-text="data.content">
                  <DocumentCopy/>
                </el-icon>
              </el-tooltip>
            </span>
          <span v-if="!readOnly">
            <span class="bar-item bg" @click="reGenerate(data.prompt)">
            <el-tooltip
                class="box-item"
                effect="dark"
                content="重新生成"
                placement="bottom"
            >
              <el-icon><Refresh/></el-icon>
            </el-tooltip>
          </span>

          <span class="bar-item bg" @click="synthesis(data.content)">
            <el-tooltip
                class="box-item"
                effect="dark"
                content="生成语音朗读"
                placement="bottom"
            >
              <i class="iconfont icon-speaker"></i>
            </el-tooltip>
          </span>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {Clock, DocumentCopy, Refresh} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";
import {dateFormat, processContent} from "@/utils/libs";
import hl from "highlight.js";
// eslint-disable-next-line no-undef,no-unused-vars
const props = defineProps({
  data: {
    type: Object,
    default: {
      icon: "",
      content: "",
      created_at: "",
      tokens: 0,
    },
  },
  readOnly: {
    type: Boolean,
    default: false
  },
  listStyle: {
    type: String,
    default: 'list',
  },
})

const mathjaxPlugin = require('markdown-it-mathjax3')
const md = require('markdown-it')({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    const codeIndex = parseInt(Date.now()) + Math.floor(Math.random() * 10000000)
    // 显示复制代码按钮
    const copyBtn = `<span class="copy-code-btn" data-clipboard-action="copy" data-clipboard-target="#copy-target-${codeIndex}">复制</span>
<textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy-target-${codeIndex}">${str.replace(/<\/textarea>/g, '&lt;/textarea>')}</textarea>`
    if (lang && hl.getLanguage(lang)) {
      const langHtml = `<span class="lang-name">${lang}</span>`
      // 处理代码高亮
      const preCode = hl.highlight(lang, str, true).value
      // 将代码包裹在 pre 中
      return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code>${copyBtn} ${langHtml}</pre>`
    }

    // 处理代码高亮
    const preCode = md.utils.escapeHtml(str)
    // 将代码包裹在 pre 中
    return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code>${copyBtn}</pre>`
  }
});
md.use(mathjaxPlugin)

const emits = defineEmits(['regen']);

if (!props.data.icon) {
  props.data.icon = "images/gpt-icon.png"
}

const synthesis = (text) => {
  console.log(text)
  ElMessage.info("语音合成功能暂不可用")
}

// 重新生成
const reGenerate = (prompt) => {
  console.log(prompt)
  emits('regen', prompt)
}
</script>

<style lang="stylus">
@import '@/assets/css/markdown/vue.css';
.chat-page,.chat-export {
  .chat-line-reply-list {
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
        width 100%
        position: relative;
        padding: 0;
        overflow: hidden;

        .content {
          min-height 20px;
          word-break break-word;
          padding: 0
          color #374151;
          font-size: var(--content-font-size);
          border-radius: 5px;
          overflow auto;

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
            display flex

            .hljs {
              border-radius 10px
              width 100%
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
            bottom 20px
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

          // 代码快

          blockquote {
            margin 0
            background-color: #ebfffe;
            padding: 0.8rem 1.5rem;
            border-left: 0.5rem solid;
            border-color: #026863;
            color: #2c3e50;
          }
        }


        .bar {
          padding 10px 10px 10px 0;

          .bar-item {
            background-color #e7e7e8;
            color #888
            padding 3px 5px;
            margin-right 10px;
            border-radius 5px;
            cursor pointer

            .el-icon {
              position relative
              top 2px;
              cursor pointer
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

  .chat-line-reply-chat {
    justify-content: center;
    padding 1.5rem;

    .chat-line-inner {
      display flex;
      width 100%
      flex-flow row

      .chat-icon {
        margin-right 20px;

        img {
          width: 36px;
          height: 36px;
          border-radius: 50%
          padding: 1px;
        }
      }

      .chat-item {
        position: relative;
        padding: 0;
        overflow: hidden;
        max-width 70%

        .content-wrapper {
          display flex
          .content {
            min-height 20px;
            word-break break-word;
            padding: 1rem
            color #374151;
            font-size: var(--content-font-size);
            overflow auto;
            background-color #F5F5F5
            border-radius: 0 10px 10px 10px;

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
              display flex

              .hljs {
                border-radius 10px
                width 100%
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
              bottom 20px
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

            // 代码快

            blockquote {
              margin 0
              background-color: #ebfffe;
              padding: 0.8rem 1.5rem;
              border-left: 0.5rem solid;
              border-color: #026863;
              color: #2c3e50;
            }
          }

        }

        .bar {
          padding 10px 10px 10px 0;

          .bar-item {
            color #888
            padding 3px 5px;
            margin-right 10px;
            border-radius 5px;

            .el-icon {
              position relative
              top 2px;
              cursor pointer
            }
          }

          .bar-item.bg {
            background-color #e7e7e8
            cursor pointer
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
