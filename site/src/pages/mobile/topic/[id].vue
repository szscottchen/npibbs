<template>
  <div class="mobile-chat-container">
    <MobileHeader :title="topic?.title || t('pages.topic.detail.topic')" :show-back="true" />

    <div v-if="topic" class="chat-content" ref="chatMessagesRef">
      <!-- 时间分割线 -->
      <div class="time-divider">{{ formatChatTime(topic.createTime) }}</div>

      <!-- 楼主消息 -->
      <div class="chat-message is-topic">
        <div class="message-avatar">
          <MyAvatar :user="topic.user" :size="36" />
        </div>
        <div class="message-content-wrapper">
          <div class="message-header">
            <span class="message-author">{{ topic.user.nickname }}</span>
            <span v-if="topic.ipLocation" class="message-location">{{ topic.ipLocation }}</span>
          </div>
          <div class="message-bubble is-topic">
            <div class="topic-title" v-if="topic.title">{{ topic.title }}</div>
            <div class="message-text" v-html="processedContent"></div>
            
            <!-- 图片列表 -->
            <div v-if="topic.imageList && topic.imageList.length" class="message-images">
              <div 
                v-for="(image, index) in topic.imageList" 
                :key="image.url"
                class="image-wrapper"
                :class="{ 'single': topic.imageList.length === 1 }"
              >
                <img 
                  v-if="isImageFile(image)" 
                  :src="image.preview" 
                  @click="previewImage(index)"
                />
                <div v-else class="file-item" @click="openFile(image.url)">
                  <i class="iconfont icon-file"></i>
                  <span class="file-name">{{ image.filename || getFileName(image.url) }}</span>
                </div>
              </div>
            </div>

            <!-- 标签 -->
            <div v-if="topic.node || topic.tags.length" class="topic-tags">
              <span v-if="topic.node" class="tag node">{{ topic.node.name }}</span>
              <span v-for="tag in topic.tags" :key="tag.id" class="tag">#{{ tag.name }}</span>
            </div>
          </div>
          
          <!-- 话题操作 -->
          <div class="message-actions">
            <span class="action-btn" :class="{ active: liked }" @click="handleLike">
              <i class="iconfont icon-like"></i>
              <span>{{ topic.likeCount || t('pages.topic.detail.like') }}</span>
            </span>
            <span class="action-btn" :class="{ active: topic.favorited }" @click="handleFavorite">
              <i class="iconfont" :class="topic.favorited ? 'icon-has-favorite' : 'icon-favorite'"></i>
              <span>{{ topic.favorited ? t('pages.topic.detail.favorited') : t('pages.topic.detail.favorite') }}</span>
            </span>
            <span class="action-btn" @click="scrollToInput">
              <i class="iconfont icon-comment"></i>
              <span>{{ topic.commentCount || t('pages.topic.detail.comment') }}</span>
            </span>
          </div>
        </div>
      </div>

      <!-- 评论消息列表 -->
      <div v-for="(group, groupIndex) in groupedComments" :key="groupIndex">
        <div class="time-divider">{{ group.timeLabel }}</div>
        <div 
          v-for="comment in group.comments" 
          :key="comment.id"
          class="chat-message"
          :class="{ 'is-me': isCurrentUser(comment.user.id) }"
          :id="`comment-${comment.id}`"
        >
          <div class="message-avatar">
            <MyAvatar :user="comment.user" :size="32" />
          </div>
          <div class="message-content-wrapper">
            <div class="message-header">
              <span class="message-author">{{ comment.user.nickname }}</span>
              <span v-if="comment.ipLocation" class="message-location">{{ comment.ipLocation }}</span>
            </div>
            <div 
              class="message-bubble"
              :class="{ 'is-me': isCurrentUser(comment.user.id) }"
              @click="handleQuoteClick(comment)"
            >
              <!-- 引用内容 -->
              <div v-if="comment.quote" class="message-quote">
                <div class="quote-content">
                  <span class="quote-author">{{ comment.quote.user.nickname }}:</span>
                  <span class="quote-text" v-html="comment.quote.content"></span>
                </div>
              </div>
              <div class="message-text" v-html="comment.content"></div>
              
              <!-- 评论图片 -->
              <div v-if="comment.imageList && comment.imageList.length" class="message-images">
                <img 
                  v-for="(image, idx) in comment.imageList" 
                  :key="idx"
                  :src="image.preview"
                  @click="previewCommentImage(comment.imageList, idx)"
                />
              </div>
            </div>
            <div class="message-actions">
              <span class="action-btn" :class="{ active: comment.liked }" @click="likeComment(comment)">
                <i class="iconfont icon-like"></i>
                <span>{{ comment.likeCount > 0 ? comment.likeCount : t('component.comment.list.like') }}</span>
              </span>
              <span class="action-btn" @click="replyTo(comment)">
                <i class="iconfont icon-comment"></i>
                <span>{{ t('component.comment.list.reply') }}</span>
              </span>
              <span
                v-if="userStore.user && userStore.user.id == topic?.user?.id && valueTypes.length > 0"
                class="action-btn value-btn"
                @click="showValue(comment)"
              >
                <i class="iconfont icon-star"></i>
                <span>价值</span>
              </span>
            </div>
            <!-- 子评论列表 -->
            <div v-if="comment.replies && comment.replies.results && comment.replies.results.length" class="comment-replies">
              <div
                v-for="reply in comment.replies.results"
                :key="reply.id"
                class="reply-item"
                :id="`comment-${reply.id}`"
              >
                <div class="reply-avatar">
                  <MyAvatar :user="reply.user" :size="24" />
                </div>
                <div class="reply-content-wrapper">
                  <div class="reply-header">
                    <span class="reply-author">{{ reply.user.nickname }}</span>
                    <span v-if="reply.ipLocation" class="reply-location">{{ reply.ipLocation }}</span>
                  </div>
                  <div class="reply-bubble">
                    <div v-if="reply.quote" class="reply-quote">
                      <span class="quote-author">{{ reply.quote.user.nickname }}:</span>
                      <span class="quote-text" v-html="reply.quote.content"></span>
                    </div>
                    <div class="reply-text" v-html="reply.content"></div>
                    <div v-if="reply.imageList && reply.imageList.length" class="reply-images">
                      <img
                        v-for="(image, idx) in reply.imageList"
                        :key="idx"
                        :src="image.preview"
                        @click="previewCommentImage(reply.imageList, idx)"
                      />
                    </div>
                  </div>
                  <div class="reply-actions">
                    <span class="action-btn" :class="{ active: reply.liked }" @click="likeComment(reply)">
                      <i class="iconfont icon-like"></i>
                      <span>{{ reply.likeCount > 0 ? reply.likeCount : t('component.comment.list.like') }}</span>
                    </span>
                    <span class="action-btn" @click="replyTo(reply)">
                      <i class="iconfont icon-comment"></i>
                      <span>{{ t('component.comment.list.reply') }}</span>
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载更多 -->
      <div v-if="hasMoreComments" class="load-more">
        <span v-if="loadingMore" class="loading-text">
          <i class="iconfont icon-loading"></i> {{ t('common.loadMore.loadMore') }}
        </span>
        <span v-else class="load-more-btn" @click="loadMoreComments">
          {{ t('component.comment.list.loadMore') }}
        </span>
      </div>

      <!-- 无更多数据 -->
      <div v-else-if="comments.length > 0" class="no-more">
        {{ t('common.loadMore.noMore') }}
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading-state">
      <i class="iconfont icon-loading"></i>
    </div>

    <!-- 新消息提示 -->
    <div v-if="newMessageCount > 0" class="new-message-tip" @click="scrollToBottom">
      <i class="iconfont icon-down"></i>
      <span>{{ newMessageCount }} {{ t('component.chat.newMessages') }}</span>
    </div>

    <!-- 价值评价对话框 -->
    <div v-if="showValueDialog" class="value-dialog-overlay" @click="showValueDialog = false">
      <div class="value-dialog" @click.stop>
        <div class="value-dialog-header">
          <span>价值评价</span>
          <i class="iconfont icon-close" @click="showValueDialog = false"></i>
        </div>
        <div class="value-dialog-content">
          <div
            v-for="(valueType, index) in valueTypes"
            :key="index"
            class="value-type-item"
            @click="submitValue(valueType, index)"
          >
            <span class="value-type-label">{{ valueType.label }}</span>
            <span class="value-type-score">+{{ valueType.score }}积分</span>
          </div>
        </div>
      </div>
    </div>

    <!-- AI总结对话框 -->
    <div v-if="showAISummary" class="ai-summary-dialog-overlay" @click="showAISummary = false">
      <div class="ai-summary-dialog" @click.stop>
        <div class="ai-summary-dialog-header">
          <span>AI总结</span>
          <div class="header-actions">
            <button 
              class="ai-generate-btn" 
              @click="generateAISummary"
              :disabled="aiSummaryGenerating"
            >
              <span v-if="!aiSummaryGenerating">生成总结</span>
              <span v-else>生成中...</span>
            </button>
            <i class="iconfont icon-close" @click="showAISummary = false"></i>
          </div>
        </div>
        <div class="ai-summary-dialog-content">
          <div v-if="aiSummaryLoading" class="ai-summary-loading">
            <i class="iconfont icon-loading"></i>
            <span>AI正在总结中...</span>
          </div>
          <div v-else-if="aiSummaryContent" class="ai-summary-text">
            {{ aiSummaryContent }}
          </div>
          <div v-else class="ai-summary-empty">
            <i class="iconfont icon-star"></i>
            <span>点击AI总结按钮生成内容总结</span>
          </div>
        </div>
        <div class="ai-summary-dialog-actions">
          <button class="ai-summary-copy-btn" @click="copyAISummary" :disabled="!aiSummaryContent">
            <i class="iconfont icon-copy"></i>
            <span>复制</span>
          </button>
          <button class="ai-summary-close-btn" @click="showAISummary = false">
            <span>关闭</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 底部输入框 -->
    <div class="chat-input-wrapper" ref="inputWrapperRef">
      <div v-if="replyToComment" class="quote-bar">
        <span>{{ t('component.comment.input.replyTo') }} {{ replyToComment.user.nickname }}</span>
        <i class="iconfont icon-close" @click="cancelReply"></i>
      </div>
      <div class="input-bar">
        <div class="input-box">
          <textarea
            ref="inputRef"
            v-model="inputContent"
            :placeholder="t('component.chat.inputPlaceholder')"
            @keydown.enter.prevent="sendMessage"
            @input="onInputChange"
            rows="1"
          ></textarea>
        </div>
        <div class="input-actions">
          <button class="ai-btn" @click="handleAISummary">
            <span>AI</span>
          </button>
          <button
            class="send-btn"
            :disabled="!canSend || sending"
            :class="{ 'can-send': canSend }"
            @click="sendMessage"
          >
            <i class="iconfont icon-send"></i>
          </button>
        </div>
      </div>

      <!-- 调试: 显示状态 <div style="font-size: 10px; color: red;">showImageUpload: {{ showImageUpload }}</div> -->
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const route = useRoute()
const { t } = useI18n()
const userStore = useUserStore()

