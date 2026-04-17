<template>
  <section v-if="topic" class="main" :class="{ 'chat-topic-page': isMobile }">
    <div v-if="isPending" class="container main-container">
      <div class="notification is-warning" style="width: 100%; margin: 20px 0">
        {{ t('pages.topic.detail.pending') }}
      </div>
    </div>

    <!-- 移动端：聊天式话题详情 -->
    <div v-if="isMobile" class="chat-container">
      <!-- 聊天头部 -->
      <div class="chat-header">
        <div class="chat-header-info">
          <h1 class="chat-title">{{ topic.title || t('pages.topic.detail.topic') }}</h1>
          <div class="chat-meta">
            <span class="node-tag" v-if="topic.node">{{ topic.node.name }}</span>
            <span class="view-count">{{ topic.viewCount }} {{ t('pages.topic.detail.view') }}</span>
          </div>
        </div>
        <topic-manage-menu v-model="topic" />
      </div>

      <!-- 聊天消息列表 -->
      <div ref="chatMessagesRef" class="chat-messages">
        <!-- 主贴消息（楼主） -->
        <div class="chat-message-group">
          <div class="time-divider">{{ formatChatTime(topic.createTime) }}</div>
          <div class="chat-message" :class="{ 'is-me': false }">
            <div class="message-avatar">
              <my-avatar :user="topic.user" :size="40" />
            </div>
            <div class="message-content-wrapper">
              <div class="message-header">
                <nuxt-link :to="`/user/${topic.user.id}`" class="message-author">
                  {{ topic.user.nickname }}
                </nuxt-link>
                <span v-if="topic.ipLocation" class="message-location">{{ topic.ipLocation }}</span>
              </div>
              <div class="message-bubble" :class="{ 'is-topic': true }">
                <div class="message-text content" v-html="processedContent" />
                <ul
                  v-if="topic.imageList && topic.imageList.length"
                  class="message-image-list"
                >
                  <li v-for="(image, index) in topic.imageList" :key="image.url">
                    <div class="image-item" v-if="isImageFile(image)">
                      <el-image
                        :src="image.preview"
                        :preview-src-list="imageUrls"
                        :initial-index="index"
                      />
                    </div>
                    <div class="file-item" v-else @click="openFile(image.url)">
                      <i class="file-icon" :class="getFileIconClass(image)"></i>
                      <span class="file-name">{{ image.filename || getFileName(image.url) }}</span>
                    </div>
                  </li>
                </ul>
                <div
                  v-if="hideContent && hideContent.exists"
                  class="hide-content"
                >
                  <div v-if="hideContent.show" class="widget has-border">
                    <div class="widget-header">
                      <span>
                        <i class="iconfont icon-lock" />
                        <span>&nbsp;{{ t('pages.topic.detail.hideContent') }}</span>
                      </span>
                    </div>
                    <div class="widget-content" v-html="hideContent.content" />
                  </div>
                  <div v-else class="hide-content-tip">
                    <i class="iconfont icon-lock" />
                    <span>{{ t('pages.topic.detail.hideContentTip') }}</span>
                  </div>
                </div>
              </div>
              <div class="message-actions">
                <span class="action-btn" :class="{ active: liked }" @click="like">
                  <i class="iconfont icon-like" />
                  <span>{{ topic.likeCount || t('pages.topic.detail.like') }}</span>
                </span>
                <span class="action-btn" :class="{ active: topic.favorited }" @click="addFavorite(topic.id)">
                  <i class="iconfont" :class="topic.favorited ? 'icon-has-favorite' : 'icon-favorite'" />
                  <span>{{ topic.favorited ? t('pages.topic.detail.favorited') : t('pages.topic.detail.favorite') }}</span>
                </span>
                <span class="action-btn" @click="scrollToBottom">
                  <i class="iconfont icon-comment" />
                  <span>{{ topic.commentCount || t('pages.topic.detail.comment') }}</span>
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 评论消息列表 -->
        <div v-for="(group, groupIndex) in groupedComments" :key="groupIndex" class="chat-message-group">
          <div class="time-divider">{{ group.timeLabel }}</div>
          <div
            v-for="comment in group.comments"
            :key="comment.id"
            class="chat-message"
            :class="{ 'is-me': isCurrentUser(comment.user.id) }"
            :id="`comment-${comment.id}`"
          >
            <div class="message-avatar">
              <my-avatar :user="comment.user" :size="36" />
            </div>
            <div class="message-content-wrapper">
              <div class="message-header">
                <nuxt-link :to="`/user/${comment.user.id}`" class="message-author">
                  {{ comment.user.nickname }}
                </nuxt-link>
                <span v-if="comment.ipLocation" class="message-location">{{ comment.ipLocation }}</span>
              </div>
              <div
                class="message-bubble"
                :class="{ 'is-me': isCurrentUser(comment.user.id) }"
                @click="handleQuoteClick(comment)"
              >
                <!-- 引用内容 -->
                <div v-if="comment.quote" class="message-quote">
                  <div class="quote-line" />
                  <div class="quote-content">
                    <span class="quote-author">{{ comment.quote.user.nickname }}:</span>
                    <span class="quote-text" v-html="comment.quote.content" />
                  </div>
                </div>
                <div class="message-text" v-html="comment.content" />
                <div
                  v-if="comment.imageList && comment.imageList.length"
                  class="message-image-list"
                >
                  <el-image
                    v-for="(image, imageIndex) in comment.imageList"
                    :key="imageIndex"
                    :src="image.preview"
                    :preview-src-list="comment.imageList.map(img => img.url)"
                    :initial-index="imageIndex"
                    fit="cover"
                  />
                </div>
              </div>
              <div class="message-actions">
                <span class="action-btn" :class="{ active: comment.liked }" @click="likeComment(comment)">
                  <i class="iconfont icon-like" />
                  <span>{{ comment.likeCount > 0 ? comment.likeCount : t('component.comment.list.like') }}</span>
                </span>
                <span class="action-btn" @click="replyTo(comment)">
                  <i class="iconfont icon-comment" />
                  <span>{{ t('component.comment.list.reply') }}</span>
                </span>
                <!-- 价值按钮（仅楼主可见） -->
                <span
                  v-if="valueTypes.length > 0"
                  class="action-btn value-btn"
                  @click="showValue(comment)"
                >
                  <i class="iconfont icon-star" />
                  <span>价值</span>
                </span>
                <span class="action-btn" @click="showCommentDetail(comment)">
                  <i class="iconfont icon-more" />
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
                    <my-avatar :user="reply.user" :size="28" />
                  </div>
                  <div class="reply-content-wrapper">
                    <div class="reply-header">
                      <nuxt-link :to="`/user/${reply.user.id}`" class="reply-author">
                        {{ reply.user.nickname }}
                      </nuxt-link>
                      <span v-if="reply.ipLocation" class="reply-location">{{ reply.ipLocation }}</span>
                    </div>
                    <div class="reply-bubble">
                      <div v-if="reply.quote" class="reply-quote">
                        <span class="quote-author">{{ reply.quote.user.nickname }}:</span>
                        <span class="quote-text" v-html="reply.quote.content"></span>
                      </div>
                      <div class="reply-text" v-html="reply.content"></div>
                      <div v-if="reply.imageList && reply.imageList.length" class="reply-image-list">
                        <el-image
                          v-for="(image, imageIndex) in reply.imageList"
                          :key="imageIndex"
                          :src="image.preview"
                          :preview-src-list="reply.imageList.map(img => img.url)"
                          :initial-index="imageIndex"
                          fit="cover"
                        />
                      </div>
                    </div>
                    <div class="reply-actions">
                      <span class="action-btn" :class="{ active: reply.liked }" @click="likeComment(reply)">
                        <i class="iconfont icon-like" />
                        <span>{{ reply.likeCount > 0 ? reply.likeCount : t('component.comment.list.like') }}</span>
                      </span>
                      <span class="action-btn" @click="replyTo(reply)">
                        <i class="iconfont icon-comment" />
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
        <div v-if="hasMoreComments" class="load-more-wrapper">
          <el-button :loading="loadingMore" @click="loadMoreComments" type="primary" text>
            {{ t('component.comment.list.loadMore') }}
          </el-button>
        </div>

        <!-- 新消息提示 -->
        <div v-if="newMessageCount > 0" class="new-message-tip" @click="scrollToBottom">
          <i class="iconfont icon-down" />
          <span>{{ newMessageCount }} {{ t('component.chat.newMessages') }}</span>
        </div>
      </div>

      <!-- 底部固定输入框 -->
      <div class="chat-input-wrapper">
        <div v-if="replyToComment" class="quote-bar">
          <span class="quote-text">{{ t('component.comment.input.replyTo') }} {{ replyToComment.user.nickname }}</span>
          <i class="iconfont icon-close" @click="cancelReply" />
        </div>
        <div class="chat-input-bar">
          <div class="input-actions">
            <i class="iconfont icon-image" @click="triggerImageUpload" />
            <i class="iconfont icon-emoji" @click="showEmojiPicker" />
          </div>
          <div class="input-box">
            <textarea
              ref="inputRef"
              v-model="inputContent"
              :placeholder="t('component.chat.inputPlaceholder')"
              @keydown.enter.prevent="sendMessage"
              @input="onInputChange"
              rows="1"
            />
          </div>
          <div class="send-btn">
            <el-button
              type="primary"
              :disabled="!canSend"
              :loading="sending"
              @click="sendMessage"
            >
              {{ t('component.chat.send') }}
            </el-button>
          </div>
        </div>
        <image-upload
          v-show="showImageUpload"
          ref="imageUploader"
          v-model="inputImages"
          class="image-upload-panel"
        />
      </div>
    </div>

    <!-- 桌面端：传统论坛风格 -->
    <div v-if="!isMobile" class="container main-container left-main size-360">
      <div class="left-container">
        <div class="main-content no-padding no-bg">
          <article class="topic-detail">
            <side-action-bar
              class="float-bar"
              entity-type="topic"
              :entity-id="topic.id"
              :liked="liked"
              :like-count="topic.likeCount"
              :comment-count="topic.commentCount"
              :favorited="topic.favorited"
            />
            <div class="topic-header">
              <div class="topic-header-left">
                <my-avatar :user="topic.user" :size="45" />
              </div>
              <div class="topic-header-center">
                <div class="topic-nickname">
                  <nuxt-link :to="`/user/${topic.user.id}`">
                    {{ topic.user.nickname }}
                  </nuxt-link>
                </div>
                <div class="topic-meta">
                  <span class="meta-item">
                    {{ t('pages.topic.detail.publishedAt') }}
                    <time>{{ usePrettyDate(topic.createTime, t) }}</time>
                  </span>
                  <span v-if="topic.ipLocation" class="meta-item">
                    {{ t('pages.topic.detail.ipLocation') }}{{ topic.ipLocation }}
                  </span>
                </div>
              </div>
              <div class="topic-header-right">
                <topic-manage-menu v-model="topic" />
              </div>
            </div>

            <!-- 内容 -->
            <div
              class="topic-content content"
              :class="{
                'topic-tweet': topic.type === 1,
              }"
            >
              <h1 v-if="topic.title" class="topic-title">
                {{ topic.title }}
                <span v-if="topic.needAHand" class="need-a-hand-tag">
                  <i class="iconfont icon-help" />
                  {{ t('pages.topic.detail.needAHand') }}
                </span>
              </h1>
              <div
              class="topic-content-detail line-numbers"
              v-html="processedContent"
            />
              <ul
                v-if="topic.imageList && topic.imageList.length"
                class="topic-image-list"
              >
                <li v-for="(image, index) in topic.imageList" :key="image.url">
                  <div class="image-item" v-if="isImageFile(image)">
                    <el-image
                      :src="image.preview"
                      :preview-src-list="imageUrls"
                      :initial-index="index"
                    />
                  </div>
                  <div class="file-item" v-else @dblclick="openFile(image.url)">
                    <i class="file-icon" :class="getFileIconClass(image)"></i>
                    <span class="file-name">{{ image.filename || getFileName(image.url) }}</span>
                  </div>
                </li>
              </ul>
              <div
                v-if="hideContent && hideContent.exists"
                class="topic-content-detail hide-content"
              >
                <div v-if="hideContent.show" class="widget has-border">
                  <div class="widget-header">
                    <span>
                      <i class="iconfont icon-lock" />
                      <span>&nbsp;{{ t('pages.topic.detail.hideContent') }}</span>
                    </span>
                  </div>
                  <div class="widget-content" v-html="hideContent.content" />
                </div>
                <div v-else class="hide-content-tip">
                  <i class="iconfont icon-lock" />
                  <span>{{ t('pages.topic.detail.hideContentTip') }}</span>
                </div>
              </div>
            </div>

            <!-- 节点、标签 -->
            <div class="topic-tags">
              <nuxt-link
                v-if="topic.node"
                :to="`/topics/node/${topic.node.id}`"
                class="topic-tag"
              >
                {{ topic.node.name }}
              </nuxt-link>
              <nuxt-link
                v-for="tag in topic.tags"
                :key="tag.id"
                :to="`/topics/tag/${tag.id}`"
                class="topic-tag"
              >
                #{{ tag.name }}
              </nuxt-link>
            </div>

            <!-- 点赞用户列表 -->
            <div v-if="likeUsers && likeUsers.length" class="topic-like-users">
              <my-avatar
                v-for="likeUser in likeUsers"
                :key="likeUser.id"
                :user="likeUser"
                :size="24"
                has-border
              />
              <span class="like-count">{{ topic.likeCount }}</span>
            </div>

            <!-- 功能按钮 -->
            <div class="topic-actions">
              <div class="action disabled">
                <i class="action-icon iconfont icon-view" />
                <div class="action-text">
                  <span>{{ t('pages.topic.detail.view') }}</span>
                  <span v-if="topic.viewCount > 0" class="action-text">
                    ({{ topic.viewCount }})
                  </span>
                </div>
              </div>
              <div class="action" @click="like">
                <i
                  class="action-icon iconfont icon-like"
                  :class="{ 'checked-icon': liked }"
                />
                <div class="action-text">
                  <span>{{ t('pages.topic.detail.like') }}</span>
                  <span v-if="topic.likeCount > 0">
                    ({{ topic.likeCount }})
                  </span>
                </div>
              </div>
              <div class="action" @click="addFavorite(topic.id)">
                <i
                  class="action-icon iconfont icon-favorite"
                  :class="{
                    'icon-has-favorite': topic.favorited,
                    'icon-favorite': !topic.favorited,
                    'checked-icon': topic.favorited,
                  }"
                />
                <div class="action-text">
                  <span>{{ t('pages.topic.detail.favorite') }}</span>
                </div>
              </div>
              <div class="action test-email-action" @click="sendTestEmail">
                <i class="action-icon iconfont icon-email" />
                <div class="action-text">
                  <span>测试邮件</span>
                </div>
              </div>
            </div>
          </article>

          <!-- 评论 -->
          <comment
            :entity-id="topic.id"
            :comment-count="topic.commentCount"
            entity-type="topic"
            @created="commentCreated"
          />
        </div>
      </div>
      <div class="right-container">
        <AISummarySidebar :topic-id="topic.id" />
      </div>
    </div>

    <!-- 评论详情弹窗（仅移动端使用） -->
    <el-dialog
      v-model="commentDetailVisible"
      :title="t('component.chat.commentDetail')"
      width="500px"
      class="comment-detail-dialog"
    >
      <div v-if="currentComment" class="comment-detail-content">
        <div class="detail-header">
          <my-avatar :user="currentComment.user" :size="40" />
          <div class="detail-info">
            <nuxt-link :to="`/user/${currentComment.user.id}`">{{ currentComment.user.nickname }}</nuxt-link>
            <span class="detail-time">{{ usePrettyDate(currentComment.createTime) }}</span>
          </div>
        </div>
        <div class="detail-body content" v-html="currentComment.content" />
        <div class="detail-actions">
          <el-button type="primary" @click="replyTo(currentComment); commentDetailVisible = false;">
            {{ t('component.comment.list.reply') }}
          </el-button>
          <el-button v-if="currentComment.liked" @click="likeComment(currentComment); commentDetailVisible = false;">
            {{ t('component.comment.list.liked') }} ({{ currentComment.likeCount }})
          </el-button>
          <el-button v-else @click="likeComment(currentComment); commentDetailVisible = false;">
            {{ t('component.comment.list.like') }}
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 价值评价弹窗 -->
    <el-dialog
      v-model="showValueDialog"
      title="价值评价"
      width="400px"
      :close-on-click-modal="false"
      class="value-dialog"
    >
      <div class="value-types">
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
      <template #footer>
        <el-button @click="showValueDialog = false">取消</el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup>
