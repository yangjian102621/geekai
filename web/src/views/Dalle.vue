<template>
  <div>
    <div class="page-dall">
      <div class="inner custom-scroll">
        <div class="sd-box">
          <h2>DALL-E 创作中心</h2>

          <div class="sd-params">
            <el-form :model="params" label-width="80px" label-position="left">
              <div class="param-line" style="padding-top: 10px">
                <el-form-item label="图片质量">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select v-model="params.quality" style="width:176px">
                        <el-option v-for="v in qualities" :label="v.name" :value="v.value" :key="v.value"/>
                      </el-select>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="图片尺寸">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select v-model="params.size" style="width:176px">
                        <el-option v-for="v in sizes" :label="v" :value="v" :key="v"/>
                      </el-select>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="图片样式">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select v-model="params.style" style="width:176px">
                        <el-option v-for="v in styles" :label="v.name" :value="v.value" :key="v.value"/>
                      </el-select>
                      <el-tooltip
                          effect="light"
                          content="生动使模型倾向于生成超真实和戏剧性的图像"
                          raw-content
                          placement="right"
                      >
                        <el-icon class="info-icon">
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-input
                    v-model="params.prompt"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    type="textarea"
                    ref="promptRef"
                    placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"
                />
              </div>

              <div class="text-info">
                <el-row :gutter="10">
                  <el-col :span="12">
                    <el-tag>每次绘图消耗{{ dallPower }}算力</el-tag>
                  </el-col>
                  <el-col :span="12">
                    <el-tag type="success">当前可用{{ power }}算力</el-tag>
                  </el-col>
                </el-row>
              </div>

            </el-form>
          </div>
          <div class="submit-btn">
            <el-button color="#47fff1" :dark="false" round @click="generate">
              立即生成
            </el-button>
          </div>
        </div>
        <div class="task-list-box">
          <div class="task-list-inner" :style="{ height: listBoxHeight + 'px' }">
            <div class="job-list-box">
              <h2>任务列表</h2>
              <task-list :list="runningJobs" />

              <h2>创作记录</h2>
              <div class="finish-job-list">
                <div v-if="finishedJobs.length > 0">
                  <v3-waterfall
                      id="waterfall"
                      :list="finishedJobs"
                      srcKey="img_thumb"
                      :gap="20"
                      :bottomGap="-10"
                      :colWidth="colWidth"
                      :distanceToScroll="100"
                      :isLoading="loading"
                      :isOver="isOver"
                      @scrollReachBottom="fetchFinishJobs()">
                    <template #default="slotProp">
                      <div class="job-item">
                        <el-image
                            v-if="slotProp.item.img_url !== ''"
                            @click="previewImg(slotProp.item)"
                            :src="slotProp.item['img_thumb']"
                            fit="cover"
                            loading="lazy">
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

                        <el-image v-else-if="slotProp.item.progress === 101">
                          <template #error>
                            <div class="image-slot">
                              <div class="err-msg-container">
                                <div class="title">任务失败</div>
                                <div class="opt">
                                  <el-popover title="错误详情" trigger="click" :width="250" :content="slotProp.item['err_msg']" placement="top">
                                    <template #reference>
                                      <el-button type="info">详情</el-button>
                                    </template>
                                  </el-popover>
                                  <el-button type="danger"  @click="removeImage(slotProp.item)">删除</el-button>
                                </div>
                              </div>
                            </div>
                          </template>
                        </el-image>

                        <el-image v-else>
                          <template #error>
                            <div class="image-slot">
                              <i class="iconfont icon-loading"></i>
                              <span>正在下载图片</span>
                            </div>
                          </template>
                        </el-image>

                        <div class="remove">
                          <el-tooltip content="删除" placement="top" effect="light">
                            <el-button type="danger" :icon="Delete" @click="removeImage(slotProp.item)" circle/>
                          </el-tooltip>
                          <el-tooltip content="分享" placement="top" effect="light" v-if="slotProp.item.publish">
                            <el-button type="warning"
                                       @click="publishImage(slotProp.item, false)"
                                       circle>
                              <i class="iconfont icon-cancel-share"></i>
                            </el-button>
                          </el-tooltip>
                          <el-tooltip content="取消分享" placement="top" effect="light" v-else>
                            <el-button type="success" @click="publishImage(slotProp.item, true)" circle>
                              <i class="iconfont icon-share-bold"></i>
                            </el-button>
                          </el-tooltip>

                          <el-tooltip content="复制提示词" placement="top" effect="light">
                            <el-button type="info" circle class="copy-prompt"
                                       :data-clipboard-text="slotProp.item.prompt">
                              <i class="iconfont icon-file"></i>
                            </el-button>
                          </el-tooltip>
                        </div>
                      </div>
                    </template>

                    <template #footer>
                      <div class="no-more-data">
                        <span>没有更多数据了</span>
                        <i class="iconfont icon-face"></i>
                      </div>
                    </template>
                  </v3-waterfall>

                </div>
                <el-empty :image-size="100" v-else/>
              </div> <!-- end finish job list-->
            </div>
          </div>
          <back-top :right="30" :bottom="30" bg-color="#0f7a71"/>
        </div><!-- end task list box -->
      </div>

    </div>

    <el-image-viewer @close="() => { previewURL = '' }" v-if="previewURL !== ''" :url-list="[previewURL]"/>
  </div>
</template>

<script setup>
import {nextTick, onMounted, onUnmounted, ref} from "vue"
import {Delete, InfoFilled, Picture} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElMessageBox, ElNotification} from "element-plus";
import Clipboard from "clipboard";
import {checkSession, getSystemInfo} from "@/store/cache";
import {useSharedStore} from "@/store/sharedata";
import TaskList from "@/components/TaskList.vue";
import BackTop from "@/components/BackTop.vue";

