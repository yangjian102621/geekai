<script lang="ts" setup>
import {onMounted} from "vue";
import {Message} from "@arco-design/web-vue";
import useSubmit from "@/composables/useSubmit";
import useRequest from "@/composables/useRequest";
import {getConfig, modelList, save} from "./api";
import SystemUploader from "./SystemUploader.vue";

const {formRef, formData: system, handleSubmit, submitting} = useSubmit({});

const [getModelOptions, modelOptions, modelOptionsLoading] = useRequest(modelList);

const rules = {
  title: [{required: true, message: "请输入网站标题"}],
  admin_title: [{required: true, message: "请输入控制台标题"}],
};

const handleSave = async () => {
  await handleSubmit(
      () =>
          save({
            key: "system",
            config: system,
          }),
      {}
  );
  Message.success("保存成功");
};

const reload = async () => {
  const {data} = await getConfig({key: "system"});
  data && Object.assign(system, data);
};

onMounted(async () => {
  getModelOptions();
  reload();
});
</script>
<template>
  <a-card :bordered="false">
    <a-form ref="formRef" :model="system" :rules="rules" auto-label-width :disabled="submitting">
      <a-form-item label="网站标题" field="title">
        <a-input v-model="system['title']"/>
      </a-form-item>
      <a-form-item label="控制台标题" field="admin_title">
        <a-input v-model="system['admin_title']"/>
      </a-form-item>
      <a-form-item label="网站Logo" field="logo">
        <SystemUploader v-model="system['logo']" placeholder="推荐图片宽高比为 1:1"/>
      </a-form-item>
      <a-form-item label="注册赠送算力" field="init_power">
        <a-input-number v-model="system['init_power']" placeholder="新用户注册赠送初始算力"/>
      </a-form-item>
      <a-form-item label="邀请用户赠送算力" field="invite_power">
        <a-input-number
            v-model="system['invite_power']"
            placeholder="邀请新用户注册赠送算力"
        />
      </a-form-item>

      <a-form-item label="VIP每月赠送算力" field="vip_month_power">
        <a-input-number v-model="system['vip_month_power']" placeholder="VIP用户每月赠送算力"/>
      </a-form-item>
      <a-form-item label="MJ绘画价格" field="mj_power">
        <a-space>
          <a-input-number v-model="system['mj_power']" placeholder=""/>
          <a-tooltip content="MidJourney 单次绘图消耗多少单位算力" position="right">
            <icon-info-circle-fill size="18"/>
          </a-tooltip>
        </a-space>
      </a-form-item>
      <a-form-item label="SD绘画价格" field="sd_power">
        <a-space>
          <a-input-number v-model="system['sd_power']" placeholder=""/>
          <a-tooltip content="Stable-Diffusion 单次绘图消耗多少单位算力" position="right">
            <icon-info-circle-fill size="18"/>
          </a-tooltip>
        </a-space>
      </a-form-item>
      <a-form-item label="DALL绘画价格" field="dall_power">
        <a-space>
          <a-input-number v-model="system['dall_power']" placeholder=""/>
          <a-tooltip content="DALL-E-3 单次绘图消耗多少单位算力" position="right">
            <icon-info-circle-fill size="18"/>
          </a-tooltip>
        </a-space>
      </a-form-item>
      <a-form-item label="开放注册" field="enabled_register">
        <a-space>
          <a-switch v-model="system['enabled_register']"/>
          <a-tooltip content="关闭注册之后只能通过管理后台添加用户" position="right">
            <icon-info-circle-fill size="18"/>
          </a-tooltip>
        </a-space>
      </a-form-item>
      <a-form-item label="注册方式" field="register_ways">
        <a-checkbox-group v-model="system['register_ways']">
          <a-checkbox value="mobile">手机注册</a-checkbox>
          <a-checkbox value="email">邮箱注册</a-checkbox>
        </a-checkbox-group>
      </a-form-item>
      <a-form-item label="启用众筹功能" field="enabled_reward">
        <a-space>
          <a-switch v-model="system['enabled_reward']"/>
          <a-tooltip content="开启众筹功能允许用户使用个人微信收款码进行收款" position="right">
            <icon-info-circle-fill size="18"/>
          </a-tooltip>
        </a-space>
      </a-form-item>
      <template v-if="system['enabled_reward']">
        <a-form-item label="众筹算力单价" field="power_price">
          <a-input-number v-model="system['power_price']" placeholder="单位算力价格，如1块10个单位算力，那便填写 0.1"/>
        </a-form-item>
        <a-form-item label="收款二维码" field="reward_img">
          <SystemUploader v-model="system['reward_img']" placeholder="众筹收款二维码地址"/>
        </a-form-item>
      </template>
      <a-form-item label="微信客服二维码" field="wechat_card_url">
        <SystemUploader v-model="system['wechat_card_url']" placeholder="微信客服二维码"/>
      </a-form-item>
      <a-form-item label="订单超时时间" field="order_pay_timeout">
        <a-space style="width: 100%">
          <a-input-number
              v-model="system['order_pay_timeout']"
              placeholder="单位：秒"
              style="width: 100%"
          />
          <a-tooltip position="right">
            <icon-info-circle-fill size="18"/>
            <template #content> 系统会定期清理超时未支付的订单<br/>默认值：900秒</template>
          </a-tooltip>
        </a-space>
      </a-form-item>
      <a-form-item label="会员充值说明" field="order_pay_info_text">
        <a-textarea
            v-model="system['order_pay_info_text']"
            :autosize="{ minRows: 3, maxRows: 10 }"
            placeholder="请输入会员充值说明文字，比如介绍会员计划"
        />
      </a-form-item>
      <a-form-item label="默认AI模型" field="default_models">
        <a-space style="width: 100%">
          <a-select
              v-model="system['default_models']"
              multiple
              :filterable="true"
              placeholder="选择AI模型，多选"
              :options="modelOptions"
              :loading="modelOptionsLoading"
              :field-names="{ value: 'value', label: 'name' }"
              style="width: 100%"
          >
          </a-select>
          <a-tooltip content="新用户注册默认开通的 AI 模型" position="right">
            <icon-info-circle-fill size="18"/>
          </a-tooltip>
        </a-space>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" :loading="submitting" @click="handleSave">保存</a-button>
      </a-form-item>
    </a-form>
  </a-card>
</template>