import { useI18n } from 'vue-i18n';
import '~/assets/css/file-icons.css';
import AISummarySidebar from '~/components/ai/AISummarySidebar.vue';

// 注册组件
const { t } = useI18n();
const route = useRoute();
const userStore = useUserStore();

// 设备检测（仅客户端）
const isMobile = ref(false);

// 话题数据
const { data: topic } = await useMyFetch(`/api/topic/${route.params.id}`);
const { data: liked } = await useMyFetch("/api/like/liked", {
  params: { entityType: "topic", entityId: route.params.id },
});
const { data: likeUsers, refresh: refreshLikeUsers } = await useMyFetch(
  `/api/topic/recentlikes/${route.params.id}`
);
const { data: hideContent, refresh: refreshHideContent } = await useMyFetch(
  `/api/topic/hide_content?topicId=${route.params.id}`
);

// 评论相关（仅移动端使用）
const comments = ref([]);
const commentCursor = ref('');
const hasMoreComments = ref(true);
const loadingMore = ref(false);
const newMessageCount = ref(0);

// 输入相关（仅移动端使用）
const inputContent = ref('');
const inputImages = ref([]);
const inputRef = ref(null);
const imageUploader = ref(null);
const sending = ref(false);
const showImageUpload = ref(false);
const replyToComment = ref(null);
const chatMessagesRef = ref(null);

