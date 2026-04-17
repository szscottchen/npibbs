<template>
  <div class="mobile-container">
    <MobileHeader title="企业微信登录" />
    
    <div class="callback-content">
      <div class="loading" v-if="loading">
        <div class="loading-icon">⏳</div>
        <div class="loading-text">{{ loadingText }}</div>
      </div>
      
      <div class="error" v-else-if="error">
        <div class="error-icon">❌</div>
        <div class="error-message">{{ error }}</div>
        <button class="retry-btn" @click="goBack">返回</button>
      </div>
    </div>
  </div>
</template>

<script setup>
const route = useRoute()
const router = useRouter()
const { $api } = useNuxtApp()
const userStore = useUserStore()

const loading = ref(true)
const loadingText = ref('正在处理登录...')
const error = ref('')

const handleCallback = async () => {
  try {
    const { code, state } = route.query

    console.log('企业微信回调参数:', { code, state })

    if (!code) {
      throw new Error('授权码不能为空')
    }

    // 验证state参数（企业微信内置浏览器可能有存储隔离，仅记录警告不阻止登录）
    const savedState = localStorage.getItem('wecom_auth_state')
    console.log('保存的state:', savedState, 'URL中的state:', state)
    
    if (!savedState) {
      console.warn('localStorage 中没有找到保存的state，企业微信内置浏览器可能有存储隔离')
      // 企业微信环境放宽验证，仅记录警告
    } else if (savedState !== state) {
      console.warn('state 不匹配，可能存在安全风险')
      throw new Error('无效的state参数')
    }

    // 清除state
    localStorage.removeItem('wecom_auth_state')

    loadingText.value = '正在验证企业微信信息...'

    // 调用后端API处理回调
    const response = await $api.get('/api/wecom/callback', {
      params: { code, state }
    })

    console.log('企业微信回调响应:', response)
    console.log('response类型:', typeof response)
    console.log('response所有键:', Object.keys(response || {}))

    // 检查后端是否返回错误（需要同时检查 errorCode 和 success 字段）
    const errorCode = response?.errorCode
    const success = response?.success
    console.log('errorCode:', errorCode, 'success:', success)
    
    if ((errorCode !== 0 && errorCode !== undefined) || success === false) {
      throw new Error(response.message || '登录失败')
    }

    // 后端返回结构：{ errorCode, message, data: { needBind, ... } 或 { token, user, redirect } }
    // 也可能直接返回 { token, user, redirect }
    let data = response?.data
    
    // 如果没有 data 字段，检查 response 本身是否有 token 或 needBind
    if (!data && (response?.token || response?.needBind !== undefined)) {
      console.log('响应直接包含数据字段，不使用 data 包装')
      data = response
    }
    
    // 如果 data 为 null 或空对象，报错
    if (!data || Object.keys(data).length === 0) {
      throw new Error(response.message || '登录响应数据为空')
    }
    
    // 调试：打印实际返回的数据
    console.log('处理后的 data:', data)
    console.log('data.needBind:', data.needBind)
    console.log('data.token:', data.token)

    // 检查是否需要绑定
    if (data.needBind) {
      // 需要绑定，跳转到绑定页面
      loadingText.value = '正在跳转到绑定页面...'

      // 保存企业微信用户信息
      if (data.wecomUserInfo) {
        sessionStorage.setItem('wecom_bind_userinfo', JSON.stringify(data.wecomUserInfo))
      }
      if (data.suggestedUsername) {
        sessionStorage.setItem('wecom_bind_username', data.suggestedUsername)
      }

      // 清除授权开始时间
      sessionStorage.removeItem('wecom_auth_start_time')
      await router.replace(data.bindUrl || '/mobile/bind/password')
    } else if (data.token) {
      // 登录成功，保存token并跳转到首页
      loadingText.value = '登录成功，正在跳转...'

      // 保存token到本地存储
      localStorage.setItem('token', data.token)

      // 保存用户信息并更新 store
      if (data.user) {
        localStorage.setItem('userInfo', JSON.stringify(data.user))
        userStore.user = data.user
      }

      // 清除授权开始时间
      sessionStorage.removeItem('wecom_auth_start_time')
      await router.replace(data.redirect || '/mobile')
    } else {
      console.error('无法识别的响应数据:', data)
      throw new Error(`登录响应数据异常: 缺少 needBind 或 token 字段，data=${JSON.stringify(data)}`)
    }

  } catch (err) {
    console.error('企业微信回调处理失败:', err)
    error.value = err.message || '登录失败，请稍后重试'
    loading.value = false
    // 清除授权开始时间，允许重新登录
    sessionStorage.removeItem('wecom_auth_start_time')
  }
}

const goBack = () => {
  // 清除授权开始时间
  sessionStorage.removeItem('wecom_auth_start_time')
  router.push('/mobile')
}

onMounted(() => {
  // 延迟执行，避免 hydration 问题
  setTimeout(() => {
    handleCallback()
  }, 100)
})
</script>

<style scoped>
.callback-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: 20px;
}

.loading, .error {
  text-align: center;
}

.loading-icon {
  font-size: 48px;
  margin-bottom: 16px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.loading-text {
  font-size: 16px;
  color: #666;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-message {
  font-size: 16px;
  color: #ff4757;
  margin-bottom: 20px;
  line-height: 1.4;
}

.retry-btn {
  padding: 10px 24px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}

.retry-btn:hover {
  background: #40a9ff;
}
</style>