<template>
  <div class="system-config" v-loading="loading">
    <div class="container">
      <el-divider content-position="center">基本设置</el-divider>
      <el-form :model="system" label-width="120px" label-position="left" ref="systemFormRef" :rules="rules">
        <el-form-item label="网站标题" prop="title">
          <el-input v-model="system['title']"/>
        </el-form-item>
        <el-form-item label="控制台标题" prop="admin_title">
          <el-input v-model="system['admin_title']"/>
        </el-form-item>
        <el-form-item label="注册赠送次数" prop="init_calls">
          <el-input v-model.number="system['init_calls']" placeholder="新用户注册赠送对话次数"/>
        </el-form-item>
        <el-alert type="info" show-icon :closable="false">
          <p>在这里维护前端聊天页面可用的 GPT 模型列表</p>
        </el-alert>
        <el-form-item label="GPT 模型" prop="models">
          <div class="models">
            <el-tag
                v-for="item in system.models"
                :key="item"
                @close="removeModel(item)"
                round
                closable
            >
              {{ item }}
            </el-tag>
            <el-button type="success" :icon="Plus" @click="addModel" size="small" circle/>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save('system')">保存</el-button>
        </el-form-item>
      </el-form>

      <el-divider content-position="center">聊天设置</el-divider>
      <el-alert type="info" show-icon :closable="false">
        <p>以下配置为新用户注册默认初始化的聊天参数，用户登录后还可以自己修改参数。</p>
      </el-alert>
      <el-form :model="chat" label-position="left" label-width="120px">
        <el-form-item label="OpenAI API 地址">
          <el-input v-model="chat['api_url']" placeholder="gpt-3/gpt-3.5-turbo/gpt-4"/>
        </el-form-item>
        <el-form-item label="默认模型">
          <el-input v-model="chat['model']" placeholder="用户默认使用的 GPT 模型"/>
        </el-form-item>
        <el-form-item label="模型温度">
          <el-input v-model="chat['temperature']" placeholder="0-1之间的小数"/>
        </el-form-item>
        <el-form-item label="Max Tokens">
          <el-input v-model="chat['max_tokens']" placeholder="回复的最大字数，最大4096"/>
        </el-form-item>
        <el-form-item label="开启聊天上下文">
          <el-switch v-model="chat['enable_context']"/>
        </el-form-item>
        <el-form-item label="保存聊天记录">
          <el-switch v-model="chat['enable_history']"/>
        </el-form-item>

        <el-form-item style="text-align: right">
          <el-button type="primary" @click="save('chat')">保存</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {nextTick, onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElMessageBox} from "element-plus";
import {Plus} from "@element-plus/icons-vue";
import {removeArrayItem} from "@/utils/libs";

const system = ref({models: []})
const chat = ref({})
const loading = ref(true)
const systemFormRef = ref(null)
const tempModel = ref('')
const models = ref([])

onMounted(() => {
  // 加载系统配置
  httpGet('/api/admin/config/get?key=system').then(res => {
    system.value = res.data
    system.value['models'].forEach(model => {
      models.value.push({
        name: model,
        edit: false
      })
    })

  }).catch(e => {
    ElMessage.error("加载系统配置失败: " + e.message)
  })

  // 加载聊天配置
  httpGet('/api/admin/config/get?key=chat').then(res => {
    chat.value = res.data
  }).catch(e => {
    ElMessage.error("加载聊天配置失败: " + e.message)
  })

  nextTick(() => {
    loading.value = false
  })
})

const rules = reactive({
  title: [{required: true, message: '请输入网站标题', trigger: 'blur',}],
  admin_title: [{required: true, message: '请输入控制台标题', trigger: 'blur',}],
  init_calls: [{required: true, message: '必须填入大于0的数组', trigger: 'blur',}],
  models: [{required: true, message: '至少保留一个 GPT 模型', trigger: 'blur',}],
})
const save = function (key) {
  systemFormRef.value.validate((valid) => {
    if (valid) {
      const data = key === 'system' ? system.value : chat.value
      httpPost('/api/admin/config/update', {key: key, config: data}).then(() => {
        ElMessage.success("操作成功！")
      }).catch(e => {
        ElMessage.error("操作失败：" + e.message)
      })
    }
  })
}

const removeModel = function (model) {
  system.value.models = removeArrayItem(system.value.models, model, (v1, v2) => {
    return v1 === v2
  })
}

// 增加 GPT 模型
const addModel = function () {
  ElMessageBox.prompt('请输入 GPT 模型名称', '新增模型', {
    confirmButtonText: '保存',
    cancelButtonText: '取消',
    inputPattern:
        /[\w+]/,
    inputErrorMessage: '请输入模型名称',

  }).then(({value}) => {
    system.value.models.push(value)
  })

}
</script>

<style lang="stylus" scoped>
.system-config {
  display flex
  justify-content center

  .container {
    width 100%
    max-width 800px;

    .el-form {
      .el-form-item__content {
        .models {
          .el-tag {
            margin-right 10px;

            .el-input {
              max-width 100px;
            }
          }

          .el-button--small {
            font-size 16px;
          }
        }
      }
    }

    .el-alert {
      margin-bottom 15px;
    }
  }

}
</style>