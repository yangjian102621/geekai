<template>
  <div class="page-mj">
    <div class="inner custom-scroll">
      <div class="mj-box">
        <h2>MidJourney 创作中心</h2>

        <div class="mj-params" :style="{ height: paramBoxHeight + 'px' }">
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
                    <el-select v-model="params.quality" placeholder="请选择" style="width:175px">
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
          <div class="extra-params">
            <el-form>
              <el-tabs v-model="activeName" class="title-tabs" @tabChange="tabChange">
                <el-tab-pane label="文生图" name="txt2img">
                  <div class="prompt-box">
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
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input v-model="params.prompt" :autosize="{ minRows: 4, maxRows: 6 }" type="textarea"
                                ref="promptRef"
                                placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"/>
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
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input v-model="params.neg_prompt" :autosize="{ minRows: 4, maxRows: 6 }" type="textarea"
                                ref="promptRef"
                                placeholder="请在此输入你不希望出现在图片上的内容，系统会自动翻译中文提示词"/>
                    </div>
                  </div>
                </el-tab-pane>
                <el-tab-pane label="图生图" name="img2img">
                  <div class="text">图生图：以某张图片为底稿参考来创作绘画，生成类似风格或类型图像，支持 PNG 和 JPG 格式图片；
                  </div>
                  <div class="param-line">
                    <div class="img-inline">
                      <div class="img-list-box">
                        <div class="img-item" v-for="imgURL in imgList">
                          <el-image :src="imgURL" fit="cover"/>
                          <el-button type="danger" :icon="Delete" @click="removeUploadImage(imgURL)" circle/>
                        </div>

                      </div>
                      <el-upload class="img-uploader" :auto-upload="true" :show-file-list="false"
                                 :http-request="uploadImg" style="--el-color-primary:#47fff1">
                        <el-icon class="uploader-icon">
                          <Plus/>
                        </el-icon>
                      </el-upload>
                    </div>
                  </div>

                  <div class="param-line" style="padding-top: 10px">
                    <el-form-item label="参考权重：">
                      <template #default>
                        <div class="form-item-inner">
                          <el-slider v-model.number="params.iw" :max="1" :step="0.01"
                                     style="width: 180px;--el-slider-main-bg-color:#47fff1"/>
                          <el-tooltip effect="light"
                                      content="使用图像权重参数--iw来调整图像 URL 与文本的重要性 <br/>权重较高时意味着图像提示将对完成的作业产生更大的影响"
                                      raw-content placement="right">
                            <el-icon>
                              <InfoFilled/>
                            </el-icon>
                          </el-tooltip>
                        </div>
                      </template>
                    </el-form-item>
                  </div>

                  <div class="prompt-box">
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
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input v-model="params.prompt" :autosize="{ minRows: 4, maxRows: 6 }" type="textarea"
                                ref="promptRef"
                                placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"/>
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
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input v-model="params.neg_prompt" :autosize="{ minRows: 4, maxRows: 6 }" type="textarea"
                                ref="promptRef"
                                placeholder="请在此输入你不希望出现在图片上的内容，系统会自动翻译中文提示词"/>
                    </div>
                  </div>
                </el-tab-pane>

                <el-tab-pane label="融图" name="blend">
                  <div class="text">请上传两张以上的图片，最多不超过五张，超过五张图片请使用图生图功能</div>
                  <div class="img-inline">
                    <div class="img-list-box">
                      <div class="img-item" v-for="imgURL in imgList">
                        <el-image :src="imgURL" fit="cover"/>
                        <el-button type="danger" :icon="Delete" @click="removeUploadImage(imgURL)" circle/>
                      </div>

                    </div>
                    <el-upload class="img-uploader" :auto-upload="true" :show-file-list="false"
                               :http-request="uploadImg" style="--el-color-primary:#47fff1">
                      <el-icon class="uploader-icon">
                        <Plus/>
                      </el-icon>
                    </el-upload>
                  </div>
                </el-tab-pane>

                <el-tab-pane label="换脸" name="swapFace">
                  <div class="text">请上传两张有脸部的图片，用左边图片的脸替换右边图片的脸</div>
                  <div class="img-inline">
                    <div class="img-list-box">
                      <div class="img-item" v-for="imgURL in imgList">
                        <el-image :src="imgURL" fit="cover"/>
                        <el-button type="danger" :icon="Delete" @click="removeUploadImage(imgURL)" circle/>
                      </div>

                    </div>
                    <el-upload class="img-uploader" :auto-upload="true" :show-file-list="false"
                               :http-request="uploadImg" style="--el-color-primary:#47fff1">
                      <el-icon class="uploader-icon">
                        <Plus/>
                      </el-icon>
                    </el-upload>
                  </div>
                </el-tab-pane>

                <el-tab-pane name="cref">
                  <template #label>
                    <el-badge value="New">
                      <span>一致性</span>
                    </el-badge>
                  </template>

                  <div class="text">注意：只有于 niji6 和 v6 模型支持一致性功能，如果选择其他模型此功能将会生成失败。</div>
                  <div class="param-line">
                    <el-form-item label="角色一致性：" prop="cref">
                      <el-input v-model="params.cref" placeholder="请输入图片URL或者上传图片"
                                style="--el-input-focus-border-color:#47fff1;--el-input-text-color:#ffffff; max-width: 500px; width: 100%"
                                size="small">
                        <template #append>
                          <el-upload
                              :auto-upload="true"
                              :show-file-list="false"
                              @click="beforeUpload('cref')"
                              :http-request="uploadImg"
                          >
                            <el-icon class="uploader-icon">
                              <UploadFilled/>
                            </el-icon>
                          </el-upload>
                        </template>
                      </el-input>
                    </el-form-item>
                  </div>

                  <div class="param-line">
                    <el-form-item label="风格一致性：" prop="sref">
                      <el-input v-model="params.sref" placeholder="请输入图片URL或者上传图片"
                                style="--el-input-focus-border-color:#47fff1; --el-input-text-color:#ffffff; max-width: 500px; width: 100%"
                                size="small">
                        <template #append>
                          <el-upload
                              :auto-upload="true"
                              :show-file-list="false"
                              @click="beforeUpload('sref')"
                              :http-request="uploadImg"
                          >
                            <el-icon class="uploader-icon">
                              <UploadFilled/>
                            </el-icon>
                          </el-upload>
                        </template>
                      </el-input>
                    </el-form-item>
                  </div>

                  <div class="param-line" style="padding-top: 10px">
                    <el-form-item label="参考权重：">
                      <template #default>
                        <div class="form-item-inner">
                          <el-slider v-model.number="params.cw" :max="100" :step="1"
                                     style="width: 180px;--el-slider-main-bg-color:#47fff1"/>
                          <el-tooltip effect="light"
                                      content="取值范围 0-100 <br/>默认值100参考原图的脸部、头发和衣服<br/>0则表示只换脸"
                                      raw-content placement="right">
                            <el-icon>
                              <InfoFilled/>
                            </el-icon>
                          </el-tooltip>
                        </div>
                      </template>
                    </el-form-item>
                  </div>

                  <div class="prompt-box">
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
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input v-model="params.prompt" :autosize="{ minRows: 4, maxRows: 6 }" type="textarea"
                                ref="promptRef"
                                placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"/>
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
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input v-model="params.neg_prompt" :autosize="{ minRows: 4, maxRows: 6 }" type="textarea"
                                ref="promptRef"
                                placeholder="请在此输入你不希望出现在图片上的内容，系统会自动翻译中文提示词"/>
                    </div>
                  </div>
                </el-tab-pane>
              </el-tabs>

              <el-row class="text-info">
                <el-tag>每次绘图消耗{{ mjPower }}算力，U/V 操作消耗{{ mjActionPower }}算力</el-tag>
                <el-tag type="success">当前可用算力：{{ power }}</el-tag>
              </el-row>

              <div class="submit-btn">
                <el-button color="#47fff1" :dark="false" @click="generate" round>立即生成</el-button>
              </div>
            </el-form>
          </div>

          <div class="job-list-box">
            <h2>任务列表</h2>
            <task-list :list="runningJobs" />

            <h2>创作记录</h2>
            <div class="finish-job-list">
              <div v-if="finishedJobs.length > 0">
                <v3-waterfall
                    id="waterfall"
                    :list="finishedJobs"
                    srcKey="thumb_url"
                    :gap="20"
                    :bottomGap="-8"
                    :colWidth="colWidth"
                    :distanceToScroll="100"
                    :isLoading="loading"
                    :isOver="isOver"
                    @scrollReachBottom="fetchFinishJobs()">
                  <template #default="slotProp">
                    <div class="job-item">
                      <el-image
                          v-if="slotProp.item.img_url !== ''"
                          :src="slotProp.item['thumb_url']"
                          :class="slotProp.item['can_opt'] ? '' : 'upscale'" @click="previewImg(slotProp.item)"
                          fit="cover" :initial-index="slotProp.index"
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

                      <div class="opt" v-if="slotProp.item['can_opt']">
                        <div class="opt-line">
                          <ul>
                            <li><a @click="upscale(1, slotProp.item)">U1</a></li>
                            <li><a @click="upscale(2, slotProp.item)">U2</a></li>
                            <li><a @click="upscale(3, slotProp.item)">U3</a></li>
                            <li><a @click="upscale(4, slotProp.item)">U4</a></li>
                            <li class="show-prompt">

                              <el-popover placement="left" title="提示词" :width="240" trigger="hover">
                                <template #reference>
                                  <el-icon>
                                    <ChromeFilled/>
                                  </el-icon>
                                </template>

                                <template #default>
                                  <div class="mj-list-item-prompt">
                                    <span>{{ slotProp.item.prompt }}</span>
                                    <el-icon class="copy-prompt-mj"
                                             :data-clipboard-text="slotProp.item.prompt">
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
                            <li><a @click="variation(1, slotProp.item)">V1</a></li>
                            <li><a @click="variation(2, slotProp.item)">V2</a></li>
                            <li><a @click="variation(3, slotProp.item)">V3</a></li>
                            <li><a @click="variation(4, slotProp.item)">V4</a></li>
                          </ul>
                        </div>
                      </div>

                      <div class="remove" v-if="slotProp.item.progress === 100">
                        <el-tooltip effect="light" content="删除任务" placement="top">
                          <el-button type="danger" :icon="Delete" @click="removeImage(slotProp.item)" circle/>
                        </el-tooltip>
                        <el-tooltip effect="light" content="取消发布" placement="top" v-if="slotProp.item.publish">
                          <el-button type="warning"
                                     @click="publishImage(slotProp.item, false)"
                                     circle>
                            <i class="iconfont icon-cancel-share"></i>
                          </el-button>
                        </el-tooltip>
                        <el-tooltip effect="light" content="发布图片" placement="top" v-else>
                          <el-button type="success" @click="publishImage(slotProp.item, true)" circle>
                            <i class="iconfont icon-share-bold"></i>
                          </el-button>
                        </el-tooltip>

                        <el-tooltip effect="light" content="复制提示词" placement="top">
                          <el-button type="success" class="copy-prompt-mj"
                                     :data-clipboard-text="slotProp.item.prompt" circle>
                            <el-icon><DocumentCopy/></el-icon>
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

    <el-image-viewer @close="() => { previewURL = '' }" v-if="previewURL !== ''" :url-list="[previewURL]"/>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from "vue"
