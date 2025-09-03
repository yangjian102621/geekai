<template>
  <div class="audio-chat-page">
    <!-- 图像特效 -->
    <div class="jimeng-create__function-panel">
      <div class="jimeng-create__param-line">
        <span class="jimeng-create__param-label">上传图片：</span>
      </div>
      <div class="jimeng-create__param-line">
        <div class="jimeng-create__upload">
          <input
            ref="imageEffectsInput"
            type="file"
            accept=".jpg,.png,.jpeg"
            @change="
              (e) =>
                jimengStore.onImageUpload({
                  file: e.target.files[0],
                  name: e.target.files[0]?.name,
                })
            "
            class="hidden"
          />
          <div @click="$refs.imageEffectsInput?.click()" class="jimeng-create__upload-content">
            <i
              v-if="!jimengStore.imageEffectsParams.image_input1.length"
              class="jimeng-create__upload-icon iconfont icon-upload"
            ></i>
            <span
              v-if="!jimengStore.imageEffectsParams.image_input1.length"
              class="jimeng-create__upload-text"
              >上传图片</span
            >
            <div v-else class="jimeng-create__upload-preview">
              <el-image
                :src="
                  jimengStore.imageEffectsParams.image_input1[0]?.url ||
                  jimengStore.imageEffectsParams.image_input1[0]?.content
                "
                fit="cover"
                class="w-32 h-32 rounded"
              />
              <button
                @click.stop="jimengStore.imageEffectsParams.image_input1 = []"
                class="jimeng-create__upload-remove-btn"
              >
                <i class="iconfont icon-close"></i>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="jimeng-create__param-line">
        <span class="jimeng-create__param-label">特效模板：</span>
      </div>
      <div class="jimeng-create__param-line">
        <CustomSelect
          v-model="jimengStore.imageEffectsParams.template_id"
          :options="
            jimengStore.imageEffectsTemplateOptions.map((opt) => ({
              label: opt.label,
              value: opt.value,
            }))
          "
          label="特效模板"
          title="选择特效模板"
          @update:modelValue="handleTemplateChange"
        />
      </div>

      <div class="jimeng-create__param-line">
        <span class="jimeng-create__param-label">输出尺寸：</span>
      </div>
      <div class="jimeng-create__param-line">
        <CustomSelect
          v-model="jimengStore.imageEffectsParams.size"
          :options="
            jimengStore.imageSizeOptions.map((opt) => ({
              label: opt.label,
              value: opt.value,
            }))
          "
          label="输出尺寸"
          title="选择尺寸"
        />
      </div>
    </div>

    <!-- 文生视频 -->
    <div
      v-if="jimengStore.activeFunction === 'text_to_video'"
      class="jimeng-create__function-panel"
    >
      <div class="jimeng-create__param-line">
        <span class="jimeng-create__param-label">提示词：</span>
      </div>
      <div class="jimeng-create__param-line">
        <textarea
          v-model="jimengStore.currentPrompt"
          placeholder="描述你想要的视频内容"
          class="jimeng-create__form-section-textarea"
          rows="4"
          maxlength="2000"
        />
        <div class="jimeng-create__form-section-counter">
          <span>{{ jimengStore.currentPrompt.length }}/2000</span>
        </div>
      </div>

      <div class="jimeng-create__param-line">
        <span class="jimeng-create__param-label">视频比例：</span>
      </div>
      <div class="jimeng-create__param-line">
        <CustomSelect
          v-model="jimengStore.textToVideoParams.aspect_ratio"
          :options="
            jimengStore.videoAspectRatioOptions.map((opt) => ({
              label: opt.label,
              value: opt.value,
            }))
          "
          label="视频比例"
          title="选择比例"
        />
      </div>
    </div>

    <!-- 图生视频 -->
    <div
      v-if="jimengStore.activeFunction === 'image_to_video'"
      class="jimeng-create__function-panel"
    >
      <div class="jimeng-create__param-line">
        <span class="jimeng-create__param-label">上传图片：</span>
      </div>
      <div class="jimeng-create__param-line">
        <div class="jimeng-create__upload">
          <input
            ref="imageToVideoInput"
            type="file"
            accept=".jpg,.png,.jpeg"
            multiple
            @change="(e) => jimengStore.handleMultipleImageUpload(e)"
            class="hidden"
          />
          <div @click="$refs.imageToVideoInput?.click()" class="jimeng-create__upload-content">
            <i
              v-if="!jimengStore.imageToVideoParams.image_urls.length"
              class="jimeng-create__upload-icon iconfont icon-upload"
            ></i>
            <span
              v-if="!jimengStore.imageToVideoParams.image_urls.length"
              class="jimeng-create__upload-text"
              >上传图片</span
            >
            <div v-else class="jimeng-create__upload-multiple">
              <div
                v-for="(image, index) in jimengStore.imageToVideoParams.image_urls"
                :key="index"
                class="jimeng-create__upload-multiple-item"
              >
                <el-image
                  :src="image?.url || image?.content"
                  fit="cover"
                  class="w-24 h-24 rounded"
                />
                <button
                  @click.stop="jimengStore.removeImage(index)"
                  class="jimeng-create__upload-remove-btn"
                >
                  <i class="iconfont icon-close"></i>
                </button>
              </div>
              <div
                v-if="jimengStore.imageToVideoParams.image_urls.length < 2"
                @click.stop="$refs.imageToVideoInput?.click()"
                class="jimeng-create__upload-multiple-add"
              >
                <i class="iconfont icon-plus"></i>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="jimeng-create__param-line">
        <span class="jimeng-create__param-label">提示词：</span>
      </div>
      <div class="jimeng-create__param-line">
        <textarea
          v-model="jimengStore.currentPrompt"
          placeholder="描述你想要的视频效果"
          class="jimeng-create__form-section-textarea"
          rows="4"
          maxlength="2000"
        />
        <div class="jimeng-create__form-section-counter">
          <span>{{ jimengStore.currentPrompt.length }}/2000</span>
        </div>
      </div>

      <div class="jimeng-create__param-line">
        <span class="jimeng-create__param-label">视频比例：</span>
      </div>
      <div class="jimeng-create__param-line">
        <CustomSelect
          v-model="jimengStore.imageToVideoParams.aspect_ratio"
          :options="
            jimengStore.videoAspectRatioOptions.map((opt) => ({
              label: opt.label,
              value: opt.value,
            }))
          "
          label="视频比例"
          title="选择比例"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
const connect = () => {}
</script>

<style scoped lang="scss">
.audio-chat-page {
  display: flex;
  flex-flow: column;
  justify-content: center;
  align-items: center;
}

canvas {
  background-color: transparent;
}
</style>