// 话题数据
const { data: topic } = await useMyFetch(`/api/topic/${route.params.id}`)
const { data: liked } = await useMyFetch("/api/like/liked", {
  params: { entityType: "topic", entityId: route.params.id },
})

// 评论相关
const comments = ref([])
const commentCursor = ref('')
const hasMoreComments = ref(true)
const loadingMore = ref(false)
const newMessageCount = ref(0)

// 输入相关
const inputContent = ref('')
const inputRef = ref(null)
const inputWrapperRef = ref(null)
const chatMessagesRef = ref(null)
const sending = ref(false)
const replyToComment = ref(null)

// 轮询
let pollingInterval = null
let lastCommentId = null

// 价值评价相关
const valueTypes = ref([])
const showValueDialog = ref(false)
const currentComment = ref(null)

// AI总结相关
const showAISummary = ref(false)
const aiSummaryContent = ref('')
const aiSummaryLoading = ref(false)
const aiSummaryGenerating = ref(false)

// 计算属性
const canSend = computed(() => {
  return inputContent.value.trim().length > 0
})

const imageUrls = computed(() => {
  if (!topic.value?.imageList) return []
  return topic.value.imageList.map(img => img.url)
})

// 按时间分组
const groupedComments = computed(() => {
  const groups = []
  let currentGroup = null
  
  comments.value.forEach(comment => {
    const timeLabel = formatChatTime(comment.createTime)
    
    if (!currentGroup || currentGroup.timeLabel !== timeLabel) {
      currentGroup = { timeLabel, comments: [] }
      groups.push(currentGroup)
    }
    currentGroup.comments.push(comment)
  })
  
  return groups
})

