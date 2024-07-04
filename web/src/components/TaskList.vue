<template>
  <div class="running-job-list">
    <div class="running-job-box" v-if="list.length > 0">
      <div class="job-item" v-for="item in list">
        <div v-if="item.progress > 0" class="job-item-inner">
          <el-image v-if="item.img_url" :src="item['img_url']" fit="cover" loading="lazy">
            <template #placeholder>
              <div class="image-slot">
                正在加载图片
              </div>
            </template>

            <template #error>
              <div class="image-slot">
                <el-icon>
                  <Picture/>
                </el-icon>
              </div>
            </template>
          </el-image>

          <div class="progress">
            <el-progress type="circle" :percentage="item.progress" :width="100"
                         color="#47fff1"/>
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
    <el-empty :image-size="100" v-else/>
  </div>
</template>

<script setup>
import {ref} from "vue";
import {CircleCloseFilled, Picture} from "@element-plus/icons-vue";
import {isImage, removeArrayItem, substr} from "@/utils/libs";
import {FormatFileSize, GetFileIcon, GetFileType} from "@/store/system";

const props = defineProps({
  list: {
    type: Array,
    default:[],
  }
})

</script>

<style scoped lang="stylus">
@import "~@/assets/css/running-job-list.styl"
</style>