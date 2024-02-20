<template>
  <div class="mobile-user-profile container">
    <van-nav-bar :title="title"/>

    <div class="content">
      <van-form @submit="save">
        <van-cell-group inset v-model="form">
          <van-field
              v-model="form.username"
              name="账号"
              label="账号"
              readonly
              disabled
          />
          <van-field label="头像">
            <template #input>
              <van-uploader v-model="fileList"
                            reupload max-count="1"
                            :deletable="false"
                            :after-read="afterRead"/>
            </template>
          </van-field>

          <van-field label="剩余对话次数">
            <template #input>
              <van-tag type="primary">{{ form.calls }}</van-tag>
            </template>
          </van-field>

          <van-field label="剩余绘图次数">
            <template #input>
              <van-tag type="primary">{{ form.img_calls }}</van-tag>
            </template>
          </van-field>

          <van-field label="累计消耗tokens">
            <template #input>
              <van-tag type="primary">{{ form.total_tokens }}</van-tag>
            </template>
          </van-field>
        </van-cell-group>
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit">
            提交
          </van-button>
        </div>
      </van-form>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {showFailToast, showNotify, showSuccessToast} from "vant";
import {httpGet, httpPost} from "@/utils/http";
import Compressor from 'compressorjs';

const title = ref('用户设置')
const form = ref({
  username: '',
  nickname: '',
  mobile: '',
  avatar: '',
  calls: 0,
  tokens: 0
})
const fileList = ref([
  {
    url: 'https://fastly.jsdelivr.net/npm/@vant/assets/leaf.jpeg',
    message: '上传中...',
  }
]);

onMounted(() => {
  httpGet('/api/user/profile').then(res => {
    form.value = res.data
    fileList.value[0].url = form.value.avatar
  }).catch((e) => {
    console.log(e.message)
    showFailToast('获取用户信息失败')
  });
})

const afterRead = (file) => {
  file.status = 'uploading';
  file.message = '上传中...';
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then((res) => {
        form.value.avatar = res.data
        file.status = 'success'
        showNotify({type: 'success', message: '上传成功'})
      }).catch((e) => {
        console.log(e.message)
        showNotify({type: 'danger', message: '上传失败'})
      })
    },
    error(err) {
      console.log(err.message);
    },
  });
};

const save = () => {
  httpPost('/api/user/profile/update', form.value).then(() => {
    showSuccessToast('保存成功')
  }).catch(() => {
    showFailToast('保存失败')
  })
}
</script>

<style lang="stylus">
.mobile-user-profile {
  .content {
    .van-field__label {
      width 100px
      text-align right
    }
  }
}
</style>