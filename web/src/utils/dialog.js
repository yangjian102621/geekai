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

// 自定义 Toast 消息系统
export function showToastMessage(message, type = 'info', duration = 3000) {
  const toast = document.createElement('div')
  toast.className = `fixed top-20 left-1/2 transform -translate-x-1/2 z-50 px-4 py-2 rounded-lg text-white font-medium ${
    type === 'error' ? 'bg-red-500' : type === 'success' ? 'bg-green-500' : 'bg-blue-500'
  } animate-fade-in`
  toast.textContent = message
  document.body.appendChild(toast)

  if (duration > 0) {
    setTimeout(() => {
      toast.classList.add('animate-fade-out')
      setTimeout(() => {
        if (document.body.contains(toast)) {
          document.body.removeChild(toast)
        }
      }, 300)
    }, duration)
  }
}
