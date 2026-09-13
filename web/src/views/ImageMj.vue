<template>
  <div class="page-mj">
    <div class="inner custom-scroll">
      <div class="mj-box">
        <header class="mj-box__head">
          <h2 class="mj-box__title">MidJourneyAI生图</h2>
        </header>

        <div class="mj-params" :style="{ height: paramBoxHeight + 'px' }">
          <el-form :model="params" label-position="top" class="mj-params-form" label-width="auto">
            <section class="mj-sec">
              <h3 class="mj-sec__title">比例</h3>
              <p class="mj-sec__hint" title="生成图的宽高比；也可在提示词末尾写 --ar 宽:高">
                点击选择比例；或在提示词中加 <code class="mj-code">--ar w:h</code>
              </p>
              <div class="mj-aspect-chips" role="list">
                <button
                  v-for="item in rates"
                  :key="item.value"
                  type="button"
                  class="mj-chip"
                  :class="{ 'mj-chip--active': item.value === params.rate }"
                  :title="'比例 ' + item.text"
                  @click="changeRate(item)"
                >
                  {{ item.text }}
                </button>
              </div>
            </section>

            <section class="mj-sec">
              <h3 class="mj-sec__title">画质</h3>
              <p class="mj-sec__hint">越高越慢，细节相对更好。</p>
              <el-form-item label="输出档位" class="mj-form-item--compact">
                <el-select v-model="params.quality" placeholder="默认" class="mj-select-full">
                  <el-option
                    v-for="item in options"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-form-item>
            </section>

            <section class="mj-sec">
              <h3 class="mj-sec__title">模型</h3>
              <div class="mj-model-cards">
                <button
                  v-for="item in models"
                  :key="item.value"
                  type="button"
                  class="mj-model-card"
                  :class="{ 'mj-model-card--active': item.value === params.model }"
                  @click="changeModel(item)"
                >
                  <div class="mj-model-card__row">
                    <span class="mj-model-card__name">{{ item.text }}</span>
                    <code class="mj-model-card__flag">{{ item.flag }}</code>
                  </div>
                  <div class="mj-model-card__badge">{{ item.badge }}</div>
                  <p class="mj-model-card__summary">{{ item.summary }}</p>
                </button>
              </div>
            </section>

            <section class="mj-sec">
              <h3 class="mj-sec__title">选项</h3>
              <div class="mj-switch-grid">
                <label class="mj-switch-row">
                  <span class="mj-switch-row__text"
                    >重复平铺
                    <el-tooltip content="生成可无缝平铺的重复图案" placement="right">
                      <i class="iconfont icon-info !text-sm"></i>
                    </el-tooltip>
                  </span>

                  <el-switch v-model="params.tile" inactive-color="#464649" />
                </label>
                <label class="mj-switch-row">
                  <span class="mj-switch-row__text"
                    >RAW 模式
                    <el-tooltip
                      content="更写实、细节更强，建议写更长提示词。动漫模型勿开。"
                      placement="right"
                    >
                      <i class="iconfont icon-info !text-sm"></i>
                    </el-tooltip>
                  </span>

                  <el-switch v-model="params.raw" inactive-color="#464649" />
                </label>
              </div>
            </section>

            <section class="mj-sec mj-sec--last">
              <h3 class="mj-sec__title">高级</h3>
              <div class="mj-slider-field" title="--chaos 0–100，越高越发散，0 最稳">
                <el-form-item class="mj-form-item--compact">
                  <template #label>
                    <span class="mr-1">创意度</span>
                    <el-tooltip content="--chaos 0–100，越高越发散，0 最稳" placement="right">
                      <i class="iconfont icon-info !text-sm"></i>
                    </el-tooltip>
                  </template>
                  <div class="w-full px-4">
                    <el-slider v-model.number="params.chaos" :max="100" :step="1" />
                  </div>
                </el-form-item>
              </div>
              <div class="mj-slider-field" title="--stylize 0–1000，越高越艺术">
                <el-form-item class="mj-form-item--compact">
                  <template #label>
                    <span class="mr-1">风格化</span>
                    <el-tooltip content="--stylize 0–1000，越高越艺术" placement="right">
                      <i class="iconfont icon-info !text-sm"></i>
                    </el-tooltip>
                  </template>
                  <div class="w-full px-4">
                    <el-slider v-model.number="params.stylize" :min="0" :max="1000" :step="1" />
                  </div>
                </el-form-item>
              </div>
              <div class="mj-slider-field" title="--seed：0 或默认随机；相同种子与描述可复现相似图">
                <el-form-item label="随机种子" class="mj-form-item--compact">
                  <el-input v-model.number="params.seed" />
                </el-form-item>
              </div>
            </section>
          </el-form>
        </div>
      </div>
      <div class="task-list-box pl-6 pr-6 pb-4 h-dvh">
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
                          <el-tooltip content="输入你想要的内容，用逗号分割" placement="right">
                            <el-icon>
                              <InfoFilled />
                            </el-icon>
                          </el-tooltip>
                          <div
                            class="flex-row justify-start items-center m-2 p-3 bg-gray-100 rounded-md text-gray-500 text-sm"
                          >
                            <span
                              >如需自定义比例，在绘画指令最后加一个空格然后加上指令(宽高比) --ar w:h
                              例如: 1 cat --ar 21:9
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>

                    <div class="param-line pt" style="position: relative">
                      <el-input
                        v-model="params.prompt"
                        :autosize="{ minRows: 4, maxRows: 6 }"
                        maxlength="1024"
                        show-word-limit
                        type="textarea"
                        ref="promptTextareaTxt"
                        v-loading="promptGenerating"
                        placeholder="请在此输入绘画提示词，您也可以点击下面的提示词助手生成绘画提示词"
                      />
                    </div>

                    <div class="flex justify-end pt-2">
                      <el-button @click="generatePrompt" type="primary" :loading="promptGenerating">
                        <span v-if="!promptGenerating">
                          <i class="iconfont icon-chuangzuo"></i>
                          生成专业绘画指令
                        </span>
                        <span v-else>生成中...</span>
                      </el-button>
                    </div>

                    <div class="param-line pt">
                      <div class="flex-row justify-between items-center">
                        <div class="flex-row justify-start items-center">
                          <span>不希望出现的内容：（可选）</span>
                          <el-tooltip
                            content="不想出现在图片上的元素(例如：树，建筑)"
                            placement="right"
                          >
                            <el-icon>
                              <InfoFilled />
                            </el-icon>
                          </el-tooltip>
                        </div>
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input
                        v-model="params.neg_prompt"
                        :autosize="{ minRows: 4, maxRows: 6 }"
                        type="textarea"
                        maxlength="2000"
                        placeholder="请在此输入你不希望出现在图片上的内容，系统会自动翻译中文提示词"
                      />
                    </div>
                  </div>
                </el-tab-pane>
                <el-tab-pane label="图生图" name="img2img">
                  <div class="text">
                    图生图：以某张图片为底稿参考来创作绘画，生成类似风格或类型图像，支持 PNG 和 JPG
                    格式图片；
                  </div>
                  <div class="param-line">
                    <ImageUpload v-model="imgList" :max-count="5" :multiple="true" />
                  </div>

                  <div class="param-line" style="padding-top: 10px">
                    <el-form-item label="参考权重：">
                      <template #default>
                        <div class="form-item-inner">
                          <el-slider
                            v-model.number="params.iw"
                            :max="1"
                            :step="0.01"
                            style="width: 180px"
                          />
                          <el-tooltip
                            content="使用图像权重参数--iw来调整图像 URL 与文本的重要性 <br/>权重较高时意味着图像提示将对完成的作业产生更大的影响"
                            raw-content
                            placement="right"
                          >
                            <el-icon>
                              <InfoFilled />
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
                          <el-tooltip content="输入你想要的内容，用逗号分割" placement="right">
                            <el-icon>
                              <InfoFilled />
                            </el-icon>
                          </el-tooltip>
                        </div>
                        <div class="flex-row justify-start items-center">
                          <span
                            >如需自定义比例，在绘画指令最后加一个空格然后加上指令(宽高比) --ar w:h
                            例如: 1 cat --ar 21:9
                          </span>
                        </div>
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input
                        v-model="params.prompt"
                        :autosize="{ minRows: 4, maxRows: 6 }"
                        type="textarea"
                        ref="promptTextareaImg"
                        v-loading="promptGenerating"
                        placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"
                      />
                    </div>

                    <el-row class="text-info">
                      <el-button
                        class="generate-btn"
                        size="small"
                        @click="generatePrompt"
                        color="#5865f2"
                        :disabled="promptGenerating"
                      >
                        <i class="iconfont icon-chuangzuo"></i>
                        <span>生成专业绘画指令</span>
                      </el-button>
                    </el-row>

                    <div class="param-line pt">
                      <div class="flex-row justify-between items-center">
                        <div class="flex-row justify-start items-center">
                          <span>不希望出现的内容：（可选）</span>
                          <el-tooltip
                            content="不想出现在图片上的元素(例如：树，建筑)"
                            placement="right"
                          >
                            <el-icon>
                              <InfoFilled />
                            </el-icon>
                          </el-tooltip>
                        </div>
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input
                        v-model="params.neg_prompt"
                        :autosize="{ minRows: 4, maxRows: 6 }"
                        type="textarea"
                        placeholder="请在此输入你不希望出现在图片上的内容，系统会自动翻译中文提示词"
                      />
                    </div>
                  </div>
                </el-tab-pane>

                <el-tab-pane label="融图" name="blend">
                  <div class="text">
                    请上传两张以上的图片，最多不超过五张，超过五张图片请使用图生图功能
                  </div>
                  <div class="param-line">
                    <ImageUpload v-model="imgList" :max-count="5" :multiple="true" />
                  </div>
                </el-tab-pane>

                <el-tab-pane label="换脸" name="swapFace">
                  <div class="text">请上传两张有脸部的图片，用左边图片的脸替换右边图片的脸</div>
                  <div class="param-line">
                    <ImageUpload v-model="imgList" :max-count="2" :multiple="true" />
                  </div>
                </el-tab-pane>

                <el-tab-pane name="cref">
                  <template #label>
                    <el-badge value="New">
                      <span>一致性</span>
                    </el-badge>
                  </template>

                  <div class="text">
                    注意：仅 <code>--niji 6</code> 与
                    <code>--v 6.1</code> 支持一致性功能；选其他模型会生成失败。
                  </div>
                  <div class="param-line cref-two-cols">
                    <div class="cref-col">
                      <label class="cref-label">角色一致性</label>
                      <ImageUpload v-model="params.cref" :max-count="1" class="cref-upload-inner" />
                    </div>
                    <div class="cref-col">
                      <label class="cref-label">风格一致性</label>
                      <ImageUpload v-model="params.sref" :max-count="1" class="cref-upload-inner" />
                    </div>
                  </div>

                  <div class="param-line" style="padding-top: 10px">
                    <el-form-item label="参考权重：">
                      <template #default>
                        <div class="form-item-inner">
                          <el-slider
                            v-model.number="params.cw"
                            :max="100"
                            :step="1"
                            style="width: 180px"
                          />
                          <el-tooltip
                            content="取值范围 0-100 <br/>默认值100参考原图的脸部、头发和衣服<br/>0则表示只换脸"
                            raw-content
                            placement="right"
                          >
                            <el-icon>
                              <InfoFilled />
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
                          <el-tooltip content="输入你想要的内容，用逗号分割" placement="right">
                            <el-icon>
                              <InfoFilled />
                            </el-icon>
                          </el-tooltip>
                          <div class="flex-row justify-start items-center">
                            <span
                              >如需自定义比例，在绘画指令最后加一个空格然后加上指令(宽高比) --ar w:h
                              例如: 1 cat --ar 21:9
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input
                        v-model="params.prompt"
                        :autosize="{ minRows: 4, maxRows: 6 }"
                        type="textarea"
                        ref="promptTextareaCref"
                        placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"
                      />
                    </div>

                    <div class="param-line pt">
                      <div class="flex-row justify-between items-center">
                        <div class="flex-row justify-start items-center">
                          <span>不希望出现的内容：（可选）</span>
                          <el-tooltip
                            content="不想出现在图片上的元素(例如：树，建筑)"
                            placement="right"
                          >
                            <el-icon>
                              <InfoFilled />
                            </el-icon>
                          </el-tooltip>
                        </div>
                      </div>
                    </div>

                    <div class="param-line pt">
                      <el-input
                        v-model="params.neg_prompt"
                        :autosize="{ minRows: 4, maxRows: 6 }"
                        type="textarea"
                        placeholder="请在此输入你不希望出现在图片上的内容，系统会自动翻译中文提示词"
                      />
                    </div>
                  </div>
                </el-tab-pane>
              </el-tabs>

              <el-row class="text-info">
                <el-text type="primary">
                  绘图 {{ mjPower }} 积分；U/V {{ mjUpscalePower }}；融图 {{ mjBlendPower }}；换脸
                  {{ mjSwapFacePower }}；局部重绘 {{ mjModalPower }} 积分； 当前可用：<el-text
                    type="warning"
                    >{{ power }}</el-text
                  >
                </el-text>
              </el-row>

              <div class="submit-btn">
                <button
                  class="px-10 py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
                  @click="generate"
                  type="button"
                >
                  <i v-if="isGenerating" class="iconfont icon-loading animate-spin"></i>
                  <i v-else class="iconfont icon-chuangzuo"></i>
                  <span>{{ isGenerating ? '创作中...' : '立即生成' }}</span>
                </button>
              </div>
            </el-form>
          </div>

          <div class="job-list-box">
            <h2 class="text-xl">任务列表</h2>
            <task-list :list="runningJobs" />
            <template v-if="finishedJobs.length > 0">
              <h2 class="text-xl">创作记录</h2>
              <div class="finish-job-list mt-3">
                <Waterfall
                  :list="finishedJobs"
                  :row-key="waterfallOptions.rowKey"
                  :gutter="waterfallOptions.gutter"
                  :has-around-gutter="waterfallOptions.hasAroundGutter"
                  :width="waterfallOptions.width"
                  :breakpoints="waterfallOptions.breakpoints"
                  :img-selector="waterfallOptions.imgSelector"
                  :background-color="waterfallOptions.backgroundColor"
                  :animation-effect="waterfallOptions.animationEffect"
                  :animation-duration="waterfallOptions.animationDuration"
                  :animation-delay="waterfallOptions.animationDelay"
                  :animation-cancel="waterfallOptions.animationCancel"
                  :lazyload="waterfallOptions.lazyload"
                  :load-props="waterfallOptions.loadProps"
                  :cross-origin="waterfallOptions.crossOrigin"
                  :align="waterfallOptions.align"
                  :is-loading="loading"
                  :is-over="isOver"
                  @afterRender="loading = false"
                >
                  <template #default="{ item, url }">
                    <div class="image-task-item">
                      <div class="image-task-media">
                        <div
                          class="image-task-preview"
                          :class="{ 'image-task-preview--failed': item.status === 'failed' }"
                        >
                          <LazyImg
                            v-if="item.status === 'success'"
                            :url="url"
                            class="image-task-image"
                            @click="previewImg(item)"
                          />
                          <img
                            v-else-if="item.status === 'failed'"
                            class="image-task-image image-task-image--failed"
                            :src="taskFailedImage"
                            title="点击查看详情"
                            @click="showDetail(item)"
                          />
                        </div>
                        <div
                          v-if="item.status === 'success'"
                          class="image-task-overlay mj-task-overlay--meta"
                          @click.stop
                        >
                          <div class="image-task-overlay-time">
                            {{ dateFormat(item.created_at) }}
                          </div>
                          <div class="image-task-tools">
                            <el-tooltip content="取消分享" placement="top" v-if="item.publish">
                              <button
                                type="button"
                                class="image-task-tool"
                                @click="publishImage(item, false)"
                              >
                                <i class="iconfont icon-cancel-share"></i>
                              </button>
                            </el-tooltip>
                            <el-tooltip content="分享" placement="top" v-else>
                              <button
                                type="button"
                                class="image-task-tool"
                                @click="publishImage(item, true)"
                              >
                                <i class="iconfont icon-share-bold"></i>
                              </button>
                            </el-tooltip>
                            <el-tooltip content="任务详情" placement="top">
                              <button
                                type="button"
                                class="image-task-tool"
                                @click="showDetail(item)"
                              >
                                <i class="iconfont icon-info text-[#6366f1]"></i>
                              </button>
                            </el-tooltip>
                            <el-tooltip content="删除" placement="top">
                              <button
                                type="button"
                                class="image-task-tool image-task-tool--danger"
                                @click="removeImage(item)"
                              >
                                <i class="iconfont icon-remove"></i>
                              </button>
                            </el-tooltip>
                          </div>
                        </div>
                        <div
                          v-else-if="item.status === 'failed'"
                          class="image-task-overlay mj-task-overlay--failed"
                          @click.stop
                        >
                          <div class="image-task-overlay-time">
                            {{ dateFormat(item.created_at) }}
                          </div>
                          <div class="image-task-tools">
                            <el-popover
                              title="错误详情"
                              trigger="click"
                              :width="260"
                              :content="item['err_msg']"
                              placement="top"
                            >
                              <template #reference>
                                <button type="button" class="image-task-tool">
                                  <i class="iconfont icon-info text-[#6366f1]"></i>
                                </button>
                              </template>
                            </el-popover>
                            <el-tooltip content="删除" placement="top">
                              <button
                                type="button"
                                class="image-task-tool image-task-tool--danger"
                                @click="removeImage(item)"
                              >
                                <i class="iconfont icon-remove"></i>
                              </button>
                            </el-tooltip>
                          </div>
                        </div>
                      </div>
                      <div v-if="item.status === 'success'" class="mj-task-opt" @click.stop>
                        <div v-if="item['can_opt']" class="mj-task-opt__uv">
                          <button
                            v-for="i in 4"
                            :key="'u' + i"
                            type="button"
                            class="mj-task-opt-btn mj-task-opt-btn--uv"
                            @click="upscale(i, item)"
                          >
                            U{{ i }}
                          </button>
                          <button
                            v-for="i in 4"
                            :key="'v' + i"
                            type="button"
                            class="mj-task-opt-btn mj-task-opt-btn--uv"
                            @click="variation(i, item)"
                          >
                            V{{ i }}
                          </button>
                        </div>
                        <div v-if="item['can_modal']" class="mj-task-opt__modal">
                          <button
                            type="button"
                            class="mj-task-opt-btn mj-task-opt-btn--modal"
                            @click="openModalDialog(item)"
                          >
                            局部重绘
                          </button>
                        </div>
                      </div>
                    </div>
                  </template>
                </Waterfall>

                <div class="flex justify-center py-10">
                  <img
                    :src="waterfallOptions.loadProps.loading"
                    class="max-w-[50px] max-h-[50px]"
                    v-if="loading"
                  />
                  <div v-else>
                    <button
                      class="px-5 py-2 rounded-full bg-purple-700 text-md text-white cursor-pointer hover:bg-purple-800 transition-all duration-300"
                      @click="fetchFinishJobs"
                      v-if="!isOver"
                    >
                      加载更多
                    </button>
                    <div class="no-more-data" v-else>
                      <span class="text-gray-500 mr-2">没有更多数据了</span>
                      <i class="iconfont icon-face"></i>
                    </div>
                  </div>
                </div>
              </div>
            </template>

            <!-- end finish job list-->
          </div>
        </div>
        <back-top :right="30" :bottom="30" />
      </div>
      <!-- end task list box -->
    </div>

    <el-image-viewer v-if="previewURL !== ''" :url-list="[previewURL]" @close="closePreview" />

    <!-- 局部重绘弹窗 -->
    <el-dialog
      v-model="modalVisible"
      title="局部重绘"
      class="modal-inpaint-dialog"
      destroy-on-close
      @closed="onModalClosed"
    >
      <div class="modal-inpaint-box">
        <p class="text-sm text-gray-600 mb-2">
          下方为原图，请在图上涂抹需要重绘的区域（半透明红色即蒙版），可调整画笔大小或清除后重画。蒙版按原图尺寸导出，与涂抹区域一致。
        </p>
        <div
          class="relative inline-block rounded overflow-hidden border border-gray-200 bg-gray-100 modal-inpaint-wrap"
          ref="modalCanvasWrapRef"
        >
          <img
            :src="modalItem?.img_url"
            ref="modalImgRef"
            class="block max-h-[75vh] w-auto max-w-full"
            crossorigin="anonymous"
            @load="onModalImageLoad"
          />
          <!-- 半透明蒙版层：与显示图同尺寸，用户在此绘制 -->
          <canvas
            ref="modalCanvasRef"
            class="absolute left-0 top-0 cursor-crosshair pointer-events-auto"
            :width="modalCanvasSize.w"
            :height="modalCanvasSize.h"
            @mousedown="modalDrawStart"
            @mousemove="modalDrawMove"
            @mouseup="modalDrawEnd"
            @mouseleave="modalDrawEnd"
          />
          <!-- 用于导出的黑白蒙版：原图 naturalWidth × naturalHeight -->
          <canvas
            ref="modalMaskCanvasRef"
            class="hidden"
            :width="modalNaturalSize.w"
            :height="modalNaturalSize.h"
          />
        </div>
        <div class="flex items-center gap-4 mt-3">
          <span class="text-sm">画笔大小：</span>
          <el-slider v-model="modalBrushSize" :min="4" :max="40" :step="2" style="width: 120px" />
          <el-button size="small" @click="modalClearMask">清除蒙版</el-button>
        </div>
        <el-form-item label="重绘提示词" required class="mt-3">
          <el-input
            v-model="modalPrompt"
            type="textarea"
            :rows="2"
            placeholder="描述你希望在该区域生成的内容"
            maxlength="1024"
            show-word-limit
          />
        </el-form-item>
      </div>
      <template #footer>
        <el-button @click="modalVisible = false">取消</el-button>
        <el-button type="primary" :loading="modalSubmitting" @click="submitModal">
          提交（消耗 {{ mjModalPower }} 积分）
        </el-button>
      </template>
    </el-dialog>

    <!-- 任务详情弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="任务详情"
      class="mj-detail-dialog"
      width="680px"
      :close-on-click-modal="false"
    >
      <div class="detail-content mj-detail-body" v-if="currentDetail">
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="任务 ID">
            <div class="mj-detail-copy-row">
              <span class="break-all font-mono text-[13px]">{{
                currentDetail.task_id || '-'
              }}</span>
              <el-tooltip v-if="currentDetail.task_id" content="复制" placement="top">
                <i
                  class="iconfont icon-copy mj-detail-copy-ico"
                  @click="copyPrompt(currentDetail.task_id)"
                />
              </el-tooltip>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="任务类型">
            {{ mjDetailTypeLabel(currentDetail.type) }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            {{ currentDetail.status || '-' }}
            <span
              v-if="currentDetail.progress != null && currentDetail.status !== 'success'"
              class="text-gray-500 ml-1"
            >
              ({{ currentDetail.progress }}%)
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="原始提示词" v-if="detailTaskPayload?.prompt">
            <div class="mj-detail-copy-row">
              <span class="break-all mj-detail-text">{{ detailTaskPayload.prompt }}</span>
              <el-tooltip content="复制" placement="top">
                <i
                  class="iconfont icon-copy mj-detail-copy-ico"
                  @click="copyPrompt(detailTaskPayload.prompt)"
                />
              </el-tooltip>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="负面提示词" v-if="detailTaskPayload?.neg_prompt">
            <span class="break-all mj-detail-text">{{ detailTaskPayload.neg_prompt }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="完整提示词">
            <div class="mj-detail-copy-row align-start">
              <span class="break-all mj-detail-text">{{ currentDetail.prompt || '—' }}</span>
              <el-tooltip v-if="currentDetail.prompt" content="复制" placement="top">
                <i
                  class="iconfont icon-copy mj-detail-copy-ico"
                  @click="copyPrompt(currentDetail.prompt)"
                />
              </el-tooltip>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="引用图片" v-if="detailTaskImages.length > 0">
            <div class="mj-detail-refimgs">
              <el-image
                v-for="(url, idx) in detailTaskImages"
                :key="'ref-' + idx"
                :src="url"
                :preview-src-list="detailTaskImages"
                fit="cover"
                class="mj-detail-refimg"
                preview-teleported
              />
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="局部重绘" v-if="detailHasMask">
            <el-text type="info">已提交蒙版（内容略）</el-text>
          </el-descriptions-item>
          <el-descriptions-item
            label="生成结果"
            v-if="currentDetail.status === 'success' && currentDetail.img_url"
          >
            <el-image
              :src="getThumbURL(currentDetail.img_url, 240, 240)"
              :preview-src-list="[currentDetail.img_url]"
              fit="cover"
              class="mj-detail-result"
              preview-teleported
            />
          </el-descriptions-item>
          <el-descriptions-item label="消耗积分">
            {{ currentDetail.power ?? 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ dateFormat(currentDetail.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item
            label="错误信息"
            v-if="currentDetail.status === 'failed' && currentDetail.err_msg"
          >
            <el-text type="danger">{{ currentDetail.err_msg }}</el-text>
          </el-descriptions-item>
          <el-descriptions-item
            label="原始载荷"
            v-if="currentDetail.task_info && !detailTaskPayload"
          >
            <span class="text-gray-500 text-xs">无法解析 JSON，以下为原始文本：</span>
            <pre class="mj-detail-raw">{{ currentDetail.task_info }}</pre>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import nodata from '@/assets/img/no-data.png'
import BackTop from '@/components/BackTop.vue'
import ImageUpload from '@/components/ImageUpload.vue'
import TaskList from '@/components/TaskList.vue'
import { checkSession, getSystemInfo } from '@/store/cache'
import { getSessionId } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import { showMessageError } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { copyObj, dateFormat, getThumbURL } from '@/utils/libs'
import { InfoFilled } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { LazyImg, Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'

const listBoxHeight = ref(0)
const paramBoxHeight = ref(0)
const loading = ref(true)
const previewURL = ref('')
const store = useSharedStore()
const waterfallOptions = store.waterfallOptions
const taskFailedImage = store.taskFailedImage

const resizeElement = function () {
  listBoxHeight.value = window.innerHeight - 30
  paramBoxHeight.value = window.innerHeight - 110
}
resizeElement()
window.onresize = () => {
  resizeElement()
}

const rates = [
  { value: '1:1', text: '1:1' },
  { value: '1:2', text: '1:2' },
  { value: '2:1', text: '2:1' },
  { value: '2:3', text: '2:3' },
  { value: '3:2', text: '3:2' },
  { value: '3:4', text: '3:4' },
  { value: '4:3', text: '4:3' },
  { value: '16:9', text: '16:9' },
  { value: '9:16', text: '9:16' },
]
const models = [
  {
    text: 'Midjourney V7',
    flag: '--v 7',
    value: ' --v 7',
    badge: '官方默认 · 全能',
    summary: '写实与国风最稳，画质与理解力强；适合绝大多数日常与商业出图。',
  },
  {
    text: 'Midjourney V8.1 Alpha',
    flag: '--v 8.1',
    value: ' --v 8.1',
    badge: '测试版 · 更快更省',
    summary: '速度明显快于 V7、成本更低，画质接近；适合批量与快速改词试错。',
  },
  {
    text: 'Niji・Journey V6',
    flag: '--niji 6',
    value: ' --niji 6',
    badge: '二次元 · 插画向',
    summary: '动漫、立绘、分镜强于通用 MJ；不做写实首选 Niji。',
  },
  {
    text: 'Midjourney V6.1',
    flag: '--v 6.1',
    value: ' --v 6.1',
    badge: '经典 · 兼容老词',
    summary: '旧 prompt 更稳、可控性高；适合不愿大改词、求稳妥的用户。',
  },
]

const options = [
  {
    value: 0,
    label: '默认',
  },
  {
    value: 0.25,
    label: '普通',
  },
  {
    value: 0.5,
    label: '清晰',
  },
  {
    value: 1,
    label: '高清',
  },
]

const router = useRouter()
const initParams = {
  task_type: 'image',
  rate: rates[0].value,
  model: models[0].value,
  chaos: 0,
  stylize: 0,
  seed: -1,
  img_arr: [],
  raw: false,
  iw: 0.7,
  prompt: router.currentRoute.value.params['prompt'] ?? '',
  neg_prompt: '',
  tile: false,
  quality: 0,
  cref: '',
  sref: '',
  cw: 0,
}
const params = ref(copyObj(initParams))

const imgList = ref([])

const activeName = ref('txt2img')

const promptTextareaTxt = ref(null)
const promptTextareaImg = ref(null)
const promptTextareaCref = ref(null)

function focusPromptInput() {
  const tab = activeName.value
  const el =
    tab === 'img2img'
      ? promptTextareaImg.value
      : tab === 'cref'
        ? promptTextareaCref.value
        : promptTextareaTxt.value
  el?.focus?.()
}

const runningJobs = ref([])
const finishedJobs = ref([])
const taskPulling = ref(true) // 任务轮询
const tastPullHandler = ref(null)
const downloadPulling = ref(false) // 图片下载轮询
const downloadPullHandler = ref(null)

const power = ref(0)

onMounted(() => {
  initData()
})

onUnmounted(() => {
  if (tastPullHandler.value) {
    clearInterval(tastPullHandler.value)
  }
  if (downloadPullHandler.value) {
    clearInterval(downloadPullHandler.value)
  }
})

// 初始化数据
const initData = () => {
  checkSession()
    .then((user) => {
      power.value = user['power']
      page.value = 0
      fetchFinishJobs()

      tastPullHandler.value = setInterval(() => {
        if (taskPulling.value) {
          fetchRunningJobs()
        }
      }, 5000)

      downloadPullHandler.value = setInterval(() => {
        if (downloadPulling.value) {
          page.value = 0
          isOver.value = false
          fetchFinishJobs()
        }
      }, 5000)
    })
    .catch(() => {})
}

const mjPower = ref(1)
const mjUpscalePower = ref(1)
const mjBlendPower = ref(1)
const mjSwapFacePower = ref(1)
const mjModalPower = ref(1)
getSystemInfo()
  .then((res) => {
    const d = res.data || {}
    const fallback = d['mj_action_power'] || 1
    mjPower.value = d['mj_power'] || 1
    mjUpscalePower.value = d['mj_upscale_power'] > 0 ? d['mj_upscale_power'] : fallback
    mjBlendPower.value = d['mj_blend_power'] > 0 ? d['mj_blend_power'] : fallback
    mjSwapFacePower.value = d['mj_swap_face_power'] > 0 ? d['mj_swap_face_power'] : fallback
    mjModalPower.value = d['mj_modal_power'] > 0 ? d['mj_modal_power'] : fallback
  })
  .catch((e) => {
    ElMessage.error('获取系统配置失败：' + e.message)
  })

// 获取运行中的任务
const fetchRunningJobs = () => {
  httpGet(`/api/mj/jobs?finish=false`)
    .then((res) => {
      const jobs = res.data.items
      const _jobs = []
      for (let i = 0; i < jobs.length; i++) {
        if (jobs[i].status === 'failed') {
          ElNotification({
            title: '任务执行失败',
            dangerouslyUseHTMLString: true,
            message: `任务ID：${jobs[i]['task_id']}<br />原因：${jobs[i]['err_msg']}`,
            type: 'error',
            duration: 0,
          })
          const refund =
            jobs[i].power != null && jobs[i].power > 0
              ? jobs[i].power
              : jobs[i].type === 'image'
                ? mjPower.value
                : mjUpscalePower.value
          power.value += refund
        }
        _jobs.push(jobs[i])
      }
      if (runningJobs.value.length !== _jobs.length) {
        page.value = 0
        downloadPulling.value = true
        fetchFinishJobs()
      }
      if (_jobs.length === 0) {
        taskPulling.value = false
      }
      runningJobs.value = _jobs
    })
    .catch((e) => {
      ElMessage.error('获取任务失败：' + e.message)
    })
}

const page = ref(0)
const pageSize = ref(15)
const isOver = ref(false)
const fetchFinishJobs = () => {
  if (isOver.value) {
    return
  }

  loading.value = true
  page.value = page.value + 1
  // 获取已完成的任务
  httpGet(`/api/mj/jobs?finish=true&page=${page.value}&page_size=${pageSize.value}`)
    .then((res) => {
      const jobs = res.data.items
      let needDownload = false
      for (let i = 0; i < jobs.length; i++) {
        if (jobs[i]['img_url'] !== '') {
          if (jobs[i].type === 'upscale' || jobs[i].type === 'swapFace') {
            jobs[i]['img_thumb'] = getThumbURL(jobs[i]['img_url'], 480, 600)
          } else {
            jobs[i]['img_thumb'] = getThumbURL(jobs[i]['img_url'], 480, 480)
          }
        } else {
          if (jobs[i].status === 'downloading') {
            needDownload = true
          }
          jobs[i]['img_thumb'] = waterfallOptions.loadProps.loading
        }
        // 如果当前是第一页，则开启图片下载轮询
        if (page.value === 1) {
          downloadPulling.value = needDownload
        }

        if (jobs[i].type !== 'upscale' && jobs[i].status === 'success') {
          jobs[i]['can_opt'] = true
        }
        // 所有已完成且有图的任务均支持局部重绘
        if (jobs[i].status === 'success' && jobs[i].img_url) {
          jobs[i]['can_modal'] = true
        }
      }

      if (jobs.length < pageSize.value) {
        isOver.value = true
      }
      if (JSON.stringify(jobs) === JSON.stringify(finishedJobs.value)) {
        loading.value = false
        return
      }

      if (page.value === 1) {
        finishedJobs.value = jobs
      } else {
        finishedJobs.value = finishedJobs.value.concat(jobs)
      }
    })
    .catch((e) => {
      ElMessage.error('获取任务失败：' + e.message)
      loading.value = false
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

// 创建绘图任务
const isGenerating = ref(false)
const generate = () => {
  if (isGenerating.value) {
    return
  }

  if (params.value.prompt === '' && params.value.task_type === 'image') {
    focusPromptInput()
    return ElMessage.error('请输入绘画提示词！')
  }
  if (params.value.model.indexOf('niji') !== -1 && params.value.raw) {
    return ElMessage.error('动漫模型不允许启用原始模式')
  }
  if (imgList.value.length !== 2 && params.value.task_type === 'swapFace') {
    return ElMessage.error('换脸操作需要上传两张图片')
  }

  const regex = /(^|\s)--ar\s+(\d+:\d+)/
  const match = regex.exec(params.value.prompt)
  if (match) {
    params.value.rate = match[2]
  }

  params.value.session_id = getSessionId()
  params.value.img_arr = imgList.value
  isGenerating.value = true
  const deductPower =
    params.value.task_type === 'blend'
      ? mjBlendPower.value
      : params.value.task_type === 'swapFace'
        ? mjSwapFacePower.value
        : mjPower.value
  httpPost('/api/mj/image', params.value)
    .then(() => {
      ElMessage.success('绘画任务推送成功，请耐心等待任务执行...')
      power.value -= deductPower
      taskPulling.value = true
      runningJobs.value.push({
        status: 'pending',
        progress: 0,
      })
      isOver.value = false
    })
    .catch((e) => {
      ElMessage.error('任务推送失败：' + e.message)
    })
    .finally(() => {
      isGenerating.value = false
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
  })
    .then(() => {
      ElMessage.success('任务推送成功，请耐心等待任务执行...')
      power.value -= mjUpscalePower.value
      taskPulling.value = true
      runningJobs.value.push({
        progress: 0,
      })
    })
    .catch((e) => {
      ElMessage.error('任务推送失败：' + e.message)
    })
}

// 局部重绘弹窗
const modalVisible = ref(false)
const modalItem = ref(null)
const modalPrompt = ref('')
const modalCanvasRef = ref(null)
const modalMaskCanvasRef = ref(null)
const modalImgRef = ref(null)
const modalCanvasWrapRef = ref(null)
const modalCanvasSize = ref({ w: 0, h: 0 }) // 显示层尺寸（与当前显示图一致）
const modalNaturalSize = ref({ w: 0, h: 0 }) // 原图尺寸，蒙版 canvas 使用
const modalBrushSize = ref(16)
const modalSubmitting = ref(false)
const modalDrawing = ref(false)
let modalCtx = null // 显示层：半透明绘制，底图可见
let modalMaskCtx = null // 导出层：黑底白字，用于 API

const openModalDialog = (item) => {
  modalItem.value = item
  modalPrompt.value = ''
  modalCanvasSize.value = { w: 0, h: 0 }
  modalNaturalSize.value = { w: 0, h: 0 }
  modalVisible.value = true
}

const onModalImageLoad = () => {
  nextTick(() => {
    const img = modalImgRef.value
    if (!img || !modalItem.value) return
    const displayW = img.offsetWidth || Math.min(img.naturalWidth || 400, 960)
    const displayH = img.offsetHeight || Math.min(img.naturalHeight || 400, 720)
    const naturalW = img.naturalWidth || displayW
    const naturalH = img.naturalHeight || displayH
    if (displayW <= 0 || displayH <= 0) return
    modalCanvasSize.value = { w: displayW, h: displayH }
    modalNaturalSize.value = { w: naturalW, h: naturalH }
    nextTick(() => {
      const canvas = modalCanvasRef.value
      const maskCanvas = modalMaskCanvasRef.value
      if (!canvas || !maskCanvas) return
      modalCtx = canvas.getContext('2d')
      modalMaskCtx = maskCanvas.getContext('2d')
      if (!modalCtx || !modalMaskCtx) return
      // 显示层：与显示图同尺寸，不填黑，绘制时用半透明红色
      modalCtx.lineCap = 'round'
      modalCtx.lineJoin = 'round'
      // 导出层：原图尺寸 naturalW×naturalH，黑底，后续绘制白色蒙版
      modalMaskCtx.fillStyle = '#000000'
      modalMaskCtx.fillRect(0, 0, naturalW, naturalH)
      modalMaskCtx.strokeStyle = '#ffffff'
      modalMaskCtx.lineCap = 'round'
      modalMaskCtx.lineJoin = 'round'
    })
  })
}

// 将显示坐标换算到原图坐标，用于在蒙版 canvas 上绘制
const displayToNatural = (displayX, displayY) => {
  const ds = modalCanvasSize.value
  const ns = modalNaturalSize.value
  if (ds.w <= 0 || ds.h <= 0) return { x: 0, y: 0, scale: 1 }
  const scaleX = ns.w / ds.w
  const scaleY = ns.h / ds.h
  return {
    x: displayX * scaleX,
    y: displayY * scaleY,
    scale: Math.max(scaleX, scaleY),
  }
}

const modalDrawStart = (e) => {
  if (!modalCtx || !modalMaskCtx) return
  modalDrawing.value = true
  const size = modalBrushSize.value
  const x = e.offsetX
  const y = e.offsetY
  const { x: maskX, y: maskY, scale } = displayToNatural(x, y)
  const maskLineWidth = size * scale
  // 显示层：半透明红色，方便对照底图
  modalCtx.strokeStyle = 'rgba(255, 80, 80, 0.55)'
  modalCtx.lineWidth = size
  modalCtx.beginPath()
  modalCtx.moveTo(x, y)
  // 导出层：原图坐标系，白色
  modalMaskCtx.strokeStyle = '#ffffff'
  modalMaskCtx.lineWidth = maskLineWidth
  modalMaskCtx.beginPath()
  modalMaskCtx.moveTo(maskX, maskY)
}

const modalDrawMove = (e) => {
  if (!modalDrawing.value || !modalCtx || !modalMaskCtx) return
  const size = modalBrushSize.value
  const x = e.offsetX
  const y = e.offsetY
  const { x: maskX, y: maskY, scale } = displayToNatural(x, y)
  const maskLineWidth = size * scale
  modalCtx.lineWidth = size
  modalCtx.lineTo(x, y)
  modalCtx.stroke()
  modalMaskCtx.lineWidth = maskLineWidth
  modalMaskCtx.lineTo(maskX, maskY)
  modalMaskCtx.stroke()
}

const modalDrawEnd = () => {
  modalDrawing.value = false
}

const modalClearMask = () => {
  const canvas = modalCanvasRef.value
  const maskCanvas = modalMaskCanvasRef.value
  if (!canvas || !maskCanvas || !modalCtx || !modalMaskCtx) return
  modalCtx.clearRect(0, 0, canvas.width, canvas.height)
  modalMaskCtx.fillStyle = '#000000'
  modalMaskCtx.fillRect(0, 0, maskCanvas.width, maskCanvas.height)
}

const getModalMaskBase64 = () => {
  const maskCanvas = modalMaskCanvasRef.value
  if (!maskCanvas || modalNaturalSize.value.w === 0) return ''
  const dataUrl = maskCanvas.toDataURL('image/png')
  return dataUrl.replace(/^data:image\/\w+;base64,/, '')
}

const submitModal = () => {
  if (!modalPrompt.value.trim()) {
    ElMessage.warning('请填写重绘提示词')
    return
  }
  if (!modalItem.value) return
  const taskId = modalItem.value.task_id || modalItem.value.message_id
  if (!taskId || !modalItem.value.channel_id) {
    ElMessage.warning('缺少原图信息，无法提交局部重绘')
    return
  }
  modalSubmitting.value = true
  const maskBase64 = getModalMaskBase64()
  const body = {
    task_id: taskId,
    channel_id: modalItem.value.channel_id,
    prompt: modalPrompt.value.trim(),
  }
  if (maskBase64) body.mask_base64 = maskBase64
  httpPost('/api/mj/modal', body)
    .then(() => {
      ElMessage.success('任务推送成功，请耐心等待任务执行...')
      power.value -= mjModalPower.value
      taskPulling.value = true
      runningJobs.value.push({ progress: 0 })
      modalVisible.value = false
    })
    .catch((e) => {
      ElMessage.error('任务推送失败：' + e.message)
    })
    .finally(() => {
      modalSubmitting.value = false
    })
}

const onModalClosed = () => {
  modalItem.value = null
  modalPrompt.value = ''
  modalCtx = null
  modalMaskCtx = null
}

// 任务详情弹窗
const detailDialogVisible = ref(false)
const currentDetail = ref(null)

/** 解析 task_info（MjTask JSON） */
function parseMjTaskInfo(raw) {
  if (raw == null || raw === '') return null
  if (typeof raw === 'object') return raw
  if (typeof raw !== 'string') return null
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

const detailTaskPayload = computed(() => parseMjTaskInfo(currentDetail.value?.task_info))

const detailTaskImages = computed(() => {
  const arr = detailTaskPayload.value?.img_arr
  return Array.isArray(arr) ? arr.filter((u) => u && String(u).trim()) : []
})

const detailHasMask = computed(() => {
  const m = detailTaskPayload.value?.mask_base64
  return typeof m === 'string' && m.length > 0
})

function mjDetailTypeLabel(type) {
  const map = {
    image: '绘图',
    upscale: '放大',
    variation: '变换',
    blend: '融图',
    swapFace: '换脸',
    modal: '局部重绘',
  }
  return map[type] || type || '-'
}

const showDetail = (item) => {
  currentDetail.value = item
  detailDialogVisible.value = true
}
const copyPrompt = (text) => {
  if (!text) return
  navigator.clipboard
    .writeText(text)
    .then(() => {
      ElMessage.success('复制成功')
    })
    .catch(() => {
      ElMessage.error('复制失败')
    })
}

const removeImage = (item) => {
  ElMessageBox.confirm('此操作将会删除任务和图片，继续操作码?', '删除提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      httpGet('/api/mj/remove', { id: item.id, user_id: item.user_id })
        .then(() => {
          ElMessage.success('任务删除成功')
          page.value = 0
          isOver.value = false
          fetchFinishJobs()
        })
        .catch((e) => {
          ElMessage.error('任务删除失败：' + e.message)
        })
    })
    .catch(() => {})
}

// 发布图片到作品墙
const publishImage = (item, action) => {
  let text = '图片发布'
  if (action === false) {
    text = '取消发布'
  }
  httpGet('/api/mj/publish', {
    id: item.id,
    action: action,
    user_id: item.user_id,
  })
    .then(() => {
      ElMessage.success(text + '成功')
      item.publish = action
      page.value = 0
      isOver.value = false
    })
    .catch((e) => {
      ElMessage.error(text + '失败：' + e.message)
    })
}

const previewImg = (item) => {
  previewURL.value = item.img_url
}

const closePreview = () => {
  previewURL.value = ''
}

// 切换菜单
const tabChange = (tab) => {
  if (tab === 'txt2img' || tab === 'img2img' || tab === 'cref') {
    params.value.task_type = 'image'
  } else {
    params.value.task_type = tab
  }
}

const promptGenerating = ref(false)
const generatePrompt = () => {
  if (params.value.prompt === '') {
    return showMessageError('请输入原始提示词')
  }
  promptGenerating.value = true
  httpPost('/api/prompt/image', { prompt: params.value.prompt })
    .then((res) => {
      params.value.prompt = res.data
      promptGenerating.value = false
    })
    .catch((e) => {
      showMessageError('生成提示词失败：' + e.message)
      promptGenerating.value = false
    })
}
</script>

<style lang="scss" scoped>
@use '../assets/css/image-mj.scss' as *;
@use '../assets/css/custom-scroll.scss' as *;

.modal-inpaint-dialog {
  :deep(.el-dialog) {
    width: min(960px, 90vw);
  }
}
.modal-inpaint-wrap {
  max-width: 100%;
}

/* 一致性参数：两列布局，每列上 label 下上传 */
.cref-two-cols {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  max-width: 520px;
}

.cref-col {
  display: flex;
  flex-direction: column;
  gap: 8px;

  .cref-label {
    font-size: 14px;
    font-weight: 500;
    color: var(--el-text-color-regular, #606266);
    margin: 0;
  }

  .cref-upload-inner {
    width: 100%;
  }
}

@media (max-width: 640px) {
  .cref-two-cols {
    grid-template-columns: 1fr;
  }
}

/* 任务详情弹窗 */
.mj-detail-dialog {
  :deep(.el-dialog__body) {
    padding-top: 6px;
  }
}

.mj-detail-body {
  max-height: min(72vh, 640px);
  overflow-y: auto;
}

.mj-detail-copy-row {
  display: flex;
  align-items: center;
  gap: 8px;

  &.align-start {
    align-items: flex-start;

    .mj-detail-copy-ico {
      margin-top: 4px;
    }
  }
}

.mj-detail-text {
  font-size: 13px;
  line-height: 1.55;
}

.mj-detail-copy-ico {
  flex-shrink: 0;
  cursor: pointer;
  color: var(--el-text-color-secondary);
  opacity: 0.85;
  transition:
    opacity 0.15s ease,
    color 0.15s ease;

  &:hover {
    opacity: 1;
    color: var(--el-color-primary);
  }
}

.mj-detail-refimgs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.mj-detail-refimg {
  width: 88px;
  height: 88px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid var(--el-border-color-lighter);
}

.mj-detail-result {
  width: 240px;
  max-width: 100%;
  height: 240px;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid var(--el-border-color-lighter);
}

.mj-detail-raw {
  display: block;
  margin-top: 8px;
  padding: 8px 10px;
  font-size: 11px;
  line-height: 1.45;
  max-height: 180px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
  border-radius: 8px;
  background: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color-lighter);
}
</style>
