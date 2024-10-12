<template>
  <el-container class="file-select-box">
    <a class="file-upload-img" @click="fetchFiles(1)">
      <i class="iconfont icon-attachment-st"></i>
    </a>
    <el-dialog
        class="file-list-dialog"
        v-model="show"
        :close-on-click-modal="true"
        :show-close="true"
        :width="800"
        title="文件管理"
    >
      <el-scrollbar ref="scrollbarRef" max-height="80vh" style="height: 100%;" @scroll="onScroll">
        <div class="file-list">
          <el-row :gutter="20">
            <el-col :span="3">
              <div class="grid-content">
                <el-upload
                    class="avatar-uploader"
                    :auto-upload="true"
                    :show-file-list="false"
                    :http-request="afterRead"
                    accept=".doc,.docx,.jpg,.png,.jpeg,.xls,.xlsx,.ppt,.pptx,.pdf,.mp4,.mp3"
                >
                  <el-icon class="avatar-uploader-icon">
                    <Plus/>
                  </el-icon>
                </el-upload>
              </div>
            </el-col>
            <el-col :span="3" v-for="file in fileData.items" :key="file.url">
              <div class="grid-content">
                <el-tooltip
                    class="box-item"
                    effect="dark"
                    :content="file.name"
                    placement="top">
                  <el-image :src="file.url" fit="cover" v-if="isImage(file.ext)" @click="insertURL(file)"/>
                  <el-image :src="GetFileIcon(file.ext)" fit="cover" v-else @click="insertURL(file)"/>
                </el-tooltip>

                <div class="opt">
                  <el-button type="danger" size="small" :icon="Delete" @click="removeFile(file)" circle/>
                </div>
              </div>
            </el-col>
          </el-row>
          <el-row justify="center" v-if="!fileData.isLastPage" @click="fetchFiles(fileData.page)">
            <el-link>加载更多</el-link>
            
          </el-row>
        </div>        
      </el-scrollbar>

    </el-dialog>
  </el-container>
</template>

<script setup>
import {reactive, ref} from "vue";
import {ElMessage} from "element-plus";
import {httpGet, httpPost} from "@/utils/http";
import {Delete, Plus} from "@element-plus/icons-vue";
import {isImage, removeArrayItem} from "@/utils/libs";
import {GetFileIcon} from "@/store/system";

const props = defineProps({
  userId: Number,
});
const emits = defineEmits(['selected']);
const show = ref(false)
const scrollbarRef = ref(null)
const fileData = reactive({
  items:[],
  page: 1,
  isLastPage: true,
})

const fetchFiles = (pageNo) => {
  if(pageNo === 1) show.value = true
  httpPost("/api/upload/list", { page: pageNo || 1, page_size: 30 }).then(res => {
    const { items, page, total_page } = res.data

    if(page === 1){
      fileData.items = items
    }else{
      fileData.items = [...fileData.items, ...items]
    }

    fileData.isLastPage = (page === total_page)

    if(!fileData.isLastPage){
      fileData.page = page + 1
    }
    
  }).catch(() => {
  })
}

// el-scrollbar 滚动回调
const onScroll = (options) => {
  const wrapRef = scrollbarRef.value.wrapRef
  scrollbarRef.value.moveY = wrapRef.scrollTop * 100 / wrapRef.clientHeight
  scrollbarRef.value.moveX = wrapRef.scrollLeft * 100 / wrapRef.clientWidth
  const poor = wrapRef.scrollHeight - wrapRef.clientHeight
  // 判断滚动到底部 自动加载数据
  if (options.scrollTop + 2 >= poor && !fileData.isLastPage) {
    fetchFiles(fileData.page)
  }
}

const afterRead = (file) => {
  const formData = new FormData();
  formData.append('file', file.file, file.name);
  // 执行上传操作
  httpPost('/api/upload', formData).then((res) => {
    fileData.items.unshift(res.data)
    ElMessage.success({message: "上传成功", duration: 500})
  }).catch((e) => {
    ElMessage.error('图片上传失败:' + e.message)
  })
};

const removeFile = (file) => {
  httpGet('/api/upload/remove?id=' + file.id).then(() => {
    fileData.items = removeArrayItem(fileData.items, file, (v1, v2) => {
      return v1.id === v2.id
    })
    ElMessage.success("文件删除成功！")
    fetchFiles(1)
  }).catch((e) => {
    ElMessage.error('文件删除失败:' + e.message)
  })
}

const insertURL = (file) => {
  show.value = false
  // 如果是相对路径，处理成绝对路径
  if (file.url.indexOf("http") === -1) {
    file.url = location.protocol + "//" + location.host + file.url
  }
  emits('selected', file)
}
</script>

<style lang="stylus">

.file-select-box {
  .file-upload-img {
    .iconfont {
      font-size: 24px;
    }
  }

  .el-dialog {

    .el-dialog__body {
      //padding 0
      overflow hidden

      .file-list {
        margin-right 10px
        .grid-content {
          margin-bottom 10px
          position relative

          .avatar-uploader {
            width 100%
            display: flex;
            justify-content: center;
            align-items: center;
            border 1px dashed #e1e1e1
            border-radius 6px

            .el-upload {
              width 100%
              height 80px
            }
          }

          .el-image {
            width 100%
            height 80px
            border 1px solid #ffffff
            border-radius 6px
            cursor pointer

            &:hover {
              border 1px solid #20a0ff
            }

          }

          .iconfont {
            color #20a0ff
            font-size 40px
          }

          .opt {
            display none
            position absolute
            top 5px
            right 5px
          }

          &:hover {
            .opt {
              display block
            }
          }
        }
      }

    }
  }
}
</style>