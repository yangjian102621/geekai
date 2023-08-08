<template>
  <div class="welcome">
    <div class="container">
      <h1 class="title">ChatGPT-PLUS</h1>

      <el-row :gutter="20">
        <el-col :span="8">
          <div class="grid-content">
            <div class="item-title">
              <div><i class="iconfont icon-quick-start"></i></div>
              <div>小试牛刀</div>
            </div>

            <div class="item-list">
              <ul>
                <li v-for="item in samples"><a @click="send(item)">{{ item }}</a></li>
              </ul>
            </div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="grid-content">
            <div class="item-title">
              <div><i class="iconfont icon-plugin"></i></div>
              <div>插件增强</div>
            </div>

            <div class="item-list">
              <ul>
                <li v-for="item in plugins"><a @click="send(item.value)">{{ item.text }}</a></li>
              </ul>
            </div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="grid-content">
            <div class="item-title">
              <div><i class="iconfont icon-control"></i></div>
              <div>能力扩展</div>
            </div>

            <div class="item-list">
              <ul>
                <li v-for="item in capabilities">{{ item }}</li>
              </ul>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script setup>

import {ref} from "vue";

const samples = ref([
  "用小学生都能听懂的术语解释什么是量子纠缠",
  "能给一位6岁男孩的生日会提供一些创造性的建议吗？",
  "如何用 Go 语言实现支持代理 Http client 请求?"
])

const plugins = ref([
  {
    value: "今日早报",
    text: "今日早报：获取当天全球的热门新闻事件列表"
  },
  {
    value: "微博热搜",
    text: "微博热搜：新浪微博热搜榜，微博当日热搜榜单"
  },
  {
    value: "今日头条",
    text: "今日头条：给用户推荐当天的头条新闻，周榜热文"
  }
])

const capabilities = ref([
  "轻松扮演翻译专家，程序员，AI 女友，文案高手...",
  "国产大语言模型支持，GLM2 模型接入中",
  "Midjourney, Stable Diffusion AI 绘画支持"
])

const emits = defineEmits(['send']);
const send = (text) => {
  emits('send', text)
}
</script>
<style scoped lang="stylus">
.welcome {
  text-align center
  display flex
  justify-content center
  margin-top 8vh

  .container {
    max-width 768px;
    width 100%

    .title {
      font-size: 2.25rem
      line-height: 2.5rem
      font-weight 600
      margin-bottom: 4rem
    }

    .grid-content {
      .item-title {
        div {
          padding 6px 10px;

          .iconfont {
            font-size 24px;
          }
        }
      }

      .item-list {
        ul {
          padding 10px;

          li {
            font-size 14px;
            padding .75rem
            border-radius 5px;
            background-color: rgba(247, 247, 248, 1);

            line-height 1.5
            color #666666

            a {
              cursor pointer
              display block
              width 100%
            }
            margin-top 10px;
          }
        }
      }
    }
  }
}
</style>