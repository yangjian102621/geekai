<template>
  <div class="list-box" ref="container">
    <div class="list-inner">
      <div
          class="list-item"
          v-for="(item, index) in items"
          :key="index"
          :style="{width:itemWidth + 'px'}"
      >
        <div class="item-inner" :style="{padding: gap/2+'px'}">
          <div class="item-wrapper">
            <slot :item="item" :index="index" :width="itemWidth"></slot>
          </div>
        </div>
      </div>
    </div>
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
    default: 12
  },
  width: {
    type: Number,
    default: 240
  }
});

const container = ref(null)
const itemWidth = ref(props.width)

onMounted(() => {
  computeSize()
})

const computeSize = () => {
  const w = container.value.offsetWidth - 10 // 减去滚动条的宽度
  let cols = Math.floor(w / props.width)
  itemWidth.value = Math.floor(w / cols)
}

window.onresize = () => {
  computeSize()
}
</script>

<style scoped lang="stylus">

.list-box {

  .list-inner {
    display flex
    flex-wrap wrap

    .list-item {
      .item-inner {
        display flex

        .item-wrapper {
          height 100%
          width 100%
          display flex
          justify-content center
        }
      }
    }
  }
}

</style>