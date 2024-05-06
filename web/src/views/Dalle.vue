<template>
  <div>
    <div class="page-dall">
      <div class="inner custom-scroll">
        <div class="sd-box">
          <h2>DALL-E 创作中心</h2>

          <div class="sd-params" :style="{ height: paramBoxHeight + 'px' }">
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
        <div class="task-list-box" @scrollend="handleScrollEnd">
          <div class="task-list-inner" :style="{ height: listBoxHeight + 'px' }">
            <div class="job-list-box">
              <h2>任务列表</h2>
              <div class="running-job-list">
                <ItemList :items="runningJobs" v-if="runningJobs.length > 0" :width="240">
                  <template #default>
                    <div class="job-item">
                      <el-image fit="cover">
                        <template #error>
                          <div class="image-slot">
                            <i class="iconfont icon-quick-start"></i>
                            <span>任务正在排队中</span>
                          </div>
                        </template>
                      </el-image>
                    </div>
                  </template>
                </ItemList>
                <el-empty :image-size="100" v-else/>
              </div>

              <h2>创作记录</h2>
              <div class="finish-job-list">
                <div v-if="finishedJobs.length > 0">
                  <ItemList :items="finishedJobs" :width="240" :gap="16">
                    <template #default="scope">
                      <div class="job-item">
                        <el-image v-if="scope.item['img_url']"
                                  :src="scope.item['img_url']+'?imageView2/1/w/240/h/240/q/75'"
                                  fit="cover"
                                  :preview-src-list="[scope.item['img_url']]"
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

                        <el-image v-else
                                  :src="scope.item['org_url']"
                                  fit="cover"
                                  :preview-src-list="[scope.item['org_url']]"
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

                        <div class="remove">
                          <el-tooltip content="删除" placement="top" effect="light">
                            <el-button type="danger" :icon="Delete" @click="removeImage($event,scope.item)" circle/>
                          </el-tooltip>
                          <el-tooltip content="分享" placement="top" effect="light" v-if="scope.item.publish">
                            <el-button type="warning"
                                       @click="publishImage($event,scope.item, false)"
                                       circle>
                              <i class="iconfont icon-cancel-share"></i>
                            </el-button>
                          </el-tooltip>
                          <el-tooltip content="取消分享" placement="top" effect="light" v-else>
                            <el-button type="success" @click="publishImage($event,scope.item, true)" circle>
                              <i class="iconfont icon-share-bold"></i>
                            </el-button>
                          </el-tooltip>

                          <el-tooltip content="复制提示词" placement="top" effect="light">
                            <el-button type="info" circle class="copy-prompt" :data-clipboard-text="scope.item.prompt">
                              <i class="iconfont icon-file"></i>
                            </el-button>
                          </el-tooltip>
                        </div>
                      </div>
                    </template>
                  </ItemList>

                  <div class="no-more-data" v-if="isOver">
                    <span>没有更多数据了</span>
                    <i class="iconfont icon-face"></i>
                  </div>
                </div>
                <el-empty :image-size="100" v-else/>
              </div> <!-- end finish job list-->
            </div>
          </div>

        </div><!-- end task list box -->
      </div>

    </div>

    <login-dialog :show="showLoginDialog" @hide="showLoginDialog =  false" @success="initData"/>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from "vue"
import {Delete, InfoFilled, Picture} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElMessageBox, ElNotification} from "element-plus";
import ItemList from "@/components/ItemList.vue";
import Clipboard from "clipboard";
import {checkSession} from "@/action/session";
import LoginDialog from "@/components/LoginDialog.vue";

const listBoxHeight = ref(window.innerHeight - 40)
const paramBoxHeight = ref(window.innerHeight - 150)
const showLoginDialog = ref(false)
const isLogin = ref(false)

window.onresize = () => {
  listBoxHeight.value = window.innerHeight - 40
  paramBoxHeight.value = window.innerHeight - 150
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

  httpGet("/api/config/get?key=system").then(res => {
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
    fetchRunningJobs()
    fetchFinishJobs(1)
    connect()
  }).catch(() => {
  });
}

const handleScrollEnd = () => {
  if (isOver.value === true) {
    return
  }
  page.value += 1
  fetchFinishJobs(page.value)
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

  // 心跳函数
  const sendHeartbeat = () => {
    clearTimeout(heartbeatHandle.value)
    new Promise((resolve, reject) => {
      if (socket.value !== null) {
        socket.value.send(JSON.stringify({type: "heartbeat", content: "ping"}))
      }
      resolve("success")
    }).then(() => {
      heartbeatHandle.value = setTimeout(() => sendHeartbeat(), 5000)
    });
  }

  const _socket = new WebSocket(host + `/api/dall/client?user_id=${userId.value}`);
  _socket.addEventListener('open', () => {
    socket.value = _socket;

    // 发送心跳消息
    sendHeartbeat()
  });

  _socket.addEventListener('message', event => {
    if (event.data instanceof Blob) {
      const reader = new FileReader();
      reader.readAsText(event.data, "UTF-8")
      reader.onload = () => {
        const message = String(reader.result)
        if (message === "FINISH") {
          page.value = 1
          fetchFinishJobs(page.value)
          isOver.value = false
        }
        fetchRunningJobs()
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
  httpGet(`/api/dall/jobs?status=0`).then(res => {
    const jobs = res.data
    const _jobs = []
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i].progress === -1) {
        ElNotification({
          title: '任务执行失败',
          dangerouslyUseHTMLString: true,
          message: `任务ID：${jobs[i]['task_id']}<br />原因：${jobs[i]['err_msg']}`,
          type: 'error',
        })
        power.value += dallPower.value
        continue
      }
      _jobs.push(jobs[i])
    }
    runningJobs.value = _jobs
  }).catch(e => {
    ElMessage.error("获取任务失败：" + e.message)
  })
}

const page = ref(1)
const pageSize = ref(15)
const isOver = ref(false)
// 获取已完成的任务
const fetchFinishJobs = (page) => {
  if (!isLogin.value) {
    return
  }
  httpGet(`/api/dall/jobs?status=1&page=${page}&page_size=${pageSize.value}`).then(res => {
    if (res.data.length < pageSize.value) {
      isOver.value = true
    }
    if (page === 1) {
      finishedJobs.value = res.data
    } else {
      finishedJobs.value = finishedJobs.value.concat(res.data)
    }
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
    showLoginDialog.value = true
    return
  }
  httpPost("/api/dall/image", params.value).then(() => {
    ElMessage.success("任务执行成功！")
    power.value -= dallPower.value
  }).catch(e => {
    ElMessage.error("任务执行失败：" + e.message)
  })
}

const removeImage = (event, item) => {
  event.stopPropagation()
  ElMessageBox.confirm(
      '此操作将会删除任务和图片，继续操作码?',
      '删除提示',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
  ).then(() => {
    httpPost("/api/dall/remove", {id: item.id, img_url: item.img_url, user_id: userId.value}).then(() => {
      ElMessage.success("任务删除成功")
      fetchFinishJobs(1)
    }).catch(e => {
      ElMessage.error("任务删除失败：" + e.message)
    })
  }).catch(() => {
  })
}

// 发布图片到作品墙
const publishImage = (event, item, action) => {
  event.stopPropagation()
  let text = "图片发布"
  if (action === false) {
    text = "取消发布"
  }
  httpPost("/api/dall/publish", {id: item.id, action: action}).then(() => {
    ElMessage.success(text + "成功")
    item.publish = action
  }).catch(e => {
    ElMessage.error(text + "失败：" + e.message)
  })
}

</script>

<style lang="stylus">
@import "@/assets/css/image-dall.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