// 方法
function isCurrentUser(userId) {
  return userStore.user && userStore.user.id === userId
}

function formatChatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return t('component.chat.justNow')
  if (minutes < 60) return t('component.chat.minutesAgo', { n: minutes })
  if (hours < 24) return t('component.chat.hoursAgo', { n: hours })
  if (days < 7) return t('component.chat.daysAgo', { n: days })
  
  return date.toLocaleDateString()
}

// 文件处理
function isImageFile(file) {
  if (typeof file === 'object' && file?.isImage !== undefined) {
    return file.isImage
  }
  const url = typeof file === 'object' ? file?.url : file
  if (!url) return false
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
  return imageExtensions.some(ext => url.toLowerCase().endsWith(ext))
}

function getFileName(url) {
  if (!url) return ''
  return url.split('/').pop()
}

function openFile(url) {
  window.open(url, '_blank')
}

function processFileLinks(content) {
  if (!content) return ''
  const fileLinkRegex = /<a\s+[^>]*href="([^"]*uploads\/files[^"]*\.([a-zA-Z0-9]+))"[^>]*>(.*?)<\/a>/gi
  return content.replace(fileLinkRegex, (match, href, ext, linkText) => {
    return `<span class="file-link-wrapper"><i class="iconfont icon-file"></i><a href="${href}" rel="nofollow">${linkText}</a></span>`
  })
}

