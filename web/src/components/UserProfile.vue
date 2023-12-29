<template>
  <div class="user-info" id="user-info">
    <el-form v-if="user.id" :model="user" label-width="150px">
      <el-row>
        <el-upload
            class="avatar-uploader"
            :auto-upload="true"
            :show-file-list="false"
            :http-request="afterRead"
        >
          <el-avatar v-if="user.avatar" :src="user.avatar" shape="circle" :size="100"/>
          <el-icon v-else class="avatar-uploader-icon">
            <Plus/>
          </el-icon>
        </el-upload>
      </el-row>
      <el-form-item label="昵称">
        {{ user['nickname'] }}
      </el-form-item>
      <el-form-item label="手机号">
        <span>{{ user.mobile }}</span>
        <el-tooltip
            class="box-item"
            effect="light"
            content="您已经是 VIP 会员"
            placement="right"
        >
          <el-image v-if="user.vip" :src="vipImg" style="height: 25px;margin-left: 10px"/>
        </el-tooltip>
      </el-form-item>
      <el-form-item label="剩余对话次数">
        <el-tag>{{ user['calls'] }}</el-tag>
      </el-form-item>
      <el-form-item label="剩余绘图次数">
        <el-tag>{{ user['img_calls'] }}</el-tag>
      </el-form-item>
      <el-form-item label="本月消耗电量">
        <el-tag type="info">{{ user['tokens'] }}</el-tag>
      </el-form-item>
      <el-form-item label="累计消耗电量">
        <el-tag type="info">{{ user['total_tokens'] }}</el-tag>
      </el-form-item>
      <el-form-item label="会员到期时间" v-if="user['expired_time']  > 0">
        <el-tag type="danger">{{ dateFormat(user['expired_time']) }}</el-tag>
      </el-form-item>

      <el-form-item label="OpenAI API KEY">
        <el-input v-model="user.chat_config['api_keys']['OpenAI']"/>
      </el-form-item>
      <el-form-item label="Azure API KEY">
        <el-input v-model="user['chat_config']['api_keys']['Azure']"/>
      </el-form-item>
      <el-form-item label="ChatGLM API KEY">
        <el-input v-model="user['chat_config']['api_keys']['ChatGLM']"/>
      </el-form-item>

      <el-row class="opt-line">
        <el-button color="#47fff1" :dark="false" round @click="save">保存</el-button>
      </el-row>
    </el-form>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {Plus} from "@element-plus/icons-vue";
import Compressor from "compressorjs";
import {dateFormat} from "@/utils/libs";
import {checkSession} from "@/action/session";

const user = ref({
  vip: false,
  username: '',
  nickname: '',
  avatar: '',
  mobile: '',
  calls: 0,
  tokens: 0,
  chat_config: {api_keys: {OpenAI: "", Azure: "", ChatGLM: ""}}
})
const vipImg = ref("/images/vip.png")

onMounted(() => {
  checkSession().then(() => {
    // 获取最新用户信息
    httpGet('/api/user/profile').then(res => {
      user.value = res.data
      user.value.chat_config.api_keys = res.data.chat_config.api_keys ?? {OpenAI: "", Azure: "", ChatGLM: ""}
    }).catch(e => {
      ElMessage.error("获取用户信息失败：" + e.message)
    });
  }).catch(e => {
    console.log(e)
  })
})

const afterRead = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then((res) => {
        user.value.avatar = res.data
        ElMessage.success({message: "上传成功", duration: 500})
      }).catch((e) => {
        ElMessage.error('图片上传失败:' + e.message)
      })
    },
    error(err) {
      console.log(err.message);
    },
  });
};

const save = () => {
  httpPost('/api/user/profile/update', user.value).then(() => {
    ElMessage.success({message: '更新成功', duration: 500})
  }).catch((e) => {
    ElMessage.error('更新失败：' + e.message)
  })
}
</script>

<style lang="stylus" scoped>
.user-info {
  padding 20px

  .el-row {
    justify-content center
    margin-bottom 10px
  }

  .opt-line {
    padding-top 20px

    .el-button {
      width 100%
    }
  }
}
</style>