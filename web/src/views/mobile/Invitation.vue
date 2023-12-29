<template>
  <div class="page-invitation" :style="{height: listBoxHeight + 'px'}">
    <div class="inner">
      <div class="title" style="padding-top: 50px">会员推广计划</div>
      <div class="share-box">
        <div class="info">
          感谢您把此应用分享给您身边的朋友，分享成功注册后您将获得 <strong>{{ inviteChatCalls }}</strong>
          次对话额度以及
          <strong>{{ inviteImgCalls }}</strong> 次AI绘画额度作为奖励。
          你可以保存下面的二维码或者直接复制分享您的专属推广链接发送给微信好友。
        </div>

        <div class="invite-qrcode">
          <el-image :src="qrImg" fit="contain"/>
        </div>
        <div class="invite-url">
          <span>{{ inviteURL }}</span>
          <el-button type="primary" plain class="copy-link" :data-clipboard-text="inviteURL">复制链接</el-button>
        </div>
      </div>

      <div class="invite-stats">
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="item-box yellow">
              <el-row :gutter="10">
                <!--                <el-col :span="10">-->
                <!--                  <div class="item-icon">-->
                <!--                    <i class="iconfont icon-role"></i>-->
                <!--                  </div>-->
                <!--                </el-col>-->
                <el-col :span="24">
                  <div class="item-info">
                    <div class="num">{{ hits }}</div>
                    <div class="text">点击量</div>
                  </div>
                </el-col>
              </el-row>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="item-box blue">
              <el-row :gutter="10">
                <!--                <el-col :span="10">-->
                <!--                  <div class="item-icon">-->
                <!--                    <i class="iconfont icon-order"></i>-->
                <!--                  </div>-->
                <!--                </el-col>-->
                <el-col :span="24">
                  <div class="item-info">
                    <div class="num">{{ regNum }}</div>
                    <div class="text">注册量</div>
                  </div>
                </el-col>
              </el-row>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="item-box green">
              <el-row :gutter="10">
                <!--                <el-col :span="10">-->
                <!--                  <div class="item-icon">-->
                <!--                    <i class="iconfont icon-chart"></i>-->
                <!--                  </div>-->
                <!--                </el-col>-->
                <el-col :span="24">
                  <div class="item-info">
                    <div class="num">{{ rate }}%</div>
                    <div class="text">转化率</div>
                  </div>
                </el-col>
              </el-row>
            </div>
          </el-col>
        </el-row>
        <div class="title" style="padding: 20px">您推荐的用户</div>
        <div class="invite-logs" style="padding-bottom: 50px">
          <invite-list v-if="isLogin"/>
        </div>
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
import InviteList from "@/components/InviteList.vue";
import {checkSession} from "@/action/session";
import {useRouter} from "vue-router";

const listBoxHeight = window.innerHeight
const inviteURL = ref("")
const qrImg = ref("")
const inviteChatCalls = ref(0)
const inviteImgCalls = ref(0)
const hits = ref(0)
const regNum = ref(0)
const rate = ref(0)
const router = useRouter()
const isLogin = ref(false)

onMounted(() => {
  checkSession().then(() => {
    isLogin.value = true
    httpGet("/api/invite/code").then(res => {
      const text = `${location.protocol}//${location.host}/register?invite_code=${res.data.code}`
      hits.value = res.data["hits"]
      regNum.value = res.data["reg_num"]
      if (hits.value > 0) {
        rate.value = ((regNum.value / hits.value) * 100).toFixed(2)
      }
      QRCode.toDataURL(text, {width: 300, height: 300, margin: 2}, (error, url) => {
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

    httpGet("/api/admin/config/get?key=system").then(res => {
      inviteChatCalls.value = res.data["invite_chat_calls"]
      inviteImgCalls.value = res.data["invite_img_calls"]
    }).catch(e => {
      ElMessage.error("获取系统配置失败：" + e.message)
    })
  }).catch(() => {
    router.push('/login')
  });

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
@import "@/assets/css/mobile/invitation.styl"
</style>
