<template>
  <div class="mobile-container">
    <MobileHeader title="话题列表" :show-back="true" />

    <div class="mobile-content">
      <!-- 筛选标签 -->
      <div class="filter-tabs">
        <div 
          v-for="tab in tabs" 
          :key="tab.key"
          class="tab-item"
          :class="{ active: currentTab === tab.key }"
          @click="switchTab(tab.key)"
        >
          {{ tab.label }}
        </div>
      </div>

      <!-- 话题列表 -->
      <div v-if="topics.length > 0" class="topic-list">
        <TopicList :topics="topics" :show-sticky="false" />
      </div>

      <!-- 加载更多 -->
      <div v-if="hasMore" class="load-more" @click="loadMore">
        <span v-if="loading">加载中...</span>
        <span v-else>加载更多</span>
      </div>

      <!-- 没有更多 -->
      <div v-if="!hasMore && topics.length > 0" class="no-more">
        没有更多了
      </div>

      <!-- 空状态 -->
      <div v-if="topics.length === 0 && !loading" class="empty-state">
        <i class="iconfont icon-empty"></i>
        <p>暂无话题</p>
      </div>
    </div>
  </div>
</template>

<script setup>
const route = useRoute()
const router = useRouter()
const { $api } = useNuxtApp()

const tabs = [
  { key: 'latest', label: '最新' },
  { key: 'hot', label: '热门' },
  { key: 'recommend', label: '推荐' }
]

// 从URL参数读取初始tab
const currentTab = ref(route.query.tab || 'latest')
const topics = ref([])
const page = ref(1)
const pageSize = 20
const hasMore = ref(true)
const loading = ref(false)

// 切换标签
const switchTab = (tab) => {
  currentTab.value = tab
  page.value = 1
  topics.value = []
  hasMore.value = true
  loadTopics()
}

// 加载话题列表
const loadTopics = async () => {
  if (loading.value || !hasMore.value) return

  loading.value = true
  try {
    let apiUrl = '/api/topic/topics'

    // 构建查询参数
    const params = {
      page: page.value,
      limit: pageSize
    }

    // 如果有节点ID
    if (route.query.nodeId) {
      params.nodeId = route.query.nodeId
    }

    // 如果有搜索关键词
    if (route.query.keyword) {
      params.keyword = route.query.keyword
    }

    const response = await $api.get(apiUrl, { params })
    
    // 处理响应数据
    const responseData = response.data
    const newTopics = responseData?.results || []

    if (page.value === 1) {
      topics.value = newTopics
    } else {
      topics.value = [...topics.value, ...newTopics]
    }

    // 判断是否还有更多
    hasMore.value = responseData?.hasMore || false
  } catch (error) {
    console.error('加载话题失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载更多
const loadMore = () => {
  if (!loading.value && hasMore.value) {
    page.value++
    loadTopics()
  }
}

// 下拉刷新
onMounted(() => {
  loadTopics()
})
</script>

<style scoped>
.mobile-container {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.mobile-content {
  padding-bottom: 16px;
}

.filter-tabs {
  display: flex;
  background: white;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 8px;
}

.tab-item {
  padding: 6px 16px;
  margin-right: 12px;
  border-radius: 16px;
  font-size: 14px;
  color: #666;
  background: #f8f8f8;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-item.active {
  background: #3273dc;
  color: white;
}

.topic-list {
  background: white;
}

.load-more {
  padding: 12px;
  text-align: center;
  font-size: 14px;
  color: #999;
  cursor: pointer;
  background: white;
  margin-top: 8px;
}

.no-more {
  padding: 12px;
  text-align: center;
  font-size: 14px;
  color: #999;
  background: white;
  margin-top: 8px;
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
  background: white;
}

.empty-state .iconfont {
  font-size: 60px;
  color: #ddd;
  margin-bottom: 16px;
}

.empty-state p {
  font-size: 14px;
  color: #999;
}
</style>
