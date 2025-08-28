<template>
  <div class="power-config form">
    <div class="container">
      <el-form
        :model="system"
        label-position="top"
        ref="systemFormRef"
        class="px-3 py-5"
        :rules="rules"
      >
        <div>
          <el-form-item label="注册赠送算力" prop="init_power">
            <el-input v-model.number="system['init_power']" placeholder="新用户注册赠送算力" />
          </el-form-item>
          <el-form-item label="邀请赠送算力" prop="invite_power">
            <el-input
              v-model.number="system['invite_power']"
              placeholder="邀请新用户注册赠送算力"
            />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                签到赠送算力
                <el-tooltip effect="dark" content="每日签到赠送算力" raw-content placement="right">
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['daily_power']" placeholder="默认值0" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ绘图算力
                <el-tooltip
                  effect="dark"
                  content="使用MidJourney画一张图消耗算力"
                  raw-content
                  placement="right"
                >
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_power']" placeholder="" />
          </el-form-item>

          <el-form-item label="Stable-Diffusion算力" prop="sd_power">
            <el-input
              v-model.number="system['sd_power']"
              placeholder="使用Stable-Diffusion画一张图消耗算力"
            />
          </el-form-item>
          <el-form-item label="Suno 算力" prop="suno_power">
            <el-input
              v-model.number="system['suno_power']"
              placeholder="使用 Suno 生成一首音乐消耗算力"
            />
          </el-form-item>
          <el-form-item label="Luma 算力" prop="luma_power">
            <el-input
              v-model.number="system['luma_power']"
              placeholder="使用 Luma 生成一段视频消耗算力"
            />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                可灵算力
                <el-tooltip
                  effect="dark"
                  content="可灵每个模型价格不一样，具体请参考：https://api.geekai.pro/models"
                  raw-content
                  placement="right"
                >
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-row :gutter="20" v-if="system['keling_powers']">
              <el-col :span="6" v-for="[key] in Object.entries(system['keling_powers'])" :key="key">
                <el-form-item :label="key" label-position="left">
                  <el-input v-model.number="system['keling_powers'][key]" size="small" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                高级语音算力
                <el-tooltip
                  effect="dark"
                  content="使用一次 OpenAI 高级语音对话消耗的算力"
                  raw-content
                  placement="right"
                >
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['advance_voice_power']" placeholder="" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                提示词算力
                <el-tooltip
                  effect="dark"
                  content="生成AI绘图提示词，歌词，视频描述消耗的算力"
                  raw-content
                  placement="right"
                >
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['prompt_power']" placeholder="" />
          </el-form-item>
        </div>

        <div style="padding: 10px">
          <el-form-item>
            <el-button type="primary" @click="save">保存</el-button>
          </el-form-item>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { copyObj } from '@/utils/libs'
import { InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

const system = ref({})
const systemFormRef = ref(null)

onMounted(() => {
  // 加载系统配置
  httpGet('/api/admin/config/get?key=system')
    .then((res) => {
      system.value = res.data
      system.value.keling_powers = system.value.keling_powers || {
        'kling-v1-6_std_5': 240,
        'kling-v1-6_std_10': 480,
        'kling-v1-6_pro_5': 420,
        'kling-v1-6_pro_10': 840,
        'kling-v1-5_std_5': 240,
        'kling-v1-5_std_10': 480,
        'kling-v1-5_pro_5': 420,
        'kling-v1-5_pro_10': 840,
        'kling-v1_std_5': 120,
        'kling-v1_std_10': 240,
        'kling-v1_pro_5': 420,
        'kling-v1_pro_10': 840,
      }
    })
    .catch((e) => {
      ElMessage.error('加载系统配置失败: ' + e.message)
    })
})

const rules = reactive({})

const save = function () {
  systemFormRef.value.validate((valid) => {
    if (valid) {
      httpPost('/api/admin/config/update/power', {
        init_power: system.value.init_power,
        invite_power: system.value.invite_power,
        daily_power: system.value.daily_power,
        mj_power: system.value.mj_power,
        sd_power: system.value.sd_power,
        dall_power: system.value.dall_power,
        suno_power: system.value.suno_power,
        luma_power: system.value.luma_power,
        keling_powers: system.value.keling_powers,
        advance_voice_power: system.value.advance_voice_power,
        prompt_power: system.value.prompt_power,
      })
        .then(() => {
          ElMessage.success('操作成功！')
        })
        .catch((e) => {
          ElMessage.error('操作失败：' + e.message)
        })
    }
  })
}
</script>

<style lang="scss" scoped>
@use '@/assets/css/admin/form.scss' as *;
@use '@/assets/css/main.scss' as *;

.power-config {
  display: flex;
  justify-content: center;
  padding: 20px;
}
</style>