const processedContent = computed(() => {
  return topic.value ? processFileLinks(topic.value.content) : ''
})

// 加载评论
async function loadComments(cursor = '') {
  try {
    const params = { entityType: 'topic', entityId: route.params.id }
    if (cursor) params.cursor = cursor

    const result = await useHttpGet('/api/comment/comments', { params })

    // 处理评论数据，将子评论挂载到父评论下
    const processComments = (results) => {
      const mainComments = []
      const repliesMap = new Map()

      // 第一遍：分离主评论和子评论
      results.forEach(comment => {
        if (comment.quoteId) {
          // 这是子评论，需要挂载到父评论下
          if (!repliesMap.has(comment.quoteId)) {
            repliesMap.set(comment.quoteId, [])
          }
          repliesMap.get(comment.quoteId).push(comment)
        } else {
          // 这是主评论
          mainComments.push(comment)
        }
      })

      // 第二遍：将子评论挂载到对应的主评论
      mainComments.forEach(comment => {
        if (repliesMap.has(comment.id)) {
          comment.replies = {
            results: repliesMap.get(comment.id),
            hasMore: false
          }
        }
      })

      return mainComments
    }

    const processedResults = processComments(result.results)

    if (cursor) {
      comments.value.push(...processedResults)
    } else {
      comments.value = processedResults
    }

    commentCursor.value = result.cursor
    hasMoreComments.value = result.hasMore

    if (comments.value.length > 0) {
      lastCommentId = comments.value[comments.value.length - 1].id
    }
  } catch (e) {
    console.error('Failed to load comments:', e)
  }
}

async function loadMoreComments() {
  if (loadingMore.value || !hasMoreComments.value) return
  loadingMore.value = true
  await loadComments(commentCursor.value)
  loadingMore.value = false
}

// 发送消息
async function sendMessage() {
  if (!canSend.value || sending.value) return
  if (!userStore.isLogin) {
    useMsgSignIn()
    return
  }

  sending.value = true
  try {
    const data = await useHttpPost('/api/comment/create', useJsonToForm({
      entityType: 'topic',
      entityId: route.params.id,
      content: inputContent.value,
      quoteId: replyToComment.value?.id || '',
    }))

    // 处理新发布的评论：子评论挂载到父评论下，主评论直接追加
    if (data.quoteId) {
      // 这是子评论，找到父评论并挂载
      const parentComment = comments.value.find(c => c.id === data.quoteId)
      if (parentComment) {
        if (!parentComment.replies) {
          parentComment.replies = { results: [], hasMore: false }
        }
        parentComment.replies.results.push(data)
      } else {
        // 如果找不到父评论，直接追加（兜底）
        comments.value.push(data)
      }
    } else {
      // 这是主评论，直接追加
      comments.value.push(data)
    }
    
    inputContent.value = ''
    replyToComment.value = null
    topic.value.commentCount++
    
    useMsgSuccess(t('component.comment.input.publishSuccess'))
    
    nextTick(() => scrollToBottom())
  } catch (e) {
    useCatchError(e)
  } finally {
    sending.value = false
  }
}

// 回复相关
function replyTo(comment) {
  replyToComment.value = comment
  inputRef.value?.focus()
}

function cancelReply() {
  replyToComment.value = null
}

function handleQuoteClick(comment) {
  if (comment.quote) {
    const el = document.getElementById(`comment-${comment.quote.id}`)
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'center' })
      el.classList.add('highlight')
      setTimeout(() => el.classList.remove('highlight'), 2000)
    }
  }
}

// 点赞
async function handleLike() {
  try {
    if (liked.value) {
      await useHttpPost('/api/like/unlike', useJsonToForm({
        entityType: 'topic',
        entityId: topic.value.id,
      }))
      liked.value = false
      topic.value.likeCount = Math.max(0, topic.value.likeCount - 1)
    } else {
      await useHttpPost('/api/like/like', useJsonToForm({
        entityType: 'topic',
        entityId: topic.value.id,
      }))
      liked.value = true
      topic.value.likeCount++
    }
  } catch (e) {
    useCatchError(e)
  }
}

