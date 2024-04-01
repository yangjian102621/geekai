<template>
  <el-container class="file-list-box">
    <el-tooltip class="box-item" effect="dark" content="打开文件管理中心">
      <el-button class="file-upload-img" @click="fetchFiles">
        <el-icon>
          <PictureFilled/>
        </el-icon>
      </el-button>
    </el-tooltip>

    <el-dialog
        v-model="show"
        :close-on-click-modal="true"
        :show-close="true"
        :width="800"
        title="文件管理"
    >

      <div class="file-list">
        <el-row :gutter="20">
          <el-col :span="3">
            <div class="grid-content">
              <el-upload
                  class="avatar-uploader"
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="afterRead"
              >
                <el-icon class="avatar-uploader-icon">
                  <Plus/>
                </el-icon>
              </el-upload>
            </div>
          </el-col>
          <el-col :span="3" v-for="file in fileList" :key="file.url">
            <div class="grid-content">
              <el-tooltip
                  class="box-item"
                  effect="dark"
                  :content="file.name"
                  placement="top">
                <el-image :src="file.url" fit="cover" v-if="isImage(file.ext)" @click="insertURL(file.url)"/>
                <el-image :src="getFileIcon(file.ext)" fit="cover" v-else @click="insertURL(file.url)"/>
              </el-tooltip>

              <div class="opt">
                <el-button type="danger" size="small" :icon="Delete" @click="removeFile(file)" circle/>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-dialog>
  </el-container>
</template>

<script setup>
import {ref} from "vue";
import {ElMessage} from "element-plus";
import {httpGet, httpPost} from "@/utils/http";
import {Delete, PictureFilled, Plus} from "@element-plus/icons-vue";
import {isImage, removeArrayItem} from "@/utils/libs";

const props = defineProps({
  userId: Number,
});
const emits = defineEmits(['selected']);
const show = ref(false)
const fileList = ref([])

const fetchFiles = () => {
  show.value = true
  httpGet("/api/upload/list").then(res => {
    fileList.value = res.data
  }).catch(() => {
  })
}

const getFileIcon = (ext) => {
  const files = {
    ".docx": "doc.png",
    ".doc": "doc.png",
    ".xls": "xls.png",
    ".xlsx": "xls.png",
    ".ppt": "ppt.png",
    ".pptx": "ppt.png",
    ".md": "md.png",
    ".pdf": "pdf.png",
    ".sql": "sql.png"
  }
  if (files[ext]) {
    return '/images/ext/' + files[ext]
  }

  return '/images/ext/file.png'
}

const afterRead = (file) => {
  const formData = new FormData();
  formData.append('file', file.file, file.name);
  // 执行上传操作
  httpPost('/api/upload', formData).then((res) => {
    fileList.value.unshift(res.data)
    ElMessage.success({message: "上传成功", duration: 500})
  }).catch((e) => {
    ElMessage.error('图片上传失败:' + e.message)
  })
};

const removeFile = (file) => {
  httpGet('/api/upload/remove?id=' + file.id).then(() => {
    fileList.value = removeArrayItem(fileList.value, file, (v1, v2) => {
      return v1.id === v2.id
    })
    ElMessage.success("文件删除成功！")
  }).catch((e) => {
    ElMessage.error('文件删除失败:' + e.message)
  })
}

const insertURL = (url) => {
  show.value = false
  // 如果是相对路径，处理成绝对路径
  if (url.indexOf("http") === -1) {
    url = location.protocol + "//" + location.host + url
  }
  emits('selected', url)
}
</script>

<style lang="stylus">

.file-list-box {
  .file-upload-img {
    padding: 8px 5px;
    border-radius: 6px;
    background: #19c37d;
    color: #fff;
    font-size: 20px;
  }

  .el-dialog {

    .el-dialog__body {
      //padding 0

      .file-list {

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