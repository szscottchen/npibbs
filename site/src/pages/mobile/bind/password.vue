<template>
  <div class="mobile-container">
    <MobileHeader title="账号绑定" />
    
    <div class="bind-content">
      <!-- 企业微信用户信息 -->
      <div class="wecom-user-info">
        <div class="avatar">
          <img :src="wecomUserInfo.avatar || '/default-avatar.png'" alt="头像" />
        </div>
        <div class="info">
          <div class="name">{{ wecomUserInfo.name }}</div>
          <div class="userid">工号：{{ wecomUserInfo.userId }}</div>
          <div class="mobile" v-if="wecomUserInfo.mobile">电话：{{ wecomUserInfo.mobile }}</div>
        </div>
      </div>
      
      <!-- 绑定说明 -->
      <div class="bind-tips">
        <div class="tips-title">绑定说明</div>
        <div class="tips-content">
          检测到您的企业微信账号尚未绑定系统账号，请输入您的账号密码进行绑定。
          绑定成功后，下次可直接通过企业微信登录。
        </div>
      </div>
      
      <!-- 密码输入表单 -->
      <div class="bind-form">
        <div class="form-item">
          <label>系统账号</label>
          <input 
            type="text" 
            v-model="bindForm.username"
            placeholder="请输入工号/用户名"
            :class="{ error: errors.username }"
          />
          <div class="error-message" v-if="errors.username">{{ errors.username }}</div>
        </div>
        
        <div class="form-item">
          <label>登录密码</label>
          <input 
            type="password" 
            v-model="bindForm.password"
            placeholder="请输入密码"
            :class="{ error: errors.password }"
          />
          <div class="error-message" v-if="errors.password">{{ errors.password }}</div>
        </div>
        
        <div class="form-actions">
          <button 
            class="bind-btn"
            :disabled="loading"
            @click="handleBind"
          >
            {{ loading ? '绑定中...' : '确认绑定' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const route = useRoute()
const router = useRouter()
const { $api } = useNuxtApp()

const loading = ref(false)
const wecomUserInfo = ref({})
const bindForm = reactive({
  username: '',
  password: ''
})
const errors = reactive({
  username: '',
  password: ''
})

// 从URL参数或会话存储中获取企业微信用户信息
const loadWeComUserInfo = () => {
  const userInfo = sessionStorage.getItem('wecom_bind_userinfo')
  if (userInfo) {
    wecomUserInfo.value = JSON.parse(userInfo)
  } else {
    // 如果没有信息，重定向到企业微信授权
    redirectToWeComAuth()
  }
  
  // 填充建议的用户名
  const suggestedUsername = sessionStorage.getItem('wecom_bind_username')
  if (suggestedUsername) {
    bindForm.username = suggestedUsername
  }
}

// 重定向到企业微信授权
const redirectToWeComAuth = () => {
  router.push('/mobile/wecom-auth')
}

// 表单验证
const validateForm = () => {
  let isValid = true
  
  // 清空错误信息
  errors.username = ''
  errors.password = ''
  
  if (!bindForm.username.trim()) {
    errors.username = '请输入账号'
    isValid = false
  }
  
  if (!bindForm.password.trim()) {
    errors.password = '请输入密码'
    isValid = false
  }
  
  return isValid
}

// 处理绑定
const handleBind = async () => {
  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    console.log('开始绑定:', {
      username: bindForm.username,
      wecomUserId: wecomUserInfo.value.userId
    })

    const response = await $api.post('/api/wecom/bind', {
      username: bindForm.username,
      password: bindForm.password,
      wecomUserId: wecomUserInfo.value.userId
    })

    console.log('绑定响应:', response)

    // 绑定成功，自动登录
    const { token, user } = response.data

    // 保存token
    localStorage.setItem('token', token)

    // 清除绑定临时信息
    sessionStorage.removeItem('wecom_bind_userinfo')
    sessionStorage.removeItem('wecom_bind_username')

    // 跳转到首页
    router.replace('/mobile')

  } catch (error) {
    console.error('绑定失败:', error)
    console.error('错误详情:', {
      message: error.message,
      response: error.response,
      data: error.response?.data
    })

    if (error.response?.data?.message) {
      const message = error.response.data.message

      if (message.includes('用户不存在') || message.includes('账号')) {
        errors.username = message
      } else if (message.includes('密码')) {
        errors.password = message
      } else {
        errors.password = message
      }
    } else {
      errors.password = '绑定失败，请稍后重试'
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadWeComUserInfo()
})
</script>

<style scoped>
.bind-content {
  padding: 16px;
}

.wecom-user-info {
  display: flex;
  align-items: center;
  padding: 16px;
  background: white;
  border-radius: 8px;
  margin-bottom: 16px;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 25px;
  overflow: hidden;
  margin-right: 12px;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.info {
  flex: 1;
}

.name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.userid, .mobile {
  font-size: 14px;
  color: #666;
  margin-bottom: 2px;
}

.bind-tips {
  padding: 12px;
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 6px;
  margin-bottom: 16px;
}

.tips-title {
  font-size: 14px;
  font-weight: 600;
  color: #0369a1;
  margin-bottom: 4px;
}

.tips-content {
  font-size: 13px;
  color: #0c4a6e;
  line-height: 1.4;
}

.bind-form {
  background: white;
  padding: 16px;
  border-radius: 8px;
}

.form-item {
  margin-bottom: 16px;
}

.form-item label {
  display: block;
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
}

.form-item input {
  width: 100%;
  height: 44px;
  padding: 0 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.form-item input.error {
  border-color: #ff4757;
}

.error-message {
  font-size: 12px;
  color: #ff4757;
  margin-top: 4px;
}

.form-actions {
  margin-top: 24px;
}

.bind-btn {
  width: 100%;
  height: 44px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.bind-btn:hover:not(:disabled) {
  background: #40a9ff;
}

.bind-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}
</style>