<template>
  <el-dialog
      class="config-dialog"
      v-model="showDialog"
      :close-on-click-modal="true"
      :before-close="close"
      style="max-width: 600px"
      title="用户设置"
  >
    <div class="user-info" id="user-info">
      <el-form v-if="form.id" :model="form" label-width="150px">
        <el-form-item label="账户">
          <span>{{ form.mobile }}</span>
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
        <el-form-item label="剩余对话次数">
          <el-tag>{{ form['calls'] }}</el-tag>
        </el-form-item>
        <el-form-item label="剩余绘图次数">
          <el-tag>{{ form['img_calls'] }}</el-tag>
        </el-form-item>
        <el-form-item label="累计消耗 Tokens">
          <el-tag type="info">{{ form['total_tokens'] }}</el-tag>
        </el-form-item>
        <el-form-item label="OpenAI API KEY">
          <el-input v-model="form.chat_config['api_keys']['OpenAI']"/>
        </el-form-item>
        <el-form-item label="Azure API KEY">
          <el-input v-model="form['chat_config']['api_keys']['Azure']"/>
        </el-form-item>
        <el-form-item label="ChatGLM API KEY">
          <el-input v-model="form['chat_config']['api_keys']['ChatGLM']"/>
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
  chat_config: {api_keys: {OpenAI: "", Azure: "", ChatGLM: ""}}
})

onMounted(() => {
  // 获取最新用户信息
  httpGet('/api/user/profile').then(res => {
    form.value = res.data
    form.value.chat_config.api_keys = res.data.chat_config.api_keys ?? {OpenAI: "", Azure: "", ChatGLM: ""}
  }).catch(e => {
    ElMessage.error("获取用户信息失败：" + e.message)
  });
})

const afterRead = (file) => {
  // console.log(file)
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then((res) => {
        form.value.avatar = res.data
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

// eslint-disable-next-line no-undef
const emits = defineEmits(['hide', 'update-user']);
const save = function () {
  httpPost('/api/user/profile/update', form.value).then(() => {
    ElMessage.success({
      message: '更新成功',
      onClose: () => emits('hide', false)
    })
    // 更新用户数据
    emits('update-user', {nickname: form.value['nickname'], avatar: form.value['avatar']});
  }).catch((e) => {
    ElMessage.error('更新失败：' + e.message)
  })
}
const close = function () {
  emits('hide', false);
}
</script>

<style lang="stylus">
.config-dialog {
  .el-dialog {
    --el-dialog-width 90%;
    max-width 800px;

    .el-dialog__body {
      overflow-y auto;

      .user-info {
        position relative;

        .el-message {
          position: absolute;
        }
      }

      .tip {
        color #c1c1c1
        font-size 12px;
      }
    }
  }
}
</style>