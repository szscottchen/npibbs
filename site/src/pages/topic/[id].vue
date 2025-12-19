<template>
  <section v-if="topic" class="main">
    <div v-if="isPending" class="container main-container">
      <div class="notification is-warning" style="width: 100%; margin: 20px 0">
        {{ t('pages.topic.detail.pending') }}
      </div>
    </div>
    <div class="container main-container left-main size-360">
      <div class="left-container">
        <div class="main-content no-padding no-bg">
          <article class="topic-detail">
            <side-action-bar
              class="float-bar"
              entity-type="topic"
              :entity-id="topic.id"
              :liked="topic.liked"
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
              <div class="action" @click="like(topic)">
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
        <user-info :user="topic.user" />
      </div>
    </div>
  </section>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

// 引入文件图标样式
import '~/assets/css/file-icons.css';

const route = useRoute();
const { t } = useI18n();

const { data: topic } = await useMyFetch(`/api/topic/${route.params.id}`);

const { data: liked } = await useMyFetch("/api/like/liked", {
  params: {
    entityType: "topic",
    entityId: route.params.id,
  },
});

const { data: likeUsers, refresh: refreshLikeUsers } = await useMyFetch(
  `/api/topic/recentlikes/${route.params.id}`
);

const { data: hideContent, refresh: refreshHideContent } = await useMyFetch(
  `/api/topic/hide_content?topicId=${route.params.id}`
);

const imageUrls = computed(() => {
  if (!topic.value.imageList || !topic.value.imageList.length) {
    return [];
  }
  const ret = [];
  for (let i = 0; i < topic.value.imageList.length; i++) {
    ret.push(topic.value.imageList[i].url);
  }
  return ret;
});

// 判断是否为图片文件
const isImageFile = (file) => {
  // 如果传入的是对象且包含isImage字段，优先使用该字段
  if (typeof file === 'object' && file !== null && file.isImage !== undefined) {
    return file.isImage;
  }
  
  // 否则按原来的逻辑检查URL扩展名
  const url = typeof file === 'object' && file !== null ? file.url : file;
  if (!url) return false;
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp'];
  const lowerUrl = url.toLowerCase();
  return imageExtensions.some(ext => lowerUrl.endsWith(ext));
};

// 获取文件图标类名
const getFileIconClass = (fileObj) => {
  if (!fileObj) return 'file-icon file-icon-default';
  
  // 如果传入的是对象且包含fileExt字段，优先使用该字段
  if (typeof fileObj === 'object' && fileObj !== null && fileObj.fileExt) {
    const ext = fileObj.fileExt.toLowerCase();
    return getIconClassByExtension(ext);
  }
  
  // 如果传入的是对象但没有fileExt字段，从URL中提取扩展名
  if (typeof fileObj === 'object' && fileObj !== null && fileObj.url) {
    const urlStr = fileObj.url.toLowerCase();
    const ext = '.' + urlStr.split('.').pop();
    return getIconClassByExtension(ext);
  }
  
  // 否则按原来的逻辑检查URL扩展名
  const lowerUrl = fileObj.toLowerCase();
  const ext = '.' + lowerUrl.split('.').pop();
  return getIconClassByExtension(ext);
};

