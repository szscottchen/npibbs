<template>
  <div class="wecom-entry">
    <div class="loading" v-if="loading">
      <div class="spinner"></div>
      <div class="text">正在跳转...</div>
    </div>
    <div class="error" v-else-if="error">
      <div class="error-icon">❌</div>
      <div class="error-text">{{ error }}</div>
    </div>
  </div>
</template>

<script setup>
const loading = ref(true)
const error = ref('')

// 直接跳转到企业微信授权页面
onMounted(async () => {
  try {
    console.log('访问企业微信入口页面')
    // 重定向到后端 API
    window.location.href = '/api/wecom/entry'
  } catch (err) {
    console.error('跳转失败:', err)
    error.value = err.message || '跳转失败，请重试'
    loading.value = false
  }
})
</script>

<style scoped>
.wecom-entry {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f5f5f5;
}

.loading, .error {
  text-align: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e0e0e0;
  border-top: 4px solid #1890ff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.text {
  font-size: 14px;
  color: #666;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-text {
  font-size: 16px;
  color: #ff4757;
}
</style>
