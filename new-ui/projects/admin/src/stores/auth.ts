import { defineStore } from 'pinia'
import { Message } from '@arco-design/web-vue'
import { userLogin, userLogout } from '@/http/login'
import router from '@/router'

const defaultState: {
  token: string
  is_super_admin?: boolean;
  permissions?: string[]
} = {
  token: null,
  is_super_admin: false,
  permissions: []
}

export const useAuthStore = defineStore({
  id: Symbol(__AUTH_KEY).toString(),
  state: () => ({ ...defaultState }),
  actions: {
    init() {
      this.$state = JSON.parse(localStorage.getItem(__AUTH_KEY));
    },
    async login(params: any) {
      try {
        const { data } = await userLogin(params)
        if (data) {
          this.$state = data;
          localStorage.setItem(__AUTH_KEY, JSON.stringify(data))
          Message.success('登录成功');
          router.replace({ name: 'home' })
          return Promise.resolve(data)
        }
      } catch (err) {
        return Promise.reject(err)
      }
    },
    async logout() {
      try {
        await userLogout()
        if (this.$state.token) {
          localStorage.removeItem(__AUTH_KEY)
          this.$reset()
        }
        Message.success('退出成功');
        router.push({ name: 'Login' })
        return Promise.resolve(true)
      } catch (err) {
        return Promise.reject(err)
      }
    }
  }
})