async function likeComment(comment) {
  try {
    if (comment.liked) {
      await useHttpPost('/api/like/unlike', useJsonToForm({
        entityType: 'comment',
        entityId: comment.id,
      }))
      comment.liked = false
      comment.likeCount = Math.max(0, comment.likeCount - 1)
    } else {
      await useHttpPost('/api/like/like', useJsonToForm({
        entityType: 'comment',
        entityId: comment.id,
      }))
      comment.liked = true
      comment.likeCount++
    }
  } catch (e) {
    useCatchError(e)
  }
}

async function handleFavorite() {
  try {
    if (topic.value.favorited) {
      await useHttpPost('/api/favorite/delete', useJsonToForm({
        entityType: 'topic',
        entityId: topic.value.id,
      }))
      topic.value.favorited = false
    } else {
      await useHttpPost('/api/favorite/add', useJsonToForm({
        entityType: 'topic',
        entityId: topic.value.id,
      }))
      topic.value.favorited = true
    }
  } catch (e) {
    useCatchError(e)
  }
}

// UI 方法
function scrollToBottom() {
  if (chatMessagesRef.value) {
    chatMessagesRef.value.scrollTop = chatMessagesRef.value.scrollHeight
    newMessageCount.value = 0
  }
}

function scrollToInput() {
  inputWrapperRef.value?.scrollIntoView({ behavior: 'smooth' })
  setTimeout(() => inputRef.value?.focus(), 300)
}

function onInputChange() {
  const textarea = inputRef.value
  if (textarea) {
    textarea.style.height = 'auto'
    textarea.style.height = Math.min(textarea.scrollHeight, 80) + 'px'
  }
}



function previewImage(index) {
  // 使用图片预览
  const images = imageUrls.value
  if (images.length > 0) {
    window.open(images[index], '_blank')
  }
}

function previewCommentImage(imageList, index) {
  const urls = imageList.map(img => img.url)
  window.open(urls[index], '_blank')
}

// 轮询新评论
function startPolling() {
  pollingInterval = setInterval(async () => {
    try {
      const result = await useHttpGet('/api/comment/comments', {
        params: { 
          entityType: 'topic', 
          entityId: route.params.id,
          afterId: lastCommentId 
        }
      })
      
      if (result.results && result.results.length > 0) {
        const newComments = result.results.filter(c =>
          !comments.value.some(existing => existing.id === c.id) &&
          !comments.value.some(existing =>
            existing.replies?.results?.some(reply => reply.id === c.id)
          )
        )

        if (newComments.length > 0) {
          // 处理新评论：子评论挂载到父评论下
          newComments.forEach(comment => {
            if (comment.quoteId) {
              // 子评论
              const parentComment = comments.value.find(c => c.id === comment.quoteId)
              if (parentComment) {
                if (!parentComment.replies) {
                  parentComment.replies = { results: [], hasMore: false }
                }
                parentComment.replies.results.push(comment)
              } else {
                comments.value.push(comment)
              }
            } else {
              // 主评论
              comments.value.push(comment)
            }
          })

          lastCommentId = result.results[result.results.length - 1].id

          const container = chatMessagesRef.value
          const isAtBottom = container.scrollHeight - container.scrollTop <= container.clientHeight + 100

          if (isAtBottom) {
            nextTick(() => scrollToBottom())
          } else {
            newMessageCount.value += newComments.length
          }
        }
      }
    } catch (e) {
      console.error('Polling error:', e)
    }
  }, 5000)
}

