<template>
  <div class="container list" v-loading="loading">

    <div class="handle-box">
      <el-button type="primary" :icon="Plus" @click="add">新增</el-button>
    </div>

    <el-row>
      <el-table :data="items" :row-key="row => row.id" table-layout="auto">
        <el-table-column prop="platform" label="所属平台"/>
        <el-table-column prop="value" label="KEY"/>
        <el-table-column prop="type" label="用途">
          <template #default="scope">
            <el-tag v-if="scope.row.type === 'chat'">聊天</el-tag>
            <el-tag v-else-if="scope.row.type === 'img'" type="success">绘图</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="use_proxy" label="使用代理">
          <template #default="scope">
            <el-switch v-model="scope.row['use_proxy']" @change="set('use_proxy',scope.row)"/>
          </template>
        </el-table-column>

        <el-table-column label="最后使用时间">
          <template #default="scope">
            <span v-if="scope.row['last_used_at']">{{ scope.row['last_used_at'] }}</span>
            <el-tag v-else>未使用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="启用状态">
          <template #default="scope">
            <el-switch v-model="scope.row['enabled']" @change="set('enabled',scope.row)"/>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button size="small" type="primary" @click="edit(scope.row)">编辑</el-button>
            <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row)" :width="200">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

    <el-dialog
        v-model="showDialog"
        :title="title"
    >
      <el-alert
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 10px; font-size:14px;">
        <p><b>注意：</b>如果是百度文心一言平台，API-KEY 为 APIKey|SecretKey，中间用竖线（|）连接</p>
        <p><b>注意：</b>如果是讯飞星火大模型，API-KEY 为 AppId|APIKey|APISecret，中间用竖线（|）连接</p>
      </el-alert>
      <el-form :model="item" label-width="120px" ref="formRef" :rules="rules">
        <el-form-item label="所属平台：" prop="platform">
          <el-select v-model="item.platform" placeholder="请选择平台" @change="changePlatform">
            <el-option v-for="item in platforms" :value="item.value" :label="item.name" :key="item.value">{{
                item.name
              }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="用途：" prop="type">
          <el-select v-model="item.type" placeholder="请选择用途" @change="changePlatform">
            <el-option v-for="item in types" :value="item.value" :label="item.name" :key="item.value">{{
                item.name
              }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="API KEY：" prop="value">
          <el-input v-model="item.value" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="API URL：" prop="api_url">
          <el-input v-model="item.api_url" autocomplete="off"
                    placeholder="如果你用了第三方的 API 中转，这里填写中转地址"/>
        </el-form-item>

        <el-form-item label="使用代理：" prop="use_proxy">
          <el-switch v-model="item.use_proxy"/>
          <el-tooltip
              effect="dark"
              content="是否使用代理访问 API URL，OpenAI 官方API需要开启代理访问"
              raw-content
              placement="right"
          >
            <el-icon>
              <InfoFilled/>
            </el-icon>
          </el-tooltip>
        </el-form-item>

        <el-form-item label="启用状态：" prop="enable">
          <el-switch v-model="item.enabled"/>
        </el-form-item>
      </el-form>

      <template #footer>
            <span class="dialog-footer">
              <el-button @click="showDialog = false">取消</el-button>
              <el-button type="primary" @click="save">提交</el-button>
            </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat, disabledDate, removeArrayItem} from "@/utils/libs";
import {InfoFilled, Plus} from "@element-plus/icons-vue";

// 变量定义
const items = ref([])
const item = ref({})
const showDialog = ref(false)
const rules = reactive({
  platform: [{required: true, message: '请选择平台', trigger: 'change',}],
  type: [{required: true, message: '请选择用途', trigger: 'change',}],
  value: [{required: true, message: '请输入 API KEY 值', trigger: 'change',}]
})
const loading = ref(true)
const formRef = ref(null)
const title = ref("")
const platforms = ref([
  {
    name: "【OpenAI】ChatGPT",
    value: "OpenAI",
    api_url: "https://api.fast-tunnel.one/v1/chat/completions",
    img_url: "https://api.openai.com/v1/images/generations"
  },
  {
    name: "【讯飞】星火大模型",
    value: "XunFei",
    api_url: "wss://spark-api.xf-yun.com/{version}/chat"
  },
  {
    name: "【清华智普】ChatGLM",
    value: "ChatGLM",
    api_url: "https://open.bigmodel.cn/api/paas/v3/model-api/{model}/sse-invoke"
  },
  {
    name: "【百度】文心一言",
    value: "Baidu",
    api_url: "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/{model}"
  },
  {
    name: "【微软】Azure",
    value: "Azure",
    api_url: "https://chat-bot-api.openai.azure.com/openai/deployments/{model}/chat/completions?api-version=2023-05-15"
  },
])
const types = ref([
  {name: "聊天", value: "chat"},
  {name: "绘画", value: "img"},
])

// 获取数据
httpGet('/api/admin/apikey/list').then((res) => {
  if (res.data) {
    // 初始化数据
    const arr = res.data;
    for (let i = 0; i < arr.length; i++) {
      arr[i].last_used_at = dateFormat(arr[i].last_used_at)
    }
    items.value = arr
  }
  loading.value = false
}).catch(() => {
  ElMessage.error("获取数据失败");
})

const add = function () {
  showDialog.value = true
  title.value = "新增 API KEY"
  item.value = {}
}

const edit = function (row) {
  showDialog.value = true
  title.value = "修改 API KEY"
  item.value = row
}

const save = function () {
  formRef.value.validate((valid) => {
    if (valid) {
      showDialog.value = false
      httpPost('/api/admin/apikey/save', item.value).then((res) => {
        ElMessage.success('操作成功！')
        if (!item.value['id']) {
          const newItem = res.data
          newItem.last_used_at = dateFormat(newItem.last_used_at)
          items.value.push(newItem)
        }
      }).catch((e) => {
        ElMessage.error('操作失败，' + e.message)
      })
    } else {
      return false
    }
  })
}

const remove = function (row) {
  httpGet('/api/admin/apikey/remove?id=' + row.id).then(() => {
    ElMessage.success("删除成功！")
    items.value = removeArrayItem(items.value, row, (v1, v2) => {
      return v1.id === v2.id
    })
  }).catch((e) => {
    ElMessage.error("删除失败：" + e.message)
  })
}

const set = (filed, row) => {
  httpPost('/api/admin/apikey/set', {id: row.id, filed: filed, value: row[filed]}).then(() => {
    ElMessage.success("操作成功！")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

const changePlatform = () => {
  let platform = null
  for (let v of platforms.value) {
    if (v.value === item.value.platform) {
      platform = v
      break
    }
  }
  if (platform !== null) {
    if (item.value.type === "chat") {
      item.value.api_url = platform.api_url
    } else if (platform.img_url) {
      item.value.api_url = platform.img_url
    }

  }
}

</script>

<style lang="stylus" scoped>
.list {

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
}

.el-form {
  .el-form-item__content {

    .el-icon {
      padding-left 10px;
    }
  }
}
</style>