import {ChromeFilled, Delete, DocumentCopy, InfoFilled, Picture, Plus, UploadFilled} from "@element-plus/icons-vue";
import Compressor from "compressorjs";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElMessageBox, ElNotification} from "element-plus";
import Clipboard from "clipboard";
import {checkSession, getSystemInfo} from "@/store/cache";
import {useRouter} from "vue-router";
import {getSessionId} from "@/store/session";
import {copyObj, removeArrayItem} from "@/utils/libs";
import {useSharedStore} from "@/store/sharedata";
import TaskList from "@/components/TaskList.vue";
import BackTop from "@/components/BackTop.vue";

const listBoxHeight = ref(0)
const paramBoxHeight = ref(0)
const loading = ref(true)
const colWidth = ref(220)
const previewURL = ref("")
const store = useSharedStore();

const resizeElement = function () {
  listBoxHeight.value = window.innerHeight - 80
  paramBoxHeight.value = window.innerHeight - 160
};
resizeElement()
window.onresize = () => {
  resizeElement()
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
  {text: "动漫风-niji4", value: " --niji 4", img: "/images/mj/nj4.jpg"},
  {text: "动漫风-niji5", value: " --niji 5", img: "/images/mj/mj-niji.png"},
  {text: "动漫风-niji5 可爱", value: " --niji 5 --style cute", img: "/images/mj/nj1.jpg"},
  {text: "动漫风-niji5 风景", value: " --niji 5 --style scenic", img: "/images/mj/nj2.jpg"},
  {text: "动漫风-niji6", value: " --niji 6", img: "/images/mj/nj3.jpg"},

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

const router = useRouter()
const initParams = {
  task_type: "image",
  rate: rates[0].value,
  model: models[0].value,
  chaos: 0,
  stylize: 0,
  seed: 0,
  img_arr: [],
  raw: false,
  iw: 0,
  prompt: router.currentRoute.value.params["prompt"] ?? "",
  neg_prompt: "",
  tile: false,
  quality: 0,
  cref: "",
  sref: "",
  cw: 0,
}
const params = ref(copyObj(initParams))

const imgList = ref([])

const activeName = ref('txt2img')

const runningJobs = ref([])
const finishedJobs = ref([])

const socket = ref(null)
const power = ref(0)
const userId = ref(0)
const isLogin = ref(false)

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

  const _socket = new WebSocket(host + `/api/mj/client?user_id=${userId.value}`);
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
        if (message === "FINISH" || message === "FAIL") {
          page.value = 0
          isOver.value = false
          fetchFinishJobs(page.value)
        }
        fetchRunningJobs()
      }
    }
  });

  _socket.addEventListener('close', () => {
    if (socket.value !== null) {
      connect()
    }
  });
}

