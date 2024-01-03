<template>
  <div class="page-mj">
    <div class="inner custom-scroll">
      <div class="mj-box">
        <h2>MidJourney 创作中心</h2>

        <div class="mj-params" :style="{ height: mjBoxHeight + 'px' }">
          <el-form :model="params" label-width="80px" label-position="left">
            <div class="param-line pt">
              <span>图片比例：</span>
              <el-tooltip effect="light" content="生成图片的尺寸比例" placement="right">
                <el-icon>
                  <InfoFilled/>
                </el-icon>
              </el-tooltip>
            </div>

            <div class="param-line pt">
              <el-row :gutter="10">
                <el-col :span="8" v-for="item in rates" :key="item.value">
                  <div class="flex-col items-center"
                       :class="item.value === params.rate ? 'grid-content active' : 'grid-content'"
                       @click="changeRate(item)">
                    <!--                    <div :class="'shape ' + item.css"></div>-->
                    <el-image class="icon" :src="item.img" fit="cover"></el-image>
                    <div class="text">{{ item.text }}</div>
                  </div>
                </el-col>
              </el-row>
            </div>

            <div class="param-line" style="padding-top: 10px">
              <el-form-item label="图片画质">
                <template #default>
                  <div class="form-item-inner flex-row items-center">
                    <el-select v-model="params.quality" placeholder="请选择">
                      <el-option v-for="item in options"
                                 :key="item.value"
                                 :label="item.label"
                                 :value="item.value">
                      </el-option>
                    </el-select>
                    <el-tooltip effect="light" content="生成的图片质量，质量越好出图越慢" placement="right">
                      <el-icon>
                        <InfoFilled/>
                      </el-icon>
                    </el-tooltip>
                  </div>
                </template>
              </el-form-item>
            </div>

            <div class="param-line pt">
              <span>模型选择：</span>
              <el-tooltip effect="light" content="MJ: 偏真实通用模型 <br/>NIJI: 偏动漫风格、适用于二次元模型" raw-content
                          placement="right">
                <el-icon>
                  <InfoFilled/>
                </el-icon>
              </el-tooltip>
            </div>
            <div class="param-line pt">
              <el-row :gutter="10">
                <el-col :span="12" v-for="item in models" :key="item.value">
                  <div :class="item.value === params.model ? 'model active' : 'model'"
                       @click="changeModel(item)">
                    <el-image :src="item.img" fit="cover"></el-image>
                    <div class="text">{{ item.text }}</div>
                  </div>
                </el-col>
              </el-row>
            </div>

            <div class="param-line">
              <el-form-item label="重复平铺">
                <template #default>
                  <div class="form-item-inner">
                    <el-switch v-model="params.tile" inactive-color="#464649" active-color="#47fff1"/>
                    <el-tooltip effect="light"
                                content="重复：--tile，参数释义：生成可用作重复平铺的图像，以创建无缝图案。" raw-content
                                placement="right">
                      <el-icon>
                        <InfoFilled/>
                      </el-icon>
                    </el-tooltip>
                  </div>
                </template>
              </el-form-item>
            </div>


            <div class="param-line">
              <el-form-item label="原始模式">
                <template #default>
                  <div class="form-item-inner">
                    <el-switch v-model="params.raw" inactive-color="#464649" active-color="#47fff1"/>
                    <el-tooltip effect="light"
                                content="启用新的RAW模式，呈现的人物写实感更加逼真，人物细节、光源、流畅度也更加接近原始作品。<br/> 同时也意味着您需要添加更长的提示。"
                                raw-content
                                placement="right">
                      <el-icon>
                        <InfoFilled/>
                      </el-icon>
                    </el-tooltip>
                  </div>
                </template>
              </el-form-item>
            </div>

            <div class="param-line" style="padding-top: 10px">
              <el-form-item label="创意度">
                <template #default>
                  <div class="form-item-inner">
                    <el-slider v-model.number="params.chaos" :max="100" :step="1"
                               style="width: 180px;--el-slider-main-bg-color:#47fff1"/>
                    <el-tooltip effect="light"
                                content="参数用法：--chaos 或--c，取值范围: 0-100 <br/> 取值越高结果越发散，反之则稳定收敛<br /> 默认值0最为精准稳定"
                                raw-content placement="right">
                      <el-icon>
                        <InfoFilled/>
                      </el-icon>
                    </el-tooltip>
                  </div>
                </template>
              </el-form-item>
            </div>

            <div class="param-line">
              <el-form-item label="风格化">
                <template #default>
                  <div class="form-item-inner">
                    <el-slider v-model.number="params.stylize" :min="0" :max="1000" :step="1"
                               style="width: 180px;--el-slider-main-bg-color:#47fff1"/>
                    <el-tooltip effect="light"
                                content="风格化：--stylize 或 --s，范围 1-1000，默认值100 <br/>高取值会产生非常艺术化但与提示关联性较低的图像"
                                raw-content placement="right">
                      <el-icon>
                        <InfoFilled/>
                      </el-icon>
                    </el-tooltip>
                  </div>
                </template>
              </el-form-item>
            </div>

            <div class="param-line">
              <el-form-item label="随机种子">
                <template #default>
                  <div class="form-item-inner">
                    <el-input v-model.number="params.seed" style="--el-input-focus-border-color:#47fff1"/>
                    <el-tooltip effect="light"
                                content="随机种子：--seed，默认值0表示随机产生 <br/>使用相同的种子参数和描述将产生相似的图像"
                                raw-content
                                placement="right">
                      <el-icon>
                        <InfoFilled/>
                      </el-icon>
                    </el-tooltip>
                  </div>
                </template>
              </el-form-item>
            </div>
          </el-form>
        </div>
      </div>
      <div class="task-list-box">
        <div class="task-list-inner" :style="{ height: listBoxHeight + 'px' }">
          <h2>AI绘画</h2>
          <el-form>
            <el-tabs v-model="activeName" class="title-tabs">
              <el-tab-pane label="图生图(可选)" name="图生图">
                <div class="text">图生图：以某张图片为底稿参考来创作绘画，生成类似风格或类型图像，支持 PNG 和 JPG 格式图片；
                </div>
                <div class="param-line pt">
                  <el-form-item label="">
                    <template #default>
                      <div class="form-item-inner flex-row items-center">
                        <el-input v-model="params.img" size="small" placeholder="请输入图片地址或者上传图片"
                                  style="width: 300px;"/>
                        <el-icon @click="params.img = ''" title="清空图片">
                          <DeleteFilled/>
                        </el-icon>
                        <el-tooltip effect="light"
                                    content="垫图：以某张图片为底稿参考来创作绘画 <br/> 支持 PNG 和 JPG 格式图片"
                                    raw-content placement="right">
                          <el-icon>
                            <InfoFilled/>
                          </el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                  </el-form-item>
                </div>

                <div class="param-line">
                  <el-upload class="img-uploader" :auto-upload="true" :show-file-list="false"
                             :http-request="afterRead" style="--el-color-primary:#47fff1">
                    <el-image v-if="params.img !== ''" :src="params.img" fit="cover"/>
                    <el-icon v-else class="uploader-icon">
                      <Plus/>
                    </el-icon>
                  </el-upload>
                </div>

                <div class="param-line" style="padding-top: 10px">
                  <el-form-item label="图像权重:">
                    <template #default>
                      <div class="form-item-inner">
                        <el-slider v-model.number="params.weight" :max="1" :step="0.01"
                                   style="width: 180px;--el-slider-main-bg-color:#47fff1"/>
                        <el-tooltip effect="light"
                                    content="使用图像权重参数--iw来调整图像 URL 与文本的重要性 <br/>权重较高时意味着图像提示将对完成的作业产生更大的影响"
                                    raw-content placement="right">
                          <el-icon style="margin-top: 9px">
                            <InfoFilled/>
                          </el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                  </el-form-item>
                </div>
              </el-tab-pane>

              <el-tab-pane label="图生文(可选)" name="图生文">
                <div class="text">图生文功能正在紧锣密鼓开发中，敬请期待...</div>
                <!--              <div class="param-line pt">-->
                <!--                <el-empty image-size="100px" description="功能建设中"/>-->
                <!--              </div>-->
              </el-tab-pane>

              <el-tab-pane label="融图(可选)" name="融图">
                <div class="text">融图功能正在紧锣密鼓开发中，敬请期待...</div>
              </el-tab-pane>
            </el-tabs>

            <div v-loading="loading" element-loading-background="rgba(122, 122, 122, 0.8)">
              <div class="param-line pt">
                <div class="flex-row justify-between items-center">
                  <div class="flex-row justify-start items-center">
                    <span>提示词：</span>
                    <el-tooltip effect="light" content="输入你想要的内容，用逗号分割" placement="right">
                      <el-icon>
                        <InfoFilled/>
                      </el-icon>
                    </el-tooltip>
                  </div>
                  <div>
                    <el-button type="primary" @click="translatePrompt">
                      <el-icon style="margin-right: 6px;font-size: 18px;">
                        <Refresh/>
                      </el-icon>
                      翻译
                    </el-button>

                    <el-tooltip
                        class="box-item"
                        effect="light"
                        raw-content
                        content="使用 AI 翻译并重写提示词，<br/>增加更多细节，风格等描述"
                        placement="top-end"
                    >
                      <el-button type="success" @click="rewritePrompt">
                        <el-icon style="margin-right: 6px;font-size: 18px;">
                          <Refresh/>
                        </el-icon>
                        翻译并重写
                      </el-button>
                    </el-tooltip>
                  </div>
                </div>
              </div>

              <div class="param-line pt">
                <el-input v-model="params.prompt" :autosize="{ minRows: 4, maxRows: 6 }" type="textarea"
                          ref="promptRef"
                          placeholder="这里输入你的英文咒语，例如：A chinese girl walking in the middle of a cobblestone street"/>
              </div>
            </div>

            <div class="param-line pt">
              <div class="flex-row justify-between items-center">
                <div class="flex-row justify-start items-center">
                  <span>不希望出现的内容：（可选）</span>
                  <el-tooltip effect="light" content="不想出现在图片上的元素(例如：树，建筑)" placement="right">
                    <el-icon>
                      <InfoFilled/>
                    </el-icon>
                  </el-tooltip>
                </div>
                <!--                <el-button type="success">-->
                <!--                  <el-icon style="margin-right: 6px;font-size: 18px;">-->
                <!--                    <Refresh/>-->
                <!--                  </el-icon>-->
                <!--                  翻译-->
                <!--                </el-button>-->
              </div>
            </div>

            <div class="param-line pt">
              <el-input v-model="params.neg_prompt" :autosize="{ minRows: 4, maxRows: 6 }" type="textarea"
                        ref="promptRef"
                        placeholder="这里输入你不希望出现在图片上的内容，元素"/>
            </div>

            <div class="submit-btn">
              <el-button color="#47fff1" :dark="false" @click="generate" round>立即生成</el-button>
              <div class="text-info">
                <el-tag type="success">绘图可用额度：{{ imgCalls }}</el-tag>
              </div>
            </div>
          </el-form>

          <h2>任务列表</h2>
          <div class="running-job-list">
            <ItemList :items="runningJobs" v-if="runningJobs.length > 0">
              <template #default="scope">
                <div class="job-item">
                  <div v-if="scope.item.progress > 0" class="job-item-inner">
                    <el-image :src="scope.item['img_url']" :zoom-rate="1.2"
                              :preview-src-list="[scope.item['img_url']]" fit="cover" :initial-index="0"
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

                    <div class="progress">
                      <el-progress type="circle" :percentage="scope.item.progress" :width="100"
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
              </template>
            </ItemList>
            <el-empty :image-size="100" v-else/>
          </div>

          <h2>创作记录</h2>
          <div class="finish-job-list">
            <ItemList :items="finishedJobs" v-if="finishedJobs.length > 0" :width="240" :gap="16">
              <template #default="scope">
                <div class="job-item">
                  <el-image
                      :src="scope.item['thumb_url']"
                      :class="scope.item.type === 'upscale' ? 'upscale' : ''" :zoom-rate="1.2"
                      :preview-src-list="[scope.item['img_url']]" fit="cover" :initial-index="scope.index"
                      loading="lazy" v-if="scope.item.progress > 0">
                    <template #placeholder>
                      <div class="image-slot">
                        正在加载图片
                      </div>
                    </template>

                    <template #error>
                      <div class="image-slot" v-if="scope.item['img_url'] === ''">
                        <i class="iconfont icon-loading"></i>
                        <span>正在下载图片</span>
                      </div>
                      <div class="image-slot" v-else>
                        <el-icon>
                          <Picture/>
                        </el-icon>
                      </div>
                    </template>
                  </el-image>

                  <div class="opt" v-if="scope.item.type !== 'upscale'">
                    <div class="opt-line">
                      <ul>
                        <li><a @click="upscale(1, scope.item)">U1</a></li>
                        <li><a @click="upscale(2, scope.item)">U2</a></li>
                        <li><a @click="upscale(3, scope.item)">U3</a></li>
                        <li><a @click="upscale(4, scope.item)">U4</a></li>
                        <li class="show-prompt">

                          <el-popover placement="left" title="提示词" :width="240" trigger="hover">
                            <template #reference>
                              <el-icon>
                                <ChromeFilled/>
                              </el-icon>
                            </template>

                            <template #default>
                              <div class="mj-list-item-prompt">
                                <span>{{ scope.item.prompt }}</span>
                                <el-icon class="copy-prompt"
                                         :data-clipboard-text="scope.item.prompt">
                                  <DocumentCopy/>
                                </el-icon>
                              </div>
                            </template>
                          </el-popover>
                        </li>
                      </ul>
                    </div>

                    <div class="opt-line">
                      <ul>
                        <li><a @click="variation(1, scope.item)">V1</a></li>
                        <li><a @click="variation(2, scope.item)">V2</a></li>
                        <li><a @click="variation(3, scope.item)">V3</a></li>
                        <li><a @click="variation(4, scope.item)">V4</a></li>
                      </ul>
                    </div>
                  </div>

                  <div class="remove">
                    <el-button type="danger" :icon="Delete" @click="removeImage(scope.item)" circle/>
                  </div>
                </div>
              </template>
            </ItemList>

            <el-empty :image-size="100" v-else/>
          </div> <!-- end finish job list-->
        </div>

      </div><!-- end task list box -->
    </div>

  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {
  ChromeFilled, Delete,
  DeleteFilled,
  DocumentCopy,
  InfoFilled,
  Picture,
  Plus,
  Refresh
} from "@element-plus/icons-vue";
import Compressor from "compressorjs";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElMessageBox, ElNotification} from "element-plus";
import ItemList from "@/components/ItemList.vue";
import Clipboard from "clipboard";
import {checkSession} from "@/action/session";
import {useRouter} from "vue-router";
import {getSessionId} from "@/store/session";

