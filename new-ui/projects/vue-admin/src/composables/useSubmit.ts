import { ref, reactive, unref } from "vue";
import { Message } from "@arco-design/web-vue";
import type { BaseResponse } from "@gpt-vue/packages/type";
function useSubmit<T extends Record<string, any> = Record<string, any>, R = any>(defaultData?: T) {
  const formRef = ref();
  const formData = reactive<T | Record<string, any>>({ ...defaultData ?? {} });
  const submitting = ref(false);

  const handleSubmit = async (api: (params?: any) => Promise<BaseResponse<R>>, params) => {
    submitting.value = true;
    try {
      const hasError = await formRef.value?.validate();
      if (!hasError) {
        const { data, message } = await api({ ...formData ?? {}, ...unref(params) });
        message && Message.success(message);
        return Promise.resolve({ formData, data });
      }
      return Promise.reject(false);
    } catch (err) {
      return Promise.reject(err);
    } finally {
      submitting.value = false;
    }
  };

  return { formRef, formData, handleSubmit, submitting };
}

export default useSubmit;
