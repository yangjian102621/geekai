<template>
  <div class="payment-content container">
    <van-nav-bar :title="title" :left-arrow="true" @click-left="onClickLeft"/>
    <div class="content">
      <div class="pay-price">
        <div class="pay-expire">支付倒计时：
          <van-count-down :time="data.expire" :auto-start="data.autoStart" format="mm:ss" @finish="finishPay"
                          ref="countDownRef"/>
        </div>
        <div class="pay-title">支付金额</div>
        <div class="pay-amount">¥<span>{{ data.amount }}</span></div>
      </div>
      <div class="pay-btn">
        <van-button round block type="primary" @click="pay">支付</van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import wx from "weixin-js-sdk";
import {ref} from "vue";
import {ElMessage} from "element-plus";
import {checkSession} from "@/action/session";
import {httpGet, httpPost} from "@/utils/http";
import {authResult, authUrl, getCode, getSignature} from "@/utils/wechatAuth";
import {isWeChat} from "@/utils/libs";
import {useRouter} from "vue-router";
import {setUserToken} from "@/store/session";
import {prevRoute} from "@/router";

const title = ref('支付')
const moment = require('moment');
const countDownRef = ref(null)
const router = useRouter()
const data = ref({
  amount: '-',
  expire: 0,
  autoStart: false,
  orderTimeout: 1800
})

const orderNo = router.currentRoute.value.query["order_no"]
const payWay = router.currentRoute.value.query["pay_way"]
if (isWeChat()) {
  checkSession().then(() => {
    httpGet("/api/config/get?key=system").then(res => {
      data.value.orderTimeout = res.data['order_pay_timeout']
      if (authResult()) {
        //判断是直接访问还是回调访问
        const queryParam = getCode()
        queryParam.login_type = "1"
        httpPost('/api/user/wxLogin', queryParam).then(() => {
          data.value.autoStart = false
          httpPost("/api/payment/queryOrder", {order_no: orderNo}).then(res => {
            const {amount, createTime} = res.data
            data.value.amount = amount
            let expire = createTime + data.value.orderTimeout - moment().unix()
            if (expire <= 0) {
              finishPay()
            } else {
              data.value.expire = Math.min(expire, data.value.orderTimeout) * 1000
              data.value.autoStart = true
              if (countDownRef.value) {
                countDownRef.value.start()
              }
            }
          }).catch(e => {
            ElMessage.error("查询支付状态失败：" + e.message)
          })
        }).catch((e) => {
          ElMessage.error('登录失败，' + e.message)
        })
      } else {
        window.location.href = authUrl(res.data['wxAppId'])
      }
    }).catch(e => {
      ElMessage.error("获取系统配置失败：" + e.message)
    })
  }).catch(() => {
    router.push('/login')
  })
} else {
  ElMessage.warning("请使用微信支付")
}

const pay = () => {
  if (!isWeChat()) {
    ElMessage.warning("请使用微信支付")
    return
  }
  httpGet(`/api/payment/doPay?order_no=${orderNo}&pay_way=${payWay}`).then(res => {
    const {nonceStr, paySign, signType, timeStamp} = res.data
    getSignature(res.data, () => {
      wx.chooseWXPay({
        timestamp: timeStamp, // 支付签名时间戳，注意微信jssdk中的所有使用timestamp字段均为小写。但最新版的支付后台生成签名使用的timeStamp字段名需大写其中的S字符
        nonceStr: nonceStr, // 支付签名随机串，不长于 32 位
        package: res.data['package'], // 统一支付接口返回的prepay_id参数值，提交格式如：prepay_id=\*\*\*）
        signType: signType, // 微信支付V3的传入RSA,微信支付V2的传入格式与V2统一下单的签名格式保持一致
        paySign: paySign, // 支付签名
        success: function () {
          // 支付成功后的回调函数
          ElMessage.success('支付成功')
          let timer = setTimeout(() => {
            clearTimeout(timer)
            router.push('/mobile/profile')
          }, 1000)
          if (countDownRef.value) {
            countDownRef.value.stop()
          }
        },
        fail: function () {
          ElMessage.error('支付失败')
        }
      })
    })
  }).catch(e => {
    ElMessage.error("查询支付状态失败：" + e.message)
  })
}
const onClickLeft = () => router.push('/mobile/profile');
const finishPay = () => {
  ElMessage.error('支付超时')
  let timer = setTimeout(() => {
    clearTimeout(timer)
    router.push('/mobile/profile')
  }, 1000)
}


</script>

<style lang="stylus">
.payment-content {
  .content {
    padding-top 60px

    .van-cell__value {
      .van-image {
        width 100%
      }
    }

    .pay-price {
      text-align center
      margin-bottom 30px
      padding 10px 15px

      .pay-expire {
        color #2778FF
        display flex
        flex-direction row
        align-items center
        justify-content center
        font-size 15px

        .van-count-down {
          color #2778FF
          font-size 15px
        }
      }

      .pay-title {
        color #333
        font-size 15px
        margin-top 30px
      }

      .pay-amount {
        color #666
        font-size 15px

        span {
          font-size 36px
          font-weight bold
        }
      }

    }

    .pay-btn {
      padding 10px 15px
    }
  }
}

</style>
