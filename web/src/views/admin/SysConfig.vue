<template>
  <div class="system-config" v-loading="loading">
    <el-form :model="form" label-width="120px">
      <el-form-item label="应用标题">
        <el-input v-model="form['title']"/>
      </el-form-item>
      <el-form-item label="控制台标题">
        <el-input v-model="form['console_title']"/>
      </el-form-item>
      <el-form-item label="代理地址">
        <el-input v-model="form['proxy_url']" placeholder="多个地址之间用逗号隔开"/>
      </el-form-item>

      <el-divider content-position="center">聊天设置</el-divider>
      <el-row>
        <el-col :span="12">
          <div class="grid-content">
            <el-form-item label="GPT模型">
              <el-input v-model="form['model']" placeholder="gpt-3/gpt-3.5-turbo/gpt-4"/>
            </el-form-item>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="grid-content">
            <el-form-item label="模型温度">
              <el-input v-model="form['temperature']" placeholder="0-1之间的小数"/>
            </el-form-item>
          </div>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <div class="grid-content">
            <el-form-item label="Max Tokens">
              <el-input v-model="form['max_tokens']" placeholder="回复的最大字数，最大4096"/>
            </el-form-item>
          </div>
        </el-col>

        <el-col :span="12">
          <div class="grid-content">
            <el-form-item label="上下文超时">
              <el-input v-model="form['chat_context_expire_time']" placeholder="单位：秒"/>
            </el-form-item>
          </div>
        </el-col>
      </el-row>

      <el-form-item label="对话上下文">
        <el-switch v-model="form['enable_context']"/>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="saveConfig">保存</el-button>
      </el-form-item>
    </el-form>

    <el-divider content-position="center">API KEY 管理</el-divider>
    <el-row class="api-key-box">
      <el-input
          v-model="apiKey"
          placeholder="输入 API KEY"
          class="input-with-select"
      >
        <template #prepend>
          <el-button type="primary">
            <el-icon>
              <Plus/>
            </el-icon>
          </el-button>
        </template>
        <template #append>
          <el-button class="new-proxy" @click="addApiKey">新增</el-button>
        </template>
      </el-input>
    </el-row>

    <el-row>
      <el-table :data="apiKeys" style="width: 100%">
        <el-table-column prop="value" label="API-KEY"/>
        <el-table-column prop="last_used" label="最后使用" width="180">
          <template #default="scope">
            <span v-if="scope.row['last_used'] > 0">{{ dateFormat(scope.row['last_used']) }}</span>
            <el-tag v-else>未使用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-popconfirm
                width="220"
                confirm-button-text="确定"
                cancel-button-text="取消"
                title="确定删除该记录吗?"
                :hide-after="0"
                @confirm="removeApiKey(scope.row.value)"
            >
              <template #reference>
                <el-button
                    size="small"
                    type="danger">删除
                </el-button>
              </template>
            </el-popconfirm>

          </template>
        </el-table-column>
      </el-table>
    </el-row>

  </div>
</template>

<script>
import {defineComponent, nextTick} from "vue";
import {Plus} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat, removeArrayItem} from "@/utils/libs";

export default defineComponent({
  name: 'SysConfig',
  components: {Plus},
  data() {
    return {
      apiKey: '',
      form: {},
      apiKeys: [],
      loading: true
    }
  },
  mounted() {
    // 获取系统配置
    httpGet('/api/admin/config/get').then((res) => {
      this.form = res.data;
    }).catch(() => {
      ElMessage.error('获取系统配置失败')
    })

    // 获取 API KEYS
    httpPost('api/admin/apikey/list').then((res) => {
      this.apiKeys = res.data
    }).catch(() => {
      ElMessage.error('获取 API KEY 失败')
    })

    nextTick(() => {
      this.loading = false
    })
  },
  computed: {
    dateFormat() {
      return dateFormat
    },
  },
  methods: {
    saveConfig: function () {
      this.form['temperature'] = parseFloat(this.form.temperature)
      this.form['chat_context_expire_time'] = parseInt(this.form.chat_context_expire_time)
      this.form['max_tokens'] = parseInt(this.form.max_tokens)
      httpPost("/api/admin/config/set", this.form).then(() => {
        ElMessage.success("保存成功");
      }).catch((e) => {
        console.log(e.message);
        ElMessage.error("保存失败");
      })
    },

    addApiKey: function () {
      if (this.apiKey.trim() === '') {
        ElMessage.error('请输入 API KEY')
        return
      }

      httpPost('api/admin/apikey/add', {api_key: this.apiKey.trim()}).then(() => {
        ElMessage.success('添加成功')
        this.apiKeys.unshift({value: this.apiKey, last_used: 0})
        this.apiKey = ''
      }).catch((e) => {
        ElMessage.error('添加失败，' + e.message)
      })
    },

    removeApiKey: function (key) {
      httpPost('api/admin/apikey/remove', {api_key: key}).then(() => {
        ElMessage.success('删除成功')
        this.apiKeys = removeArrayItem(this.apiKeys, key, function (v1, v2) {
          return v1.value === v2
        })
      }).catch((e) => {
        ElMessage.error('删除失败，' + e.message)
      })
    }
  }
})
</script>

<style lang="stylus" scoped>
.system-config {

  .api-key-box {
    padding-bottom: 10px;
  }

}
</style>