const clipboard = ref(null)
onMounted(() => {
  initData()
  clipboard.value = new Clipboard('.copy-prompt-mj');
  clipboard.value.on('success', () => {
    ElMessage.success("复制成功！");
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！');
  })
})

onUnmounted(() => {
  clipboard.value.destroy()
  if (socket.value !== null) {
    socket.value.close()
    socket.value = null
  }
})

// 初始化数据
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

const mjPower = ref(1)
const mjActionPower = ref(1)
getSystemInfo().then(res => {
  mjPower.value = res.data["mj_power"]
  mjActionPower.value = res.data["mj_action_power"]
}).catch(e => {
  ElMessage.error("获取系统配置失败：" + e.message)
})

// 获取运行中的任务
const fetchRunningJobs = () => {
  if (!isLogin.value) {
    return
  }

  httpGet(`/api/mj/jobs?finish=false`).then(res => {
    const jobs = res.data
    const _jobs = []
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i].progress === 101) {
        ElNotification({
          title: '任务执行失败',
          dangerouslyUseHTMLString: true,
          message: `任务ID：${jobs[i]['task_id']}<br />原因：${jobs[i]['err_msg']}`,
          type: 'error',
          duration: 0,
        })
        if (jobs[i].type === 'image') {
          power.value += mjPower.value
        } else {
          power.value += mjActionPower.value
        }
      }
      _jobs.push(jobs[i])
    }
    runningJobs.value = _jobs
  }).catch(e => {
    ElMessage.error("获取任务失败：" + e.message)
  })
}

