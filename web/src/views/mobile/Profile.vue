<template>
  <div class="app-background" theme="dark">
    <div class="member">
      <div class="title" style="padding-top: 50px">
        会员中心
      </div>
      <div class="inner" :style="{height: listBoxHeight + 'px'}">
        <div class="product-box">
          <ItemList :items="list" v-if="list.length > 0" :gap="30">
            <template #default="scope">
              <div class="product-item" :style="{width: scope.width+'px'}">
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

                  <div class="pay-way">
                    <el-button type="primary" @click="alipay(scope.item)" size="small" v-if="payWays['alipay']">
                      <i class="iconfont icon-alipay"></i> 支付宝
                    </el-button>
                    <el-button type="success" @click="huPiPay(scope.item)" size="small" v-if="payWays['hupi']">
                      <span v-if="payWays['hupi']['name'] === 'wechat'"><i
                          class="iconfont icon-wechat-pay"></i> 微信</span>
                      <span v-else><i class="iconfont icon-alipay"></i> 支付宝</span>
                    </el-button>
                  </div>
                </div>
              </div>
            </template>
          </ItemList>
        </div>
        <el-row>
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
              <el-col :span="24" style="padding-bottom: 50px;">
                <el-button type="danger" round @click="logout">退出登录</el-button>
              </el-col>
            </el-row>
          </div>
        </el-row>
      </div>
      <login-dialog :show="showLoginDialog" @hide="showLoginDialog = false"/>
      <password-dialog v-if="isLogin" :show="showPasswordDialog" width="100%" @hide="showPasswordDialog = false"
                       @logout="logout"/>
      <bind-mobile v-if="isLogin" :show="showBindMobileDialog" width="100%" :mobile="user.mobile"
                   @hide="showBindMobileDialog = false"/>
      <reward-verify v-if="isLogin" :show="showRewardVerifyDialog" @hide="showRewardVerifyDialog = false"/>
      <el-dialog
          v-model="showRewardDialog"
          :show-close="true"
          width="100%"
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
            <count-down :second="orderTimeout" @timeout="refreshPayCode" ref="countDownRef"/>
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
            <span class="text">请打开手机{{ payName }}扫码支付</span>
          </div>
        </div>
      </el-dialog>
    </div>
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
const countDownRef = ref(null)
const orderTimeout = ref(1800)
const loading = ref(true)
const orderPayInfoText = ref("")

const payWays = ref({})
const payName = ref("支付宝")
const curPay = ref("alipay") // 当前支付方式


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
    orderPayInfoText.value = res.data['order_pay_info_text']
    if (res.data['order_pay_timeout'] > 0) {
      orderTimeout.value = res.data['order_pay_timeout']
    }
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })

  httpGet("/api/payment/payWays").then(res => {
    payWays.value = res.data
  }).catch(e => {
    ElMessage.error("获取支付方式失败：" + e.message)
  })
})

// refresh payment qrcode
const refreshPayCode = () => {
  if (curPay.value === 'alipay') {
    alipay()
  } else if (curPay.value === 'hupi') {
    huPiPay()
  }
}

const genPayQrcode = () => {
  loading.value = true
  text.value = ""
  httpPost("/api/payment/qrcode", {
    pay_way: curPay.value,
    product_id: curPayProduct.value.id,
    user_id: user.value.id
  }).then(res => {
    showPayDialog.value = true
    qrcode.value = res.data['image']
    activeOrderNo.value = res.data['order_no']
    queryOrder(activeOrderNo.value)
    loading.value = false
    // 重置计数器
    if (countDownRef.value) {
      countDownRef.value.resetTimer()
    }
  }).catch(e => {
    ElMessage.error("生成支付订单失败：" + e.message)
  })
}

const alipay = (row) => {
  payName.value = "支付宝"
  curPay.value = "alipay"
  if (!user.value.id) {
    showLoginDialog.value = true
    return
  }
  if (row) {
    curPayProduct.value = row
  }
  genPayQrcode()
}

// 虎皮椒支付
const huPiPay = (row) => {
  payName.value = payWays.value["hupi"]["name"] === "wechat" ? '微信' : '支付宝'
  curPay.value = "hupi"
  if (!user.value.id) {
    showLoginDialog.value = true
    return
  }
  if (row) {
    curPayProduct.value = row
  }
  genPayQrcode()

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
@import "@/assets/css/mobile/profile.styl"
</style>
