<template>
  <div>
    <div
      class="custom-scroll min-h-full bg-slate-50 dark:bg-[#0f1117]"
      v-loading="loading"
      element-loading-background="rgba(255,255,255,.6)"
      :element-loading-text="loadingText"
    >
      <div class="mx-auto max-w-7xl space-y-6 px-4 py-8">
        <!-- Page header：左占位 + 居中标题副标题 + 右上「AI 生成」（与截图对齐） -->
        <div class="flex items-start justify-between gap-2 pt-1">
          <div class="w-12 shrink-0 sm:w-16" aria-hidden="true" />
          <div class="min-w-0 flex-1 text-center">
            <h1 class="m-0 text-[2rem] font-bold leading-snug text-indigo-600 dark:text-indigo-400">
              会员中心
            </h1>
            <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">
              升级算力套餐，让创作效率持续提升
            </p>
          </div>
        </div>

        <!-- 中间两列：左用户信息 | 右账号操作（截图 1:1） -->
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2 md:items-stretch">
          <!-- 左：用户资料 + 可用积分 -->
          <div
            class="flex h-full min-h-0 flex-col rounded-2xl border border-gray-100 bg-white p-6 shadow-sm dark:border-[#2d3448] dark:bg-[#1e2433] dark:shadow-none"
          >
            <div
              v-if="isLogin && user"
              class="flex w-full flex-1 items-center justify-between gap-5 min-h-0"
            >
              <div class="flex items-start gap-4">
                <el-upload
                  class="shrink-0 cursor-pointer"
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="afterRead"
                  accept=".png,.jpg,.jpeg,.bmp"
                >
                  <div class="flex flex-col items-center gap-1">
                    <el-tooltip content="点击上传头像" placement="top" v-if="user.avatar">
                      <el-avatar
                        :src="user.avatar"
                        shape="circle"
                        :size="56"
                        class="block transition-opacity hover:opacity-[0.85]"
                      />
                    </el-tooltip>
                    <div v-else class="flex flex-col items-center gap-1">
                      <div
                        class="flex h-14 w-14 cursor-pointer items-center justify-center rounded-full bg-indigo-50 text-indigo-600 ring-1 ring-inset ring-indigo-100 transition-colors hover:bg-indigo-100 dark:bg-indigo-500/15 dark:text-indigo-400 dark:ring-indigo-500/30 dark:hover:bg-indigo-500/20"
                      >
                        <el-icon class="text-2xl"><Plus /></el-icon>
                      </div>
                      <span class="text-[10px] text-gray-400 dark:text-slate-500">用户头像</span>
                    </div>
                  </div>
                </el-upload>
                <div class="flex min-w-0 flex-1 flex-col gap-1 pt-0.5">
                  <span
                    class="truncate text-base font-bold text-slate-900 hover:underline dark:text-slate-100"
                    :class="{ 'cursor-pointer': isLogin }"
                    title="点击修改用户名与昵称"
                    @click="isLogin && openNicknameDialog()"
                    >{{ displayName }}</span
                  >
                </div>
              </div>
              <div
                class="flex shrink-0 flex-col !items-start rounded-xl border border-indigo-200/80 bg-indigo-50/90 px-4 py-3.5 dark:border-indigo-500/35 dark:bg-indigo-500/10 min-w-[200px]"
              >
                <span class="text-sm font-medium text-indigo-600 dark:text-indigo-400"
                  >可用积分</span
                >
                <span
                  class="text-2xl font-bold tabular-nums text-indigo-600 dark:text-indigo-400"
                  >{{ userScore }}</span
                >
              </div>
            </div>
            <div v-else class="flex min-h-[120px] flex-1 items-center justify-center">
              <span class="text-sm text-slate-500 dark:text-slate-500">请先登录以查看账户详情</span>
            </div>
          </div>

          <!-- 右：快捷操作（与左侧卡片等高） -->
          <div
            class="flex h-full min-h-0 flex-col rounded-2xl border border-gray-100 bg-white p-6 shadow-sm dark:border-[#2d3448] dark:bg-[#1e2433] dark:shadow-none"
          >
            <div class="flex min-h-0 flex-1 w-full items-center justify-center">
              <div class="grid w-full grid-cols-2 items-center gap-3 sm:grid-cols-5 sm:gap-3">
                <button
                  v-for="act in accountActions"
                  :key="act.key"
                  type="button"
                  class="group flex h-[64px] max-w-[100px] w-full cursor-pointer flex-col items-center justify-center gap-1.5 self-center rounded-xl border border-gray-200 bg-white px-2 py-2 transition-colors hover:border-indigo-500 hover:bg-indigo-50/80 dark:border-[#2d3448] dark:bg-[#1e2433] dark:hover:border-indigo-400 dark:hover:bg-indigo-500/12"
                  @click="handleAccountAction(act.key)"
                >
                  <i
                    class="iconfont text-[1.375rem] text-slate-500 transition-colors group-hover:text-indigo-600 dark:text-slate-400 dark:group-hover:text-indigo-400"
                    :class="act.icon"
                  ></i>
                  <span
                    class="text-center text-[0.8125rem] font-medium leading-tight text-slate-800 transition-colors group-hover:text-indigo-600 dark:text-slate-300 dark:group-hover:text-indigo-400"
                    >{{ act.label }}</span
                  >
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 精选套餐 -->
        <section
          class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm dark:border-[#2d3448] dark:bg-[#1e2433] dark:shadow-none"
        >
          <div class="mb-6">
            <h2 class="m-0 text-lg font-bold text-slate-900 dark:text-slate-100">精选套餐</h2>
            <p class="mt-1.5 text-sm text-slate-500 dark:text-slate-400">
              多档算力套餐覆盖高频使用场景，灵活支付立即生效。
            </p>
          </div>

          <div
            v-if="list.length > 0"
            class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
          >
            <div
              class="relative flex flex-col gap-4 rounded-2xl border border-gray-100 bg-white p-5 shadow-sm transition-[border-color,box-shadow] hover:border-indigo-300 hover:shadow-md dark:border-[#2d3448] dark:bg-[#1e2433] dark:shadow-none dark:hover:border-indigo-500/50 dark:hover:shadow-lg dark:hover:shadow-black/20"
              v-for="item in list"
              :key="item.id || item.name"
            >
              <div
                class="absolute right-3 top-3 inline-flex items-center rounded-full border border-indigo-200/80 bg-indigo-50 px-2.5 py-0.5 text-[0.6875rem] font-semibold tracking-wide text-indigo-600 dark:border-indigo-500/30 dark:bg-indigo-500/15 dark:text-indigo-400"
              >
                热销
              </div>
              <div class="flex items-center gap-4">
                <el-image
                  class="h-14 w-14 shrink-0 overflow-hidden rounded-xl border border-gray-100 dark:border-[#2d3448]"
                  :src="vipImg"
                  fit="cover"
                />
                <div class="min-w-0">
                  <h3 class="m-0 text-[0.9375rem] font-bold text-slate-900 dark:text-slate-100">
                    {{ item.name }}
                  </h3>
                  <p class="mt-0.5 text-sm text-slate-500 dark:text-slate-400">
                    算力值：{{ item.power }}
                  </p>
                </div>
              </div>
              <div class="flex items-baseline gap-0.5 text-indigo-600 dark:text-indigo-400">
                <span class="text-base font-semibold">¥</span>
                <span class="text-[2rem] font-extrabold leading-none">{{ item.price }}</span>
              </div>
              <div v-if="item.features" class="flex flex-col gap-2">
                <div
                  class="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400"
                  v-for="feature in item.features"
                  :key="feature"
                >
                  <i
                    class="iconfont icon-check shrink-0 text-sm text-green-500 dark:text-green-400"
                  ></i>
                  <span>{{ feature }}</span>
                </div>
              </div>
              <div class="mt-auto grid grid-cols-2 gap-2">
                <el-button
                  class="h-10 w-full text-[0.8125rem] font-medium"
                  type="success"
                  @click="wxPay(item)"
                >
                  <i class="iconfont icon-wechat-pay mr-1"></i>
                  <span>微信支付</span>
                </el-button>
                <el-button
                  class="h-10 w-full text-[0.8125rem] font-medium"
                  color="#1677FF"
                  @click="alipay(item)"
                >
                  <i class="iconfont icon-alipay mr-1"></i>
                  <span>支付宝</span>
                </el-button>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无数据" :image="nodata" />
        </section>

        <!-- 消费账单 -->
        <section
          class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm dark:border-[#2d3448] dark:bg-[#1e2433] dark:shadow-none"
        >
          <div class="mb-6">
            <h2 class="m-0 text-lg font-bold text-slate-900 dark:text-slate-100">消费账单</h2>
            <p class="mt-1.5 text-sm text-slate-500 dark:text-slate-400">
              实时同步最新订单状态，轻松掌握消费明细。
            </p>
          </div>
          <div class="min-h-[80px] rounded-xl border border-gray-100 dark:border-[#2d3448]">
            <user-order v-if="isLogin" :key="userOrderKey" />
            <div v-else class="py-6 text-center text-sm text-slate-500 dark:text-slate-500">
              登录后可查看完整订单记录
            </div>
          </div>
        </section>

        <el-dialog
          v-model="showNicknameDialog"
          title="修改用户名与昵称"
          width="400px"
          :close-on-click-modal="true"
          @close="showNicknameDialog = false"
        >
          <el-form label-width="70px">
            <el-form-item label="用户名">
              <el-input
                v-model="usernameInput"
                placeholder="请输入用户名"
                maxlength="30"
                show-word-limit
                clearable
              />
            </el-form-item>
            <el-form-item label="昵称">
              <el-input
                v-model="nicknameInput"
                placeholder="请输入昵称"
                maxlength="20"
                show-word-limit
                clearable
              />
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="showNicknameDialog = false">取消</el-button>
            <el-button type="primary" @click="saveNickname">保存</el-button>
          </template>
        </el-dialog>

        <password-dialog
          v-if="isLogin"
          :show="showPasswordDialog"
          @hide="showPasswordDialog = false"
        />

        <el-dialog
          v-model="showBindMobileDialog"
          title="绑定手机"
          width="400px"
          :close-on-click-modal="true"
          @close="showBindMobileDialog = false"
        >
          <bind-mobile @hide="showBindMobileDialog = false" />
        </el-dialog>

        <el-dialog
          v-model="showBindEmailDialog"
          title="绑定邮箱"
          width="400px"
          :close-on-click-modal="true"
          @close="showBindEmailDialog = false"
        >
          <bind-email @hide="showBindEmailDialog = false" />
        </el-dialog>

        <el-dialog
          v-model="showRedeemVerifyDialog"
          title="卡密兑换"
          width="450px"
          :close-on-click-modal="true"
          @close="showRedeemVerifyDialog = false"
        >
          <redeem-verify @hide="redeemCallback" />
        </el-dialog>
      </div>
    </div>

    <el-dialog v-model="showQrCode" :show-close="true" class="!w-[334px] !h-[368px]">
      <template #header>
        <div class="flex items-center justify-center text-base text-slate-600 dark:text-slate-400">
          <span>{{ title }}</span>
        </div>
      </template>
      <div ref="qrContainerRef" class="flex items-center justify-center">
        <el-image :src="qrImg" class="h-[300px] w-[300px]" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import nodata from '@/assets/img/no-data.png'
import BindEmail from '@/components/BindEmail.vue'
import BindMobile from '@/components/BindMobile.vue'
import PasswordDialog from '@/components/PasswordDialog.vue'
import RedeemVerify from '@/components/RedeemVerify.vue'
import UserOrder from '@/components/UserOrder.vue'
import { checkSession, getSystemInfo } from '@/store/cache'
import { useSharedStore } from '@/store/sharedata'
import { closeLoading, showLoading } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { isMobile } from '@/utils/libs'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import Compressor from 'compressorjs'
import QRCode from 'qrcode'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const list = ref([])
const vipImg = ref('/images/menu/member.png')
const showPasswordDialog = ref(false)
const showBindMobileDialog = ref(false)
const showBindEmailDialog = ref(false)
const showRedeemVerifyDialog = ref(false)
const showNicknameDialog = ref(false)
const nicknameInput = ref('')
const usernameInput = ref('')
const user = ref(null)
const isLogin = ref(false)
const orderTimeout = ref(1800)
const loading = ref(true)
const loadingText = ref('加载中...')

const store = useSharedStore()
const userOrderKey = ref(0)
const showQrCode = ref(false)
const qrImg = ref('')
const title = ref('')
const handler = ref(null)

/** 右侧账号操作按钮（文案 + iconfont 类名 + 点击 key） */
const accountActions = [
  { key: 'nickname', label: '修改资料', icon: 'icon-user-fill' },
  { key: 'email', label: '绑定邮箱', icon: 'icon-email' },
  { key: 'mobile', label: '绑定手机', icon: 'icon-mobile' },
  { key: 'password', label: '修改密码', icon: 'icon-password' },
  { key: 'redeem', label: '卡密兑换', icon: 'icon-redeem' },
]

const displayName = computed(() => {
  if (!user.value) {
    return '未登录'
  }
  return user.value.nickname || user.value.username || user.value.email || '未设置'
})

const userScore = computed(() => {
  if (!user.value || typeof user.value.power === 'undefined' || user.value.power === null) {
    return '--'
  }
  return user.value.power
})

onMounted(() => {
  checkSession()
    .then((_user) => {
      user.value = _user
      isLogin.value = true
      // 获取最新用户信息
      httpGet('/api/user/profile')
        .then((res) => {
          user.value = { ...user.value, ...res.data }
        })
        .catch((e) => {
          console.error('获取用户信息失败：' + e.message)
        })
    })
    .catch(() => {
      store.setShowLoginDialog(true)
    })

  httpGet('/api/product/list')
    .then((res) => {
      list.value = res.data
      loading.value = false
    })
    .catch((e) => {
      ElMessage.error('获取产品套餐失败：' + e.message)
    })

  getSystemInfo()
    .then((res) => {
      if (res.data['order_pay_timeout'] > 0) {
        orderTimeout.value = res.data['order_pay_timeout']
      }
    })
    .catch((e) => {
      ElMessage.error('获取系统配置失败：' + e.message)
    })
})

// 修改昵称
const openNicknameDialog = () => {
  nicknameInput.value = user.value?.nickname || ''
  usernameInput.value = user.value?.username || ''
  showNicknameDialog.value = true
}

const handleAccountAction = (key) => {
  switch (key) {
    case 'nickname':
      openNicknameDialog()
      break
    case 'email':
      showBindEmailDialog.value = true
      break
    case 'mobile':
      showBindMobileDialog.value = true
      break
    case 'password':
      showPasswordDialog.value = true
      break
    case 'redeem':
      showRedeemVerifyDialog.value = true
      break
    default:
      break
  }
}

const saveNickname = () => {
  httpPost('/api/user/profile/update', {
    nickname: nicknameInput.value,
    username: usernameInput.value,
  })
    .then(() => {
      if (user.value) {
        user.value.nickname = nicknameInput.value
        user.value.username = usernameInput.value
      }
      ElMessage.success({ message: '保存成功', duration: 500 })
      showNicknameDialog.value = false
    })
    .catch((e) => {
      ElMessage.error('保存失败：' + e.message)
    })
}

// 头像上传处理函数
const afterRead = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData()
      formData.append('file', result, result.name)
      // 执行上传操作
      httpPost('/api/upload', formData)
        .then((res) => {
          user.value.avatar = res.data.url
          // 自动保存用户信息
          httpPost('/api/user/profile/update', { avatar: res.data.url })
            .then(() => {
              ElMessage.success({ message: '头像上传成功', duration: 500 })
            })
            .catch((e) => {
              ElMessage.error('保存失败：' + e.message)
            })
        })
        .catch((e) => {
          ElMessage.error('图片上传失败:' + e.message)
        })
    },
    error(err) {
      console.log(err.message)
      ElMessage.error('图片处理失败')
    },
  })
}

