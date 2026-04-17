<template>
  <div class="mobile-container">
    <MobileHeader title="企业微信登录" />
    
    <div class="auth-content">
      <div class="loading" v-if="loading">
        <div class="loading-icon">⏳</div>
        <div class="loading-text">正在跳转企业微信授权...</div>
      </div>
      
      <div class="error" v-else-if="error">
        <div class="error-icon">❌</div>
        <div class="error-message">{{ error }}</div>
        <button class="retry-btn" @click="retryAuth">重试</button>
      </div>
    </div>
  </div>
</template>

<script setup>
const route = useRoute()
const router = useRouter()
const { $api } = useNuxtApp()

const loading = ref(true)
const error = ref('')
const isRedirecting = ref(false)

const initAuth = async () => {
  // 防止重复跳转
  if (isRedirecting.value) {
    console.log('正在跳转中，跳过重复调用')
    return
  }

  // 检查是否已经在授权流程中（5分钟内）
  const authStartTime = sessionStorage.getItem('wecom_auth_start_time')
  if (authStartTime) {
    const elapsed = Date.now() - parseInt(authStartTime)
    if (elapsed < 5 * 60 * 1000) { // 5分钟内
      console.log('授权流程正在进行中，跳过重复调用')
      return
    }
  }

  try {
    isRedirecting.value = true

    // 1. 获取企业微信授权URL
    console.log('开始请求企业微信授权...')
    const response = await $api.get('/api/wecom/entry')
    console.log('企业微信授权响应:', response)

    // 检查后端是否返回错误
    if (response?.errorCode !== 0 && response?.errorCode !== undefined) {
      throw new Error(response.message || '获取授权失败')
    }

    // 后端返回结构：{ errorCode, message, data: { authUrl, state } }
    const responseData = response?.data || {}
    const { authUrl, state } = responseData

    if (!authUrl) {
      throw new Error('授权URL为空，请检查企业微信配置')
    }

    // 2. 保存state用于回调验证（使用localStorage，因为企业微信跳转会丢失sessionStorage）
    localStorage.setItem('wecom_auth_state', state)

    // 3. 记录授权开始时间，防止重复跳转
    sessionStorage.setItem('wecom_auth_start_time', Date.now().toString())

    // 4. 替换授权URL中的回调地址为前端页面
    const redirectUri = encodeURIComponent(window.location.origin + '/mobile/wecom-callback')
    const modifiedAuthUrl = authUrl.replace(/redirect_uri=[^&]+/, `redirect_uri=${redirectUri}`)
    console.log('修改后的授权URL:', modifiedAuthUrl)

    // 5. 跳转到企业微信授权页面
    window.location.href = modifiedAuthUrl

  } catch (err) {
    console.error('获取企业微信授权失败:', err)
    error.value = err.message || '获取授权失败，请稍后重试'
    loading.value = false
    isRedirecting.value = false
    // 清除授权开始时间，允许重试
    sessionStorage.removeItem('wecom_auth_start_time')
  }
}

const retryAuth = () => {
  loading.value = true
  error.value = ''
  // 清除授权开始时间
  sessionStorage.removeItem('wecom_auth_start_time')
  initAuth()
}

// 页面加载时自动开始授权流程
onMounted(() => {
  // 延迟执行，避免 hydration 问题导致的重复调用
  setTimeout(() => {
    initAuth()
  }, 100)
})
</script>

<style scoped>
.auth-content {
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