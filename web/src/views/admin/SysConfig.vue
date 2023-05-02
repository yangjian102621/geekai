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

      <el-row>
        <el-col :span="12">
          <div class="grid-content">
            <el-form-item label="Max Tokens">
              <el-input v-model="form['max_tokens']" placeholder="回复的最大字数，最大4096"/>
            </el-form-item>
          </div>
        </el-col>

        <el-col :span="12">
          <div class="grid-content">
            <el-form-item label="上下文超时">
              <el-input v-model="form['chat_context_expire_time']" placeholder="默认60min"/>
            </el-form-item>
          </div>
        </el-col>
      </el-row>

      <el-form-item label="对话上下文">
        <el-switch v-model="form['enable_context']" />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="save">保存</el-button>
      </el-form-item>
    </el-form>

    <el-divider content-position="center">API KEY 管理</el-divider>
    <el-row class="api-key-box">
      <el-button type="primary" @click="save">
        <el-icon class="el-icon--right"><Plus /></el-icon> 新增
      </el-button>
    </el-row>

    <el-row>
      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="key" label="API-KEY" />
        <el-table-column prop="last_used" label="最后使用" width="180">
          <template #default="scope">
            <span>{{scope.row['last_used']}}</span>
            <el-tag>未使用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.$index, scope.row)"
            >Edit</el-button
            >
            <el-button
                size="small"
                type="danger"
                @click="handleDelete(scope.$index, scope.row)"
            >Delete</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </el-row>

  </div>
</template>

<script>
import {defineComponent} from "vue";
import {Plus} from "@element-plus/icons-vue";

export default defineComponent({
  name: 'SysConfig',
  components: {Plus},
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

  .api-key-box {
    padding-bottom: 10px;
    justify-content: end;

    .el-icon--right {
      margin-right 5px;
    }
  }

}
</style>