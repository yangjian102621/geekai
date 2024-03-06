import { defineStore } from 'pinia'
import { Message } from '@arco-design/web-vue'
import { userLogin, userLogout } from '@/http/login'
import router from '@/router'

export const useAuthStore = defineStore({
  id: __AUTH_KEY,
  state: () => ({ token: null }),
  actions: {
    init() {
      this.$state.token = localStorage.getItem(__AUTH_KEY);
    },
    async login(params) {
      try {
        const { data } = await userLogin(params)
        if (data) {
          this.$state.token = data;
          localStorage.setItem(__AUTH_KEY, data)
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
          this.$restate.token = null
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
