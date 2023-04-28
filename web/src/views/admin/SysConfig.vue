<template>
  <div class="system-config">
    <el-form :model="form" label-width="120px">
      <el-form-item label="应用标题">
        <el-input v-model="form['title']"/>
      </el-form-item>
      <el-form-item label="控制台标题">
        <el-input v-model="form['console_title']"/>
      </el-form-item>
      <el-form-item label="代理地址">
        <el-input v-model="form['console_title']" placeholder="多个地址之间用逗号隔开"/>
      </el-form-item>

      <el-divider content-position="center">聊天设置</el-divider>
      <el-row>
        <el-col :span="12">
          <div class="grid-content">
            <el-form-item label="GPT模型">
              <el-input v-model="form['console_title']" placeholder="目前只支持 GPT-3 和 GPT-3.5"/>
            </el-form-item>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="grid-content">
            <el-form-item label="模型温度">
              <el-input v-model="form['console_title']" placeholder="0-1之间的小数"/>
            </el-form-item>
          </div>
        </el-col>
      </el-row>

      <el-form-item>
        <el-button type="primary" @click="save">Create</el-button>
      </el-form-item>
    </el-form>

  </div>
</template>

<script>
import {defineComponent} from "vue";

export default defineComponent({
  name: 'SysConfig',
  data() {
    return {
      title: "系统管理",
      form: {}
    }
  },
  methods: {
    save: function () {

    }
  }
})
</script>

<style lang="stylus" scoped>
.system-config {

}
</style>