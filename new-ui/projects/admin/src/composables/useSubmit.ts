import { ref, reactive, unref } from "vue";
import type { BaseResponse } from "@chatgpt-plus/packages/type";
function useSubmit<T extends Record<string, any> = Record<string, any>, R = any>(defaultData?: T) {
  const formRef = ref();
  const formData = reactive<T | Record<string, any>>({ ...defaultData ?? {} });
  const submitting = ref(false);

  const handleSubmit = async (api: (params?: any) => Promise<BaseResponse<R>>, params) => {
    submitting.value = true;
    try {
      const hasError = await formRef.value?.validate();
      if (hasError) return Promise.reject({ validateErrors: hasError });

      const { data, code, message } = await api({ ...formData ?? {}, ...unref(params) });
      if (code) {
        return Promise.reject({ requestErrors: message })
      }

      return Promise.resolve({ formData, data });
    } catch (err) {
      return Promise.reject({ errors: err });
    } finally {
      submitting.value = false;
    }
  };

  return { formRef, formData, handleSubmit, submitting };
}

export default useSubmit;
