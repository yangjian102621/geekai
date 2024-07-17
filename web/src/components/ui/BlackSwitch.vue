<template>

  <el-switch v-model="model" :size="size"
             @change="$emit('update:value', $event)"
             style="--el-switch-on-color:#555555;--el-color-white:#0E0808"/>

</template>

<script>
export default {
  name: 'BlackSwitch',
  props: {
    value : Boolean,
    size: {
      type: String,
      default: 'default',
    }
  },
  data() {
    return {
      model: this.value
    }
  }
}
</script>
