<template>
  <div class="mobile-sd">
    <van-form @submit="generate">
      <van-cell-group inset>
        <div>
          <van-field
              v-model="quality"
              is-link
              label="图片质量"
              placeholder="选择图片质量"
              @click="showQualityPicker = true"
          />
          <van-popup v-model:show="showQualityPicker" position="bottom" teleport="#app">
            <van-picker
                :columns="qualities"
                @cancel="showQualityPicker = false"
                @confirm="qualityConfirm"
            />
          </van-popup>
        </div>

        <div>
          <van-field
              v-model="size"
              is-link
              label="图片尺寸"
              placeholder="选择图片尺寸"
              @click="showSizePicker = true"
          />
          <van-popup v-model:show="showSizePicker" position="bottom" teleport="#app">
            <van-picker
                :columns="sizes"
                @cancel="showSizePicker = false"
                @confirm="sizeConfirm"
            />
          </van-popup>
        </div>

        <div>
          <van-field
              v-model="style"
              is-link
              label="图片样式"
              placeholder="选择图片样式"
              @click="showStylePicker = true"
          />
          <van-popup v-model:show="showStylePicker" position="bottom" teleport="#app">
            <van-picker
                :columns="styles"
                @cancel="showStylePicker = false"
                @confirm="styleConfirm"
            />
          </van-popup>
        </div>

        <van-field
            v-model="params.prompt"
            rows="3"
            autosize
            type="textarea"
            placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"
        />

        <div class="text-line pt-6">
          <el-tag>绘图消耗{{ dallPower }}算力，当前算力：{{ power }}</el-tag>
        </div>

        <div class="text-line">
          <van-button round block type="primary" native-type="submit">
            立即生成
          </van-button>
        </div>
      </van-cell-group>
    </van-form>

    <h3>任务列表</h3>
    <div class="running-job-list">
      <van-empty v-if="runningJobs.length ===0"
                 image="https://fastly.jsdelivr.net/npm/@vant/assets/custom-empty-image.png"
                 image-size="80"
                 description="暂无记录"
      />
      <van-grid :gutter="10" :column-num="3" v-else>
        <van-grid-item v-for="item in runningJobs" :key="item.id">
          <div v-if="item.progress > 0">
            <van-image src="/images/img-holder.png"></van-image>
            <div class="progress">
              <van-circle
                  v-model:current-rate="item.progress"
                  :rate="item.progress"
                  :speed="100"
                  :text="item.progress+'%'"
                  :stroke-width="60"
                  size="90px"
              />
            </div>
          </div>

          <div v-else class="task-in-queue">
            <span class="icon"><i class="iconfont icon-quick-start"></i></span>
            <span class="text">排队中</span>
          </div>

        </van-grid-item>
      </van-grid>
    </div>

    <h3>创作记录</h3>
    <div class="finish-job-list">
      <van-empty v-if="finishedJobs.length ===0"
                 image="https://fastly.jsdelivr.net/npm/@vant/assets/custom-empty-image.png"
                 image-size="80"
                 description="暂无记录"
      />

      <van-list v-else
                v-model:error="error"
                v-model:loading="loading"
                :finished="finished"
                error-text="请求失败，点击重新加载"
                finished-text="没有更多了"
                @load="onLoad"
      >
        <van-grid :gutter="10" :column-num="2">
          <van-grid-item v-for="item in finishedJobs" :key="item.id">
            <div class="failed" v-if="item.progress === 101">
              <div class="title">任务失败</div>
              <div class="opt">
                <van-button size="small" @click="showErrMsg(item)">详情</van-button>
                <van-button type="danger" @click="removeImage($event,item)" size="small">删除</van-button>
              </div>
            </div>
            <div class="job-item" v-else>
              <van-image
                  :src="item['img_url']"
                  :class="item['can_opt'] ? '' : 'upscale'"
                  lazy-load
                  @click="imageView(item)"
                  fit="cover">
                <template v-slot:loading>
                  <van-loading type="spinner" size="20"/>
                </template>
              </van-image>

              <div class="remove">
                <el-button type="danger" :icon="Delete" @click="removeImage($event, item)" circle/>
                <el-button type="warning" v-if="item.publish" @click="publishImage($event,item, false)"
                           circle>
                  <i class="iconfont icon-cancel-share"></i>
                </el-button>
                <el-button type="success" v-else @click="publishImage($event, item, true)" circle>
                  <i class="iconfont icon-share-bold"></i>
                </el-button>
                <el-button type="primary" @click="showPrompt(item)" circle>
                  <i class="iconfont icon-prompt"></i>
                </el-button>
              </div>
            </div>
          </van-grid-item>
        </van-grid>
      </van-list>

    </div>
    <button style="display: none" class="copy-prompt-dall" :data-clipboard-text="prompt" id="copy-btn-dall">复制</button>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from "vue"
