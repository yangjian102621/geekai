<template>
  <div class="page-invitation">
    <div class="inner">
      <h2>会员推广计划</h2>
      <div class="share-box">
        <div class="info">
          我们非常欢迎您把此应用分享给您身边的朋友，分享成功注册后您将获得 {{ inviteChatCalls }} 次对话额度以及
          {{ inviteImgCalls }} 次AI绘画额度作为奖励。
          你可以保存下面的二维码或者直接复制分享您的专属推广链接发送给微信好友。
        </div>

        <div class="invite-qrcode">
          <el-image :src="qrImg"/>
        </div>

        <div class="invite-url">
          <span>{{ inviteURL }}</span>
          <el-button type="primary" plain class="copy-link" :data-clipboard-text="inviteURL">复制链接</el-button>
        </div>
      </div>

      <h2>您推荐用户</h2>

      <div class="invite-logs">
        <el-empty :image-size="100"/>
      </div>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import QRCode from "qrcode";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import Clipboard from "clipboard";

const inviteURL = ref("")
const qrImg = ref("")
const inviteChatCalls = ref(0)
const inviteImgCalls = ref(0)
const users = ref([])


onMounted(() => {
  httpGet("/api/invite/code").then(res => {
    const text = `${location.protocol}//${location.host}/register?invite_code=${res.data.code}`
    QRCode.toDataURL(text, {width: 400, height: 400, margin: 2}, (error, url) => {
      if (error) {
        console.error(error)
      } else {
        qrImg.value = url;
      }
    });
    inviteURL.value = text
  }).catch(e => {
    ElMessage.error("获取邀请码失败：" + e.message)
  })

  // 复制链接
  const clipboard = new Clipboard('.copy-link');
  clipboard.on('success', () => {
    ElMessage.success('复制成功！');
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })
})
</script>

<style lang="stylus" scoped>
.page-invitation {
  display: flex;
  justify-content: center;
  background-color: #282c34;
  height 100vh

  .inner {
    max-width 800px
    width 100%
    color #e1e1e1

    h2 {
      color #ffffff;
    }

    .share-box {
      .info {
        line-height 1.5
        border 1px solid #444444
        border-radius 10px
        padding 10px
      }

      .invite-qrcode {
        padding 20px
        text-align center
      }

      .invite-url {
        padding 15px
        display flex
        justify-content space-between
        border 1px solid #444444
        border-radius 10px

        span {
          position relative
          font-family 'Microsoft YaHei', '微软雅黑', Arial, sans-serif
          top 5px
        }
      }
    }

  }

}
</style>
