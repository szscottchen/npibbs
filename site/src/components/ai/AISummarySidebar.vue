<template>
  <div class="ai-summary-sidebar">
    <div class="ai-summary-header">
      <h3>{{ $t('ai_summary.title') }}</h3>
      <button 
        @click="generateSummary" 
        :disabled="generating || !canGenerate"
        class="generate-btn"
        :class="{ 'is-loading': generating }"
      >
        <span v-if="!generating">{{ $t('ai_summary.generate') }}</span>
        <span v-else>{{ $t('ai_summary.generating') }}</span>
      </button>
    </div>

    <div v-if="summary" class="summary-content">
      <pre class="summary-text">{{ summary }}</pre>
      <div class="summary-actions">
        <button @click="copySummary" class="copy-btn">
          {{ $t('ai_summary.copy') }}
        </button>
      </div>
    </div>

    <div v-else-if="generating" class="loading-state">
      <div class="loading-spinner"></div>
      <p>{{ $t('ai_summary.generating') }}</p>
    </div>

    <div v-else class="empty-state">
      <p>{{ $t('ai_summary.empty') }}</p>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  topicId: {
    type: [String, Number],
    required: true
  }
})

const { $toast } = useNuxtApp()

// 状态管理
const summary = ref('')
const generating = ref(false)
const canGenerate = computed(() => props.topicId && !generating.value)

// 生成总结
const generateSummary = async () => {
  if (!canGenerate.value) return

  generating.value = true

  try {
    const response = await $fetch(`/api/ai-sum/topic/${props.topicId}/generate-summary`, {
      method: 'POST',
      body: {}
    })

    console.log('生成总结响应:', response)

    // 处理 web.JsonData 包装格式: { data: { content: "..." } }
    const data = response.data || response
    if (data && data.content) {
      summary.value = data.content
      $toast.success($t('ai_summary.generate_success'))
    } else {
      console.error('响应格式不正确:', response)
      $toast.error($t('ai_summary.generate_failed'))
    }
  } catch (error) {
    console.error('生成AI总结失败:', error)
    const errorMsg = error?.message || error?.statusMessage || $t('ai_summary.generate_failed')
    $toast.error(errorMsg)
  } finally {
    generating.value = false
  }
}

// 复制总结
const copySummary = async () => {
  if (!summary.value) return
  
  try {
    await navigator.clipboard.writeText(summary.value)
    $toast.success($t('ai_summary.copy_success'))
  } catch (error) {
    console.error('复制失败:', error)
    $toast.error($t('ai_summary.copy_failed'))
  }
}

// 组件挂载时尝试获取已有的总结
onMounted(async () => {
  try {
    const response = await $fetch(`/api/ai-sum/topic/${props.topicId}/summary`)
    console.log('获取总结响应:', response)
    // 处理 web.JsonData 包装格式: { data: { content: "..." } }
    const data = response.data || response
    if (data && data.content) {
      summary.value = data.content
    }
  } catch (error) {
    // 404 或其他错误表示没有总结，静默处理
    console.log('暂无AI总结')
  }
})
</script>

<style scoped>
.ai-summary-sidebar {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 20px;
}

.ai-summary-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  min-height: 28px;
  flex-shrink: 0;
}

.ai-summary-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.generate-btn {
  background: #007bff;
  color: white;
  border: none;
  padding: 4px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: background-color 0.2s;
  height: 28px;
  line-height: 1;
}

.generate-btn:hover:not(:disabled) {
  background: #0056b3;
}

.generate-btn:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

.generate-btn.is-loading {
  opacity: 0.7;
}

.summary-content {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  padding: 12px;
}

.summary-text {
  font-size: 12px;
  line-height: 1.6;
  color: #495057;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
  margin: 0 0 12px 0;
  padding: 0;
  background: transparent;
  border: none;
  display: block;
  overflow-wrap: break-word;
}

.summary-actions {
  display: flex;
  justify-content: flex-end;
}

.copy-btn {
  background: #28a745;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.copy-btn:hover {
  background: #218838;
}

.loading-state {
  text-align: center;
  padding: 20px;
  color: #6c757d;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #f3f3f3;
  border-top: 2px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 8px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 20px;
  color: #6c757d;
}
</style>