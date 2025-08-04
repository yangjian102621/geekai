/**
 * Util lib functions
 */
import { isMobile } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import {
  closeToast,
  showConfirmDialog,
  showFailToast,
  showLoadingToast,
  showSuccessToast,
  showToast,
} from 'vant'

export function showLoginDialog(router) {
  showConfirmDialog({
    title: '登录',
    message: '此操作需要登录才能进行，前往登录？',
  })
    .then(() => {
      router.push('/login')
    })
    .catch(() => {
      // on cancel
    })
}

export function showMessageOK(message) {
  if (isMobile()) {
    showSuccessToast(message)
  } else {
    ElMessage.success(message)
  }
}

export function showMessageInfo(message) {
  if (isMobile()) {
    showToast(message)
  } else {
    ElMessage.info(message)
  }
}

export function showMessageError(message) {
  if (isMobile()) {
    showFailToast({ message: message })
  } else {
    ElMessage.error(message)
  }
}

export function showLoading(message = '正在处理...') {
  showLoadingToast({ message: message, forbidClick: true, duration: 0 })
}

export function closeLoading() {
  closeToast()
}
