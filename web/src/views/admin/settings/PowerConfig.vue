<template>
  <div class="power-config form">
    <div class="container">
      <el-form :model="firstUseGift" label-position="top" class="px-3 py-5" v-loading="giftLoading">
        <el-form-item label="首次客户端登录赠送订阅">
          <el-switch v-model="firstUseGift.enabled" />
        </el-form-item>
        <el-form-item label="赠送套餐 ID">
          <el-input-number v-model="firstUseGift.plan_id" :min="0" :max="2147483647" :precision="0" />
        </el-form-item>
        <p class="text-sm text-gray-500 mb-4">每个账号仅赠送一次；修改套餐仅影响后续领取。</p>
        <el-button type="primary" :loading="giftSaving" :disabled="!giftLoaded" @click="saveFirstUseGift">
          保存订阅赠送配置
        </el-button>
      </el-form>
      <el-divider />
      <el-form :model="rewardQuota" label-position="top" class="px-3 py-5" v-loading="quotaLoading">
        <el-form-item label="签到、邀请积分追加 Token 额度">
          <el-switch v-model="rewardQuota.enabled" />
        </el-form-item>
        <el-form-item label="每获得 1 积分，追加的平台额度">
          <el-input-number v-model="rewardQuota.quota_per_point" :min="1" :max="1000000000" :precision="0" />
        </el-form-item>
        <p class="text-sm text-gray-500 mb-4">积分继续保留。配置仅影响之后的新奖励；未绑定 Token 账号的奖励会在绑定后发放。</p>
        <el-button type="primary" :loading="quotaSaving" :disabled="!quotaLoaded" @click="saveRewardQuota">
          保存额度换算配置
        </el-button>
      </el-form>
      <el-divider />
      <el-form
        :model="system"
        label-position="top"
        ref="systemFormRef"
        class="px-3 py-5"
        :rules="rules"
      >
        <div>
          <el-form-item label="注册赠送积分" prop="init_power">
            <el-input v-model.number="system['init_power']" placeholder="新用户注册赠送积分" />
          </el-form-item>
          <el-form-item label="邀请赠送积分" prop="invite_power">
            <el-input
              v-model.number="system['invite_power']"
              placeholder="邀请新用户注册赠送积分"
            />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                签到赠送积分
                <el-tooltip effect="dark" content="每日签到赠送积分" raw-content placement="right">
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
                MJ绘图积分
                <el-tooltip
                  effect="dark"
                  content="使用MidJourney画一张图消耗积分"
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
                MJ 放大（Upscale）积分
                <el-tooltip effect="dark" content="MJ 放大、变换（V1-V4）操作消耗积分" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_upscale_power']" placeholder="未配置时使用 MJ 操作积分" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ 混合（Blend）积分
                <el-tooltip effect="dark" content="MJ 融图操作消耗积分" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_blend_power']" placeholder="未配置时使用 MJ 操作积分" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ 换脸积分
                <el-tooltip effect="dark" content="MJ 换脸操作消耗积分" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_swap_face_power']" placeholder="未配置时使用 MJ 操作积分" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ 局部重绘积分
                <el-tooltip effect="dark" content="MJ 局部重绘（Inpaint）操作消耗积分" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_modal_power']" placeholder="未配置时使用 MJ 操作积分" />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="label-title">
                MJ 操作积分（回退默认值）
                <el-tooltip effect="dark" content="上述 MJ 分项未配置时使用的默认值" raw-content placement="right">
                  <el-icon><InfoFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input v-model.number="system['mj_action_power']" placeholder="" />
          </el-form-item>

          <el-form-item label="Suno 积分" prop="suno_power">
            <el-input
              v-model.number="system['suno_power']"
              placeholder="使用 Suno 生成一首音乐消耗积分"
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
const firstUseGift = ref({ enabled: false, plan_id: 0 })
const giftLoading = ref(true)
const giftLoaded = ref(false)
const giftSaving = ref(false)
const rewardQuota = ref({ enabled: false, quota_per_point: 1000 })
const quotaLoading = ref(true)
const quotaLoaded = ref(false)
const quotaSaving = ref(false)

onMounted(() => {
  httpGet('/api/admin/config/get?key=reward_quota')
    .then((res) => {
      rewardQuota.value = { enabled: false, quota_per_point: 1000, ...res.data }
      quotaLoaded.value = true
    })
    .catch((e) => ElMessage.error('加载额度换算配置失败：' + e.message))
    .finally(() => { quotaLoading.value = false })
  httpGet('/api/admin/config/get?key=first_use_gift')
    .then((res) => {
      firstUseGift.value = { enabled: false, plan_id: 0, ...res.data }
      giftLoaded.value = true
    })
    .catch((e) => ElMessage.error('加载订阅赠送配置失败：' + e.message))
    .finally(() => { giftLoading.value = false })
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

const saveRewardQuota = async () => {
  if (!Number.isInteger(rewardQuota.value.quota_per_point) || rewardQuota.value.quota_per_point < 1 || rewardQuota.value.quota_per_point > 1000000000) {
    ElMessage.error('每积分额度须为 1 至 1000000000 的整数')
    return
  }
  quotaSaving.value = true
  try {
    await httpPost('/api/admin/config/update/reward_quota', rewardQuota.value)
    ElMessage.success('额度换算配置已保存')
  } catch (e) {
    ElMessage.error('保存额度换算配置失败：' + e.message)
  } finally {
    quotaSaving.value = false
  }
}

const saveFirstUseGift = async () => {
  if (!Number.isInteger(firstUseGift.value.plan_id) || firstUseGift.value.plan_id < 0 ||
      (firstUseGift.value.enabled && firstUseGift.value.plan_id === 0)) {
    ElMessage.error('启用赠送时，请填写有效的套餐 ID')
    return
  }
  giftSaving.value = true
  try {
    await httpPost('/api/admin/config/update/first_use_gift', firstUseGift.value)
    ElMessage.success('订阅赠送配置已保存')
  } catch (e) {
    ElMessage.error('保存订阅赠送配置失败：' + e.message)
  } finally {
    giftSaving.value = false
  }
}

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