const listBoxHeight = ref(window.innerHeight - 40)
const mjBoxHeight = ref(window.innerHeight - 150)
window.onresize = () => {
  listBoxHeight.value = window.innerHeight - 40
  mjBoxHeight.value = window.innerHeight - 150
}
const rates = [
  {css: "square", value: "1:1", text: "1:1", img: "/images/mj/rate_1_1.png"},
  {css: "size1-2", value: "1:2", text: "1:2", img: "/images/mj/rate_1_2.png"},
  {css: "size2-1", value: "2:1", text: "2:1", img: "/images/mj/rate_2_1.png"},
  {css: "size2-3", value: "2:3", text: "2:3", img: "/images/mj/rate_3_4.png"},
  {css: "size3-2", value: "3:2", text: "3:2", img: "/images/mj/rate_4_3.png"},
  {css: "size3-4", value: "3:4", text: "3:4", img: "/images/mj/rate_3_4.png"},
  {css: "size4-3", value: "4:3", text: "4:3", img: "/images/mj/rate_4_3.png"},
  {css: "size16-9", value: "16:9", text: "16:9", img: "/images/mj/rate_16_9.png"},
  {css: "size9-16", value: "9:16", text: "9:16", img: "/images/mj/rate_9_16.png"},
]
const models = [
  {text: "写实模式MJ-6.0", value: " --v 6", img: "/images/mj/mj-v6.png"},
  {text: "优质模式MJ-5.2", value: " --v 5.2", img: "/images/mj/mj-v5.2.png"},
  {text: "优质模式MJ-5.1", value: " --v 5.1", img: "/images/mj/mj-v5.1.jpg"},
  {text: "虚幻模式MJ-5", value: " --v 5", img: "/images/mj/mj-v5.jpg"},
  {text: "真实模式MJ-4", value: " --v 4", img: "/images/mj/mj-v4.jpg"},
  {text: "动漫风niji5 原始", value: " --niji 5", img: "/images/mj/mj-niji.png"},
  {text: "动漫风niji5 可爱", value: " --niji 5 --style cute", img: "/images/mj/nj1.jpg"},
  {text: "动漫风niji5 风景", value: " --niji 5 --style scenic", img: "/images/mj/nj2.jpg"},
  {text: "动漫风niji5 表现力", value: " --niji 5 --style expressive", img: "/images/mj/nj3.jpg"},
  {text: "动漫风niji4", value: " --niji 4", img: "/images/mj/nj4.jpg"},
]

