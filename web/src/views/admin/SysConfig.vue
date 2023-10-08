<template>
  <div class="system-config" v-loading="loading">
    <div class="container">
      <el-divider content-position="center">基本设置</el-divider>
      <el-form :model="system" label-width="150px" label-position="right" ref="systemFormRef" :rules="rules">
        <el-form-item label="网站标题" prop="title">
          <el-input v-model="system['title']"/>
        </el-form-item>
        <el-form-item label="控制台标题" prop="admin_title">
          <el-input v-model="system['admin_title']"/>
        </el-form-item>
        <el-form-item label="注册赠送对话次数" prop="user_init_calls">
          <el-input v-model.number="system['user_init_calls']" placeholder="新用户注册赠送对话次数"/>
        </el-form-item>
        <el-form-item label="注册赠送绘图次数" prop="init_img_calls">
          <el-input v-model.number="system['init_img_calls']" placeholder="新用户注册赠送绘图次数"/>
        </el-form-item>
        <el-form-item label="开放注册服务" prop="enabled_register">
          <el-switch v-model="system['enabled_register']"/>
        </el-form-item>
        <el-form-item label="短信验证服务" prop="enabled_msg">
          <el-switch v-model="system['enabled_msg']"/>
        </el-form-item>
        <el-form-item label="开放AI绘画" prop="enabled_draw">
          <el-switch v-model="system['enabled_draw']"/>
        </el-form-item>
        <el-form-item label="收款二维码" prop="reward_img">
          <el-input v-model="system['reward_img']" placeholder="众筹收款二维码地址">
            <template #append>
              <el-upload
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="uploadRewardImg"
              >
                <el-icon class="uploader-icon">
                  <UploadFilled/>
                </el-icon>
              </el-upload>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save('system')">保存</el-button>
        </el-form-item>
      </el-form>

      <el-divider content-position="center">模型通用配置</el-divider>

      <el-form :model="chat" label-position="right" label-width="150px" ref="chatFormRef" :rules="rules">
        <el-form-item label="开启聊天上下文">
          <el-switch v-model="chat['enable_context']"/>
        </el-form-item>
        <el-form-item label="保存聊天记录">
          <el-switch v-model="chat['enable_history']"/>
        </el-form-item>
        <el-form-item label="会话上下文深度">
          <el-input-number v-model="chat['context_deep']" :min="0" :max="10"/>
          <div class="tip" style="margin-top: 10px;">会话上下文深度：在老会话中继续会话，默认加载多少条聊天记录作为上下文。如果设置为
            0
            则不加载聊天记录，仅仅使用当前角色的上下文。该配置参数最好设置为 2 的整数倍。
          </div>
        </el-form-item>

        <el-divider content-position="center">OpenAI</el-divider>
        <el-form-item label="API 地址" prop="open_ai.api_url">
          <el-input v-model="chat['open_ai']['api_url']" placeholder="支持变量，{model} => 模型名称"/>
        </el-form-item>
        <el-form-item label="模型创意度">
          <el-slider v-model="chat['open_ai']['temperature']" :max="2" :step="0.1"/>
          <div class="tip">值越大 AI 回答越发散，值越小回答越保守，建议保持默认值</div>
        </el-form-item>
        <el-form-item label="最大响应长度">
          <el-input v-model.number="chat['open_ai']['max_tokens']" placeholder="回复的最大字数，最大4096"/>
        </el-form-item>

        <el-divider content-position="center">Azure</el-divider>
        <el-form-item label="API 地址" prop="azure.api_url">
          <el-input v-model="chat['azure']['api_url']" placeholder="支持变量，{model} => 模型名称"/>
        </el-form-item>
        <el-form-item label="模型创意度">
          <el-slider v-model="chat['azure']['temperature']" :max="2" :step="0.1"/>
          <div class="tip">值越大 AI 回答越发散，值越小回答越保守，建议保持默认值</div>
        </el-form-item>
        <el-form-item label="最大响应长度">
          <el-input v-model.number="chat['azure']['max_tokens']" placeholder="回复的最大字数，最大4096"/>
        </el-form-item>

        <el-divider content-position="center">ChatGLM</el-divider>
        <el-form-item label="API 地址" prop="chat_gml.api_url">
          <el-input v-model="chat['chat_gml']['api_url']" placeholder="支持变量，{model} => 模型名称"/>
        </el-form-item>
        <el-form-item label="模型创意度">
          <el-slider v-model="chat['chat_gml']['temperature']" :max="2" :step="0.1"/>
          <div class="tip">值越大 AI 回答越发散，值越小回答越保守，建议保持默认值</div>
        </el-form-item>
        <el-form-item label="最大响应长度">
          <el-input v-model.number="chat['chat_gml']['max_tokens']" placeholder="回复的最大字数，最大4096"/>
        </el-form-item>

        <el-form-item style="text-align: right">
          <el-button type="primary" @click="save('chat')">保存</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import Compressor from "compressorjs";