import {Delete} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import Clipboard from "clipboard";
import {checkSession, getClientId, getSystemInfo} from "@/store/cache";
import {useRouter} from "vue-router";
import {getSessionId} from "@/store/session";
import {
  showConfirmDialog,
  showDialog,
  showFailToast,
  showImagePreview,
  showNotify,
  showSuccessToast,
  showToast
} from "vant";
import {showLoginDialog} from "@/utils/libs";
import {useSharedStore} from "@/store/sharedata";

const listBoxHeight = ref(window.innerHeight - 40)
const mjBoxHeight = ref(window.innerHeight - 150)
const isLogin = ref(false)

window.onresize = () => {
  listBoxHeight.value = window.innerHeight - 40
  mjBoxHeight.value = window.innerHeight - 150
}

const qualities = [
  {text: "标准", value: "standard"},
  {text: "高清", value: "hd"},
]
const sizes = [
  {text:"1024x1024",value:"1024x1024"},
  {text:"1792x1024",value:"1792x1024"},
  {text: "1024x1792",value:"1024x1792"},
]
const styles = [
  {text: "生动", value: "vivid"},
  {text: "自然", value: "natural"}
]
const params = ref({
  client_id: getClientId(),
  quality: qualities[0].value,
  size: sizes[0].value,
  style: styles[0].value,
  prompt: ""
})
const quality = ref(qualities[0].text)
const size = ref(sizes[0].text)
const style = ref(styles[0].text)

const showQualityPicker = ref(false)
const showStylePicker = ref(false)
const showSizePicker = ref(false)


const runningJobs = ref([])
const finishedJobs = ref([])
const router = useRouter()
const power = ref(0)
const dallPower = ref(0) // 画一张 DALL 图片消耗算力

const userId = ref(0)
const store = useSharedStore()
const clipboard = ref(null)
const prompt = ref('')
onMounted(() => {
  initData()
  clipboard.value = new Clipboard(".copy-prompt-dall");
  clipboard.value.on('success', () => {
    showNotify({type: 'success', message: '复制成功', duration: 1000})
  })
  clipboard.value.on('error', () => {
    showNotify({type: 'danger', message: '复制失败', duration: 2000})
  })

  getSystemInfo().then(res => {
    dallPower.value = res.data.dall_power
  }).catch(e => {
    showNotify({type: "danger", message: "获取系统配置失败：" + e.message})
  })

  store.addMessageHandler("dall", (data) => {
    if (data.channel !== "dall" || data.clientId !== getClientId()) {
      return
    }
    if (data.body === "FINISH" || data.body === "FAIL") {
      page.value = 1
      fetchFinishJobs(1)
    }
    fetchRunningJobs()
  })

})

onUnmounted(() => {
  clipboard.value.destroy()
  store.removeMessageHandler("dall")
})


const initData = () => {
  checkSession().then(user => {
    power.value = user['power']
    isLogin.value = true
    fetchRunningJobs()
    fetchFinishJobs(1)
  }).catch(() => {
    loading.value = false
  });
}

const fetchRunningJobs = () => {
  // 获取运行中的任务
  httpGet(`/api/dall/jobs?finish=0`).then(res => {
    runningJobs.value = res.data.items
  }).catch(e => {
    showNotify({type: "danger", message: "获取任务失败：" + e.message})
  })
}