function stopPolling() {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

// 加载价值类型
async function loadValueTypes() {
  try {
    const config = await useHttpGet('/api/sys-config/configs')
    if (config && config.valueTypes && config.valueTypes.length > 0) {
      valueTypes.value = config.valueTypes
    }
  } catch (e) {
    console.error('Failed to load value types:', e)
  }
}

// 显示价值评价对话框
function showValue(comment) {
  if (!userStore.user) {
    useMsgSignIn()
    return
  }
  // 只有主贴发布人可以评价价值
  if (userStore.user.id !== topic.value?.user?.id) {
    useMsgError('只有主贴发布人可以评价价值')
    return
  }
  currentComment.value = comment
  showValueDialog.value = true
}

// 提交价值评价
async function submitValue(valueType, index) {
  if (!currentComment.value) return
  try {
    await useHttpPost(
      '/api/comment/value',
      useJsonToForm({
        commentId: currentComment.value.id,
        valueType: index
      })
    )
    useMsgSuccess(`评价成功：${valueType.label} (+${valueType.score}积分)`)
    showValueDialog.value = false
    currentComment.value = null
  } catch (e) {
    useCatchError(e)
  }
}

// AI总结
async function handleAISummary() {
  if (!topic.value) return

  if (!userStore.isLogin) {
    useMsgSignIn()
    return
  }

  aiSummaryLoading.value = true
  try {
    const topicId = topic.value.id

    // 先尝试获取已有总结
    const getResponse = await $fetch(`/api/ai-sum/topic/${topicId}/summary`)
    const existingData = getResponse.data || getResponse

    if (existingData && existingData.content) {
      // 已有总结，直接显示
      aiSummaryContent.value = existingData.content
      showAISummary.value = true
      return
    }

    // 没有已有总结，生成新的
    const genResponse = await $fetch(`/api/ai-sum/topic/${topicId}/generate-summary`, {
      method: 'POST',
      body: {}
    })

    const data = genResponse.data || genResponse

    if (data && data.content) {
      aiSummaryContent.value = data.content
      showAISummary.value = true
      useMsgSuccess('AI总结生成成功')
    } else if (data && data.success) {
      // 生成中，需要重新获取
      await new Promise(resolve => setTimeout(resolve, 2000))
      const retryResponse = await $fetch(`/api/ai-sum/topic/${topicId}/summary`)
      const retryData = retryResponse.data || retryResponse
      if (retryData && retryData.content) {
        aiSummaryContent.value = retryData.content
        showAISummary.value = true
        useMsgSuccess('AI总结生成成功')
      } else {
        useMsgError('AI总结生成中，请稍后重试')
      }
    } else {
      useMsgError('AI总结生成失败')
    }
  } catch (e) {
    console.error('AI总结错误:', e)
    useMsgError(e?.message || 'AI总结服务暂时不可用')
  } finally {
    aiSummaryLoading.value = false
  }
}

// 生成AI总结（用于重新生成，包含新评论）
async function generateAISummary() {
  if (!topic.value) return
  
  if (!userStore.isLogin) {
    useMsgSignIn()
    return
  }
  
  aiSummaryGenerating.value = true
  try {
    const topicId = topic.value.id
    
    // 调用生成总结接口
    const genResponse = await $fetch(`/api/ai-sum/topic/${topicId}/generate-summary`, {
      method: 'POST',
      body: {}
    })
    
    const data = genResponse.data || genResponse
    
    if (data && data.content) {
      aiSummaryContent.value = data.content
      useMsgSuccess('AI总结生成成功')
    } else if (data && data.success) {
      // 生成中，需要等待后重新获取
      await new Promise(resolve => setTimeout(resolve, 2000))
      const retryResponse = await $fetch(`/api/ai-sum/topic/${topicId}/summary`)
      const retryData = retryResponse.data || retryResponse
      if (retryData && retryData.content) {
        aiSummaryContent.value = retryData.content
        useMsgSuccess('AI总结生成成功')
      } else {
        useMsgError('AI总结生成中，请稍后重试')
      }
    } else {
      useMsgError('AI总结生成失败')
    }
  } catch (e) {
    console.error('AI总结生成错误:', e)
    useMsgError(e?.message || 'AI总结服务暂时不可用')
  } finally {
    aiSummaryGenerating.value = false
  }
}

// 复制AI总结内容
async function copyAISummary() {
  if (!aiSummaryContent.value) return
  
  try {
    await navigator.clipboard.writeText(aiSummaryContent.value)
    useMsgSuccess('AI总结内容已复制到剪贴板')
  } catch (e) {
    // 如果clipboard API不可用，使用备用方法
    const textArea = document.createElement('textarea')
    textArea.value = aiSummaryContent.value
    document.body.appendChild(textArea)
    textArea.select()
    try {
      document.execCommand('copy')
      useMsgSuccess('AI总结内容已复制到剪贴板')
    } catch (err) {
      useMsgError('复制失败，请手动复制内容')
    }
    document.body.removeChild(textArea)
  }
}

// 生命周期
onMounted(() => {
  loadComments()
  loadValueTypes()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})

useHead({
  title: useTopicSiteTitle(topic.value),
})
</script>

<style scoped>
.mobile-chat-container {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.chat-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  padding-bottom: 80px;
}

/* 时间分割线 */
.time-divider {
  text-align: center;
  font-size: 12px;
  color: #999;
  margin: 16px 0;
}

/* 消息项 */
.chat-message {
  display: flex;
  margin-bottom: 16px;
  animation: messageSlideIn 0.3s ease;
}

.chat-message.is-topic {
  margin-bottom: 20px;
}

.chat-message.is-me {
  flex-direction: row-reverse;
}

@keyframes messageSlideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-avatar {
  flex-shrink: 0;
  margin: 0 8px;
}

.message-content-wrapper {
  display: flex;
  flex-direction: column;
  max-width: 75%;
  min-width: 0;
}

.chat-message.is-me .message-content-wrapper {
  align-items: flex-end;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
  font-size: 12px;
}

.message-author {
  color: #666;
  font-weight: 500;
}

.message-location {
  color: #999;
  font-size: 11px;
}

/* 消息气泡 */
.message-bubble {
  background: white;
  padding: 10px 12px;
  border-radius: 12px;
  color: #333;
  word-break: break-word;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.message-bubble.is-topic {
  background: #e6f3ff;
  border: 1px solid #b3d9ff;
}

.message-bubble.is-me {
  background: #95ec69;
  color: #000;
  border-radius: 12px 4px 12px 12px;
}

.topic-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #333;
}

.message-text {
  font-size: 15px;
  line-height: 1.6;
}

.message-text :deep(img) {
  max-width: 100%;
  border-radius: 8px;
}

/* 引用 */
.message-quote {
  margin-bottom: 8px;
  padding: 8px;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 8px;
  border-left: 3px solid #07c160;
}

.quote-author {
  color: #07c160;
  font-weight: 500;
  font-size: 13px;
}

.quote-text {
  color: #666;
  font-size: 13px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 图片列表 */
.message-images {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
  margin-top: 8px;
}

.message-images img {
  width: 100%;
  aspect-ratio: 1;
  object-fit: cover;
  border-radius: 6px;
  cursor: pointer;
}

.image-wrapper.single {
  grid-column: span 2;
}

.image-wrapper.single img {
  aspect-ratio: 16/9;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px;
  background: #f5f5f5;
  border-radius: 6px;
  font-size: 13px;
}

.file-item .iconfont {
  font-size: 24px;
  color: #666;
}

.file-name {
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 标签 */
.topic-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px dashed rgba(0, 0, 0, 0.1);
}

.tag {
  padding: 2px 8px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.05);
  font-size: 12px;
  color: #666;
}

.tag.node {
  background: #07c160;
  color: white;
}

/* 消息操作 */
.message-actions {
  display: flex;
  gap: 12px;
  margin-top: 6px;
  padding: 0 4px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
  cursor: pointer;
}

.action-btn .iconfont {
  font-size: 14px;
}

.action-btn.active {
  color: #ff7875;
}

/* 加载更多 */
.load-more,
.no-more {
  text-align: center;
  padding: 16px;
  font-size: 13px;
  color: #999;
}

.load-more-btn {
  cursor: pointer;
  color: #07c160;
}

.loading-text .iconfont {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 新消息提示 */
.new-message-tip {
  position: fixed;
  bottom: 70px;
  left: 50%;
  transform: translateX(-50%);
  background: #07c160;
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  animation: bounceIn 0.3s ease;
  z-index: 100;
}

@keyframes bounceIn {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

/* 底部输入框 */
.chat-input-wrapper {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  border-top: 1px solid #eee;
  z-index: 100;
}

.quote-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: #f5f5f5;
  font-size: 13px;
  color: #666;
}

.quote-bar .iconfont {
  font-size: 16px;
  color: #999;
}

.input-bar {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  padding: 10px 12px;
}

.input-box {
  flex: 1;
  min-width: 0;
}

.input-box textarea {
  width: 100%;
  min-height: 36px;
  max-height: 80px;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 18px;
  background: #f5f5f5;
  font-size: 15px;
  resize: none;
  outline: none;
}

.input-box textarea:focus {
  border-color: #07c160;
  background: white;
}

.input-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.input-actions .iconfont {
  font-size: 24px;
  color: #666;
}

.ai-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 10px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  box-shadow: 0 2px 6px rgba(102, 126, 234, 0.3);
}

.ai-btn:active {
  transform: scale(0.95);
}

.send-btn {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 50%;
  background: #ddd;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.send-btn.can-send {
  background: #07c160;
}

.send-btn .iconfont {
  font-size: 18px;
  color: white;
}

/* AI总结对话框 */
.ai-summary-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 16px;
}

