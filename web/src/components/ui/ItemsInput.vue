<template>
  <!-- 多项目输入组件 -->
  <div class="items-input-box">
    <el-tag
        v-for="tag in tags"
        :key="tag"
        closable
        :disable-transitions="false"
        @close="handleClose(tag)"
    >
      {{ tag }}
    </el-tag>
    <el-input
        v-if="inputVisible"
        ref="InputRef"
        v-model="inputValue"
        class="w-20"
        size="small"
        @keyup.enter="handleInputConfirm"
        @blur="handleInputConfirm"
    />
    <el-button v-else class="button-new-tag" size="small" @click="showInput">
      + 新增
    </el-button>
  </div>
</template>
<script setup>

import {nextTick, ref, watch} from "vue";
// eslint-disable-next-line no-undef
const props = defineProps({
  value : {
    type: Array,
    default: () => []
  },
});
// eslint-disable-next-line no-undef
const emits = defineEmits(['update:value']);
const tags = ref(props.value)
const inputValue = ref('')
const inputVisible = ref(false)
const InputRef = ref(null)

watch(() => props.value, (newValue) => {
  tags.value = newValue
})

const handleClose = (tag) => {
  tags.value.splice(tags.value.indexOf(tag), 1)
}

const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    InputRef.value?.input?.focus()
  })
}

const handleInputConfirm = () => {
  if (inputValue.value) {
    tags.value.push(inputValue.value)
  }
  inputVisible.value = false
  inputValue.value = ''
  emits('update:value', tags.value)
}
</script>

<style scoped lang="stylus">
.items-input-box {
  display flex

  .el-tag {
    display flex
    margin-right 6px
  }
}

</style>