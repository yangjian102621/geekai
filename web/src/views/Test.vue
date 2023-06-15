<template>
  <div class="role-list">


    <el-form :model="form1" label-width="120px" ref="formRef" :rules="rules">
      <el-form-item label="角色名称：" prop="name">
        <el-input
            v-model="form1.name"
            autocomplete="off"
        />
      </el-form-item>

      <el-form-item label="角色标志：" prop="key">
        <el-input
            v-model="form1.key"
            autocomplete="off"
        />
      </el-form-item>

      <el-form-item label="角色图标：" prop="icon">
        <el-input
            v-model="form1.icon"
            autocomplete="off"
        />
      </el-form-item>

      <el-form-item label="打招呼信息：" prop="hello_msg">
        <el-input
            v-model="form1.hello_msg"
            autocomplete="off"
        />
      </el-form-item>

      <el-form-item label="上下文信息：" prop="context">
        <template #default>
          <el-table :data="form1.context" :border="childBorder" size="small">
            <el-table-column label="对话角色" width="120">
              <template #default="scope">
                <el-input
                    v-model="scope.row.role"
                    autocomplete="off"
                />
              </template>
            </el-table-column>
            <el-table-column label="对话内容">
              <template #header>
                <div class="context-msg-key">
                  <span>对话内容</span>
                  <span class="fr">
                      <el-button type="primary" @click="addContext" size="small">
                      <el-icon>
                        <Plus/>
                      </el-icon>
                      增加一行
                    </el-button>
                    </span>
                </div>
              </template>

              <template #default="scope">
                <div class="context-msg-content">
                  <el-input
                      v-model="scope.row.content"
                      autocomplete="off"
                  />
                  <span><el-icon @click="removeContext(scope.$index)"><RemoveFilled/></el-icon></span>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </template>
      </el-form-item>

      <el-form-item label="启用状态">
        <el-switch v-model="form1.enable"/>
      </el-form-item>
    </el-form>

    <span class="dialog-footer">
              <el-button @click="showDialog = false">取消</el-button>
              <el-button type="primary" @click="doUpdate">保存</el-button>
            </span>
  </div>
</template>

<script setup>

import {Plus, RemoveFilled} from "@element-plus/icons-vue";
import {reactive, ref} from "vue";
import {httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";

const showDialog = ref(false)
const childBorder = ref(true)
const form1 = ref({context: []})
// const form2 = ref({context: []})
const formRef = ref(null)

const rules = reactive({
  name: [{required: true, message: '请输入用户名', trigger: 'change',}],
  key: [{required: true, message: '请输入角色标识', trigger: 'change',}],
  icon: [{required: true, message: '请输入角色图标', trigger: 'change',}],
  hello_msg: [{required: true, message: '请输入打招呼信息', trigger: 'change',}]
})


// 编辑
const doUpdate = function () {
  formRef.value.validate((valid) => {
    if (valid) {
      showDialog.value = false
      httpPost('/api/admin/chat-roles/set', form1.value).then(() => {
        ElMessage.success('更新角色成功')
        // 更新当前数据行
      }).catch((e) => {
        ElMessage.error('更新角色失败，' + e.message)
      })
    } else {
      return false
    }
  })
}

const addContext = function () {
  form1.value.context.push({role: '', content: ''})
}

const removeContext = function (index) {
  form1.value.context.splice(index, 1);
}

</script>

<style lang="stylus" scoped>
.role-list {
  .opt-box {
    padding-bottom: 10px;

    .el-icon {
      margin-right 5px;
    }
  }

  .context-msg-key {
    .fr {
      float right

      .el-icon {
        margin-right 5px
      }
    }
  }

  .context-msg-content {
    display flex

    .el-icon {
      font-size: 20px;
      margin-top 5px;
      margin-left 5px;
      cursor pointer
    }
  }
}
</style>