// 弹窗相关
const commentDetailVisible = ref(false);
const currentComment = ref(null);

// 价值评价相关
const showValueDialog = ref(false);
const currentValueComment = ref(null);
const valueTypes = ref([]);

// 轮询
let pollingInterval = null;
let lastCommentId = null;

// 计算属性
const isPending = computed(() => topic.value?.status === 2);

const canSend = computed(() => {
  return inputContent.value.trim().length > 0 || inputImages.value.length > 0;
});

const imageUrls = computed(() => {
  if (!topic.value?.imageList) return [];
  return topic.value.imageList.map(img => img.url);
});

// 按时间分组的评论（仅移动端使用）
const groupedComments = computed(() => {
  const groups = [];
  let currentGroup = null;

  comments.value.forEach(comment => {
    const timeLabel = formatChatTime(comment.createTime);

    if (!currentGroup || currentGroup.timeLabel !== timeLabel) {
      currentGroup = { timeLabel, comments: [] };
      groups.push(currentGroup);
    }
    currentGroup.comments.push(comment);
  });

  return groups;
});

// 设备检测
onMounted(() => {
  if (import.meta.client && typeof navigator !== 'undefined') {
    const ua = navigator.userAgent;
    // 只检测设备类型，不包含企业微信（避免企业微信电脑端被误判为移动端）
    const isMobileDevice = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(ua);
    // 只有在真正的移动设备或移动浏览器才使用移动端布局
    isMobile.value = isMobileDevice && !/wxwork/i.test(ua);

    // 如果是移动端，加载评论
    console.log('UserAgent:', ua);
    console.log('isMobile:', isMobile.value);
    // 总是加载价值类型配置（用于调试）
    loadValueTypes();
    if (isMobile.value) {
      loadComments();
      startPolling();

      // 监听滚动显示/隐藏新消息提示
      if (chatMessagesRef.value) {
        chatMessagesRef.value.addEventListener('scroll', () => {
          const container = chatMessagesRef.value;
          const isAtBottom = container.scrollHeight - container.scrollTop <= container.clientHeight + 50;
          if (isAtBottom) {
            newMessageCount.value = 0;
          }
        });
      }
    }
  }
});

onUnmounted(() => {
  stopPolling();
});

useHead({
  title: useTopicSiteTitle(topic.value),
});

// ========== 通用方法 ==========

// 判断是否为图片文件
function isImageFile(file) {
  if (typeof file === 'object' && file !== null && file.isImage !== undefined) {
    return file.isImage;
  }
  const url = typeof file === 'object' && file !== null ? file.url : file;
  if (!url) return false;
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp'];
  const lowerUrl = url.toLowerCase();
  return imageExtensions.some(ext => lowerUrl.endsWith(ext));
}

// 获取文件图标类名
function getFileIconClass(fileObj) {
  if (!fileObj) return 'file-icon file-icon-default';

  if (typeof fileObj === 'object' && fileObj !== null && fileObj.fileExt) {
    const ext = fileObj.fileExt.toLowerCase();
    return getIconClassByExtension(ext);
  }

  if (typeof fileObj === 'object' && fileObj !== null && fileObj.url) {
    const urlStr = fileObj.url.toLowerCase();
    const ext = '.' + urlStr.split('.').pop();
    return getIconClassByExtension(ext);
  }

  const lowerUrl = fileObj.toString().toLowerCase();
  const ext = '.' + lowerUrl.split('.').pop();
  return getIconClassByExtension(ext);
}

