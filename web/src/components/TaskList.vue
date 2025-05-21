<template>
  <div class="running-job-list pt-4 pb-4">
    <div class="running-job-box" v-if="list.length > 0">
      <div class="job-item" v-for="item in list" :key="item.id">
        <div v-if="item.progress > 0" class="job-item-inner">
          <div class="progress" v-if="item.progress > 0">
            <el-progress type="circle" :percentage="item.progress" :width="100" color="#47fff1" />
          </div>
        </div>
        <el-image fit="cover" v-else>
          <template #error>
            <div class="image-slot">
              <i class="iconfont icon-quick-start"></i>
              <span>任务正在排队中</span>
            </div>
          </template>
        </el-image>
      </div>
    </div>
    <el-empty :image-size="100" v-else :image="nodata" description="暂无任务" />
  </div>
</template>

<script setup>
import nodata from "@/assets/img/no-data.png";

// eslint-disable-next-line no-undef
const props = defineProps({
  list: {
    type: Array,
    default: [],
  },
});
</script>

<style scoped lang="stylus">
@import "../assets/css/running-job-list.styl"
</style>
