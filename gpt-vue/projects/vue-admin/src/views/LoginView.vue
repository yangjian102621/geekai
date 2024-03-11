<script lang="ts" setup>
import { onMounted, reactive } from "vue";
import { useAuthStore } from "@/stores/auth";
import { captcha } from "@/http/login";
import useState from "@/composables/useState";
import useRequest from "@/composables/useRequest";

// 表单
function useFormData() {
  const formData = reactive({
    username: "",
    password: "",
    captcha: "",
  });
  const rules = {
    username: [{ required: true, message: "请输入您的账号" }],
    password: [{ required: true, message: "请输入您的密码" }],
    captcha: [{ required: true, message: "请输入验证码" }],
  };
  const authStore = useAuthStore();
  const [loginRequest, _, submitting] = useRequest(authStore.login);
  return { formData, loginRequest, submitting, rules };
}

// 验证码
function useCaptcha() {
  const captchaImage = reactive({
    pic_path: "",
    captcha_id: "",
  });
  const getCaptchaImage = async () => {
    const { data } = await captcha();
    Object.assign(captchaImage, data);
  };
  onMounted(getCaptchaImage);
  return { captchaImage, getCaptchaImage };
}

// 记住密码
function useRemeberPWD(formData) {
  const storageKey = "r-f";
  const [isRemember, setIsRemember] = useState(false);
  const onIsRememberChange = (v) => {
    if (v) {
      const value = {
        username: formData.username,
        password: formData.password,
      };
      localStorage.setItem(storageKey, JSON.stringify(value));
    } else {
      localStorage.removeItem(storageKey);
    }
    setIsRemember(v);
  };
  onMounted(() => {
    const getter = localStorage.getItem(storageKey);
    if (getter) {
      setIsRemember(true);
      Object.assign(formData, JSON.parse(getter));
    }
  });
  return { isRemember, onIsRememberChange };
}

const { formData, loginRequest, submitting, rules } = useFormData();
const { captchaImage, getCaptchaImage } = useCaptcha();
const { isRemember, onIsRememberChange } = useRemeberPWD(formData);

// 表单提交
async function handleSubmit({ errors }: any) {
  if (errors) return;
  try {
    await loginRequest({
      ...formData,
      captcha_id: captchaImage.captcha_id,
    });
  } catch (err) {
    getCaptchaImage();
  }
}
</script>
<template>
  <div class="bg">
    <div class="content">
      <!-- 左侧图片 -->
      <span class="left">
        <img src="/left-img.png" alt="" style="width: 468px" />
      </span>
      <!-- 表单 -->
      <div class="right-content">
        <div class="form-box">
          <div class="title">ChatGPT Plus Admin</div>
          <a-form
            ref="formRef"
            :model="formData"
            class="form"
            size="medium"
            auto-label-width
            :label-col-props="{ span: 0 }"
            :wrapper-col-props="{ span: 24 }"
            :rules="rules"
            @submit="handleSubmit"
          >
            <a-space direction="vertical" style="width: 100%">
              <a-form-item field="username" label="账号">
                <a-input v-model="formData.username" placeholder="请输入您的账号" class="input" />
              </a-form-item>
              <a-form-item field="password" label="密码">
                <a-input-password
                  v-model="formData.password"
                  placeholder="请输入您的密码"
                  class="input"
                />
              </a-form-item>
              <a-form-item field="captcha" label="验证码">
                <a-input v-model="formData.captcha" placeholder="请输入验证码" class="input">
                  <template #append>
                    <img
                      class="captcha-image"
                      :src="captchaImage.pic_path"
                      alt="验证码"
                      title="点击刷新验证码"
                      @click="getCaptchaImage()"
                    />
                  </template>
                </a-input>
              </a-form-item>
              <a-form-item hide-label>
                <a-checkbox :model-value="isRemember" @change="onIsRememberChange"
                  >记住密码</a-checkbox
                >
              </a-form-item>
            </a-space>
            <a-form-item hide-label>
              <a-button
                :loading="submitting"
                html-type="submit"
                long
                type="primary"
                class="sign-in-btn"
              >
                登录
              </a-button>
            </a-form-item>
          </a-form>
        </div>
      </div>
    </div>
  </div>
</template>
<style lang="less" scoped>
.bg {
  width: 100%;
  height: 100vh;
  background: linear-gradient(133deg, #ffffff 0%, #dde8fe 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: auto;
}

.content {
  width: 1080px;
  height: 557px;
  display: flex;
}

.left {
  width: 540px;
  height: 557px;
  background: #2670fe;
  border-radius: 16px 0 0 16px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.right-content {
  width: 540px;
  height: 557px;
  background: #ffffff;
  border-radius: 0 16px 16px 0;
  display: flex;
  justify-content: center;
}

.title {
  font-size: 24px;
  font-weight: 800;
  color: #333333;
  line-height: 35px;
  letter-spacing: 1px;
  text-align: center;
  margin-bottom: 40px;
}

.form-box {
  display: flex;
  flex-direction: column;
  width: 438px;
  padding: 58px 0;
  height: 100%;
  box-sizing: border-box;
}

.form {
  flex: 1;
  justify-content: space-between;
}

.input {
  border-radius: 10px 10px 10px 10px;
  height: 60px;
  border: 2px solid #e5e6eb;
}

.captcha-image {
  //width: 100%;
  //height: 100%;
  cursor: pointer;
}

.sign-in-btn {
  height: 60px;
  font-weight: 500;
  line-height: 33px;
  font-size: 28px;
  border-radius: 10px 10px 10px 10px;
}
</style>
