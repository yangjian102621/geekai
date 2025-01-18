<template>
  <div class="chat-line chat-line-prompt-list" v-if="listStyle === 'list'">
    <div class="chat-line-inner">
        <div class="chat-icon">
          <img :src="data.icon" alt="User"/>
        </div>

        <div class="chat-item">
          <div v-if="files.length > 0" class="file-list-box">
            <div v-for="file in files">
              <div class="image" v-if="isImage(file.ext)">
                <el-image :src="file.url" fit="cover"/>
              </div>
              <div class="item" v-else>
                <div class="icon">
                  <el-image :src="GetFileIcon(file.ext)" fit="cover"  />
                </div>
                <div class="body">
                  <div class="title">
                    <el-link :href="file.url" target="_blank" style="--el-font-weight-primary:bold">{{file.name}}</el-link>
                  </div>
                  <div class="info">
                    <span>{{GetFileType(file.ext)}}</span>
                    <span>{{FormatFileSize(file.size)}}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="content" v-html="content"></div>
          <div class="bar" v-if="data.created_at > 0">
            <span class="bar-item"><el-icon><Clock/></el-icon> {{ dateFormat(data.created_at) }}</span>
            <span class="bar-item">tokens: {{ finalTokens }}</span>
          </div>
        </div>
      </div>
  </div>

  <div class="chat-line chat-line-prompt-chat" v-else>
    <div class="chat-line-inner">
      <div class="chat-icon">
        <img :src="data.icon" alt="User"/>
      </div>

      <div class="chat-item">

        <div v-if="files.length > 0" class="file-list-box">
          <div v-for="file in files">
            <div class="image" v-if="isImage(file.ext)">
              <el-image :src="file.url" fit="cover"/>
            </div>
            <div class="item" v-else>
              <div class="icon">
                <el-image :src="GetFileIcon(file.ext)" fit="cover"  />
              </div>
              <div class="body">
                <div class="title">
                  <el-link :href="file.url" target="_blank" style="--el-font-weight-primary:bold">{{file.name}}</el-link>
                </div>
                <div class="info">
                  <span>{{GetFileType(file.ext)}}</span>
                  <span>{{FormatFileSize(file.size)}}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="content-wrapper">
          <div class="content" v-html="content"></div>
        </div>
        <div class="bar" v-if="data.created_at > 0">
          <span class="bar-item"><el-icon><Clock/></el-icon> {{ dateFormat(data.created_at) }}</span>
<!--          <span class="bar-item">tokens: {{ finalTokens }}</span>-->
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {Clock} from "@element-plus/icons-vue";
import {httpPost} from "@/utils/http";
import hl from "highlight.js";
import {dateFormat, isImage, processPrompt} from "@/utils/libs";
import {FormatFileSize, GetFileIcon, GetFileType} from "@/store/system";

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
const props = defineProps({
  data: {
    type: Object,
    default: {
      content: '',
      created_at: '',
      tokens: 0,
      model: '',
      icon: '',
    },
  },
  listStyle: {
    type: String,
    default: 'list',
  },
})
const finalTokens = ref(props.data.tokens)
const content =ref(processPrompt(props.data.content))
const files = ref([])

onMounted(() => {
  processFiles()
})

const processFiles = () => {
  if (!props.data.content) {
    return
  }

  const linkRegex = /(https?:\/\/\S+)/g;
  const links = props.data.content.match(linkRegex);
  if (links) {
    httpPost("/api/upload/list", {urls: links}).then(res => {
      files.value = res.data.items

      for (let link of links) {
        if (isExternalImg(link, files.value)) {
          files.value.push({url:link, ext: ".png"})
        }
      }
    }).catch(() => {
    })

    for (let link of links) {
      content.value = content.value.replace(link,"")
    }

  }
  content.value = md.render(content.value.trim())
}
const isExternalImg = (link, files) => {
  return isImage(link) && !files.find(file => file.url === link)
}
</script>

<style lang="stylus">
@import '@/assets/css/markdown/vue.css';
.chat-page,.chat-export {
  .chat-line-prompt-list {
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
        width 100%
        padding: 0 5px 0 0;
        overflow: hidden;

        .file-list-box {
          display flex
          flex-flow column
          .image {
            display flex
            flex-flow row
            margin-right 10px
            position relative

            .el-image {
              border 1px solid #e3e3e3
              border-radius 10px
              margin-bottom 10px
            }
          }
          .item {
            display flex
            flex-flow row
            border-radius 10px
            background-color #ffffff
            border 1px solid #e3e3e3
            padding 6px
            margin-bottom 10px

            .icon {
              .el-image {
                width 40px
                height 40px
              }
            }
            .body {
              margin-left 8px
              font-size 14px
              .title {
                font-weight bold
                line-height 24px
                color #0D0D0D
              }
              .info {
                color #B4B4B4

                span {
                  margin-right 10px
                }
              }
            }
          }
        }

        .content {
          word-break break-word;
          padding: 0;
          color #374151;
          font-size: var(--content-font-size);
          border-radius: 5px;
          overflow: auto;

          img {
            max-width: 600px;
            border-radius: 10px;
            margin 10px 0
          }

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
          padding 10px 10px 10px 0;

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

  .chat-line-prompt-chat {
    background-color #ffffff;
    justify-content: center;
    width 100%
    padding-bottom: 1.5rem;
    padding-top: 1.5rem;

    .chat-line-inner {
      display flex;
      width 100%;
      padding 0 25px;
      flex-flow row-reverse

      .chat-icon {
        margin-left 20px;

        img {
          width: 36px;
          height: 36px;
          border-radius: 50%;
          padding: 1px;
        }
      }

      .chat-item {
        padding: 0;
        overflow: hidden;
        max-width 60%

        .file-list-box {
          display flex
          flex-flow column
          .image {
            display flex
            flex-flow row
            margin-right 10px
            position relative

            .el-image {
              border 1px solid #e3e3e3
              border-radius 10px
              margin-bottom 10px
            }
          }
          .item {
            display flex
            flex-flow row
            border-radius 10px
            background-color #ffffff
            border 1px solid #e3e3e3
            padding 6px
            margin-bottom 10px

            .icon {
              .el-image {
                width 40px
                height 40px
              }
            }
            .body {
              margin-left 8px
              font-size 14px
              .title {
                font-weight bold
                line-height 24px
                color #0D0D0D
              }
              .info {
                color #B4B4B4

                span {
                  margin-right 10px
                }
              }
            }
          }
        }


        .content-wrapper {
          display flex
          flex-flow row-reverse
          .content {
              word-break break-word;
              padding: 1rem
              color #222222;
              font-size: var(--content-font-size);
              overflow: auto;
              background-color #98e165
              border-radius: 10px 0 10px 10px;

              img {
                max-width: 600px;
                border-radius: 10px;
                margin 10px 0
              }

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
            }
          }
        }
      }
    }


  }

}
</style>