import {ElMessage} from "element-plus";
import {UploadFilled} from "@element-plus/icons-vue";

const system = ref({models: []})
const chat = ref({
  open_ai: {api_url: "", temperature: 1, max_tokens: 1024},
  azure: {api_url: "", temperature: 1, max_tokens: 1024},
  chat_gml: {api_url: "", temperature: 1, max_tokens: 1024},
  context_deep: 0,
  enable_context: true,
  enable_history: true,
})
const loading = ref(true)
const systemFormRef = ref(null)
const chatFormRef = ref(null)

onMounted(() => {
  // 加载系统配置
  httpGet('/api/admin/config/get?key=system').then(res => {
    system.value = res.data
  }).catch(e => {
    ElMessage.error("加载系统配置失败: " + e.message)
  })

  // 加载聊天配置
  httpGet('/api/admin/config/get?key=chat').then(res => {
    // chat.value = res.data
    if (res.data.open_ai) {
      chat.value.open_ai = res.data.open_ai
    }
    if (res.data.azure) {
      chat.value.azure = res.data.azure
    }
    if (res.data.chat_gml) {
      chat.value.chat_gml = res.data.chat_gml
    }
    chat.value.context_deep = res.data.context_deep
    chat.value.enable_context = res.data.enable_context
    chat.value.enable_history = res.data.enable_history
    loading.value = false
  }).catch(e => {
    ElMessage.error("加载聊天配置失败: " + e.message)
  })

})

const rules = reactive({
  title: [{required: true, message: '请输入网站标题', trigger: 'blur',}],
  admin_title: [{required: true, message: '请输入控制台标题', trigger: 'blur',}],
  user_init_calls: [{required: true, message: '请输入赠送对话次数', trigger: 'blur'}],
  user_img_calls: [{required: true, message: '请输入赠送绘图次数', trigger: 'blur'}],
  open_ai: {api_url: [{required: true, message: '请输入 API URL', trigger: 'blur'}]},
  azure: {api_url: [{required: true, message: '请输入 API URL', trigger: 'blur'}]},
  chat_gml: {api_url: [{required: true, message: '请输入 API URL', trigger: 'blur'}]},
})
const save = function (key) {
  if (key === 'system') {
    systemFormRef.value.validate((valid) => {
      if (valid) {
        httpPost('/api/admin/config/update', {key: key, config: system.value}).then(() => {
          ElMessage.success("操作成功！")
        }).catch(e => {
          ElMessage.error("操作失败：" + e.message)
        })
      }
    })
  } else if (key === 'chat') {
    chatFormRef.value.validate((valid) => {
      if (valid) {
        httpPost('/api/admin/config/update', {key: key, config: chat.value}).then(() => {
          ElMessage.success("操作成功！")
        }).catch(e => {
          ElMessage.error("操作失败：" + e.message)
        })
      }
    })
  }
}

// 图片上传
const uploadRewardImg = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then((res) => {
        system.value['reward_img'] = res.data
        ElMessage.success('上传成功')
      }).catch((e) => {
        ElMessage.error('上传失败:' + e.message)
      })
    },
    error(err) {
      console.log(err.message);
    },
  });
};

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

        .tip-text {
          padding-left 10px;
        }

        .tip {
          color #c1c1c1
          font-size 12px;
          line-height 1.5;
        }

        .uploader-icon {
          font-size 24px
          position relative
          top 3px
        }
      }
    }

    .el-alert {
      margin-bottom 15px;
    }
  }

}
</style>