.ai-summary-dialog {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 400px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.ai-summary-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.ai-summary-dialog-header .header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ai-summary-dialog-header .iconfont {
  font-size: 20px;
  color: #999;
  cursor: pointer;
  padding: 4px;
}

.ai-generate-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
  height: 32px;
  line-height: 1;
}

.ai-generate-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.ai-generate-btn:not(:disabled):active {
  background: #5a6fd8;
  transform: scale(0.98);
}

.ai-summary-dialog-content {
  flex: 1;
  padding: 16px 20px;
  overflow-y: auto;
  min-height: 200px;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  justify-content: flex-start;
}

.ai-summary-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #666;
}

.ai-summary-loading .iconfont {
  font-size: 32px;
  color: #667eea;
  animation: rotate 1s linear infinite;
}

.ai-summary-text {
  font-size: 14px;
  line-height: 1.7;
  color: #333;
  text-align: left;
  width: 100%;
  white-space: pre-wrap;
  word-break: break-word;
}

.ai-summary-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #999;
  text-align: center;
}

.ai-summary-empty .iconfont {
  font-size: 48px;
  color: #ddd;
}

.ai-summary-dialog-actions {
  display: flex;
  gap: 12px;
  padding: 16px 20px 20px;
  border-top: 1px solid #f0f0f0;
}

