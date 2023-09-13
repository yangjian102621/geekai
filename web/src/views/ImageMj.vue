<template>
  <div class="page-mj">
    <div class="inner">
      <div class="mj-box">
        <h2>MidJourney 创作中心</h2>

        <div class="mj-params">
          <div class="param-line">
            <span>图片比例：</span>
            <el-tooltip
                effect="light"
                content="生成图片的尺寸比例"
                placement="right"
            >
              <el-icon>
                <InfoFilled/>
              </el-icon>
            </el-tooltip>
          </div>

          <div class="param-line">
            <el-row :gutter="10">
              <el-col :span="8" v-for="item in rates" :key="item.value">
                <div :class="item.value === params.rate?'grid-content active':'grid-content'"
                     @click="changeRate(item)">
                  <div :class="'shape '+item.css"></div>
                  <div class="text">{{ item.text }}</div>
                </div>
              </el-col>
            </el-row>
          </div>

          <div class="param-line">
            <span>模型选择：</span>
            <el-tooltip
                effect="light"
                content="MJ: 偏真实通用模型 <br/>NIJI: 偏动漫风格、适用于二次元模型"
                raw-content
                placement="right"
            >
              <el-icon>
                <InfoFilled/>
              </el-icon>
            </el-tooltip>
          </div>
          <div class="param-line">
            <el-row :gutter="10">
              <el-col :span="12" v-for="item in models" :key="item.value">
                <div :class="item.value === params.model?'grid-content active':'grid-content'"
                     @click="changeModel(item)">
                  <div class="img">
                    <el-image src="https://ai.2021it.com/assets/mj-8c02cbcc.png"></el-image>
                  </div>
                  <div class="text">{{ item.text }}</div>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </div>
      <div class="task-list-box" :style="{ height: listBoxHeight + 'px' }">
        <h2>任务列表</h2>

        <h2>创作记录</h2>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref} from "vue"
import {InfoFilled} from "@element-plus/icons-vue";

const listBoxHeight = window.innerHeight - 20
const rates = [
  {css: "horizontal", value: "16:9", text: "横图"},
  {css: "square", value: "1:1", text: "方图"},
  {css: "vertical", value: "9:16", text: "竖图"},
]
const models = [
  {text: "标准模型", value: "--v 5.2"},
  {text: "动漫模型", value: "--niji 5"},
]
const params = ref({
  rate: rates[0].value,
  model: models[0].value
})
const changeRate = (item) => {
  params.value.rate = item.value
}
const changeModel = (item) => {
  params.value.model = item.value
}
</script>

<style lang="stylus" scoped>
.page-mj {
  background-color: #282c34;

  .inner {
    display: flex;

    .mj-box {
      margin 10px
      background-color #262626
      border 1px solid #454545
      min-width 300px
      max-width 300px
      padding 10px
      border-radius 10px
      color #ffffff;
      font-size 14px

      h2 {
        font-weight: bold;
        font-size 20px
        text-align center
        color #47fff1
      }

      .mj-params {
        margin-top 10px

        .param-line {
          padding 6px 10px 6px 10px

          .el-icon {
            position relative
            top 3px
          }

          .grid-content {
            background-color #383838
            border-radius 5px
            padding 8px 14px
            display flex
            cursor pointer

            &:hover {
              background-color #585858
            }

            .shape {
              width 16px
              height 16px
              margin-right 5px
              border 1px solid #C4C4C4
              border-radius 3px
            }

            .shape.vertical {
              width 12px
              height 20px
            }

            .shape.horizontal {
              height 12px
              width 20px
              position relative
              top 3px
            }

            .img {
              
            }
          }


          .grid-content.active {
            color #47fff1
            background-color #585858

            .shape {
              border 1px solid #47fff1
            }
          }
        }
      }
    }


    .task-list-box {
      width 100%
      padding 10px
      color #ffffff
      overflow-x hidden
      overflow-y auto

      h2 {
        font-size 20px
      }
    }
  }

}
</style>
