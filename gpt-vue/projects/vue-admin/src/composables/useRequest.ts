import { ref } from "vue";
import type { Ref } from "vue";
import type { BaseResponse } from "@gpt-vue/packages/type";

type Request<T> = (params?: any) => Promise<BaseResponse<T>>
function useRequest<T>(request: Request<T>) {
  const result = ref<T>()
  const loading = ref(false)

  const requestData = async (params?: any) => {
    try {
      const res = await request(params)
      result.value = res.data
      return Promise.resolve(res)
    } catch (err) {
      return Promise.reject(err)
    } finally {
      loading.value = false
    }
  }

  return [requestData, result, loading] as [Request<T>, Ref<T>, Ref<boolean>]
}

export default useRequest