<template>
  <div class="member custom-scroll">
    <div class="title">
      会员充值中心
    </div>
    <div class="inner" :style="{height: listBoxHeight + 'px'}">
      <el-row :gutter="20">
        <el-col :span="7">
          <div class="user-profile">
            <user-profile/>

            <el-row class="user-opt" :gutter="20">
              <el-col :span="12">
                <el-button type="primary" @click="showPasswordDialog = true">修改密码</el-button>
              </el-col>
              <el-col :span="12">
                <el-button type="primary" @click="showBindMobileDialog = true">绑定手机号</el-button>
              </el-col>
              <el-col :span="12">
                <el-button type="primary" v-if="enableReward" @click="showRewardDialog = true">加入众筹</el-button>
              </el-col>
              <el-col :span="12">
                <el-button type="primary" v-if="enableReward" @click="showRewardVerifyDialog = true">众筹核销
                </el-button>
              </el-col>

              <el-col :span="24" style="padding-top: 30px">
                <el-button type="danger" round @click="logout">退出登录</el-button>
              </el-col>
            </el-row>
          </div>
        </el-col>

        <el-col :span="17">
          <div class="product-box">
            <div class="info">
              <el-alert type="info" show-icon :closable="false" effect="dark">
                <strong>说明:</strong> 成为本站会员后每月有500次对话额度，50次 AI 绘画额度，限制下月1号解除，若在期间超过次数后可单独购买点卡。
                当月充值的点卡有效期可以延期到下个月底。
              </el-alert>
            </div>

            <ItemList :items="list" v-if="list.length > 0" :gap="30" :width="240">
              <template #default="scope">
                <div class="product-item" :style="{width: scope.width+'px'}" @click="orderPay(scope.item)">
                  <div class="image-container">
                    <el-image :src="vipImg" fit="cover"/>
                  </div>
                  <div class="product-title">
                    <span class="name">{{ scope.item.name }}</span>
                  </div>
                  <div class="product-info">
                    <div class="info-line">
                      <span class="label">商品原价：</span>
                      <span class="price">￥{{ scope.item.price }}</span>
                    </div>
                    <div class="info-line">
                      <span class="label">促销立减：</span>
                      <span class="price">￥{{ scope.item.discount }}</span>
                    </div>
                    <div class="info-line">
                      <span class="label">有效期：</span>
                      <span class="expire" v-if="scope.item.days > 0">{{ scope.item.days }}天</span>
                      <span class="expire" v-else>当月有效</span>
                    </div>
                  </div>
                </div>
              </template>
            </ItemList>

            <h2 class="headline">消费账单</h2>

            <div class="user-order">
              <user-order v-if="isLogin"/>
            </div>
          </div>
        </el-col>
      </el-row>

    </div>

    <login-dialog :show="showLoginDialog" @hide="showLoginDialog = false"/>

    <password-dialog v-if="isLogin" :show="showPasswordDialog" @hide="showPasswordDialog = false"
                     @logout="logout"/>

    <bind-mobile v-if="isLogin" :show="showBindMobileDialog" :mobile="user.mobile"
                 @hide="showBindMobileDialog = false"/>

    <reward-verify v-if="isLogin" :show="showRewardVerifyDialog" @hide="showRewardVerifyDialog = false"/>

    <el-dialog
        v-model="showRewardDialog"
        :show-close="true"
        width="400px"
        title="参与众筹"
    >
      <el-alert type="info" :closable="false">
        <div style="font-size: 14px">您好，众筹 9.9元，就可以兑换 100 次对话，以此来覆盖我们的 OpenAI
          账单和服务器的费用。<strong
              style="color: #f56c6c">由于本人没有开通微信支付，付款后请凭借转账单号,点击【众筹核销】按钮手动核销。</strong>
        </div>
      </el-alert>
      <div style="text-align: center;padding-top: 10px;">
        <el-image v-if="enableReward" :src="rewardImg"/>
      </div>
    </el-dialog>

    <el-dialog
        v-model="showPayDialog"
        :close-on-click-modal="false"
        :show-close="true"
        :width="400"
        title="充值订单支付">
      <div class="pay-container">
        <div class="count-down">
          <count-down :second="orderTimeout" @timeout="orderPay" ref="countDown"/>
        </div>

        <div class="pay-qrcode" v-loading="loading">
          <el-image :src="qrcode"/>
        </div>

        <div class="tip success" v-if="text !== ''">
          <el-icon>
            <SuccessFilled/>
          </el-icon>
          <span class="text">{{ text }}</span>
        </div>
        <div class="tip" v-else>
          <el-icon>
            <InfoFilled/>
          </el-icon>
          <span class="text">请打开手机支付宝扫码支付</span>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {ElMessage} from "element-plus";
import {httpGet, httpPost} from "@/utils/http";
import ItemList from "@/components/ItemList.vue";
import {InfoFilled, SuccessFilled} from "@element-plus/icons-vue";
import LoginDialog from "@/components/LoginDialog.vue";
import {checkSession} from "@/action/session";
import UserProfile from "@/components/UserProfile.vue";
import PasswordDialog from "@/components/PasswordDialog.vue";
import BindMobile from "@/components/BindMobile.vue";
import RewardVerify from "@/components/RewardVerify.vue";
import {useRouter} from "vue-router";
import {removeUserToken} from "@/store/session";
import UserOrder from "@/components/UserOrder.vue";
import CountDown from "@/components/CountDown.vue";