const page = ref(0)
const pageSize = ref(15)
const isOver = ref(false)
const fetchFinishJobs = () => {
  if (!isLogin.value || isOver.value) {
    return
  }

  loading.value = true
  page.value = page.value + 1
  // 获取已完成的任务
  httpGet(`/api/mj/jobs?finish=true&page=${page.value}&page_size=${pageSize.value}`).then(res => {
    const jobs = res.data
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i]['img_url'] !== "") {
        if (jobs[i].type === 'upscale' || jobs[i].type === 'swapFace') {
          jobs[i]['thumb_url'] = jobs[i]['img_url'] + '?imageView2/1/w/480/h/600/q/75'
        } else {
          jobs[i]['thumb_url'] = jobs[i]['img_url'] + '?imageView2/1/w/480/h/480/q/75'
        }
      } else {
        jobs[i]['thumb_url'] = '/images/img-placeholder.jpg'
      }

      if ((jobs[i].type === 'image' || jobs[i].type === 'variation') && jobs[i].progress === 100) {
        jobs[i]['can_opt'] = true
      }
    }
    if (jobs.length < pageSize.value) {
      isOver.value = true
    }
    if (page.value === 1) {
      finishedJobs.value = jobs
    } else {
      finishedJobs.value = finishedJobs.value.concat(jobs)
    }
    loading.value = false
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