// 根据文件扩展名获取图标类名
function getIconClassByExtension(ext) {
  const extMap = {
    '.doc': 'file-icon file-icon-word',
    '.docx': 'file-icon file-icon-word',
    '.xls': 'file-icon file-icon-excel',
    '.xlsx': 'file-icon file-icon-excel',
    '.csv': 'file-icon file-icon-excel',
    '.ppt': 'file-icon file-icon-ppt',
    '.pptx': 'file-icon file-icon-ppt',
    '.pdf': 'file-icon file-icon-pdf',
    '.txt': 'file-icon file-icon-text',
    '.log': 'file-icon file-icon-text',
    '.md': 'file-icon file-icon-text',
    '.zip': 'file-icon file-icon-zip',
    '.rar': 'file-icon file-icon-zip',
    '.7z': 'file-icon file-icon-zip',
    '.tar': 'file-icon file-icon-zip',
    '.gz': 'file-icon file-icon-zip',
    '.jpg': 'file-icon file-icon-image',
    '.jpeg': 'file-icon file-icon-image',
    '.png': 'file-icon file-icon-image',
    '.gif': 'file-icon file-icon-image',
    '.bmp': 'file-icon file-icon-image',
    '.webp': 'file-icon file-icon-image',
    '.svg': 'file-icon file-icon-image',
    '.mp3': 'file-icon file-icon-audio',
    '.wav': 'file-icon file-icon-audio',
    '.flac': 'file-icon file-icon-audio',
    '.aac': 'file-icon file-icon-audio',
    '.ogg': 'file-icon file-icon-audio',
    '.mp4': 'file-icon file-icon-video',
    '.avi': 'file-icon file-icon-video',
    '.mkv': 'file-icon file-icon-video',
    '.mov': 'file-icon file-icon-video',
    '.wmv': 'file-icon file-icon-video',
    '.flv': 'file-icon file-icon-video',
    '.html': 'file-icon file-icon-code',
    '.htm': 'file-icon file-icon-code',
    '.css': 'file-icon file-icon-code',
    '.js': 'file-icon file-icon-code',
    '.ts': 'file-icon file-icon-code',
    '.json': 'file-icon file-icon-code',
    '.xml': 'file-icon file-icon-code',
    '.py': 'file-icon file-icon-code',
    '.java': 'file-icon file-icon-code',
    '.cpp': 'file-icon file-icon-code',
    '.c': 'file-icon file-icon-code',
    '.php': 'file-icon file-icon-code',
    '.go': 'file-icon file-icon-code',
    '.rb': 'file-icon file-icon-code',
    '.sql': 'file-icon file-icon-code',
  };

  return extMap[ext] || 'file-icon file-icon-default';
}

// 处理内容中的文件链接
function processFileLinks(content) {
  if (!content) return '';
  const fileLinkRegex = /<a\s+[^>]*href="([^"]*uploads\/files[^"]*\.([a-zA-Z0-9]+))"[^>]*>(.*?)<\/a>/gi;
  return content.replace(fileLinkRegex, (_, href, ext, linkText) => {
    const fileExt = '.' + ext.toLowerCase();
    const iconClass = getFileIconClass(fileExt);
    return `<span class="file-link-wrapper">
      <i class="${iconClass}"></i>
      <a href="${href}" rel="nofollow" title="${linkText}">${linkText}</a>
    </span>`;
  });
}

// 处理后的帖子内容
const processedContent = computed(() => {
  return topic.value ? processFileLinks(topic.value.content) : '';
});

// 获取文件名
function getFileName(url) {
  if (!url) return '';
  const parts = url.split('/');
  return parts[parts.length - 1];
}

// 双击打开文件
function openFile(url) {
  window.open(url, '_blank');
}

// 点赞
async function like() {
  try {
    if (liked.value) {
      await useHttpPost("/api/like/unlike", useJsonToForm({
        entityType: "topic",
        entityId: topic.value.id,
      }));
      liked.value = false;
      topic.value.likeCount = topic.value.likeCount > 0 ? topic.value.likeCount - 1 : 0;
      useMsgSuccess(t('pages.topic.detail.likeSuccess'));
      await refreshLikeUsers();
    } else {
      await useHttpPost("/api/like/like", useJsonToForm({
        entityType: "topic",
        entityId: topic.value.id,
      }));
      liked.value = true;
      topic.value.likeCount++;
      useMsgSuccess(t('pages.topic.detail.likeSuccess'));
      await refreshLikeUsers();
    }
  } catch (e) {
    useCatchError(e);
  }
}

// 收藏
async function addFavorite(topicId) {
  try {
    if (topic.value.favorited) {
      await useHttpPost("/api/favorite/delete", useJsonToForm({
        entityType: "topic",
        entityId: topicId,
      }));
      topic.value.favorited = false;
      useMsgSuccess(t('pages.topic.detail.favoriteSuccess'));
    } else {
      await useHttpPost("/api/favorite/add", useJsonToForm({
        entityType: "topic",
        entityId: topicId,
      }));
      topic.value.favorited = true;
      useMsgSuccess(t('pages.topic.detail.favoriteSuccess'));
    }
  } catch (e) {
    useCatchError(e);
  }
}

// 评论创建回调（桌面端）
async function commentCreated() {
  refreshHideContent();
}

// ========== 移动端专用方法 ==========

function isCurrentUser(userId) {
  return userStore.user && userStore.user.id === userId;
}

function formatChatTime(time) {
  const date = new Date(time);
  const now = new Date();
  const diff = now - date;
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);

  if (minutes < 1) return t('component.chat.justNow');
  if (minutes < 60) return t('component.chat.minutesAgo', { n: minutes });
  if (hours < 24) return t('component.chat.hoursAgo', { n: hours });
  if (days < 7) return t('component.chat.daysAgo', { n: days });

  return date.toLocaleDateString();
}

// 处理评论数据，将子评论挂载到父评论下
function processComments(results) {
  if (!results || results.length === 0) return [];

  const mainComments = [];
  const repliesMap = new Map();

  // 第一遍：分离主评论和子评论
  results.forEach(comment => {
    if (comment.quoteId) {
      // 这是子评论，需要挂载到父评论下
      if (!repliesMap.has(comment.quoteId)) {
        repliesMap.set(comment.quoteId, []);
      }
      repliesMap.get(comment.quoteId).push(comment);
    } else {
      // 这是主评论
      mainComments.push(comment);
    }
  });

  // 第二遍：将子评论挂载到对应的主评论
  mainComments.forEach(comment => {
    if (repliesMap.has(comment.id)) {
      comment.replies = {
        results: repliesMap.get(comment.id),
        hasMore: false
      };
    }
  });

  return mainComments;
}

// 加载评论
async function loadComments(cursor = '') {
  try {
    const params = { entityType: 'topic', entityId: route.params.id };
    if (cursor) params.cursor = cursor;

    const result = await useHttpGet('/api/comment/comments', { params });

    // 处理评论数据，将子评论挂载到父评论下
    const processedResults = processComments(result.results || []);

    if (cursor) {
      comments.value.push(...processedResults);
    } else {
      comments.value = processedResults;
    }

    commentCursor.value = result.cursor;
    hasMoreComments.value = result.hasMore;

    if (comments.value.length > 0) {
      lastCommentId = comments.value[comments.value.length - 1].id;
    }
  } catch (e) {
    console.error('Failed to load comments:', e);
  }
}

