<template>
  <div class="custom-scroll">
    <div class="page-invitation">
      <div class="inner">
        <h2>会员推广计划</h2>
        <div class="share-box">
          <div class="info">
            我们非常欢迎您把此应用分享给您身边的朋友，分享成功注册后您将获得 <strong>{{ inviteChatCalls }}</strong>
            次对话额度以及
            <strong>{{ inviteImgCalls }}</strong> 次AI绘画额度作为奖励。
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

        <div class="invite-stats">
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="item-box yellow">
                <el-row :gutter="10">
                  <el-col :span="10">
                    <div class="item-icon">
                      <i class="iconfont icon-role"></i>
                    </div>
                  </el-col>
                  <el-col :span="14">
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
                  <el-col :span="10">
                    <div class="item-icon">
                      <i class="iconfont icon-order"></i>
                    </div>
                  </el-col>
                  <el-col :span="14">
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
                  <el-col :span="10">
                    <div class="item-icon">
                      <i class="iconfont icon-chart"></i>
                    </div>
                  </el-col>
                  <el-col :span="14">
                    <div class="item-info">
                      <div class="num">{{ rate }}%</div>
                      <div class="text">转化率</div>
                    </div>
                  </el-col>
                </el-row>
              </div>
            </el-col>
          </el-row>
        </div>

        <h2>您推荐用户</h2>

        <div class="invite-logs">
          <invite-list v-if="isLogin"/>
        </div>
      </div>
    </div>

    <login-dialog :show="showLoginDialog" @hide="showLoginDialog =  false" @success="initData"/>
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
import LoginDialog from "@/components/LoginDialog.vue";

const inviteURL = ref("")
const qrImg = ref("")
const inviteChatCalls = ref(0)
const inviteImgCalls = ref(0)
const hits = ref(0)
const regNum = ref(0)
const rate = ref(0)
const isLogin = ref(false)
const showLoginDialog = ref(false)

onMounted(() => {
  initData()

  // 复制链接
  const clipboard = new Clipboard('.copy-link');
  clipboard.on('success', () => {
    ElMessage.success('复制成功！');
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })
})

const initData = () => {
  checkSession().then(() => {
    isLogin.value = true
    httpGet("/api/invite/code").then(res => {
      const text = `${location.protocol}//${location.host}/register?invite_code=${res.data.code}`
      hits.value = res.data["hits"]
      regNum.value = res.data["reg_num"]
      if (hits.value > 0) {
        rate.value = ((regNum.value / hits.value) * 100).toFixed(2)
      }
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

    httpGet("/api/config/get?key=system").then(res => {
      inviteChatCalls.value = res.data["invite_chat_calls"]
      inviteImgCalls.value = res.data["invite_img_calls"]
    }).catch(e => {
      ElMessage.error("获取系统配置失败：" + e.message)
    })
  }).catch(() => {
    showLoginDialog.value = true
  });
}
</script>

<style lang="stylus" scoped>
@import "@/assets/css/custom-scroll.styl"
.page-invitation {
  display: flex;
  justify-content: center;
  background-color: #282c34;
  height 100vh
  overflow-x hidden
  overflow-y visible

  .inner {
    display flex
    flex-flow column
    max-width 1000px
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

        strong {
          color #f56c6c
        }
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

    .invite-stats {
      padding 30px 10px

      .item-box {
        border-radius 10px
        padding 0 10px

        .el-col {
          height 140px
          display flex
          align-items center
          justify-content center

          .iconfont {
            font-size 60px
          }

          .item-info {
            font-size 18px

            .text, .num {
              padding 3px 0
              text-align center
            }

            .num {
              font-size 40px
            }
          }
        }
      }

      .yellow {
        background-color #ffeecc
        color #D68F00
      }

      .blue {
        background-color #D6E4FF
        color #1062FE
      }

      .green {
        background-color #E7F8EB
        color #2D9F46
      }
    }


    .invite-logs {
      padding-bottom 20px
    }
  }

}
</style>