const listBoxHeight = ref(0)
// const paramBoxHeight = ref(0)
const isLogin = ref(false)
const loading = ref(true)
const colWidth = ref(220)
const isOver = ref(false)
const previewURL = ref("")
const store = useSharedStore();

const resizeElement = function () {
  listBoxHeight.value = window.innerHeight - 90
  // paramBoxHeight.value = window.innerHeight - 110
};
resizeElement()
window.onresize = () => {
  resizeElement()
}
const qualities = [
  {name: "标准", value: "standard"},
  {name: "高清", value: "hd"},
]
const sizes = ["1024x1024", "1792x1024", "1024x1792"]
const styles = [
  {name: "生动", value: "vivid"},
  {name: "自然", value: "natural"}
]
const params = ref({
  quality: "standard",
  size: "1024x1024",
  style: "vivid",
  prompt: ""
})

const finishedJobs = ref([])
const runningJobs = ref([])
const power = ref(0)
const dallPower = ref(0) // 画一张 SD 图片消耗算力
const clipboard = ref(null)
const userId = ref(0)
onMounted(() => {
  initData()
  clipboard.value = new Clipboard('.copy-prompt');
  clipboard.value.on('success', () => {
    ElMessage.success("复制成功！");
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！');
  })

  getSystemInfo().then(res => {
    dallPower.value = res.data["dall_power"]
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })
})

onUnmounted(() => {
  clipboard.value.destroy()
  if (socket.value !== null) {
    socket.value.close()
    socket.value = null
  }
})

const initData = () => {
  checkSession().then(user => {
    power.value = user['power']
    userId.value = user.id
    isLogin.value = true

    page.value = 0
    fetchRunningJobs()
    fetchFinishJobs()
    connect()
  }).catch(() => {
  });
}

const socket = ref(null)
const heartbeatHandle = ref(null)
const connect = () => {
  let host = process.env.VUE_APP_WS_HOST
  if (host === '') {
    if (location.protocol === 'https:') {
      host = 'wss://' + location.host;
    } else {
      host = 'ws://' + location.host;
    }
  }

  const _socket = new WebSocket(host + `/api/dall/client?user_id=${userId.value}`);
  _socket.addEventListener('open', () => {
    socket.value = _socket;
  });

  _socket.addEventListener('message', event => {
    if (event.data instanceof Blob) {
      const reader = new FileReader();
      reader.readAsText(event.data, "UTF-8")
      reader.onload = () => {
        const message = String(reader.result)
        if (message === "FINISH" || message === "FAIL") {
          page.value = 0
          isOver.value = false
          fetchFinishJobs(page.value)
        }
        nextTick(() => fetchRunningJobs())
      }
    }
  });

  _socket.addEventListener('close', () => {
    if (socket.value !== null) {
      connect()
    }
  })
}

const fetchRunningJobs = () => {
  if (!isLogin.value) {
    return
  }
  // 获取运行中的任务
  httpGet(`/api/dall/jobs?finish=false`).then(res => {
    runningJobs.value = res.data
  }).catch(e => {
    ElMessage.error("获取任务失败：" + e.message)
  })
}

const page = ref(1)
const pageSize = ref(15)
// 获取已完成的任务
const fetchFinishJobs = () => {
  if (!isLogin.value) {
    return
  }

  loading.value = true
  page.value = page.value + 1

  httpGet(`/api/dall/jobs?finish=true&page=${page.value}&page_size=${pageSize.value}`).then(res => {
    if (res.data.length < pageSize.value) {
      isOver.value = true
    }
    const imageList = res.data
    for (let i = 0; i < imageList.length; i++) {
      imageList[i]["img_thumb"] = imageList[i]["img_url"] + "?imageView2/4/w/300/h/0/q/75"
    }
    if (page.value === 1) {
      finishedJobs.value = imageList
    } else {
      finishedJobs.value = finishedJobs.value.concat(imageList)
    }

    loading.value = false
  }).catch(e => {
    ElMessage.error("获取任务失败：" + e.message)
  })
}


// 创建绘图任务
const promptRef = ref(null)
const generate = () => {
  if (params.value.prompt === '') {
    promptRef.value.focus()
    return ElMessage.error("请输入绘画提示词！")
  }

  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }
  httpPost("/api/dall/image", params.value).then(() => {
    ElMessage.success("任务执行成功！")
    power.value -= dallPower.value
  }).catch(e => {
    ElMessage.error("任务执行失败：" + e.message)
  })
}

const removeImage = (item) => {
  ElMessageBox.confirm(
      '此操作将会删除任务和图片，继续操作码?',
      '删除提示',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
  ).then(() => {
    httpGet("/api/dall/remove", {id: item.id}).then(() => {
      ElMessage.success("任务删除成功")
      page.value = 0
      isOver.value = false
      fetchFinishJobs()
    }).catch(e => {
      ElMessage.error("任务删除失败：" + e.message)
    })
  }).catch(() => {
  })
}

const previewImg = (item) => {
  previewURL.value = item.img_url
}

// 发布图片到作品墙
const publishImage = (item, action) => {
  let text = "图片发布"
  if (action === false) {
    text = "取消发布"
  }
  httpGet("/api/dall/publish", {id: item.id, action: action}).then(() => {
    ElMessage.success(text + "成功")
    item.publish = action
    page.value = 0
    isOver.value = false
  }).catch(e => {
    ElMessage.error(text + "失败：" + e.message)
  })
}

</script>

<style lang="stylus">
@import "@/assets/css/image-dall.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
