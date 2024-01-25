<template>
  <div class="chat-export" v-loading="loading">
    <div class="chat-box" id="chat-box">
      <div class="title">
        <h2>{{ chatTitle }}</h2>
        <el-button type="success" @click="exportChat" :icon="Promotion">
          导出 PDF 文档
        </el-button>
      </div>

      <div v-for="item in chatData" :key="item.id">
        <chat-prompt
            v-if="item.type==='prompt'"
            :icon="item.icon"
            :created-at="dateFormat(item['created_at'])"
            :tokens="item['tokens']"
            :content="item.content"/>
        <chat-reply v-else-if="item.type==='reply'"
                    :icon="item.icon"
                    :org-content="item.orgContent"
                    :created-at="dateFormat(item['created_at'])"
                    :tokens="item['tokens']"
                    :content="item.content"/>
        <chat-mid-journey v-else-if="item.type==='mj'"
                          :content="item.content"
                          :icon="item.icon"
                          :created-at="dateFormat(item['created_at'])"/>
      </div>
    </div><!-- end chat box -->
  </div>
</template>
<script setup>

import {dateFormat} from "@/utils/libs";
import ChatReply from "@/components/ChatReply.vue";
import ChatPrompt from "@/components/ChatPrompt.vue";
import {nextTick, ref} from "vue";
import {useRouter} from "vue-router";
import {httpGet} from "@/utils/http";
import 'highlight.js/styles/a11y-dark.css'
import hl from "highlight.js";
import {ElMessage} from "element-plus";
import {Promotion} from "@element-plus/icons-vue";
import ChatMidJourney from "@/components/ChatMidJourney.vue";

const chatData = ref([])
const router = useRouter()
const chatId = router.currentRoute.value.query['chat_id']
const loading = ref(true)
const chatTitle = ref('')

httpGet('/api/chat/history?chat_id=' + chatId).then(res => {
  const data = res.data
  if (!data) {
    loading.value = false
    return
  }

  const md = require('markdown-it')({breaks: true});
  for (let i = 0; i < data.length; i++) {
    if (data[i].type === "prompt") {
      chatData.value.push(data[i]);
      continue;
    } else if (data[i].type === "mj") {
      data[i].content = JSON.parse(data[i].content)
      data[i].content.content = md.render(data[i].content?.content)
      chatData.value.push(data[i]);
      continue;
    }

    data[i].orgContent = data[i].content;
    data[i].content = md.render(data[i].content);
    chatData.value.push(data[i]);
  }

  nextTick(() => {
    hl.configure({ignoreUnescapedHTML: true})
    const blocks = document.querySelector("#chat-box").querySelectorAll('pre code');
    blocks.forEach((block) => {
      hl.highlightElement(block)
    })
  })
  loading.value = false
}).catch(e => {
  ElMessage.error('加载聊天记录失败：' + e.message);
})

httpGet('/api/chat/detail?chat_id=' + chatId).then(res => {
  chatTitle.value = res.data.title
}).catch(e => {
  ElMessage.error("加载会失败： " + e.message)
})

const exportChat = () => {
  window.print()
}
</script>
<style lang="stylus">
.chat-export {
  display flex
  justify-content center
  padding 0 20px

  .chat-box {
    width 800px;
    // 变量定义
    --content-font-size: 16px;
    --content-color: #c1c1c1;

    font-family: 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
    padding: 0 0 50px 0;

    .title {
      text-align center
    }


    .chat-line {
      font-size: 14px;
      display: flex;
      align-items: flex-start;

      .chat-line-inner {
        .content {
          padding-top: 0
          font-size 16px;

          p:first-child {
            margin-top 0
          }
        }
      }
    }

    .chat-line-reply {
      padding-top: 1.5rem;

      .chat-line-inner {
        display flex

        .copy-reply {
          display none
        }

        .bar-item {
          background-color: #f7f7f8;
          color: #888;
          padding: 3px 5px;
          margin-right: 10px;
          border-radius: 5px;
        }

        .chat-icon {
          margin-right: 20px

          img {
            width 30px
            height 30px
            border-radius: 10px;
            padding: 1px
          }
        }

        .chat-item {
          img {
            max-width 90%
          }
        }
      }
    }
  }
}
</style>