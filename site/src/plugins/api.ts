export default defineNuxtPlugin((nuxtApp) => {
  // 创建 $api 方法，支持 GET 和 POST 请求
  const $api = {
    get: (url: string, options: any = {}) => {
      return $fetch(url, {
        method: 'GET',
        ...options
      })
    },
    post: (url: string, body: any, options: any = {}) => {
      return $fetch(url, {
        method: 'POST',
        body,
        ...options
      })
    },
    put: (url: string, body: any, options: any = {}) => {
      return $fetch(url, {
        method: 'PUT',
        body,
        ...options
      })
    },
    delete: (url: string, options: any = {}) => {
      return $fetch(url, {
        method: 'DELETE',
        ...options
      })
    }
  }

  nuxtApp.provide('api', $api)
})
