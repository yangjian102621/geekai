<template>
  <div class="chat-export-mobile">
    <div class="chat-box">
      <van-sticky :offset-top="0" position="top">
        <van-nav-bar left-arrow left-text="返回" @click-left="router.back()">
          <template #title>
            <van-dropdown-menu>
              <van-dropdown-item :title="title">
                <van-cell center title="角色"> {{ role }}</van-cell>
                <van-cell center title="模型">{{ model }}</van-cell>
              </van-dropdown-item>
            </van-dropdown-menu>
          </template>
        </van-nav-bar>
      </van-sticky>

      <div class="chat-list-wrapper">
        <div id="message-list-box" class="message-list-box">
          <van-list
              v-model:error="error"
              :finished="finished"
              error-text="请求失败，点击重新加载"
              @load="onLoad"
          >
            <van-cell v-for="item in chatData" :key="item" :border="false" class="message-line">
              <chat-prompt
                  v-if="item.type==='prompt'"
                  :content="item.content"
                  :created-at="dateFormat(item['created_at'])"
                  :icon="item.icon"
                  :tokens="item['tokens']"/>
              <chat-reply v-else-if="item.type==='reply'"
                          :content="item.content"
                          :created-at="dateFormat(item['created_at'])"
                          :icon="item.icon"
                          :org-content="item.orgContent"
                          :tokens="item['tokens']"/>
            </van-cell>
          </van-list>
        </div>
      </div>

    </div><!-- end chat box -->
  </div>
</template>
<script setup>

import {dateFormat, processContent} from "@/utils/libs";
import ChatReply from "@/components/mobile/ChatReply.vue";
import ChatPrompt from "@/components/mobile/ChatPrompt.vue";
import {nextTick, ref} from "vue";
import {useRouter} from "vue-router";
import {httpGet} from "@/utils/http";
import 'highlight.js/styles/a11y-dark.css'
import hl from "highlight.js";

const chatData = ref([])
const router = useRouter()
const chatId = router.currentRoute.value.query['chat_id']
const title = router.currentRoute.value.query['title']
const role = router.currentRoute.value.query['role']
const model = router.currentRoute.value.query['model']
const finished = ref(false)
const error = ref(false)

const md = require('markdown-it')({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    const codeIndex = parseInt(Date.now()) + Math.floor(Math.random() * 10000000)
    // 显示复制代码按钮
    const copyBtn = `<span class="copy-code-mobile" data-clipboard-action="copy" data-clipboard-target="#copy-target-${codeIndex}">复制</span>
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

const onLoad = () => {
  httpGet('/api/chat/history?chat_id=' + chatId).then(res => {
    // 加载状态结束
    finished.value = true;
    const data = res.data
    if (data && data.length > 0) {
      for (let i = 0; i < data.length; i++) {
        if (data[i].type === "prompt") {
          chatData.value.push(data[i]);
          continue;
        }

        data[i].orgContent = data[i].content;
        data[i].content = md.render(processContent(data[i].content))
        chatData.value.push(data[i]);
      }

      nextTick(() => {
        hl.configure({ignoreUnescapedHTML: true})
        const blocks = document.querySelector("#message-list-box").querySelectorAll('pre code');
        blocks.forEach((block) => {
          hl.highlightElement(block)
        })
      })
    }
  }).catch(() => {
    error.value = true
  })

};

</script>
<style lang="stylus">
.chat-export-mobile {
  background #F5F5F5;
  height 100vh

  .chat-box {
    font-family: 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;

    .message-list-box {
      background #F5F5F5;
      padding-top 50px
      padding-bottom: 10px

      .van-cell {
        background none
      }
    }


    .van-nav-bar__title {
      .van-dropdown-menu__title {
        margin-right 10px
      }

      .van-cell__title {
        text-align left
      }
    }

    .van-nav-bar__right {
      .van-icon {
        font-size 20px
      }
    }
  }
}
</style>