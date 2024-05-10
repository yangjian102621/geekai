<template>
  <el-form label-width="150px" label-position="right">
    <el-tabs type="border-card">
      <el-tab-pane label="MJ-PLUS">
        <div v-if="mjPlusConfigs">
          <el-form-item label="网站标题">
            <el-input v-model="sdConfigs"/>
          </el-form-item>
        </div>
        <el-empty v-else></el-empty>
        <el-row style="justify-content: center">
          <el-button round>
            <el-icon><Plus /></el-icon>
            <span>新增配置</span>
          </el-button>
        </el-row>

      </el-tab-pane>

      <el-tab-pane label="MJ-PROXY">
        <div v-if="mjProxyConfigs">
          <el-form-item label="注册赠送算力">
            <el-input v-model.number="sdConfigs" placeholder="新用户注册赠送算力"/>
          </el-form-item>
        </div>
        <el-empty v-else />
      </el-tab-pane>

      <el-tab-pane label="Stable-Diffusion">
        <el-form-item label="注册赠送算力">
          <el-input v-model.number="sdConfigs" placeholder="新用户注册赠送算力"/>
        </el-form-item>
      </el-tab-pane>
    </el-tabs>

    <div style="padding: 10px;">
      <el-form-item>
        <el-button type="primary" @click="save('system')">保存</el-button>
      </el-form-item>
    </div>
  </el-form>
</template>

<script setup>
import {ref} from "vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {Plus} from "@element-plus/icons-vue";

// 变量定义
const sdConfigs = ref([])
const mjPlusConfigs = ref([])
const mjProxyConfigs = ref([])

httpGet("/api/admin/config/get/draw").then(res => {
  sdConfigs.value = res.data.sd
  mjPlusConfigs.value = res.data.mj_plus
  mjProxyConfigs.value = res.data.mj_proxy
}).catch(e =>{
  ElMessage.error("获取AI绘画配置失败："+e.message)
})
</script>

<style lang="stylus" scoped>
.menu {

  .opt-box {
    padding-bottom: 10px;
    display flex;
    justify-content flex-end

    .el-icon {
      margin-right: 5px;
    }
  }

  .menu-icon {
    width 36px
    height 36px
  }

  .el-select {
    width: 100%
  }

}
</style>