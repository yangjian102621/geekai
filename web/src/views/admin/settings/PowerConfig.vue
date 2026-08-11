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
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ 放大（Upscale）算力
                <el-tooltip effect="dark" content="MJ 放大、变换（V1-V4）操作消耗算力" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_upscale_power']" placeholder="未配置时使用 MJ 操作算力" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ 混合（Blend）算力
                <el-tooltip effect="dark" content="MJ 融图操作消耗算力" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_blend_power']" placeholder="未配置时使用 MJ 操作算力" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ 换脸算力
                <el-tooltip effect="dark" content="MJ 换脸操作消耗算力" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_swap_face_power']" placeholder="未配置时使用 MJ 操作算力" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ 局部重绘算力
                <el-tooltip effect="dark" content="MJ 局部重绘（Inpaint）操作消耗算力" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_modal_power']" placeholder="未配置时使用 MJ 操作算力" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ 操作算力（回退默认值）
                <el-tooltip effect="dark" content="上述 MJ 分项未配置时使用的默认值" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_action_power']" placeholder="" />
          </el-form-item>

          <el-form-item label="Suno 算力" prop="suno_power">
            <el-input
              v-model.number="system['suno_power']"
              placeholder="使用 Suno 生成一首音乐消耗算力"
            />
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
        mj_action_power: system.value.mj_action_power,
        mj_upscale_power: system.value.mj_upscale_power,
        mj_blend_power: system.value.mj_blend_power,
        mj_swap_face_power: system.value.mj_swap_face_power,
        mj_modal_power: system.value.mj_modal_power,
        suno_power: system.value.suno_power,
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
