<template>
  <div class="audio-chat-page bg-gray-50">
    <div class="flex m-3">
      <el-select v-model="currentFunction" placeholder="请选择功能" popper-class="custom-select">
        <template #prefix>
          <i
            class="iconfont !text-xl"
            :class="functions.find((f) => f.key === currentFunction).icon"
          ></i>
        </template>
        <el-option v-for="f in functions" :value="f.key" :key="f.key" :label="f.name">
          <div class="flex items-center space-x-2">
            <i class="iconfont !text-xl" :class="f.icon"></i>
            <span>{{ f.name }}</span>
          </div>
        </el-option>
      </el-select>
    </div>

    <div class="p-3">
      <param-builder-mobile
        v-model="formData"
        :items="params[currentFunction]"
        :progress="progress[currentFunction]"
      />

      <!-- 调试信息 -->
      <div class="mt-6 p-4 bg-gray-100 rounded-lg">
        <h3 class="text-lg font-bold mb-2">调试信息</h3>
        <div class="text-sm">
          <p><strong>当前功能:</strong> {{ currentFunction }}</p>
          <p><strong>表单数据:</strong></p>
          <pre class="bg-white p-2 rounded mt-1 overflow-auto">{{
            JSON.stringify(formData, null, 2)
          }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import ParamBuilderMobile from '@/components/mobile/ParamBuilderMobile.vue'
import { JimengFunctions, JimengParams } from '@/store/data/jimeng_params'
import { ref } from 'vue'

const functions = JimengFunctions
const params = JimengParams
const formData = ref({})
const currentFunction = ref('image')
const progress = ref({
  image: 100,
  video: 100,
  virtualHuman: 38,
  actionTransfer: 65,
})
</script>

<style scoped lang="scss">
.custom-select {
  .el-select-dropdown__item.is-selected {
    background-color: var(--el-fill-color-light);
  }
}
</style>
