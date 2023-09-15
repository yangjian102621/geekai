<template>
  <div class="chat-line chat-line-mj" v-loading="loading">
    <div class="chat-line-inner">
      <div class="chat-icon">
        <img :src="icon" alt="User"/>
      </div>

      <div class="chat-item">
        <div class="content">
          <div class="text" v-html="data.html"></div>
          <div class="images" v-if="data.image?.url !== ''">
            <el-image :src="data.image?.url"
                      :zoom-rate="1.2"
                      :preview-src-list="[data.image?.url]"
                      fit="cover"
                      :initial-index="0" loading="lazy">
              <template #placeholder>
                <div class="image-slot"
                     :style="{height: height+'px', lineHeight:height+'px'}">
                  正在加载图片<span class="dot">...</span></div>
              </template>

              <template #error>
                <div class="image-slot">
                  <el-icon>
                    <Picture/>
                  </el-icon>
                </div>
              </template>
            </el-image>
          </div>
        </div>

        <div class="opt" v-if="data.showOpt &&data.image?.hash !== ''">
          <div class="opt-line">
            <ul>
              <li><a @click="upscale(1)">U1</a></li>
              <li><a @click="upscale(2)">U2</a></li>
              <li><a @click="upscale(3)">U3</a></li>
              <li><a @click="upscale(4)">U4</a></li>
            </ul>
          </div>

          <div class="opt-line">
            <ul>
              <li><a @click="variation(1)">V1</a></li>
              <li><a @click="variation(2)">V2</a></li>
              <li><a @click="variation(3)">V3</a></li>
              <li><a @click="variation(4)">V4</a></li>
            </ul>
          </div>
        </div>

        <div class="bar" v-if="createdAt !== ''">
          <span class="bar-item"><el-icon><Clock/></el-icon> {{ createdAt }}</span>
          <span class="bar-item">tokens: {{ tokens }}</span>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import {ref, watch} from "vue";
import {Clock, Picture} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";
import {httpPost} from "@/utils/http";
import {getSessionId} from "@/store/session";

const props = defineProps({
  content: Object,
  icon: String,
  chatId: String,
  roleId: Number,
  createdAt: String
});

const data = ref(props.content)
const tokens = ref(0)
const cacheKey = "img_placeholder_height"
const item = localStorage.getItem(cacheKey);
const loading = ref(false)
const height = ref(0)
if (item) {
  height.value = parseInt(item)
}
if (data.value["image"]?.width > 0) {
  height.value = 350 * data.value["image"]?.height / data.value["image"]?.width
  localStorage.setItem(cacheKey, height.value)
}
data.value["showOpt"] = data.value["content"]?.indexOf("- Image #") === -1;
// console.log(data.value)

watch(() => props.content, (newVal) => {
  data.value = newVal;
});
const emits = defineEmits(['disable-input', 'disable-input']);
const upscale = (index) => {
  send('/api/mj/upscale', index)
}

const variation = (index) => {
  send('/api/mj/variation', index)
}

const send = (url, index) => {
  loading.value = true
  emits('disable-input')
  httpPost(url, {
    index: index,
    src: "chat",
    message_id: data.value?.["message_id"],
    message_hash: data.value?.["image"]?.hash,
    session_id: getSessionId(),
    prompt: data.value?.["prompt"],
    chat_id: props.chatId,
    role_id: props.roleId,
    icon: props.icon,
  }).then(() => {
    ElMessage.success("任务推送成功，请耐心等待任务执行...")
    loading.value = false
  }).catch(e => {
    ElMessage.error("任务推送失败：" + e.message)
    emits('disable-input')
  })
}
</script>

<style lang="stylus">
.chat-line-mj {
  background-color #ffffff;
  justify-content: center;
  width 100%
  padding-bottom: 1.5rem;
  padding-top: 1.5rem;
  border-bottom: 1px solid #d9d9e3;

  .chat-line-inner {
    display flex;
    width 100%;
    max-width 900px;
    padding-left 10px;

    .chat-icon {
      margin-right 20px;

      img {
        width: 36px;
        height: 36px;
        border-radius: 10px;
        padding: 1px;
      }
    }

    .chat-item {
      position: relative;
      padding: 0 5px 0 0;
      overflow: hidden;

      .content {
        word-break break-word;
        padding: 6px 10px;
        color #374151;
        font-size: var(--content-font-size);
        border-radius: 5px;
        overflow: auto;

        .text {
          p:first-child {
            margin-top 0
          }
        }

        .images {
          max-width 350px;

          .el-image {
            border-radius 10px;

            .image-slot {
              color #c1c1c1
              width 350px
              text-align center
              border-radius 10px;
              border 1px solid #e1e1e1
            }
          }
        }
      }

      .opt {
        .opt-line {
          margin 6px 0

          ul {
            display flex
            flex-flow row
            padding-left 10px

            li {
              margin-right 10px

              a {
                padding 6px 0
                width 64px
                text-align center
                border-radius 5px
                display block
                cursor pointer
                background-color #4E5058
                color #ffffff

                &:hover {
                  background-color #6D6F78
                }
              }
            }
          }
        }
      }

      .bar {
        padding 10px;

        .bar-item {
          background-color #f7f7f8;
          color #888
          padding 3px 5px;
          margin-right 10px;
          border-radius 5px;

          .el-icon {
            position relative
            top 2px;
          }
        }
      }
    }
  }


}
</style>