const loading = ref(false)
const finished = ref(false)
const error = ref(false)
const page = ref(0)
const pageSize = ref(10)
// 获取已完成的任务
const fetchFinishJobs = (page) => {
  loading.value = true
  httpGet(`/api/dall/jobs?finish=1&page=${page}&page_size=${pageSize.value}`).then(res => {
    const jobs = res.data.items
    if (jobs.length < pageSize.value) {
      finished.value = true
    }
    const _jobs = []
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i].progress === -1) {
        jobs[i]['thumb_url'] = jobs[i]['img_url'] + '?imageView2/1/w/480/h/600/q/75'
      }
      _jobs.push(jobs[i])
    }
    if (page === 1) {
      finishedJobs.value = _jobs
    } else {
      finishedJobs.value = finishedJobs.value.concat(_jobs)
    }
    loading.value = false
  }).catch(e => {
    loading.value = false
    showNotify({type: "danger", message: "获取任务失败：" + e.message})
  })
}

const onLoad = () => {
  page.value += 1
  fetchFinishJobs(page.value)
}

// 创建绘图任务
const promptRef = ref(null)
const generate = () => {
  if (!isLogin.value) {
    return showLoginDialog(router)
  }

  if (params.value.prompt === '') {
    promptRef.value.focus()
    return showToast("请输入绘画提示词！")
  }

  if (!params.value.seed) {
    params.value.seed = -1
  }
  params.value.session_id = getSessionId()
  httpPost("/api/dall/image", params.value).then(() => {
    showSuccessToast("绘画任务推送成功，请耐心等待任务执行...")
    power.value -= dallPower.value
    fetchRunningJobs()
  }).catch(e => {
    showFailToast("任务推送失败：" + e.message)
  })
}

const showPrompt = (item) => {
  prompt.value = item.prompt
  showConfirmDialog({
    title: "绘画提示词",
    message: item.prompt,
    confirmButtonText: "复制",
    cancelButtonText: "关闭",
  }).then(() => {
    document.querySelector('#copy-btn-dall').click()
  }).catch(() => {
  });
}

const showErrMsg = (item) => {
  showDialog({
    title: '错误详情',
    message: item['err_msg'],
  }).then(() => {
    // on close
  });
}

const removeImage = (event, item) => {
  event.stopPropagation()
  showConfirmDialog({
    title: '标题',
    message:
        '此操作将会删除任务和图片，继续操作码?',
  }).then(() => {
    httpGet("/api/dall/remove", {id: item.id, user_id: item.user_id}).then(() => {
      showSuccessToast("任务删除成功")
      fetchFinishJobs(1)
    }).catch(e => {
      showFailToast("任务删除失败：" + e.message)
    })
  }).catch(() => {
    showToast("您取消了操作")
  });
}

// 发布图片到作品墙
const publishImage = (event, item, action) => {
  event.stopPropagation()
  let text = "图片发布"
  if (action === false) {
    text = "取消发布"
  }
  httpGet("/api/dall/publish", {id: item.id, action: action, user_id: item.user_id}).then(() => {
    showSuccessToast(text + "成功")
    item.publish = action
  }).catch(e => {
    showFailToast(text + "失败：" + e.message)
  })
}

const imageView = (item) => {
  showImagePreview([item['img_url']]);
}


const qualityConfirm = (item) => {
  params.value.quality = item.selectedOptions[0].value;
  quality.value = item.selectedOptions[0].text
  showQualityPicker.value = false
}

const styleConfirm = (item) => {
  params.value.style = item.selectedOptions[0].value;
  style.value = item.selectedOptions[0].text
  showStylePicker.value = false
}

const sizeConfirm =(item) => {
  params.value.size = item.selectedOptions[0].value
  size.value=item.selectedOptions[0].text
  showSizePicker.value =false
}

</script>

<style lang="stylus">
@import "@/assets/css/mobile/image-sd.styl"
</style>