const options = [
  {
    value: 0,
    label: '默认'
  },
  {
    value: 0.25,
    label: '普通'
  },
  {
    value: 0.5,
    label: '清晰'
  },
  {
    value: 1,
    label: '高清'
  },
]

const params = ref({
  rate: rates[0].value,
  model: models[0].value,
  chaos: 0,
  stylize: 100,
  seed: 0,
  raw: false,
  img: "",
  weight: 0.25,
  prompt: "",
  neg_prompt: "",
  tile: false,
  quality: 0
})

const activeName = ref('图生图')

const runningJobs = ref([])
const finishedJobs = ref([])
const router = useRouter()

const socket = ref(null)
const imgCalls = ref(0)
const loading = ref(false)
const userId = ref(0)

const rewritePrompt = () => {
  loading.value = true
  httpPost("/api/prompt/rewrite", {"prompt": params.value.prompt}).then(res => {
    params.value.prompt = res.data
    loading.value = false
  }).catch(e => {
    loading.value = false
    ElMessage.error("翻译失败：" + e.message)
  })
}

const translatePrompt = () => {
  loading.value = true
  httpPost("/api/prompt/translate", {"prompt": params.value.prompt}).then(res => {
    params.value.prompt = res.data
    loading.value = false
  }).catch(e => {
    loading.value = false
    ElMessage.error("翻译失败：" + e.message)
  })
}