const listBoxHeight = window.innerHeight - 97
const list = ref([])
const showLoginDialog = ref(false)
const showPayDialog = ref(false)
const vipImg = ref("/images/vip.png")
const enableReward = ref(false) // 是否启用众筹功能
const rewardImg = ref('/images/reward.png')
const qrcode = ref("")
const showPasswordDialog = ref(false);
const showBindMobileDialog = ref(false);
const showRewardDialog = ref(false);
const showRewardVerifyDialog = ref(false);
const text = ref("")
const user = ref(null)
const isLogin = ref(false)
const router = useRouter()
const curPayProduct = ref(null)
const activeOrderNo = ref("")
const countDown = ref(null)
const orderTimeout = ref(1800)
const loading = ref(true)


onMounted(() => {
  checkSession().then(_user => {
    user.value = _user
    isLogin.value = true
    httpGet("/api/product/list").then((res) => {
      list.value = res.data
    }).catch(e => {
      ElMessage.error("获取产品套餐失败：" + e.message)
    })
  }).catch(() => {
    router.push("/login")
  })

  httpGet("/api/admin/config/get?key=system").then(res => {
    rewardImg.value = res.data['reward_img']
    enableReward.value = res.data['enabled_reward']
    if (res.data['order_pay_timeout'] > 0) {
      orderTimeout.value = res.data['order_pay_timeout']
    }
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })
})

const orderPay = (row) => {
  if (!user.value.id) {
    showLoginDialog.value = true
    return
  }
  if (row) {
    curPayProduct.value = row
  }
  loading.value = true
  text.value = ""
  httpPost("/api/payment/alipay/qrcode", {product_id: curPayProduct.value.id, user_id: user.value.id}).then(res => {
    showPayDialog.value = true
    qrcode.value = res.data['image']
    activeOrderNo.value = res.data['order_no']
    queryOrder(activeOrderNo.value)
    loading.value = false
    // 重置计数器
    if (countDown.value) {
      countDown.value.resetTimer()
    }
  }).catch(e => {
    ElMessage.error("生成支付订单失败：" + e.message)
  })
}

const queryOrder = (orderNo) => {
  httpPost("/api/payment/query", {order_no: orderNo}).then(res => {
    if (res.data.status === 1) {
      text.value = "扫码成功，请在手机上进行支付！"
      queryOrder(orderNo)
    } else if (res.data.status === 2) {
      text.value = "支付成功，正在刷新页面"
      setTimeout(() => location.reload(), 500)
    } else {
      // 如果当前订单没有过期，继续等待订单的下一个状态
      if (activeOrderNo.value === orderNo) {
        queryOrder(orderNo)
      }
    }
  }).catch(e => {
    ElMessage.error("查询支付状态失败：" + e.message)
  })
}

const logout = function () {
  httpGet('/api/user/logout').then(() => {
    removeUserToken();
    router.push('/login');
  }).catch(() => {
    ElMessage.error('注销失败！');
  })
}

</script>

<style lang="stylus">
@import "@/assets/css/custom-scroll.styl"
.member {
  background-color: #282c34;
  height 100vh

  .el-dialog {
    .el-dialog__body {
      padding-top 10px

      .pay-container {
        .count-down {
          display flex
          justify-content center
        }

        .pay-qrcode {
          display flex
          justify-content center

          .el-image {
            width 360px;
            height 360px;
          }
        }

        .tip {
          display flex
          justify-content center

          .el-icon {
            font-size 24px
          }

          .text {
            font-size: 16px
            margin-left 10px
          }
        }

        .tip.success {
          color #07c160
        }
      }
    }
  }

  .title {
    text-align center
    background-color #25272d
    font-size 24px
    color #ffffff
    padding 10px
    border-bottom 1px solid #3c3c3c
  }

  .inner {
    color #ffffff
    padding 15px 0 15px 15px;
    overflow-x hidden
    overflow-y visible

    .user-profile {
      padding 10px 20px 20px 20px
      background-color #393F4A
      color #ffffff
      border-radius 10px
      height 100vh

      .el-form-item__label {
        color #ffffff
        justify-content start
      }

      .user-opt {
        .el-col {
          padding 10px

          .el-button {
            width 100%
          }
        }
      }
    }


    .product-box {
      .info {
        .el-alert__description {
          font-size 14px !important
          margin 0
        }
        padding 10px 20px 20px 0
      }

      .list-box {
        .product-item {
          border 1px solid #666666
          border-radius 6px
          overflow hidden
          cursor pointer
          transition: all 0.3s ease; /* 添加过渡效果 */

          .image-container {
            display flex
            justify-content center

            .el-image {
              padding 6px

              .el-image__inner {
                border-radius 10px
              }
            }
          }

          .product-title {
            display flex
            padding 10px

            .name {
              width 100%
              text-align center
              font-size 16px
              font-weight bold
              color #47fff1
            }
          }

          .product-info {
            padding 10px 20px
            font-size 14px
            color #999999

            .info-line {
              display flex
              width 100%
              padding 5px 0

              .label {
                display flex
                width 100%
              }

              .price, .expire {
                display flex
                width 90px
                justify-content right
              }

              .price {
                color #f56c6c
              }

              .expire {
                color #409eff
              }
            }

          }

          &:hover {
            box-shadow: 0 0 10px rgba(71, 255, 241, 0.6); /* 添加阴影效果 */
            transform: translateY(-10px); /* 向上移动10像素 */
          }
        }
      }

      .headline {
        padding 0 20px 20px 0
      }

      .user-order {
        padding 0 20px 20px 0
      }
    }
  }

}
</style>
