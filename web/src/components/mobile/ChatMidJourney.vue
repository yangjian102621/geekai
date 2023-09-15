<template>
  <div class="mobile-message-mj">
    <div class="chat-icon">
      <van-image :src="icon"/>
    </div>

    <div class="chat-item">
      <div class="triangle"></div>
      <div class="content-box">
        <div class="content">
          <div class="content-inner">
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
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import {ref, watch} from "vue";
import {Picture} from "@element-plus/icons-vue";
import {httpPost} from "@/utils/http";
import {getSessionId} from "@/store/session";
import {showNotify} from "vant";

const props = defineProps({
  content: Object,
  icon: String,
  chatId: String,
  roleId: Number,
  createdAt: String
});

const data = ref(props.content)
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
    key: data.value?.["key"],
    prompt: data.value?.["prompt"],
    chat_id: props.chatId,
    role_id: props.roleId,
    icon: props.icon,
  }).then(() => {
    showNotify({type: "success", message: "任务推送成功，请耐心等待任务执行..."})
    loading.value = false
  }).catch(e => {
    showNotify({type: "danger", message: "任务推送失败：" + e.message})
    emits('disable-input')
  })
}
</script>

<style lang="stylus">
.mobile-message-mj {
  display flex
  justify-content: flex-start;

  .chat-icon {
    margin-right 5px

    .van-image {
      width 32px

      img {
        border-radius 5px
      }
    }
  }

  .chat-item {
    display: inline-block;
    position: relative;
    padding: 0 0 0 5px;
    overflow: hidden;

    .triangle {
      width: 0;
      height: 0;
      border-top: 5px solid transparent;
      border-bottom: 5px solid transparent;
      border-right: 5px solid #fff;
      position: absolute;
      left: 0;
      top: 13px;
    }

    .content-box {

      display flex
      flex-direction row

      .content {
        text-align left
        width 100%
        overflow-x auto
        min-height 20px;
        word-break break-word;
        padding: 5px 10px;
        color #444444
        background-color: #ffffff;
        font-size: 16px
        border-radius: 5px;

        .content-inner {
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
                  padding 3px 0
                  width 50px
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

      }
    }

  }
}


.van-theme-dark {
  .mobile-message-reply {
    .chat-item {
      .triangle {
        border-right: 5px solid #404042;
      }

      .content-box {
        .content {
          color #c1c1c1
          background-color: #404042;

          p > code {
            color #c1c1c1
            background-color #2b2b2b
          }
        }
      }

    }
  }

}
</style>