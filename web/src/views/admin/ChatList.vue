<template>
  <div class="container chat-list">
    <el-tabs v-model="activeName" @tab-change="handleChange">
      <el-tab-pane label="对话列表" name="chat" v-loading="data.chat.loading">
        <div class="handle-box">
          <el-input v-model.number="data.chat.query.user_id" placeholder="账户ID" class="handle-input mr10"
                    @keyup="searchChat($event)"></el-input>
          <el-input v-model="data.chat.query.title" placeholder="对话标题" class="handle-input mr10"
                    @keyup="searchChat($event)"></el-input>
          <el-date-picker
              v-model="data.chat.query.created_at"
              type="daterange"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="margin-right: 10px;width: 200px; position: relative;top:3px;"
          />
          <el-button type="primary" :icon="Search" @click="fetchChatData">搜索</el-button>
        </div>

        <el-row>
          <el-table :data="data.chat.items" :row-key="row => row.id" table-layout="auto">
            <el-table-column prop="user_id" label="账户ID"/>
            <el-table-column prop="username" label="账户"/>
            <el-table-column label="图标">
              <template #default="scope">
                <el-avatar :size="30" :src="scope.row.role.icon"/>
              </template>
            </el-table-column>
            <el-table-column label="角色">
              <template #default="scope">
                <span>{{ scope.row.role.name }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="model" label="模型"/>
            <el-table-column prop="title" label="标题"/>
            <el-table-column prop="msg_num" label="消息数量"/>
            <el-table-column prop="token" label="消耗算力"/>

            <el-table-column label="创建时间">
              <template #default="scope">
                <span>{{ dateFormat(scope.row['created_at']) }}</span>
              </template>
            </el-table-column>

            <el-table-column label="操作" width="180">
              <template #default="scope">
                <el-button size="small" type="primary" @click="showMessages(scope.row)">查看</el-button>
                <el-popconfirm title="确定要删除当前记录吗?" @confirm="removeChat(scope.row)">
                  <template #reference>
                    <el-button size="small" type="danger">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-row>

        <div class="pagination">
          <el-pagination v-if="data.chat.total > 0" background
                         layout="total,prev, pager, next"
                         :hide-on-single-page="true"
                         v-model:current-page="data.chat.page"
                         v-model:page-size="data.chat.pageSize"
                         @current-change="fetchChatData()"
                         :total="data.chat.total"/>

        </div>
      </el-tab-pane>
      <el-tab-pane label="消息记录" name="message">
        <div class="handle-box">
          <el-input v-model.number="data.message.query.user_id" placeholder="账户ID" class="handle-input mr10"
                    @keyup="searchMessage($event)"></el-input>
          <el-input v-model="data.message.query.content" placeholder="消息内容" class="handle-input mr10"
                    @keyup="searchMessage($event)"></el-input>
          <el-input v-model="data.message.query.model" placeholder="模型" class="handle-input mr10"
                    @keyup="searchMessage($event)"></el-input>
          <el-date-picker
              v-model="data.message.query.created_at"
              type="daterange"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="margin-right: 10px;width: 200px; position: relative;top:3px;"
          />
          <el-button type="primary" :icon="Search" @click="fetchMessageData">搜索</el-button>
        </div>

        <el-row>
          <el-table :data="data.message.items" :row-key="row => row.id" table-layout="auto">
            <el-table-column prop="user_id" label="账户ID"/>
            <el-table-column prop="username" label="账户"/>
            <el-table-column label="角色">
              <template #default="scope">
                <el-avatar :size="30" :src="scope.row.icon"/>
              </template>
            </el-table-column>
            <el-table-column prop="model" label="模型"/>

            <el-table-column label="消息内容">
              <template #default="scope">
                <el-text style="width: 200px" truncated @click="showContent(scope.row.content)">
                  {{ scope.row.content }}
                </el-text>
              </template>
            </el-table-column>

            <el-table-column prop="token" label="算力"/>

            <el-table-column label="创建时间">
              <template #default="scope">
                <span>{{ dateFormat(scope.row['created_at']) }}</span>
              </template>
            </el-table-column>

            <el-table-column label="操作" width="180">
              <template #default="scope">
                <el-button size="small" type="primary" @click="showContent(scope.row.content)">查看</el-button>
                <el-popconfirm title="确定要删除当前记录吗?" @confirm="removeMessage(scope.row)">
                  <template #reference>
                    <el-button size="small" type="danger">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-row>

        <div class="pagination">
          <el-pagination v-if="data.message.total > 0" background
                         layout="total,prev, pager, next"
                         :hide-on-single-page="true"
                         v-model:current-page="data.message.page"
                         v-model:page-size="data.message.pageSize"
                         @current-change="fetchMessageData()"
                         :total="data.message.total"/>

        </div>
      </el-tab-pane>
    </el-tabs>


    <el-dialog
        v-model="showContentDialog"
        title="消息详情"
    >
      <div v-html="dialogContent" style="overflow: auto; max-height: 300px"></div>
    </el-dialog>

    <el-dialog
        v-model="showChatItemDialog"
        title="对话详情"
    >
      <div class="chat-box common-layout">
        <div v-for="item in messages" :key="item.id">
          <chat-prompt
              v-if="item.type==='prompt'"
              :icon="item.icon"
              :created-at="dateFormat(item['created_at'])"
              :tokens="item['tokens']"
              :model="item.model"
              :content="item.content"/>
          <chat-reply v-else-if="item.type==='reply'"
                      :icon="item.icon"
                      :org-content="item.orgContent"
                      :created-at="dateFormat(item['created_at'])"
                      :tokens="item['tokens']"
                      :content="item.content"/>
        </div>
      </div><!-- end chat box -->
    </el-dialog>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat, processContent, removeArrayItem} from "@/utils/libs";