async function loadMoreComments() {
  if (loadingMore.value || !hasMoreComments.value) return;
  loadingMore.value = true;
  await loadComments(commentCursor.value);
  loadingMore.value = false;
}

// 发送测试邮件 - 完全拷贝server\cmd\test\main.go的功能
async function sendTestEmail() {
  try {
    console.log('测试邮件按钮被点击了！');
    console.log('开始发送测试邮件...');
    
    // 调用后端API发送测试邮件
    const result = await useHttpPost('/api/topic/test-email', {
      smtpHost: 'mail01.tkmold.com',
      smtpPort: 25,
      fromEmail: 'aibbs@tkmold.com',
      password: 'TK.com2026',
      toEmail: 'scottchen@tkmold.com',
      subject: 'Go邮件测试',
      body: 'Go邮件来了'
    });
    
    if (result) {
      alert('测试邮件发送成功！');
      console.log('测试邮件发送成功！');
    } else {
      alert('测试邮件发送失败，请查看控制台日志');
      console.error('测试邮件发送失败');
    }
  } catch (error) {
    alert('测试邮件发送失败：' + error.message);
    console.error('测试邮件发送失败：', error);
  }
}

// 新增：发送测试邮件（无认证）
async function sendTestEmailNoAuth() {
  try {
    console.log('测试邮件(无认证)按钮被点击了！');
    console.log('开始发送测试邮件(无认证)...');
    
    // 调用后端API发送测试邮件（无认证）
    const result = await useHttpPost('/api/topic/test-email-no-auth', {
      smtpHost: 'mail01.tkmold.com',
      smtpPort: 25,
      fromEmail: 'aibbs@tkmold.com',
      toEmail: 'scottchen@tkmold.com',
      subject: 'Go邮件测试(无认证)',
      body: 'Go邮件来了(无认证)'
    });
    
    if (result) {
      alert('测试邮件(无认证)发送成功！');
      console.log('测试邮件(无认证)发送成功！');
    } else {
      alert('测试邮件(无认证)发送失败，请查看控制台日志');
      console.error('测试邮件(无认证)发送失败');
    }
  } catch (error) {
    alert('测试邮件(无认证)发送失败：' + error.message);
    console.error('测试邮件(无认证)发送失败：', error);
  }
}

// 临时调试函数
function debugFunction() {
  console.log('临时调试按钮被点击了！');
  console.log('当前topic信息：', topic.value);
  console.log('当前用户：', userStore.user);
  console.log('isMobile：', isMobile.value);
  console.log('valueTypes：', valueTypes.value);
  alert('调试信息已打印到控制台，请查看！');
}

// 新增：在移动端聊天输入框上方插入测试邮件按钮
function addTestEmailButton() {
  // 仅在移动端且已登录时显示
  if (!isMobile.value || !userStore.user) return;

  nextTick(() => {
    const inputBar = document.querySelector('.chat-input-bar');
    if (!inputBar) return;

    // 避免重复插入
    if (document.querySelector('.test-email-bar')) return;

    const testBar = document.createElement('div');
    testBar.className = 'test-email-bar';
    testBar.innerHTML = `
      <button class="test-email-btn" @click="sendTestEmail">
        <i class="iconfont icon-email"></i>
        发送测试邮件
      </button>
      <button class="test-email-btn no-auth" @click="sendTestEmailNoAuth">
        <i class="iconfont icon-email"></i>
        无认证测试邮件
      </button>
    `;

    // 插入到输入框上方
    inputBar.parentNode.insertBefore(testBar, inputBar);

    // 绑定点击事件
    testBar.querySelector('.test-email-btn:not(.no-auth)').addEventListener('click', sendTestEmail);
    testBar.querySelector('.test-email-btn.no-auth').addEventListener('click', sendTestEmailNoAuth);
  });
}

// 监听登录状态变化，动态插入按钮
watchEffect(() => {
  addTestEmailButton();
});

// 发送消息
async function sendMessage() {
  if (!canSend.value || sending.value) return;
  if (!userStore.isLogin) {
    useMsgSignIn();
    return;
  }

  sending.value = true;
  try {
    const data = await useHttpPost('/api/comment/create', useJsonToForm({
      entityType: 'topic',
      entityId: route.params.id,
      content: inputContent.value,
      imageList: inputImages.value.length ? JSON.stringify(inputImages.value) : '',
      quoteId: replyToComment.value?.id || '',
    }));

    if (!comments.value) {
      comments.value = [];
    }

    // 处理新发布的评论：子评论挂载到父评论下，主评论直接追加
    if (data.quoteId) {
      // 这是子评论，找到父评论并挂载
      const parentComment = comments.value.find(c => c.id === data.quoteId);
      if (parentComment) {
        if (!parentComment.replies) {
          parentComment.replies = { results: [], hasMore: false };
        }
        parentComment.replies.results.push(data);
      } else {
        // 如果找不到父评论，直接追加（兜底）
        comments.value.push(data);
      }
    } else {
      // 这是主评论，直接追加
      comments.value.push(data);
    }

    inputContent.value = '';
    inputImages.value = [];
    replyToComment.value = null;
    showImageUpload.value = false;
    topic.value.commentCount++;

    refreshHideContent();

    useMsgSuccess(t('component.comment.input.publishSuccess'));

    nextTick(() => scrollToBottom());
  } catch (e) {
    useCatchError(e);
  } finally {
    sending.value = false;
  }
}

// 回复相关
function replyTo(comment) {
  replyToComment.value = comment;
  inputRef.value?.focus();
}

function cancelReply() {
  replyToComment.value = null;
}

function handleQuoteClick(comment) {
  if (comment.quote) {
    const el = document.getElementById(`comment-${comment.quote.id}`);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'center' });
      el.classList.add('highlight');
      setTimeout(() => el.classList.remove('highlight'), 2000);
    }
  }
}

// 点赞评论
async function likeComment(comment) {
  try {
    if (comment.liked) {
      await useHttpPost('/api/like/unlike', useJsonToForm({
        entityType: 'comment',
        entityId: comment.id,
      }));
      comment.liked = false;
      comment.likeCount = Math.max(0, comment.likeCount - 1);
    } else {
      await useHttpPost('/api/like/like', useJsonToForm({
        entityType: 'comment',
        entityId: comment.id,
      }));
      comment.liked = true;
      comment.likeCount++;
    }
  } catch (e) {
    useCatchError(e);
  }
}

// UI 方法
function scrollToBottom() {
  if (chatMessagesRef.value) {
    chatMessagesRef.value.scrollTop = chatMessagesRef.value.scrollHeight;
    newMessageCount.value = 0;
  }
}

function triggerImageUpload() {
  showImageUpload.value = !showImageUpload.value;
  if (showImageUpload.value) {
    imageUploader.value?.onClick();
  }
}

function showEmojiPicker() {
  useMsgInfo(t('component.chat.emojiComingSoon'));
}

function onInputChange() {
  const textarea = inputRef.value;
  if (textarea) {
    textarea.style.height = 'auto';
    textarea.style.height = Math.min(textarea.scrollHeight, 120) + 'px';
  }
}