const imgKey = ref("")
const beforeUpload = (key) => {
  imgKey.value = key
}

// 图片上传
const uploadImg = (file) => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }

  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then((res) => {
        if (imgKey.value === '') {
          imgList.value.push(res.data.url)
        } else { // 单张图片上传
          params.value[imgKey.value] = res.data.url
          imgKey.value = ''
        }
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
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }

  if (params.value.prompt === '' && params.value.task_type === "image") {
    promptRef.value.focus()
    return ElMessage.error("请输入绘画提示词！")
  }
  if (params.value.model.indexOf("niji") !== -1 && params.value.raw) {
    return ElMessage.error("动漫模型不允许启用原始模式")
  }
  if (imgList.value.length !== 2 && params.value.task_type === "swapFace") {
    return ElMessage.error("换脸操作需要上传两张图片")
  }
  params.value.session_id = getSessionId()
  params.value.img_arr = imgList.value
  httpPost("/api/mj/image", params.value).then(() => {
    ElMessage.success("绘画任务推送成功，请耐心等待任务执行...")
    power.value -= mjPower.value
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
    power.value -= mjActionPower.value
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
    httpGet("/api/mj/remove", {id: item.id, user_id: item.user_id}).then(() => {
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

// 发布图片到作品墙
const publishImage = (item, action) => {
  let text = "图片发布"
  if (action === false) {
    text = "取消发布"
  }
  httpGet("/api/mj/publish", {id: item.id, action: action,user_id: item.user_id}).then(() => {
    ElMessage.success(text + "成功")
    item.publish = action
    page.value = 0
    isOver.value = false
    item.publish = action
  }).catch(e => {
    ElMessage.error(text + "失败：" + e.message)
  })
}

const previewImg = (item) => {
  previewURL.value = item.img_url

}

// 切换菜单
const tabChange = (tab) => {
  if (tab === "txt2img" || tab === "img2img" || tab === "cref") {
    params.value.task_type = "image"
  } else {
    params.value.task_type = tab
  }
}

// 删除已上传图片
const removeUploadImage = (url) => {
  imgList.value = removeArrayItem(imgList.value, url)
}

</script>

<style lang="stylus">
@import "@/assets/css/image-mj.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