const connect = () => {
  let host = process.env.VUE_APP_WS_HOST
  if (host === '') {
    if (location.protocol === 'https:') {
      host = 'wss://' + location.host;
    } else {
      host = 'ws://' + location.host;
    }
  }
  const _socket = new WebSocket(host + `/api/mj/client?user_id=${userId.value}`);
  _socket.addEventListener('open', () => {
    socket.value = _socket;
  });

  _socket.addEventListener('message', event => {
    if (event.data instanceof Blob) {
      fetchRunningJobs(userId.value)
      fetchFinishJobs(userId.value)
    }
  });

  _socket.addEventListener('close', () => {
    connect()
  });
}

onMounted(() => {
  checkSession().then(user => {
    imgCalls.value = user['img_calls']
    userId.value = user.id

    fetchRunningJobs(userId.value)
    fetchFinishJobs(userId.value)
    connect()

  }).catch(() => {
    router.push('/login')
  });

  const clipboard = new Clipboard('.copy-prompt');
  clipboard.on('success', () => {
    ElMessage.success("复制成功！");
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })
})

// 获取运行中的任务
const fetchRunningJobs = (userId) => {
  httpGet(`/api/mj/jobs?status=0&user_id=${userId}`).then(res => {
    const jobs = res.data
    const _jobs = []
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i].progress === -1) {
        ElNotification({
          title: '任务执行失败',
          message: "任务ID：" + jobs[i]['task_id'],
          type: 'error',
        })
        imgCalls.value += 1
        continue
      }
      _jobs.push(jobs[i])
    }
    runningJobs.value = _jobs
  }).catch(e => {
    ElMessage.error("获取任务失败：" + e.message)
  })
}

