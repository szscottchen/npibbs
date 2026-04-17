<template>
  <div class="mobile-container">
    <div class="mobile-content">
      <!-- 搜索栏 -->
      <div class="search-bar">
        <input 
          type="text" 
          placeholder="搜索话题..." 
          v-model="searchKeyword"
          @keyup.enter="handleSearch"
        />
        <div class="search-icon" @click="handleSearch">
          <i class="iconfont icon-search"></i>
        </div>
      </div>

      <!-- 筛选标签 -->
      <div class="filter-tabs">
        <div 
          v-for="tab in tabs" 
          :key="tab.id"
          class="tab-item"
          :class="{ active: currentTab === tab.id }"
          @click="switchTab(tab)"
        >
          {{ tab.name }}
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
const configStore = useConfigStore()
const { config } = storeToRefs(configStore)

// 从URL参数读取初始tab
const currentTab = ref(route.query.nodeId ? parseInt(route.query.nodeId) : 0)
const topics = ref([])
const page = ref(1)
const pageSize = 20
const hasMore = ref(true)
const loading = ref(false)
const searchKeyword = ref('')

// 标签定义，与桌面端 TopicsNav 一致
const tabs = ref([
  { id: 0, name: '最新' },
  { id: -2, name: '关注' },
  { id: -1, name: '推荐' },
  { id: -3, name: '呼叫支援' }
])

// 切换标签
const switchTab = (tab) => {
  currentTab.value = tab.id
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
    // 根据节点ID选择API路径，与桌面端保持一致
    let apiUrl = '/api/topic/topics'
    let params = {
      page: page.value
    }

    if (currentTab.value === -3) {
      // 呼叫支援使用特殊接口
      apiUrl = '/api/topic/need_a_hand'
      // 不需要 nodeId 参数
    } else {
      // 其他都使用 /api/topic/topics，带上 nodeId 参数
      params.nodeId = currentTab.value
    }

    const response = await $api.get(apiUrl, { params })
    const responseData = response.data
    const newTopics = responseData?.results || []

    if (page.value === 1) {
      topics.value = newTopics
      cursor.value = responseData?.cursor || null
    } else {
      topics.value = [...topics.value, ...newTopics]
      cursor.value = responseData?.cursor || null
    }

    hasMore.value = responseData?.hasMore || false
  } catch (error) {
    console.error('加载话题失败:', error)
  } finally {
    loading.value = false
  }
}

// 游标
const cursor = ref(null)

// 加载更多
const loadMore = () => {
  if (!loading.value && hasMore.value) {
    page.value++
    loadTopics()
  }
}

// 搜索
const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push(`/mobile/topics?keyword=${encodeURIComponent(searchKeyword.value)}`)
  }
}



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

.search-bar {
  position: relative;
  padding: 12px 16px;
  background: white;
  border-bottom: 1px solid #eee;
}

.search-bar input {
  width: 100%;
  height: 40px;
  padding: 0 40px 0 12px;
  border: 1px solid #ddd;
  border-radius: 20px;
  font-size: 15px;
  background: #f8f8f8;
}

.search-bar input:focus {
  outline: none;
  border-color: #3273dc;
  background: white;
}

.search-icon {
  position: absolute;
  right: 28px;
  top: 50%;
  transform: translateY(-50%);
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3273dc;
  cursor: pointer;
}

.search-icon .iconfont {
  font-size: 18px;
}

.filter-tabs {
  display: flex;
  background: white;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
}

.tab-item {
  flex: 1;
  padding: 8px 0;
  text-align: center;
  font-size: 14px;
  color: #666;
  background: #f8f8f8;
  margin: 0 4px;
  border-radius: 16px;
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