// 根据文件扩展名获取图标类名
const getIconClassByExtension = (ext) => {
  const extMap = {
    // Word文档
    '.doc': 'file-icon file-icon-word',
    '.docx': 'file-icon file-icon-word',
    
    // Excel表格
    '.xls': 'file-icon file-icon-excel',
    '.xlsx': 'file-icon file-icon-excel',
    '.csv': 'file-icon file-icon-excel',
    
    // PowerPoint演示文稿
    '.ppt': 'file-icon file-icon-ppt',
    '.pptx': 'file-icon file-icon-ppt',
    
    // PDF文档
    '.pdf': 'file-icon file-icon-pdf',
    
    // 文本文件
    '.txt': 'file-icon file-icon-text',
    '.log': 'file-icon file-icon-text',
    '.md': 'file-icon file-icon-text',
    
    // 压缩文件
    '.zip': 'file-icon file-icon-zip',
    '.rar': 'file-icon file-icon-zip',
    '.7z': 'file-icon file-icon-zip',
    '.tar': 'file-icon file-icon-zip',
    '.gz': 'file-icon file-icon-zip',
    
    // 图片文件
    '.jpg': 'file-icon file-icon-image',
    '.jpeg': 'file-icon file-icon-image',
    '.png': 'file-icon file-icon-image',
    '.gif': 'file-icon file-icon-image',
    '.bmp': 'file-icon file-icon-image',
    '.webp': 'file-icon file-icon-image',
    '.svg': 'file-icon file-icon-image',
    
    // 音频文件
    '.mp3': 'file-icon file-icon-audio',
    '.wav': 'file-icon file-icon-audio',
    '.flac': 'file-icon file-icon-audio',
    '.aac': 'file-icon file-icon-audio',
    '.ogg': 'file-icon file-icon-audio',
    
    // 视频文件
    '.mp4': 'file-icon file-icon-video',
    '.avi': 'file-icon file-icon-video',
    '.mkv': 'file-icon file-icon-video',
    '.mov': 'file-icon file-icon-video',
    '.wmv': 'file-icon file-icon-video',
    '.flv': 'file-icon file-icon-video',
    
    // 代码文件
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
};

// 处理内容中的文件链接，添加图标
const processFileLinks = (content) => {
  if (!content) return '';
  
  // 匹配所有uploads/files路径的文件链接
  const fileLinkRegex = /<a\s+[^>]*href="([^"]*uploads\/files[^"]*\.([a-zA-Z0-9]+))"[^>]*>(.*?)<\/a>/gi;
  
  return content.replace(fileLinkRegex, (match, href, ext, linkText) => {
    // 获取文件扩展名
    const fileExt = '.' + ext.toLowerCase();
    const iconClass = getFileIconClass(fileExt);
    
    // 构建带图标的链接
    return `<span class="file-link-wrapper">
      <i class="${iconClass}"></i>
      <a href="${href}" rel="nofollow" title="${linkText}">${linkText}</a>
    </span>`;
  });
};

// 处理后的帖子内容
const processedContent = computed(() => {
  return topic.value ? processFileLinks(topic.value.content) : '';
});

// 获取文件名
const getFileName = (url) => {
  if (!url) return '';
  const parts = url.split('/');
  return parts[parts.length - 1];
};

// 双击打开文件
const openFile = (url) => {
  window.open(url, '_blank');
};

useHead({
  title: useTopicSiteTitle(topic.value),
});

const isPending = computed(() => {
  return topic.value?.status === 2;
});

async function like() {
  try {
    if (liked.value) {
      await useHttpPost(
        "/api/like/unlike",
        useJsonToForm({
          entityType: "topic",
          entityId: topic.value.id,
        })
      );
      liked.value = false;
      topic.value.likeCount =
        topic.value.likeCount > 0 ? topic.value.likeCount - 1 : 0;

      useMsgSuccess(t('pages.topic.detail.likeSuccess'));
      await refreshLikeUsers();
    } else {
      await useHttpPost(
        "/api/like/like",
        useJsonToForm({
          entityType: "topic",
          entityId: topic.value.id,
        })
      );
      liked.value = true;
      topic.value.likeCount++;

      useMsgSuccess(t('pages.topic.detail.likeSuccess'));
      await refreshLikeUsers();
    }
  } catch (e) {
    useCatchError(e);
  }
}

async function addFavorite(topicId) {
  try {
    if (topic.value.favorited) {
      await useHttpPost(
        "/api/favorite/delete",
        useJsonToForm({
          entityType: "topic",
          entityId: topicId,
        })
      );
      topic.value.favorited = false;
      useMsgSuccess(t('pages.topic.detail.favoriteSuccess'));
    } else {
      await useHttpPost(
        "/api/favorite/add",
        useJsonToForm({
          entityType: "topic",
          entityId: topicId,
        })
      );
      topic.value.favorited = true;
      useMsgSuccess(t('pages.topic.detail.favoriteSuccess'));
    }
  } catch (e) {
    useCatchError(e);
  }
}

async function commentCreated() {
  refreshHideContent();
}
</script>

<style lang="scss" scoped>
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
          // transform-style: preserve-3d;

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
</style>
