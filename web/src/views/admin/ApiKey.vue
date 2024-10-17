<template>
  <div class="container list" v-loading="loading">

    <div class="handle-box">
      <el-select v-model="query.type" placeholder="类型" class="handle-input">
        <el-option
            v-for="item in types"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>

      <el-button :icon="Search" @click="fetchData">搜索</el-button>
      <el-button type="primary" :icon="Plus" @click="add">新增</el-button>
      <a href="https://api.chat-plus.net" target="_blank" style="margin-left: 10px">
        <el-button type="success" :icon="ShoppingCart" @click="add" plain>购买API-KEY</el-button>
      </a>
    </div>

    <el-row>
      <el-table :data="items" :row-key="row => row.id" table-layout="auto">
        <el-table-column prop="name" label="名称"/>
        <el-table-column prop="value" label="API KEY">
          <template #default="scope">
            <span>{{ substr(scope.row.value, 20) }}</span>
            <el-icon class="copy-key" :data-clipboard-text="scope.row.value">
              <DocumentCopy/>
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="api_url" label="API URL">
          <template #default="scope">
            <span>{{ substr(scope.row.api_url, 30) }}</span>
            <el-icon class="copy-key" :data-clipboard-text="scope.row.api_url">
              <DocumentCopy/>
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型">
          <template #default="scope">
            {{getTypeName(scope.row.type)}}
          </template>
        </el-table-column>
        <el-table-column prop="proxy_url" label="代理地址"/>

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
        :close-on-click-modal="false"
        :title="title"
    >
      <el-form :model="item" label-width="120px" ref="formRef" :rules="rules">
        <el-form-item label="名称：" prop="name">
          <el-input v-model="item.name" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="类型：" prop="type">
          <el-select v-model="item.type" placeholder="请选择类型" @change="changeType">
            <el-option v-for="item in types" :value="item.value" :label="item.label" :key="item.value">{{
                item.label
              }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="API KEY：" prop="value">
          <el-input v-model="item.value" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="API URL：" prop="api_url">
          <el-input v-model="item.api_url" autocomplete="off"
                    placeholder="只填 BASE URL 即可，如：https://api.openai.com 或者 wss://api.openai.com"/>
        </el-form-item>

<!--        <el-form-item label="代理地址：" prop="proxy_url">-->
<!--          <el-input v-model="item.proxy_url" autocomplete="off"/>-->
<!--          <div class="info">如果想要通过代理来访问 API，请填写代理地址，如：http://127.0.0.1:7890</div>-->
<!--        </el-form-item>-->

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
import {onMounted, onUnmounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat, removeArrayItem, substr} from "@/utils/libs";
import {DocumentCopy, Plus, ShoppingCart, Search} from "@element-plus/icons-vue";
import ClipboardJS from "clipboard";

// 变量定义
const items = ref([])
const query = ref({type: ''})
const item = ref({
  enabled: true, api_url: ""
})
const showDialog = ref(false)
const rules = reactive({
  name: [{required: true, message: '请输入名称', trigger: 'change',}],
  type: [{required: true, message: '请选择用途', trigger: 'change',}],
  value: [{required: true, message: '请输入 API KEY 值', trigger: 'change',}]
})

const loading = ref(true)
const formRef = ref(null)
const title = ref("")
const types = ref([
  {label: "对话", value:"chat"},
  {label: "Midjourney", value:"mj"},
  {label: "Stable-Diffusion", value:"sd"},
  {label: "DALL-E", value:"dalle"},
  {label: "Suno文生歌", value:"suno"},
  {label: "Luma视频", value:"luma"},
  {label: "Realtime API", value:"realtime"},
])
const isEdit = ref(false)
const clipboard = ref(null)
onMounted(() => {
  clipboard.value = new ClipboardJS('.copy-key');
  clipboard.value.on('success', () => {
    ElMessage.success('复制成功！');
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！');
  })

  fetchData()
})

onUnmounted(() => {
  clipboard.value.destroy()
})

const changeType = (event) => {
  if (isEdit.value) {
    return
  }
  
  if (event === 'realtime') {
    item.value.api_url = "wss://api.geekai.pro"
  } else {
    item.value.api_url = "https://api.geekai.pro"
  }
}

const getTypeName = (type) => {
  for (let v of types.value) {
    if (v.value === type) {
      return v.label
    }
  }
  return ""
}

// 获取数据

const fetchData = () => {
  httpGet('/api/admin/apikey/list', query.value).then((res) => {
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
}

const add = function () {
  showDialog.value = true
  title.value = "新增 API KEY"
  isEdit.value = false
}

const edit = function (row) {
  showDialog.value = true
  title.value = "修改 API KEY"
  item.value = row
  isEdit.value = true
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

</script>

<style lang="stylus" scoped>
.list {
  .handle-box {
    margin-bottom 20px
    .handle-input {
      max-width 150px;
      margin-right 10px;
    }
  }

  .copy-key {
    margin-left 5px
    cursor pointer
  }

  .el-select {
    width: 100%
  }

  .pagination {
    padding 20px 0
    display flex
    justify-content righ
  }
}

.el-form {
  .el-form-item__content {
    .info {
      color #999999
    }
    .el-icon {
      padding-left: 10px;
    }
  }
}
</style>