function showCommentDetail(comment) {
  currentComment.value = comment;
  commentDetailVisible.value = true;
}

// 加载价值类型配置
async function loadValueTypes() {
  console.log('loadValueTypes called');
  try {
    const config = await useHttpGet('/api/sys-config/configs');
    console.log('loadValueTypes config:', config);
    if (config && config.valueTypes && config.valueTypes.length > 0) {
      valueTypes.value = config.valueTypes;
      console.log('valueTypes loaded:', valueTypes.value);
    } else {
      console.log('No valueTypes in config');
    }
  } catch (e) {
    console.error('Failed to load value types:', e);
  }
}

// 判断当前用户是否是主贴发布人
function isTopicOwner() {
  console.log('isTopicOwner check:', {
    hasUser: !!userStore.user,
    hasTopic: !!topic.value,
    hasTopicUser: !!(topic.value && topic.value.user),
    userId: userStore.user?.id,
    topicUserId: topic.value?.user?.id,
    valueTypesCount: valueTypes.value.length
  });
  if (!userStore.user || !topic.value || !topic.value.user) {
    return false;
  }
  // 使用宽松相等，与桌面端保持一致
  const result = userStore.user.id == topic.value.user.id;
  console.log('isTopicOwner result:', result, userStore.user.id, topic.value.user.id);
  return result;
}

// 显示价值评价对话框
function showValue(comment) {
  if (!userStore.user) {
    useMsgSignIn();
    return;
  }

  // 只有主贴发布人可以评价价值
  if (!isTopicOwner()) {
    useMsgError('只有主贴发布人可以评价价值');
    return;
  }

  currentValueComment.value = comment;
  showValueDialog.value = true;
}

// 提交价值评价
async function submitValue(valueType, index) {
  try {
    await useHttpPost('/api/comment/value', useJsonToForm({
      commentId: currentValueComment.value.id,
      valueType: index
    }));

    useMsgSuccess(`评价成功：${valueType.label} (+${valueType.score}积分)`);
    showValueDialog.value = false;
    currentValueComment.value = null;
  } catch (e) {
    useCatchError(e);
  }
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
      });

      if (result.results && result.results.length > 0) {
        const newComments = result.results.filter(c =>
          !comments.value.some(existing => existing.id === c.id) &&
          !comments.value.some(existing =>
            existing.replies?.results?.some(reply => reply.id === c.id)
          )
        );

        if (newComments.length > 0) {
          // 处理新评论：子评论挂载到父评论下
          newComments.forEach(comment => {
            if (comment.quoteId) {
              // 子评论
              const parentComment = comments.value.find(c => c.id === comment.quoteId);
              if (parentComment) {
                if (!parentComment.replies) {
                  parentComment.replies = { results: [], hasMore: false };
                }
                parentComment.replies.results.push(comment);
              } else {
                comments.value.push(comment);
              }
            } else {
              // 主评论
              comments.value.push(comment);
            }
          });

          lastCommentId = result.results[result.results.length - 1].id;

          const container = chatMessagesRef.value;
          const isAtBottom = container.scrollHeight - container.scrollTop <= container.clientHeight + 100;

          if (isAtBottom) {
            nextTick(() => scrollToBottom());
          } else {
            newMessageCount.value += newComments.length;
          }
        }
      }
    } catch (e) {
      console.error('Polling error:', e);
    }
  }, 5000);
}

function stopPolling() {
  if (pollingInterval) {
    clearInterval(pollingInterval);
    pollingInterval = null;
  }
}
</script>

<style lang="scss" scoped>
// ========== 移动端聊天式风格 ==========
.chat-topic-page {
  height: 100vh;
  overflow: hidden;
}

.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  max-width: 900px;
  margin: 0 auto;
  background: var(--bg-color);
}

// 头部
.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-color);
  flex-shrink: 0;

  .chat-header-info {
    flex: 1;
    min-width: 0;
  }

  .chat-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-color);
    margin: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .chat-meta {
    display: flex;
    gap: 12px;
    margin-top: 4px;
    font-size: 12px;
    color: var(--text-color3);

    .node-tag {
      padding: 2px 8px;
      background: var(--bg-color2);
      border-radius: 10px;
    }
  }
}

// 消息列表
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  scroll-behavior: smooth;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background-color: var(--border-color);
    border-radius: 3px;
  }
}

.chat-message-group {
  margin-bottom: 20px;
}

.time-divider {
  text-align: center;
  font-size: 12px;
  color: var(--text-color3);
  margin: 16px 0;
  position: relative;

  &::before,
  &::after {
    content: '';
    position: absolute;
    top: 50%;
    width: 30%;
    height: 1px;
    background: var(--border-color);
  }

  &::before { left: 0; }
  &::after { right: 0; }
}

