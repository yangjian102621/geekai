<template>
  <div class="list-box" ref="containerRef">
    <el-row :gutter="gap">
      <el-col v-for="item in items" :key="item.id" :span="span" :style="{marginBottom:gap+'px'} ">
        <slot :item="item"></slot>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
// 列表组件
import {onMounted, ref} from "vue";

// eslint-disable-next-line no-undef
const props = defineProps({
  items: {
    type: Array,
    required: true
  },
  gap: {
    type: Number,
    default: 10
  },
  width: {
    type: Number,
    default: 240
  }
});

const containerRef = ref(null)
const span = ref(12)

onMounted(() => {
  calcSpan()
})

const calcSpan = () => {
  let cols = Math.floor(containerRef.value.offsetWidth / props.width)
  if (cols >= 12) {
    span.value = 1
    return
  }
  console.log(cols)
  while (cols > 1) {
    if (24 % cols === 0) {
      span.value = 24 / cols
      return
    }
    cols -= 1
  }
  span.value = 12
}
window.onresize = () => calcSpan()
</script>

<style scoped lang="stylus">

.list-box {
}

</style>