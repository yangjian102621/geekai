<template>
  <div class="waterfall" ref="container">
    <div class="waterfall-inner">
      <div
          class="waterfall-item"
          v-for="(item, index) in items"
          :key="index"
          :style="{width:itemWidth + 'px', height:height+'px', marginBottom: margin*2+'px'}"
      >
        <div :style="{marginLeft: margin+'px', marginRight: margin+'px'}">
          <div class="item-wrapper">
            <slot :item="item" :index="index"></slot>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 列表组件
import {onMounted, ref} from "vue";

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
  },
  height: {
    type: Number,
    default: 240
  }
});

const container = ref(null)
const itemWidth = ref(props.width)
const margin = ref(props.gap)

onMounted(() => {
  computeSize()
})

const computeSize = () => {
  const w = container.value.offsetWidth - 10 // 减去滚动条的宽度
  let cols = Math.floor(w / props.width)
  itemWidth.value = w / cols
  while (itemWidth.value < props.width && cols > 1) {
    cols -= 1
    itemWidth.value = w / cols
  }

  if (props.gap > 0) {
    margin.value = props.gap / 2
  }
}

window.onresize = () => {
  computeSize()
}
</script>

<style scoped lang="stylus">

.waterfall {

  .waterfall-inner {
    display flex
    flex-wrap wrap

    .waterfall-item {

      div {
        display flex
        height 100%
        overflow hidden

        .item-wrapper {
          height 100%
          width 100%
        }
      }
    }
  }
}

</style>