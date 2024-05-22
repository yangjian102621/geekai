<template>
  <el-form label-width="150px" label-position="right" class="draw-config">
    <el-tabs type="border-card">
      <el-tab-pane label="MJ-PLUS">
        <div v-if="mjPlusConfigs">
          <div class="config-item" v-for="(v,k) in mjPlusConfigs">
            <el-form-item label="是否启用">
              <el-switch v-model="v['Enabled']"/>
            </el-form-item>
            <el-form-item label="API 地址">
              <el-input v-model="v['ApiURL']" placeholder="API 地址"/>
            </el-form-item>
            <el-form-item label="API 令牌">
              <el-input v-model="v['ApiKey']" placeholder="API KEY"/>
            </el-form-item>
            <el-form-item label="绘画模式">
              <el-select v-model="v['Mode']" placeholder="请选择模式">
                <el-option v-for="item in mjModels" :value="item.value" :label="item.name" :key="item.value">{{
                    item.name
                  }}
                </el-option>
              </el-select>
            </el-form-item>

            <el-button class="remove" type="danger" :icon="Delete" circle @click="removeItem(mjPlusConfigs,k)"/>
          </div>
        </div>
        <el-empty v-else></el-empty>

        <el-row style="justify-content: center; padding: 10px">
          <el-button round @click="addConfig(mjPlusConfigs)">
            <el-icon><Plus /></el-icon>
            <span>新增配置</span>
          </el-button>
        </el-row>

      </el-tab-pane>

      <el-tab-pane label="MJ-PROXY">
        <div v-if="mjProxyConfigs">
          <div class="config-item" v-for="(v,k) in mjProxyConfigs">
            <el-form-item label="是否启用">
              <el-switch v-model="v['Enabled']"/>
            </el-form-item>
            <el-form-item label="API 地址">
              <el-input v-model="v['ApiURL']" placeholder="API 地址"/>
            </el-form-item>
            <el-form-item label="API 令牌">
              <el-input v-model="v['ApiKey']" placeholder="API KEY"/>
            </el-form-item>

            <el-button class="remove" type="danger" :icon="Delete" circle @click="removeItem(mjProxyConfigs,k)"/>
          </div>
        </div>
        <el-empty v-else />

        <el-row style="justify-content: center; padding: 10px">
          <el-button round @click="addConfig(mjProxyConfigs)">
            <el-icon>
              <Plus/>
            </el-icon>
            <span>新增配置</span>
          </el-button>
        </el-row>
      </el-tab-pane>

      <el-tab-pane label="Stable-Diffusion">
        <div v-if="sdConfigs">
          <div class="config-item" v-for="(v,k) in sdConfigs">
            <el-form-item label="是否启用">
              <el-switch v-model="v['Enabled']"/>
            </el-form-item>
            <el-form-item label="API 地址">
              <el-input v-model="v['ApiURL']" placeholder="API 地址"/>
            </el-form-item>
            <el-form-item label="API 令牌">
              <el-input v-model="v['ApiKey']" placeholder="API KEY"/>
            </el-form-item>
            <el-form-item label="模型">
              <el-input v-model="v['Model']" placeholder="绘画模型"/>
            </el-form-item>
            <el-button class="remove" type="danger" :icon="Delete" circle @click="removeItem(sdConfigs,k)"/>
          </div>
        </div>
        <el-empty v-else/>

        <el-row style="justify-content: center; padding: 10px">
          <el-button round @click="addConfig(sdConfigs)">
            <el-icon>
              <Plus/>
            </el-icon>
            <span>新增配置</span>
          </el-button>
        </el-row>
      </el-tab-pane>
    </el-tabs>

    <div style="padding: 10px;">
      <el-form-item>
        <el-button type="primary" @click="saveConfig()">保存</el-button>
      </el-form-item>
    </div>
  </el-form>
</template>

<script setup>
import {ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {Delete, Plus} from "@element-plus/icons-vue";

// 变量定义
const sdConfigs = ref([])
const mjPlusConfigs = ref([])
const mjProxyConfigs = ref([])
const mjModels = ref([
  {name: "慢速（Relax）", value: "relax"},
  {name: "快速（Fast）", value: "fast"},
  {name: "急速（Turbo）", value: "turbo"},
])

httpGet("/api/admin/config/get/app").then(res => {
  sdConfigs.value = res.data.sd
  mjPlusConfigs.value = res.data.mj_plus
  mjProxyConfigs.value = res.data.mj_proxy
}).catch(e =>{
  ElMessage.error("获取配置失败："+e.message)
})

const addConfig = (configs) => {
  configs.push({
    Enabled: true,
    ApiKey: '',
    ApiURL: '',
    Mode: 'fast'
  })
}

const saveConfig = () => {
  httpPost('/api/admin/config/update/draw', {
    'sd': sdConfigs.value,
    'mj_plus': mjPlusConfigs.value,
    'mj_proxy': mjProxyConfigs.value
  }).then(() => {
    ElMessage.success("配置更新成功")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

const removeItem = (arr, k) => {
  arr.splice(k, 1)
}
</script>

<style lang="stylus" scoped>
.draw-config {

  .config-item {
    position relative
    padding 15px 10px 10px 10px
    border 1px solid var(--el-border-color)
    border-radius 10px
    margin-bottom 10px

    .remove {
      position absolute
      right 15px
      top 15px
    }
  }
}
</style>