<template>
  <div v-for="article in articles" :key="article.id" class="article-item">
    <div class="article-item-main">
      <div class="article-info">
        <nuxt-link
          class="article-title"
          :to="'/article/' + article.id"
          target="_blank"
          >{{ article.title }}</nuxt-link
        >

        <div class="article-summary">
          {{ article.summary }}
        </div>
        
        <!-- 文件附件列表 -->
        <div v-if="article.attachments && article.attachments.length > 0" class="article-attachments">
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
      </div>

      <div class="article-meta">
        <div class="article-meta-left">
          <span class="article-meta-item">
            <nuxt-link :to="'/user/' + article.user.id" class="article-author">
              <span>{{ article.user.nickname }}</span>
            </nuxt-link>
            <time :datetime="useFormatDate(article.createTime)"
              >{{ $t("component.articleList.publishedAt") }}
              {{ usePrettyDate(article.createTime) }}</time
            >
          </span>
        </div>

        <div class="article-meta-right">
          <div v-if="article.tags && article.tags.length > 0">
            <nuxt-link
              v-for="tag in article.tags"
              :key="tag.id"
              class="article-tag"
              :to="'/articles/tag/' + tag.id"
              >{{ tag.name }}</nuxt-link
            >
          </div>
        </div>
      </div>
    </div>
    <div v-if="article.cover" class="article-item-cover">
      <img :src="article.cover.url" />
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  articles: {
    type: Array,
    default() {
      return [];
    },
    required: false,
  },
});

// 获取文件图标类名
const getFileIconClass = (ext) => {
  const extMap = {
    '.doc': 'iconfont icon-file-word',
    '.docx': 'iconfont icon-file-word',
    '.xls': 'iconfont icon-file-excel',
    '.xlsx': 'iconfont icon-file-excel',
    '.ppt': 'iconfont icon-file-ppt',
    '.pptx': 'iconfont icon-file-ppt',
    '.pdf': 'iconfont icon-file-pdf',
    '.txt': 'iconfont icon-file-text',
    '.zip': 'iconfont icon-file-zip',
    '.rar': 'iconfont icon-file-zip',
  };
  
  // 如果没有对应的文件图标，使用通用的file图标
  return extMap[ext] || 'iconfont icon-file';
};

// 打开文件
const openFile = (url) => {
  window.open(url, '_blank');
};

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};
</script>

<style lang="scss" scoped>
.article-item {
  padding: 8px 8px;
  transition: background 0.5s;
  border-radius: 3px;
  // background: var(--bg-color);
  line-height: 20px;

  &:not(:last-child) {
    margin-bottom: 6px;
  }

  display: flex;
  align-items: center;

  .article-item-main {
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: space-between;

    .article-info {
      .article-title {
        font-size: 18px;
        line-height: 24px;
        font-weight: 500;
        color: var(--text-color);
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .article-summary {
        font-size: 14px;
        color: var(--text-color2);
        overflow: hidden;
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 3;
        line-clamp: 3;
        text-align: justify;
        padding-top: 4px;
        word-break: break-all;
        text-overflow: ellipsis;
      }
    }

    .article-meta {
      display: flex;
      // justify-content: space-between;
      align-items: center;
      font-size: 13px;
      padding-top: 4px;

      .article-meta-left {
        .article-meta-item {
          padding: 0 6px 0 0;
          color: var(--text-color3);

          .article-author {
            font-weight: bold;
            padding: 0 3px;
          }
        }
      }

      .article-meta-right {
        margin-left: 10px;

        @media screen and (max-width: 768px) {
          & {
            display: none;
          }
        }

        .article-tag {
          padding: 2px 8px;
          justify-content: center;
          align-items: center;
          border-radius: 12.5px;
          margin-right: 10px;
          background: var(--bg-color2);
          border: 1px solid var(--border-color2);
          color: var(--text-color3);
          font-size: 12px;

          &:hover {
            color: var(--text-link-color);
            background: var(--bg-color);
          }
        }
      }
    }
  }

  .article-item-cover {
    display: flex;
    margin-left: 6px;
    img {
      min-width: 140px;
      min-height: 90px;
      width: 140px;
      height: 90px;
      object-fit: contain;
      background-color: #f5f5f5;

      @media screen and (max-width: 768px) {
        & {
          min-width: 110px;
          min-height: 80px;
          width: 110px;
          height: 80px;
        }
      }
    }
  }
  
  .article-attachments {
    margin-top: 8px;
    
    .attachment-item {
      display: flex;
      align-items: center;
      padding: 6px 8px;
      margin-bottom: 4px;
      background: var(--bg-color2);
      border: 1px solid var(--border-color2);
      border-radius: 4px;
      cursor: pointer;
      transition: background 0.3s ease;
      
      &:hover {
        background: var(--bg-color3);
      }
      
      .attachment-icon {
        font-size: 16px;
        margin-right: 8px;
        min-width: 16px;
      }
      
      .attachment-name {
        flex: 1;
        font-size: 13px;
        color: var(--text-color);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      
      .attachment-size {
        font-size: 12px;
        color: var(--text-color3);
        margin-left: 8px;
      }
    }
  }
}

/* 文件图标颜色 */
.icon-file-word {
  color: #2b579a;
}

.icon-file-excel {
  color: #217346;
}

.icon-file-ppt {
  color: #d24726;
}

.icon-file-pdf {
  color: #f40f02;
}

.icon-file-text {
  color: #25a9f6;
}

.icon-file-zip {
  color: #ffb300;
}

.icon-file {
  color: #666;
}
</style>