// 消息项
.chat-message {
  display: flex;
  margin-bottom: 16px;
  animation: messageSlideIn 0.3s ease;

  &.is-me {
    flex-direction: row-reverse;

    .message-content-wrapper {
      align-items: flex-end;
    }

    .message-header {
      justify-content: flex-end;
    }

    .message-bubble {
      background: #95ec69;
      color: #000;
      border-radius: 16px 4px 16px 16px;

      &.is-topic {
        background: var(--bg-color2);
        color: var(--text-color);
        border-radius: 4px 16px 16px 16px;
      }

      .message-quote {
        background: rgba(0, 0, 0, 0.05);
      }
    }

    .message-actions {
      justify-content: flex-end;
    }
  }

  &.highlight {
    animation: highlightPulse 2s ease;
  }
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

@keyframes highlightPulse {
  0%, 100% { background: transparent; }
  50% { background: rgba(var(--color-primary-rgb), 0.1); }
}

.message-avatar {
  flex-shrink: 0;
  margin: 0 12px;
}

.message-content-wrapper {
  display: flex;
  flex-direction: column;
  max-width: 70%;
  min-width: 0;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 12px;

  .message-author {
    color: var(--text-color2);
    font-weight: 500;

    &:hover {
      color: var(--text-link-color);
    }
  }

  .message-location {
    color: var(--text-color3);
  }
}

.message-bubble {
  background: var(--bg-color2);
  padding: 12px 16px;
  border-radius: 4px 16px 16px 16px;
  color: var(--text-color);
  word-break: break-word;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

  &.is-topic {
    background: #e6f3ff;
    border: 1px solid #b3d9ff;
  }

  .message-quote {
    display: flex;
    margin-bottom: 8px;
    padding: 8px;
    background: rgba(0, 0, 0, 0.03);
    border-radius: 8px;
    cursor: pointer;

    .quote-line {
      width: 3px;
      background: var(--color-primary);
      border-radius: 2px;
      margin-right: 8px;
      flex-shrink: 0;
    }

    .quote-content {
      flex: 1;
      min-width: 0;
    }

    .quote-author {
      color: var(--color-primary);
      font-weight: 500;
    }

    .quote-text {
      color: var(--text-color3);
      font-size: 13px;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }

  .message-text {
    line-height: 1.6;
    white-space: pre-wrap;

    :deep(img) {
      max-width: 100%;
      border-radius: 8px;
    }
  }
}

.message-image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;

  :deep(.el-image) {
    width: 80px;
    height: 80px;
    border-radius: 8px;
    cursor: pointer;
    object-fit: cover;
  }
}

.message-actions {
  display: flex;
  gap: 12px;
  margin-top: 6px;
  padding: 0 4px;

  .action-btn {
    font-size: 12px;
    color: var(--text-color3);
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
    transition: color 0.2s;

    i {
      font-size: 14px;
    }

    &:hover,
    &.active {
      color: var(--color-primary);
    }

    &.active {
      font-weight: 500;
    }

    &.test-email-btn {
      background: linear-gradient(135deg, #ff6b6b, #ff8e53);
      color: #fff;
      padding: 4px 8px;
      border-radius: 12px;
      font-weight: 600;
      box-shadow: 0 2px 6px rgba(255,107,107,0.3);

      &:hover {
        background: linear-gradient(135deg, #ff5252, #ff7a00);
        color: #fff;
        box-shadow: 0 4px 10px rgba(255,82,82,0.4);
      }

      i {
        font-size: 12px;
      }
    }
  }
}

// 加载更多
.load-more-wrapper {
  text-align: center;
  padding: 20px;
}

// 新消息提示
.new-message-tip {
  position: absolute;
  bottom: 80px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--color-primary);
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 13px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  animation: bounceIn 0.3s ease;

  i {
    font-size: 12px;
  }
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

// 底部输入框
。chat-input-wrapper {
  background: var(--bg-color);
  border-top: 1px solid var(--border-color);
  flex-shrink: 0;
  position: relative;

  // 测试邮件按钮栏
  .test-email-bar {
    padding: 8px 16px;
    background: var(--bg-color2);
    border-bottom: 1px solid var(--border-color);

    .test-email-btn {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 6px 12px;
      background: linear-gradient(135deg, #ff6b6b, #ff8e53);
      color: #fff;
      border: none;
      border-radius: 16px;
      font-size: 13px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.3s ease;
      box-shadow: 0 2px 6px rgba(255,107,107,0.3);
      margin-right: 8px;

      &:hover {
        background: linear-gradient(135deg, #ff5252, #ff7a00);
        box-shadow: 0 4px 10px rgba(255,82,82,0.4);
      }

      i {
        font-size: 14px;
      }

      &.no-auth {
        background: linear-gradient(135deg, #6b8cff, #53c1ff);
        box-shadow: 0 2px 6px rgba(107,140,255,0.3);

        &:hover {
          background: linear-gradient(135deg, #5270ff, #00a7ff);
          box-shadow: 0 4px 10px rgba(82,112,255,0.4);
        }
      }
    }
  }

  .quote-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 16px;
    background: var(--bg-color2);
    font-size: 13px;
    color: var(--text-color2);

    .quote-text {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    i {
      cursor: pointer;
      color: var(--text-color3);

      &:hover {
        color: var(--text-link-color);
      }
    }
  }

  .chat-input-bar {
    display: flex;
    align-items: flex-end;
    gap: 12px;
    padding: 12px 16px;

    .input-actions {
      display: flex;
      gap: 12px;

      i {
        font-size: 24px;
        color: var(--text-color3);
        cursor: pointer;

        &:hover {
          color: var(--text-link-color);
        }
      }
    }

    .input-box {
      flex: 1;
      min-width: 0;

      textarea {
        width: 100%;
        min-height: 40px;
        max-height: 120px;
        padding: 10px 14px;
        border: 1px solid var(--border-color);
        border-radius: 20px;
        background: var(--bg-color2);
        color: var(--text-color);
        font-size: 14px;
        resize: none;
        outline: none;
        transition: all 0.2s;

        &:focus {
          border-color: var(--color-primary);
          background: var(--bg-color);
        }

        &::placeholder {
          color: var(--text-color4);
        }
      }
    }

    .send-btn {
      .el-button {
        padding: 10px 20px;
        border-radius: 20px;
      }
    }
  }

  .image-upload-panel {
    padding: 0 16px 12px;
  }
}

// 图片列表（话题）
.topic-image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;

  li {
    list-style: none;

    .image-item {
      :deep(.el-image) {
        border-radius: 8px;
        overflow: hidden;
      }
    }

    .file-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px 12px;
      background: var(--bg-color);
      border-radius: 8px;
      cursor: pointer;

      .file-icon {
        font-size: 24px;
      }

      .file-name {
        font-size: 13px;
        color: var(--text-color);
      }
    }
  }
}

// 隐藏内容
.hide-content {
  margin-top: 10px;

  .hide-content-tip {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px;
    background: var(--bg-color);
    border-radius: 8px;
    color: var(--color-primary);
    font-size: 13px;
  }
}

// 评论详情弹窗
.comment-detail-dialog {
  .comment-detail-content {
    .detail-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 16px;

      .detail-info {
        a {
          display: block;
          color: var(--text-color2);
          font-weight: 500;
        }

        .detail-time {
          font-size: 12px;
          color: var(--text-color3);
        }
      }
    }

    .detail-body {
      line-height: 1.8;
      margin-bottom: 20px;
    }

    .detail-actions {
      display: flex;
      gap: 12px;
    }
  }
}

// 响应式
@media screen and (max-width: 768px) {
  .chat-container {
    max-width: 100%;
  }

  .message-content-wrapper {
    max-width: 80%;
  }

  .chat-input-bar {
    padding: 8px 12px;

    .input-actions i {
      font-size: 20px;
    }
  }
}

// ========== 桌面端传统论坛风格 ==========
.topic-detail {
  margin-bottom: 20px;
  background-color: var(--bg-color);
  border-radius: var(--border-radius);

  .float-bar {
    position: fixed;
    margin-left: -58px;
    top: 300px;

    @media screen and (max-width: 1300px) {
      display: none;
    }
  }

  .topic-header,
  .topic-content,
  .topic-tags,
  .topic-like-users,
  .topic-actions {
    margin: 0 16px 16px 16px;
  }

  .topic-header {
    display: flex;

    .topic-header-left {
      margin: 10px 10px 0 0;
    }

    .topic-header-center {
      margin: 10px 10px 0 0;
      width: 100%;

      .topic-nickname a {
        color: var(--text-color2);
        font-size: 16px;
        font-weight: bold;
        overflow: hidden;
      }

      .topic-meta {
        position: relative;
        font-size: 12px;
        line-height: 24px;
        color: var(--text-color3);
        margin-top: 5px;

        span.meta-item {
          font-size: 12px;

          &:not(:last-child) {
            margin-right: 8px;
          }
        }
      }
    }

    .topic-header-right {
      margin-top: 10px;
      min-width: max-content;
    }
  }

  .topic-content {
    font-size: 15px;
    color: var(--text-color);
    white-space: normal;
    word-break: break-all;
    word-wrap: break-word;
    padding-top: 0 !important;

    .topic-title {
      font-weight: 700;
      font-size: 20px;
      word-wrap: break-word;
      word-break: normal;
      border-bottom: 1px solid var(--border-color4);
      padding-bottom: 10px;
    }

    .topic-content-detail {
      font-size: 16px;
      line-height: 24px;
      word-wrap: break-word;
      -webkit-font-smoothing: antialiased;
    }

    &.topic-tweet {
      .topic-content-detail {
        white-space: pre-line;
      }
    }

    .topic-image-list {
      margin-left: 0;
      margin-top: 10px;

      li {
        cursor: pointer;
        border: 1px dashed var(--border-color2);
        text-align: center;

        display: inline-block;
        vertical-align: middle;
        margin: 0 8px 8px 0;
        background-color: var(--bg-color2);
        background-size: 32px 32px;
        background-position: 50%;
        background-repeat: no-repeat;
        overflow: hidden;
        position: relative;

        .image-item {
          display: block;
          overflow: hidden;

          & > img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            transition: all 0.5s ease-out 0.1s;

            &:hover {
              transform: matrix(1.04, 0, 0, 1.04, 0, 0);
              backface-visibility: hidden;
            }
          }
        }

        .file-item {
          width: 120px;
          height: 120px;
          margin-right: 10px;
          margin-bottom: 10px;
          border-radius: 4px;
          overflow: hidden;
          cursor: pointer;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          background: var(--bg-color2);
          border: 1px solid var(--border-color);

          .file-icon {
            font-size: 32px;
            margin-bottom: 8px;
            font-family: "iconfont" !important;
            font-style: normal;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;

            &.icon-file-word {
              color: #2b579a;
            }

            &.icon-file-excel {
              color: #217346;
            }

            &.icon-file-ppt {
              color: #d24726;
            }

            &.icon-file-pdf {
              color: #f40f02;
            }

            &.icon-file-text {
              color: #25a9f6;
            }

            &.icon-file-zip {
              color: #ffb300;
            }

            &.icon-file {
              color: #666;
            }
          }

          .file-name {
            font-size: 12px;
            text-align: center;
            padding: 0 4px;
            word-break: break-all;
            color: var(--text-color);
          }
        }

        /* 只有一个图片时 */
        &:first-child:nth-last-child(1) {
          width: 210px;
          height: 210px;
          line-height: 210px;

          .image-item {
            width: 210px;
            height: 210px;
          }
        }

        /* 只有两个图片时 */
        &:first-child:nth-last-child(2),
        &:first-child:nth-last-child(2) ~ li {
          width: 180px;
          height: 180px;
          line-height: 180px;

          .image-item {
            width: 180px;
            height: 180px;
          }
        }

        /*大于两个图片时*/
        &:first-child:nth-last-child(n + 3),
        &:first-child:nth-last-child(n + 3) ~ li {
          width: 120px;
          height: 120px;
          line-height: 120px;

          .image-item {
            width: 120px;
            height: 120px;
          }
        }
      }
    }

    .hide-content {
      margin: 20px 0;

      .widget-header {
        span {
          font-weight: 500;
        }
      }

      .hide-content-tip {
        border: 1px solid var(--border-hover-color);
        border-radius: 2px;
        padding: 6px 12px;
        font-size: 14px;
        color: #3273dc;
      }
    }
  }

  .topic-tags {
    .topic-tag {
      padding: 2px 8px;
      justify-content: center;
      align-items: center;
      border-radius: 12.5px;
      margin-right: 10px;
      background: var(--bg-color2);
      border: 1px solid var(--border-color);
      color: var(--text-color3);
      font-size: 12px;

      &:hover {
        color: var(--text-link-color);
        background: var(--bg-color);
      }
    }
  }

  .topic-like-users {
    width: 80%;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    margin-bottom: 10px;

    .avatar-a {
      margin-right: -3px;
    }

    .like-count {
      margin-left: 8px;
      font-size: 14px;
      color: var(--text-color);
    }
  }

  .topic-actions {
    padding: 10px 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-top: 1px solid var(--border-color4);

    .action {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      background: var(--bg-color);
      cursor: pointer;
      color: var(--text-color3);
      font-size: 14px;

      .checked-icon {
        color: var(--color-red);
      }

      &.disabled {
        cursor: not-allowed;

        &:hover {
          color: var(--text-color3);

          > .action-icon {
            fill: var(--text-color3);
          }
        }
      }

      > .action-icon {
        fill: #8590a6;
      }

      .action-text {
        color: var(--text-color);
        margin-left: 5px;
      }

      &:hover {
        color: var(--text-link-color);

        > .action-icon {
          fill: var(--text-link-color);
        }
      }

      &.test-email-action {
        background: linear-gradient(135deg, #ff6b6b, #ff8e53);
        color: #fff;
        box-shadow: 0 2px 6px rgba(255,107,107,0.3);

        &:hover {
          background: linear-gradient(135deg, #ff5252, #ff7a00);
          box-shadow: 0 4px 10px rgba(255,82,82,0.4);
        }

        > .action-icon {
          fill: #fff;
        }

        .action-text {
          color: #fff;
        }
      }
    }
  }

  .need-a-hand-tag {
    display: inline-flex;
    align-items: center;
    background: #ff6b6b;
    color: white;
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: normal;
    margin-left: 8px;

    .iconfont {
      margin-right: 4px;
      font-size: 14px;
    }
  }
}

// 价值评价弹窗样式
.value-dialog {
  .value-types {
    .value-type-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 16px;
      margin-bottom: 8px;
      border: 1px solid var(--border-color);
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.3s ease;

      &:hover {
        background-color: var(--bg-color2);
        border-color: var(--color-primary);
      }

      .value-type-label {
        font-size: 14px;
        color: var(--text-color);
      }

      .value-type-score {
        font-size: 14px;
        color: var(--color-primary);
        font-weight: 500;
      }
    }
  }
}

// 子评论列表样式
.comment-replies {
  margin-top: 10px;
  margin-left: 50px;
  background: var(--bg-color2);
  border-radius: 8px;
  padding: 12px;
}

.reply-item {
  display: flex;
  margin-bottom: 12px;

  &:last-child {
    margin-bottom: 0;
  }
}

.reply-avatar {
  flex-shrink: 0;
  margin-right: 10px;
}

.reply-content-wrapper {
  flex: 1;
  min-width: 0;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 12px;

  .reply-author {
    color: var(--text-color2);
    font-weight: 500;

    &:hover {
      color: var(--text-link-color);
    }
  }

  .reply-location {
    color: var(--text-color3);
  }
}

.reply-bubble {
  background: var(--bg-color);
  padding: 10px 12px;
  border-radius: 8px;
  color: var(--text-color);
  word-break: break-word;

  .reply-quote {
    margin-bottom: 6px;
    padding: 6px 8px;
    background: rgba(0, 0, 0, 0.03);
    border-radius: 6px;
    border-left: 2px solid var(--color-primary);
    font-size: 12px;

    .quote-author {
      color: var(--color-primary);
      font-weight: 500;
    }

    .quote-text {
      color: var(--text-color3);
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }

  .reply-text {
    font-size: 14px;
    line-height: 1.5;

    :deep(img) {
      max-width: 100%;
      border-radius: 6px;
    }
  }
}

.reply-image-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
  margin-top: 8px;

  :deep(.el-image) {
    border-radius: 6px;
    overflow: hidden;
    aspect-ratio: 1;
  }
}

.reply-actions {
  display: flex;
  gap: 16px;
  margin-top: 8px;
  padding: 0 2px;

  .action-btn {
    font-size: 12px;
    color: var(--text-color3);
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;

    &:hover {
      color: var(--text-link-color);
    }

    &.active {
      color: var(--color-primary);
    }

    .iconfont {
      font-size: 14px;
    }
  }
}
</style>