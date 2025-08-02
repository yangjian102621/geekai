<template>
  <div class="mobile-feedback">
    <!-- 顶部导航 -->
    <van-nav-bar title="意见反馈" left-arrow @click-left="router.back()" fixed placeholder />

    <!-- 反馈表单 -->
    <div class="feedback-content">
      <!-- 反馈类型 -->
      <van-cell-group title="反馈类型">
        <van-radio-group v-model="feedbackType" direction="horizontal">
          <van-radio name="bug">问题反馈</van-radio>
          <van-radio name="feature">功能建议</van-radio>
          <van-radio name="other">其他</van-radio>
        </van-radio-group>
      </van-cell-group>

      <!-- 反馈内容 -->
      <van-cell-group title="反馈内容">
        <van-field
          v-model="feedbackContent"
          type="textarea"
          placeholder="请详细描述您遇到的问题或建议..."
          :rows="6"
          maxlength="500"
          show-word-limit
          autosize
        />
      </van-cell-group>

      <!-- 联系方式 -->
      <van-cell-group title="联系方式（选填）">
        <van-field v-model="contactInfo" placeholder="邮箱或手机号，方便我们回复您" clearable />
      </van-cell-group>

      <!-- 图片上传 -->
      <van-cell-group title="上传截图（选填）">
        <van-uploader
          v-model="fileList"
          :max-count="3"
          :after-read="afterRead"
          :before-delete="beforeDelete"
          upload-text="上传图片"
        />
      </van-cell-group>

      <!-- 提交按钮 -->
      <div class="submit-section">
        <van-button type="primary" block :loading="submitting" @click="submitFeedback">
          提交反馈
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { showSuccessToast, showToast } from 'vant'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const feedbackType = ref('bug')
const feedbackContent = ref('')
const contactInfo = ref('')
const fileList = ref([])
const submitting = ref(false)

// 上传图片后的处理
const afterRead = (file) => {
  // 这里可以处理图片上传逻辑
  console.log('上传文件:', file)
}

// 删除图片前的处理
const beforeDelete = (file, detail) => {
  // 这里可以处理图片删除逻辑
  console.log('删除文件:', file)
  return true
}

// 提交反馈
const submitFeedback = async () => {
  if (!feedbackContent.value.trim()) {
    showToast('请输入反馈内容')
    return
  }

  submitting.value = true

  try {
    const feedbackData = {
      type: feedbackType.value,
      content: feedbackContent.value,
      contact: contactInfo.value,
      images: fileList.value.map((file) => file.url || file.content),
      timestamp: new Date().toISOString(),
    }

    // 暂时使用本地存储保存反馈数据
    // 后续可以对接后端API
    const existingFeedback = JSON.parse(localStorage.getItem('userFeedback') || '[]')
    existingFeedback.push(feedbackData)
    localStorage.setItem('userFeedback', JSON.stringify(existingFeedback))

    showSuccessToast('反馈提交成功，感谢您的建议！')

    // 清空表单
    feedbackContent.value = ''
    contactInfo.value = ''
    fileList.value = []

    // 返回上一页
    setTimeout(() => {
      router.back()
    }, 1500)
  } catch (error) {
    console.error('提交反馈失败:', error)
    showToast('提交失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="scss" scoped>
.mobile-feedback {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.feedback-content {
  padding: 16px;
}

.van-cell-group {
  margin-bottom: 16px;
  border-radius: 8px;
  overflow: hidden;
}

.van-radio-group {
  padding: 16px;
  display: flex;
  gap: 24px;
}

.van-field {
  padding: 16px;
}

.submit-section {
  margin-top: 32px;
  padding: 0 16px;
}

.van-uploader {
  padding: 16px;
}

// 自定义样式
:deep(.van-cell-group__title) {
  padding: 16px 16px 8px;
  font-size: 14px;
  font-weight: 500;
  color: #323233;
}

:deep(.van-field__control) {
  min-height: 120px;
}

:deep(.van-radio__label) {
  font-size: 14px;
}
</style>
