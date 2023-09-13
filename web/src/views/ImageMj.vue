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
              <el-col :span="8">
                <div class="grid-content">
                  <div class="shape vertical"></div>
                  <div class="text">竖图</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="grid-content">
                  <div class="shape square"></div>
                  <div class="text">方图</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="grid-content">
                  <div class="shape horizontal"></div>
                  <div class="text">横图</div>
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
const params = ref({
  rate: '16:9'
})
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
            padding 5px 6px
            display flex
            cursor pointer

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
