<template>
  <div>
    <div class="page-mark-map">
      <div class="inner custom-scroll">
        <div class="mark-map-box">
          <h2>思维导图创作中心</h2>

          <div class="mark-map-params" :style="{ height: leftBoxHeight + 'px' }">
            <el-form :model="params" label-width="80px" label-position="left">
              <div class="param-line">
                你的需求？
              </div>
              <div class="param-line">
                <el-input
                    v-model="prompt"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    type="textarea"
                    placeholder="请给AI输入提示词，让AI帮你完善"
                />
              </div>

              <div class="param-line">
                <el-button color="#47fff1" :dark="false" round @click="generate">智能生成思维导图</el-button>
              </div>

              <div class="param-line">
                使用已有内容生成？
              </div>
              <div class="param-line">
                <el-input
                    v-model="prompt"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    type="textarea"
                    placeholder="请用markdown语法输入您想要生成思维导图的内容！"
                />
              </div>

              <div class="param-line">
                <el-button color="#47fff1" :dark="false" round @click="generate">直接生成（免费）</el-button>
              </div>


              <div class="text-info">
               <el-row :gutter="10">
                 <el-col :span="12">
                   <el-tag>每次生成消耗1算力</el-tag>
                 </el-col>
                 <el-col :span="12">
                   <el-tag type="success">当前可用算力：{{ power }}</el-tag>
                 </el-col>
               </el-row>
              </div>

            </el-form>
          </div>
        </div>

        <div class="right-box">
          <h2>思维导图</h2>
          <div class="body">
            <svg ref="svgRef" :style="{ height: leftBoxHeight + 'px' }"/>
          </div>
        </div><!-- end task list box -->
      </div>

    </div>

    <login-dialog :show="showLoginDialog" @hide="showLoginDialog =  false" @success="initData"/>
  </div>
</template>

<script setup>
import LoginDialog from "@/components/LoginDialog.vue";
import {ref, onMounted, onUpdated} from 'vue';
import {Markmap} from 'markmap-view';
import {loadJS, loadCSS} from 'markmap-common';
import {Transformer} from 'markmap-lib';

const leftBoxHeight = ref(window.innerHeight - 105)
const rightBoxHeight = ref(window.innerHeight - 105)

const prompt = ref("")
const text = ref(`# Geek-AI 助手

* 完整的开源系统，前端应用和后台管理系统皆可开箱即用。
* 基于 Websocket 实现，完美的打字机体验。
* 内置了各种预训练好的角色应用,轻松满足你的各种聊天和应用需求。
* 支持 OPenAI，Azure，文心一言，讯飞星火，清华 ChatGLM等多个大语言模型。
* 支持 MidJourney / Stable Diffusion AI 绘画集成，开箱即用。
* 支持使用个人微信二维码作为充值收费的支付渠道，无需企业支付通道。
* 已集成支付宝支付功能，微信支付，支持多种会员套餐和点卡购买功能。
* 集成插件 API 功能，可结合大语言模型的 function 功能开发各种强大的插件。
`)
const showLoginDialog = ref(false)
const isLogin = ref(false)
const power = ref(0)
const transformer = new Transformer();
const {scripts, styles} = transformer.getAssets()
loadCSS(styles);
loadJS(scripts);


const svgRef = ref(null)
const markMap = ref(null)

onMounted(() => {
  markMap.value = Markmap.create(svgRef.value)
  update()
});

const update = () => {
  const {root} = transformer.transform(text.value)
  markMap.value.setData(root)
  markMap.value.fit()
}

onUpdated(update)

window.onresize = () => {
  leftBoxHeight.value = window.innerHeight - 145
}

</script>

<style lang="stylus">
@import "@/assets/css/mark-map.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
