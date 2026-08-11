<template>
  <div class="mobile-image container">
    <CustomTabs :model-value="activeName" @update:model-value="activeName = $event" class="my-tab">
      <CustomTabPane name="mj" label="MJ" v-if="activeMenu.mj">
        <image-mj />
      </CustomTabPane>
      <CustomTabPane name="image" label="AI图像生成" v-if="activeMenu.image">
        <image-page />
      </CustomTabPane>
    </CustomTabs>
  </div>
</template>

<script setup>
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { httpGet } from '@/utils/http'
import ImagePage from '@/views/mobile/pages/Image.vue'
import ImageMj from '@/views/mobile/pages/ImageMj.vue'
import { onMounted, ref } from 'vue'

const activeName = ref('')
const menus = ref([])
const activeMenu = ref({
  mj: false,
  image: false,
})

onMounted(() => {
  httpGet('/api/menu/list').then((res) => {
    menus.value = res.data
    activeMenu.value = {
      mj: menus.value.some((item) => item.url === '/mj'),
      image: menus.value.some((item) => item.url === '/image'),
    }
  })
})
</script>

<style lang="scss">
.mobile-image {
  .my-tab {
    .van-tab__panel {
      padding: 10px;
    }
  }
}
</style>
