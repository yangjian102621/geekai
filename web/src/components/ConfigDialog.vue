<template>
  <el-dialog
      v-model="showDialog"
      :close-on-click-modal="false"
      :before-close="close"
      :top="50+'px'"
      title="用户设置"
  >
    <div class="user-info" id="user-info">
      <el-form v-if="form.id" :model="form" label-width="120px">
        <el-form-item label="昵称">
          <el-input v-model="form.nickname"/>
        </el-form-item>
        <el-form-item label="头像">
          <el-upload
              class="avatar-uploader"
              :auto-upload="true"
              :show-file-list="false"
              :http-request="afterRead"
          >
            <el-avatar v-if="form.avatar" :src="form.avatar" shape="square" :size="100"/>
            <el-icon v-else class="avatar-uploader-icon">
              <Plus/>
            </el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="form.username" readonly disabled/>
        </el-form-item>

        <el-form-item label="绑定手机号">
          <el-input v-model="form.mobile" readonly disabled/>
        </el-form-item>

        <el-form-item label="聊天上下文">
          <el-switch v-model="form.chat_config.enable_context"/>
        </el-form-item>
        <el-form-item label="聊天记录">
          <el-switch v-model="form.chat_config.enable_history"/>
        </el-form-item>
        <el-form-item label="Model">
          <el-select v-model="form.chat_config.model" placeholder="默认会话模型">
            <el-option
                v-for="item in models"
                :key="item"
                :label="item.toUpperCase()"
                :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="MaxTokens">
          <el-input v-model.number="form.chat_config.max_tokens"/>
        </el-form-item>
        <el-form-item label="Temperature">
          <el-input v-model.number="form.chat_config.temperature"/>
        </el-form-item>
        <el-form-item label="剩余调用次数">
          <el-tag>{{ form['calls'] }}</el-tag>
        </el-form-item>
        <el-form-item label="消耗 Tokens">
          <el-tag type="info">{{ form['tokens'] }}</el-tag>
        </el-form-item>
        <el-form-item label="API KEY">
          <el-input v-model="form['chat_config']['api_key']"/>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">关闭</el-button>
        <el-button type="primary" @click="save">
          保存
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>


import {computed, onMounted, ref} from "vue"
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {Plus} from "@element-plus/icons-vue";
import Compressor from "compressorjs";

// eslint-disable-next-line no-undef
const props = defineProps({
  show: Boolean,
  user: Object,
  models: Array,
});

const showDialog = computed(() => {
  return props.show
})
const form = ref({
  username: '',
  nickname: '',
  avatar: '',
  mobile: '',
  calls: 0,
  tokens: 0,
  chat_configs: {}
})

onMounted(() => {
  // 获取最新用户信息
  httpGet('/api/user/profile').then(res => {
    form.value = res.data
  }).catch(() => {
    ElMessage.error({message: '获取用户信息失败', appendTo: '#user-info'})
  });
})

const afterRead = (file) => {
  console.log(file)
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then((res) => {
        form.value.avatar = res.data
        ElMessage.success({message: '上传成功', appendTo: '#user-info', duration: 1000})
      }).catch((e) => {
        console.log(e.message)
        ElMessage.error({message: '上传失败', appendTo: '#user-info'})
      })
    },
    error(err) {
      console.log(err.message);
    },
  });
};

// eslint-disable-next-line no-undef
const emits = defineEmits(['hide', 'update-user']);
const save = function () {
  httpPost('/api/user/profile/update', form.value).then(() => {
    ElMessage.success({
      message: '更新成功',
      appendTo: '#user-info',
      onClose: () => emits('hide', false)
    })
    // 更新用户数据
    emits('update-user', {nickname: form.value['nickname'], avatar: form.value['avatar']});
  }).catch(() => {
    ElMessage.error({
      message: '更新失败',
      appendTo: document.getElementById('user-info')
    })
  })
}
const close = function () {
  emits('hide', false);
}
</script>

<style lang="stylus">
.el-dialog {
  --el-dialog-width 90%;
  max-width 800px;

  .el-dialog__body {
    padding-top 10px;
    max-height 600px;
    overflow-y auto;

    .user-info {
      position relative;

      .el-message {
        position: absolute;
      }
    }
  }
}
</style>