<template>
  <section class="main">
    <div v-if="isPending" class="container main-container">
      <div class="notification is-warning" style="width: 100%; margin: 20px 0">
        {{ t("pages.article.detail.pending") }}
      </div>
    </div>
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <article class="article-detail">
          <div class="article-header">
            <div class="article-title-wrapper">
              <h1 class="article-title">
                {{ article.title }}
              </h1>
              <div class="article-manage-menu">
                <article-manage-menu :article="article" />
              </div>
            </div>
            <div class="article-meta">
              <span class="article-meta-item">
                <nuxt-link
                  :to="'/user/' + article.user.id"
                  class="article-author"
                  >{{ article.user.nickname }}</nuxt-link
                >{{ t("pages.article.detail.publishedAt") }}
                <time :datetime="useFormatDate(article.createTime)">{{
                  usePrettyDate(article.createTime)
                }}</time>
              </span>
            </div>
          </div>

          <div
            class="article-content content line-numbers"
            v-html="processedContent"
          ></div>

          <!-- 附件列表 -->
          <div v-if="article.attachments && article.attachments.length > 0" class="article-attachments">
            <h3 class="attachments-title">{{ t('pages.article.detail.attachments') }}</h3>
            <div 
              v-for="(attachment, index) in article.attachments" 
              :key="index" 
              class="attachment-item"
              @click="openFile(attachment.url)"
            >
              <i :class="getFileIconClass(attachment.fileExt)" class="attachment-icon"></i>
              <span class="attachment-name">{{ attachment.filename }}</span>
              <span class="attachment-size">{{ formatFileSize(attachment.size) }}</span>
            </div>
          </div>

          <!--节点、标签-->
          <div class="article-tags">
            <nuxt-link
              v-for="tag in article.tags"
              :key="tag.id"
              :to="'/articles/tag/' + tag.id"
              class="article-tag"
              >#{{ tag.name }}</nuxt-link
            >
          </div>
        </article>

        <!-- 评论 -->
        <comment
          :entity-id="article.id"
          :comment-count="article.commentCount"
          entity-type="article"
        />
      </div>
      <div class="right-container">
        <user-info :user="article.user" />
      </div>
    </div>
  </section>
</template>

<script setup>
const { t } = useI18n();
const route = useRoute();
const { data: article, error } = await useMyFetch(
  `/api/article/${route.params.id}`
);

// 引入文件图标样式
import '~/assets/css/file-icons.css';

if (error.value) {
  throw createError({
    statusCode: 500,
    message: error.value.message || t("pages.redirect.error"),
  });
}

useHead({
  title: useSiteTitle(article.value.title),
});

const isPending = computed(() => {
  return article.value.status === 2;
});

// 获取文件图标类名（参照Topic系统）
const getFileIconClass = (fileExt) => {
  if (!fileExt) return 'file-icon file-icon-default';
  
  // 确保扩展名以点开头
  const ext = fileExt.startsWith('.') ? fileExt.toLowerCase() : '.' + fileExt.toLowerCase();
  
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
  
  // 如果没有对应的文件图标，使用通用的file图标
  return extMap[ext] || 'file-icon file-icon-default';
};

// 处理内容中的文件链接，添加图标（参照Topic系统）
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

// 处理后的文章内容
const processedContent = computed(() => {
  return article.value ? processFileLinks(article.value.content) : '';
});

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

// 打开文件
const openFile = (url) => {
  window.open(url, '_blank');
};
</script>

<style lang="scss" scoped>
.article-detail {
  margin-bottom: 12px;
  padding: 12px;
  border-radius: var(--border-radius);
  background-color: var(--bg-color);
  overflow: hidden;

  .article-title {
    a {
      font-size: 18px;
      line-height: 30px;
      font-weight: 500;
      color: var(--text-color);
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  .article-header {
    padding: 10px 0;
    border-bottom: 1px solid var(--border-color);

    .article-title-wrapper {
      display: flex;
      .article-title {
        width: 100%;
        color: var(--text-color);
        font-weight: normal;
        overflow: hidden;
        text-overflow: ellipsis;
        font-size: 18px;
        line-height: 30px;
      }
      .article-manage-menu {
        min-width: max-content;
      }
    }

    .article-meta {
      display: inline-block;
      font-size: 13px;
      padding-top: 6px;

      .article-meta-item {
        padding: 0 6px 0 0;
        color: var(--text-color3);

        a {
          color: var(--text-link-color);

          &.article-author {
            font-weight: bold;
            padding: 0 3px;
          }
        }
      }
    }
  }

  .article-content {
    font-size: 15px;
    margin-top: 10px;
    margin-bottom: 10px;

    a.article-share-summary {
      color: var(--text-color);
    }
  }

  .article-tags {
    margin-top: 10px;
    .article-tag {
      height: 25px;
      padding: 0 8px;
      display: inline-flex;
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
        border: 1px solid var(--border-hover-color);
      }
    }
  }

  .article-attachments {
    margin-top: 20px;
    padding: 16px;
    background: var(--bg-color2);
    border-radius: 8px;
    border: 1px solid var(--border-color);

    .attachments-title {
      margin: 0 0 12px 0;
      font-size: 16px;
      font-weight: 600;
      color: var(--text-color);
    }

    .attachment-item {
      display: flex;
      align-items: center;
      padding: 8px 12px;
      margin-bottom: 8px;
      background: var(--bg-color);
      border: 1px solid var(--border-color2);
      border-radius: 6px;
      cursor: pointer;
      transition: background 0.3s ease;

      &:hover {
        background: var(--bg-color3);
      }

      .attachment-icon {
        width: 20px;
        height: 20px;
        margin-right: 12px;
        flex-shrink: 0;
      }

      .attachment-name {
        flex: 1;
        font-size: 14px;
        color: var(--text-color);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .attachment-size {
        font-size: 12px;
        color: var(--text-color3);
        margin-left: 8px;
        flex-shrink: 0;
      }
    }
  }
}</style>
