<template>
  <el-container class="chat-file-list">
    <div v-for="file in fileList">
      <div class="image" v-if="isImage(file.ext)">
        <el-image :src="file.url" fit="cover"/>
        <div class="action">
          <el-icon @click="removeFile(file)"><CircleCloseFilled /></el-icon>
        </div>
      </div>
      <div class="item" v-else>
        <div class="icon">
          <el-image :src="GetFileIcon(file.ext)" fit="cover"  />
        </div>
        <div class="body">
          <div class="title">
            <el-link :href="file.url" target="_blank" style="--el-font-weight-primary:bold">{{substr(file.name, 30)}}</el-link>
          </div>
          <div class="info">
            <span>{{GetFileType(file.ext)}}</span>
            <span>{{FormatFileSize(file.size)}}</span>
          </div>
        </div>
        <div class="action">
          <el-icon @click="removeFile(file)"><CircleCloseFilled /></el-icon>
        </div>
      </div>
    </div>
  </el-container>
</template>

<script setup>
import {ref} from "vue";
import {CircleCloseFilled} from "@element-plus/icons-vue";
import {isImage, removeArrayItem, substr} from "@/utils/libs";
import {FormatFileSize, GetFileIcon, GetFileType} from "@/store/system";

const props = defineProps({
  files: {
    type: Array,
    default:[],
  }
})
const emits = defineEmits(['removeFile']);
const fileList = ref(props.files)


const removeFile = (file) => {
  fileList.value = removeArrayItem(fileList.value, file, (v1,v2) => v1.url===v2.url)
  emits('removeFile', file)
}

</script>

<style scoped lang="stylus">

.chat-file-list {
  display flex
  flex-flow row
  .image {
    display flex
    flex-flow row
    margin-right 10px
    position relative

    .el-image {
      height 56px
      width 56px
      border 1px solid #e3e3e3
      border-radius 10px
    }
  }
  .item {
    position relative
    display flex
    flex-flow row
    border-radius 10px
    background-color #ffffff
    border 1px solid #e3e3e3
    padding 6px
    margin-right 10px

    .icon {
      .el-image {
        width 40px
        height 40px
      }
    }
    .body {
      margin-left 5px
      font-size 14px
      .title {
        line-height 24px
        color #0D0D0D
      }
      .info {
        color #B4B4B4

        span {
          margin-right 10px
        }
      }
    }
  }

  .action {
    position absolute
    top -8px
    right -8px
    color #da0d54
    cursor pointer
    font-size 20px
  }
}

</style>