const selectedPid = ref(0)
const wxPay = (product) => {
  selectedPid.value = product.id
  title.value = '请打开微信扫码支付'
  generateOrder('wxpay')
}

const alipay = (product) => {
  selectedPid.value = product.id
  title.value = '请打开支付宝扫码支付'
  generateOrder('alipay')
}

const generateOrder = (payWay) => {
  showLoading('正在生成支付订单...')
  // 生成支付订单
  httpPost('/api/payment/create', {
    pid: selectedPid.value,
    pay_way: payWay,
    domain: `${window.location.protocol}//${window.location.host}`,
    device: isMobile() ? 'mobile' : 'pc',
  })
    .then((res) => {
      closeLoading()

      if (isMobile()) {
        window.location.href = res.data.pay_url
      } else {
        QRCode.toDataURL(res.data.pay_url, { width: 300, height: 300, margin: 2 }, (error, url) => {
          if (error) {
            console.error(error)
          } else {
            qrImg.value = url
          }
        })
        // 查询订单状态
        if (handler.value) {
          clearTimeout(handler.value)
        }
        handler.value = setTimeout(() => queryOrder(res.data.order_no), 3000)
        showQrCode.value = true
      }
    })
    .catch((e) => {
      closeLoading()
      ElMessage.error('生成支付订单失败：' + e.message)
    })
}

const queryOrder = async (orderNo) => {
  const res = await httpGet('/api/order/query?order_no=' + orderNo)
  if (res?.data.status === 2) {
    // 订单支付成功
    clearTimeout(handler.value)
    ElMessage.success('支付成功')
    showQrCode.value = false
    // 更新用户积分
    if (user.value) {
      const currentScore = Number(user.value.scores) || 0
      user.value.scores = currentScore + Number(res.data.credit || 0)
    }
  } else {
    handler.value = setTimeout(() => queryOrder(orderNo), 3000)
  }
}

const redeemCallback = (success) => {
  showRedeemVerifyDialog.value = false

  if (success) {
    userOrderKey.value += 1
  }
}

// 组件卸载时清理定时器
onUnmounted(() => {
  if (handler.value) {
    clearTimeout(handler.value)
    handler.value = null
  }
})
</script>