.ai-summary-copy-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.ai-summary-copy-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.ai-summary-copy-btn:not(:disabled):active {
  background: #5a6fd8;
  transform: scale(0.98);
}

.ai-summary-close-btn {
  flex: 1;
  padding: 12px;
  background: #f5f5f5;
  color: #666;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.ai-summary-close-btn:active {
  background: #e8e8e8;
  transform: scale(0.98);
}

/* 加载状态 */
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.loading-state .iconfont {
  font-size: 40px;
  color: #07c160;
  animation: rotate 1s linear infinite;
}

/* 高亮动画 */
.highlight {
  animation: highlightPulse 2s ease;
}

@keyframes highlightPulse {
  0%, 100% { background: transparent; }
  50% { background: rgba(7, 193, 96, 0.1); }
}

/* 价值评价对话框 */
.value-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}

.value-dialog {
  background: white;
  border-radius: 12px;
  width: 80%;
  max-width: 300px;
  overflow: hidden;
}

.value-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #eee;
  font-size: 16px;
  font-weight: 500;
}

.value-dialog-header .iconfont {
  font-size: 20px;
  color: #999;
  cursor: pointer;
}

.value-dialog-content {
  padding: 12px;
}

.value-type-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  margin-bottom: 8px;
  border: 1px solid #eee;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.value-type-item:last-child {
  margin-bottom: 0;
}

.value-type-item:active {
  background: #f5f5f5;
  border-color: #07c160;
}

.value-type-label {
  font-size: 15px;
  color: #333;
}

.value-type-score {
  font-size: 14px;
  color: #07c160;
  font-weight: 500;
}

.action-btn.value-btn {
  color: #faad14;
}

.action-btn.value-btn:active {
  color: #d48806;
}

/* 子评论列表 */
.comment-replies {
  margin-top: 8px;
  margin-left: 32px;
  background: #f5f5f5;
  border-radius: 8px;
  padding: 10px;
}

.reply-item {
  display: flex;
  margin-bottom: 10px;
}

.reply-item:last-child {
  margin-bottom: 0;
}

.reply-avatar {
  flex-shrink: 0;
  margin-right: 8px;
}

.reply-content-wrapper {
  flex: 1;
  min-width: 0;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
}

.reply-author {
  font-size: 12px;
  color: #666;
}

.reply-location {
  font-size: 11px;
  color: #999;
}

.reply-bubble {
  background: #f8f9fa;
  padding: 8px 10px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
}

.reply-quote {
  margin-bottom: 6px;
  padding: 6px;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 6px;
  border-left: 2px solid #07c160;
  font-size: 12px;
}

.reply-quote .quote-author {
  color: #07c160;
  font-weight: 500;
}

.reply-quote .quote-text {
  color: #666;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.reply-text {
  color: #333;
  word-break: break-all;
}

.reply-images {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 4px;
  margin-top: 6px;
}

.reply-images img {
  width: 100%;
  aspect-ratio: 1;
  object-fit: cover;
  border-radius: 4px;
  cursor: pointer;
}

.reply-actions {
  display: flex;
  gap: 12px;
  margin-top: 6px;
  padding: 0 2px;
}

.reply-actions .action-btn {
  font-size: 11px;
  color: #999;
}
</style>
