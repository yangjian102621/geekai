<template>
  <div class="test-3d">
    <h1>3D预览功能测试</h1>

    <div class="test-container">
      <h2>测试1: 默认立方体（无模型URL）</h2>
      <div class="preview-wrapper">
        <ThreeDPreview />
      </div>
    </div>

    <div class="test-container">
      <h2>测试2: 带模型URL（如果有的话）</h2>
      <div class="preview-wrapper">
        <ThreeDPreview v-if="testModelUrl" :model-url="testModelUrl" :model-type="testModelType" />
        <div v-else class="no-model">
          <p>没有测试模型URL</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import ThreeDPreview from '@/components/ThreeDPreview.vue'
import { ref } from 'vue'

// 测试用的模型URL（可以替换为实际的模型文件）
const testModelUrl = ref('https://img.r9it.com/R03TQZ7PZ386RGL7PTMNGFOHAJW15WYF.glb')
const testModelType = ref('glb')
</script>

<style lang="scss" scoped>
.test-3d {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.test-container {
  margin-bottom: 40px;

  h2 {
    margin-bottom: 16px;
    color: #333;
  }
}

.preview-wrapper {
  width: 100%;
  height: 500px;
  border: 2px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
}

.no-model {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  color: #666;
}
</style>