import {Search} from "@element-plus/icons-vue";
import 'highlight.js/styles/a11y-dark.css'
import hl from "highlight.js";
import ChatPrompt from "@/components/ChatPrompt.vue";
import ChatReply from "@/components/ChatReply.vue";

// 变量定义
const data = ref({
  "chat": {
    items: [],
    query: {title: "", created_at: [], page: 1, page_size: 15},
    total: 0,
    page: 1,
    pageSize: 15,
    loading: true
  },
  "message": {
    items: [],
    query: {title: "", created_at: [], page: 1, page_size: 15},
    total: 0,
    page: 1,
    pageSize: 15,
    loading: true
  }
})
const items = ref([])
const query = ref({title: "", created_at: []})
const total = ref(0)
const page = ref(1)
const pageSize = ref(15)
const loading = ref(true)
const activeName = ref("chat")

onMounted(() => {
  fetchChatData()
})

const handleChange = (tab) => {
  if (tab === "chat") {
    fetchChatData()
  } else if (tab === "message") {
    fetchMessageData()
  }
}

// 搜索对话
const searchChat = (evt) => {
  if (evt.keyCode === 13) {
    fetchChatData()
  }
}

// 搜索消息
const searchMessage = (evt) => {
  if (evt.keyCode === 13) {
    fetchMessageData()
  }
}

// 获取数据
const fetchChatData = () => {
  const d = data.value.chat
  d.query.page = d.page
  d.query.page_size = d.pageSize
  httpPost('/api/admin/chat/list', d.query).then((res) => {
    if (res.data) {
      d.items = res.data.items
      d.total = res.data.total
      d.page = res.data.page
      d.pageSize = res.data.page_size
    }
    d.loading = false
  }).catch(e => {
    ElMessage.error("获取数据失败：" + e.message);
  })
}

const fetchMessageData = () => {
  const d = data.value.message
  d.query.page = d.page
  d.query.page_size = d.pageSize
  httpPost('/api/admin/chat/message', d.query).then((res) => {
    if (res.data) {
      d.items = res.data.items
      d.total = res.data.total
      d.page = res.data.page
      d.pageSize = res.data.page_size
    }
    d.loading = false
  }).catch(e => {
    ElMessage.error("获取数据失败：" + e.message);
  })
}

const removeChat = function (row) {
  httpGet('/api/admin/chat/remove?chat_id=' + row.chat_id).then(() => {
    ElMessage.success("删除成功！")
    fetchChatData()
  }).catch((e) => {
    ElMessage.error("删除失败：" + e.message)
  })
}

const removeMessage = function (row) {
  httpGet('/api/admin/chat/message/remove?id=' + row.id).then(() => {
    ElMessage.success("删除成功！")
    fetchMessageData()
  }).catch((e) => {
    ElMessage.error("删除失败：" + e.message)
  })
}

const latexPlugin = require('markdown-it-latex2img')
const mathjaxPlugin = require('markdown-it-mathjax')
const md = require('markdown-it')({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hl.getLanguage(lang)) {
      // 处理代码高亮
      const preCode = hl.highlight(lang, str, true).value
      // 将代码包裹在 pre 中
      return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code></pre>`
    }

    // 处理代码高亮
    const preCode = md.utils.escapeHtml(str)
    // 将代码包裹在 pre 中
    return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code></pre>`
  }
});
md.use(latexPlugin)
md.use(mathjaxPlugin)

const showContentDialog = ref(false)
const dialogContent = ref("")
const showContent = (content) => {
  showContentDialog.value = true
  dialogContent.value = md.render(processContent(content))
}

const showChatItemDialog = ref(false)
const messages = ref([])
const showMessages = (row) => {
  showChatItemDialog.value = true
  messages.value = []
  httpGet('/api/admin/chat/history?chat_id=' + row.chat_id).then(res => {
    const data = res.data
    for (let i = 0; i < data.length; i++) {
      data[i].orgContent = data[i].content;
      data[i].content = md.render(processContent(data[i].content))
      messages.value.push(data[i]);
    }
  }).catch(e => {
    // TODO: 显示重新加载按钮
    ElMessage.error('加载聊天记录失败：' + e.message);
  })
}
</script>

<style lang="stylus" scoped>
.chat-list {

  .handle-box {
    .handle-input {
      max-width 150px;
      margin-right 10px;
    }
  }

  .opt-box {
    padding-bottom: 10px;
    display flex;
    justify-content flex-end

    .el-icon {
      margin-right: 5px;
    }
  }

  .el-select {
    width: 100%
  }

  .chat-box {
    overflow-y: auto;
    overflow-x hidden

    // 变量定义
    --content-font-size: 16px;
    --content-color: #c1c1c1;

    font-family: 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
    height 90vh

    .chat-line {
      // 隐藏滚动条

      ::-webkit-scrollbar {
        width: 0;
        height: 0;
        background-color: transparent;
      }

      font-size: 14px;
      display: flex;
      align-items: flex-start;

    }
  }

}
</style>