const fetchFinishJobs = (userId) => {
  // 获取已完成的任务
  httpGet(`/api/mj/jobs?status=1&user_id=${userId}`).then(res => {
    const jobs = res.data
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i]['use_proxy']) {
        jobs[i]['thumb_url'] = jobs[i]['img_url'] + '?x-oss-process=image/quality,q_60&format=webp'
      } else {
        if (jobs[i].type === 'upscale') {
          jobs[i]['thumb_url'] = jobs[i]['img_url'] + '?imageView2/1/w/480/h/600/q/75'
        } else {
          jobs[i]['thumb_url'] = jobs[i]['img_url'] + '?imageView2/1/w/480/h/480/q/75'
        }
      }
    }
    finishedJobs.value = jobs
  }).catch(e => {
    ElMessage.error("获取任务失败：" + e.message)
  })
}

// 切换图片比例
const changeRate = (item) => {
  params.value.rate = item.value
}
// 切换模型
const changeModel = (item) => {
  params.value.model = item.value
}

// 图片上传
const afterRead = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then((res) => {
        params.value.img = res.data
        ElMessage.success('上传成功')
      }).catch((e) => {
        ElMessage.error('上传失败:' + e.message)
      })
    },
    error(err) {
      console.log(err.message);
    },
  });
};
// 创建绘图任务
const promptRef = ref(null)
const generate = () => {
  if (params.value.prompt === '') {
    promptRef.value.focus()
    return ElMessage.error("请输入绘画提示词！")
  }
  if (params.value.model.indexOf("niji") !== -1 && params.value.raw) {
    return ElMessage.error("动漫模型不允许启用原始模式")
  }
  params.value.session_id = getSessionId()
  httpPost("/api/mj/image", params.value).then(() => {
    ElMessage.success("绘画任务推送成功，请耐心等待任务执行...")
    imgCalls.value -= 1
  }).catch(e => {
    ElMessage.error("任务推送失败：" + e.message)
  })
}

// 图片放大任务
const upscale = (index, item) => {
  send('/api/mj/upscale', index, item)
}

// 图片变换任务
const variation = (index, item) => {
  send('/api/mj/variation', index, item)
}

const send = (url, index, item) => {
  httpPost(url, {
    index: index,
    channel_id: item.channel_id,
    message_id: item.message_id,
    message_hash: item.hash,
    session_id: getSessionId(),
    prompt: item.prompt,
  }).then(() => {
    ElMessage.success("任务推送成功，请耐心等待任务执行...")
    imgCalls.value -= 1
  }).catch(e => {
    ElMessage.error("任务推送失败：" + e.message)
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
    httpPost("/api/mj/remove", {id: item.id, img_url: item.img_url, user_id: userId.value}).then(() => {
      ElMessage.success("任务删除成功")
    }).catch(e => {
      ElMessage.error("任务删除失败：" + e.message)
    })
  }).catch(() => {
  })
}

</script>

<style lang="stylus">
@import "@/assets/css/image-mj.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
