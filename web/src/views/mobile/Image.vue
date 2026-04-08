<template>
  <div class="mobile-image container">
    <CustomTabs :model-value="activeName" @update:model-value="activeName = $event" class="my-tab">
      <CustomTabPane name="mj" label="MJ" v-if="activeMenu.mj">
        <image-mj />
      </CustomTabPane>
      <CustomTabPane name="sd" label="SD" v-if="activeMenu.sd">
        <image-sd />
      </CustomTabPane>
      <CustomTabPane name="dall" label="DALL" v-if="activeMenu.dall">
        <image-dall />
      </CustomTabPane>
    </CustomTabs>
  </div>
</template>

<script setup>
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { httpGet } from '@/utils/http'
import ImageDall from '@/views/mobile/pages/ImageDall.vue'
import ImageMj from '@/views/mobile/pages/ImageMj.vue'
import ImageSd from '@/views/mobile/pages/ImageSd.vue'
import { onMounted, ref } from 'vue'

const activeName = ref('')
const menus = ref([])
const activeMenu = ref({
  mj: false,
  sd: false,
  dall: false,
})

onMounted(() => {
  httpGet('/api/menu/list').then((res) => {
    menus.value = res.data
    activeMenu.value = {
      mj: menus.value.some((item) => item.url === '/mj'),
      sd: menus.value.some((item) => item.url === '/sd'),
      dall: menus.value.some((item) => item.url === '/dalle'),
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
