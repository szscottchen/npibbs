<template>
  <ul class="topic-list">
    <li v-for="topic in topics" :key="topic.id" class="topic-item">
      <div class="topic-header">
        <div class="topic-header-main">
          <my-avatar :user="topic.user" />
          <div class="topic-userinfo">
            <a :href="`/user/${topic.user.id}`" class="topic-nickname">
              {{ topic.user.nickname }}
            </a>
            <div class="topic-time">
              {{ $t("component.topicList.publishedAt") }}&nbsp;{{
                usePrettyDate(topic.createTime)
              }}
            </div>
          </div>
        </div>
        <div class="topic-title-container">
          <h1 class="topic-title">
            <nuxt-link :to="`/topic/${topic.id}`" target="_blank">
              {{ topic.title }}
            </nuxt-link>
            <span v-if="topic.needAHand" class="need-a-hand-tag">
              <i class="iconfont icon-help" />
              {{ t('component.topicList.needAHand') }}
            </span>
          </h1>
        </div>
        <div class="topic-header-right">
          <span v-if="showSticky && topic.sticky" class="topic-sticky-icon">{{
            $t("component.topicList.sticky")
          }}</span>
        </div>
      </div>
      <div class="topic-content" :class="{ 'topic-tweet': topic.type === 1 }">
        <template v-if="topic.type === 0">
          <nuxt-link
            :to="`/topic/${topic.id}`"
            class="topic-summary"
            target="_blank"
          >
            {{ topic.summary }}
          </nuxt-link>
        </template>
        <template v-if="topic.type === 1">
          <nuxt-link
            v-if="topic.content"
            :to="`/topic/${topic.id}`"
            class="topic-summary"
            target="_blank"
          >
            {{ topic.content }}
          </nuxt-link>
          <ul
            v-if="topic.imageList && topic.imageList.length"
            class="topic-image-list"
          >
            <li v-for="(image, index) in topic.imageList" :key="index">
              <nuxt-link
                v-if="image.isImage"
                :to="`/topic/${topic.id}`"
                class="image-item"
                target="_blank"
              >
                <img :src="image.preview" />
              </nuxt-link>
              <div 
                v-else 
                class="file-item"
                @click="openFile(image.url)"
              >
                <i :class="getFileIconClass(image.fileExt)" class="file-icon"></i>
                <span class="file-ext">{{ image.fileExt.replace('.', '') }}</span>
              </div>
            </li>
          </ul>
        </template>
      </div>
      <div class="topic-bottom">
        <div class="topic-tags">
          <nuxt-link
            v-if="topic.node"
            class="topic-tag"
            target="_blank"
            :to="`/topics/node/${topic.node.id}`"
            :alt="topic.node.name"
          >
            <img v-if="topic.node.logo" :src="topic.node.logo" />
            <span>{{ topic.node.name }}</span>
          </nuxt-link>
        </div>

        <div class="topic-actions">
          <div
            class="btn EASE"
            :class="{ liked: topic.liked }"
            @click="like(topic)"
          >
            <i class="iconfont icon-like" />
            <span v-if="topic.likeCount > 0">{{ topic.likeCount }}</span>
          </div>
          <div class="btn EASE" @click="toTopicDetail(topic.id)">
            <i class="iconfont icon-comment" />
            <span v-if="topic.commentCount > 0">{{ topic.commentCount }}</span>
          </div>
          <!-- <div class="btn EASE" @click="toTopicDetail(topic.id)">
            <i class="iconfont icon-view" />
            <span v-if="topic.viewCount > 0">{{ topic.viewCount }}</span>
          </div> -->
        </div>
      </div>
    </li>
  </ul>
</template>
<script setup>
const { t } = useI18n();

const props = defineProps({
  topics: {
    type: Array,
    default() {
      return [];
    },
    required: false,
  },
  showAvatar: {
    type: Boolean,
    default: true,
  },
  showSticky: {
    type: Boolean,
    default: false,
  },
});

const like = async (topic) => {
  try {
    if (topic.liked) {
      await useHttpPost(
        "/api/like/unlike",
        useJsonToForm({
          entityType: "topic",
          entityId: topic.id,
        })
      );
      topic.liked = false;
      topic.likeCount = topic.likeCount > 0 ? topic.likeCount - 1 : 0;
      useMsgSuccess(t("component.topicList.unlikeSuccess"));
    } else {
      await useHttpPost(
        "/api/like/like",
        useJsonToForm({
          entityType: "topic",
          entityId: topic.id,
        })
      );
      topic.liked = true;
      topic.likeCount++;
      useMsgSuccess(t("component.topicList.likeSuccess"));
    }
  } catch (e) {
    useCatchError(e);
  }
};

const toTopicDetail = (topicId) => {
  useLinkTo(`/topic/${topicId}`);
};

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


</script>
<style lang="scss" scoped>
.topic-list {
  .topic-item {
    padding: 10px 20px;
    position: relative;
    overflow: hidden;
    border-radius: 3px;

    &:not(:last-child):after {
      position: absolute;
      content: "";
      bottom: 0;
      left: 32px;
      right: 32px;
      height: 1px;
      background: var(--border-color2);  /* 加深颜色 */
    }

    .topic-header {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .topic-header-main {
        display: flex;
        align-items: center;

        .topic-userinfo {
          margin-left: 10px;
          .topic-nickname {
            font-weight: 500;
            font-size: 14px;
            color: var(--text-color);
          }

          .topic-time {
            margin-top: 3px;
            font-size: 10px;  /* 减小字体大小 */
            color: var(--text-color3);
          }
        }
      }

      .topic-title-container {
        flex: 1;
        text-align: left;
        margin: 0 10px;
        align-self: flex-start;
        
        .topic-title {
          margin: 0;
          font-size: 16px;
          font-weight: 600;
          
          a {
            color: #3273dc;  /* 蓝色 */
            
            &:hover {
              text-decoration: underline;
            }
          }
        }
      }

      .topic-header-right {
        .topic-sticky-icon {
          font-size: 13px;
          line-height: 13px;
          color: #ff7827;
          background: #ffe7d9;
          border-radius: 2px;
          padding: 3px 6px;
          white-space: nowrap;
        }
      }
    }

    .topic-content {
      margin-top: 4px;
      .topic-summary {
        font-size: 14px;
        margin-bottom: 4px;
        width: 100%;
        text-decoration: none;
        color: var(--text-color3);
        word-wrap: break-word;

        overflow: hidden;
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 3;
        line-clamp: 3;
        text-align: justify;
        word-break: break-all;
        text-overflow: ellipsis;
      }

      &.topic-tweet {
        .topic-summary {
          color: var(--text-color);
          white-space: pre-line;
        }
      }

      .topic-image-list {
      display: grid;
      gap: 6px;
      margin-top: 10px;
      width: 100%;

      .image-item {
        display: block;
        border-radius: 4px;
        overflow: hidden;
        border: 1px solid var(--border-color2);

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.3s ease;
          aspect-ratio: 1/1;
          object-fit: cover;

          &:hover {
            transform: scale(1.05);
          }
        }
      }

      .file-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border: 1px solid var(--border-color2);
        border-radius: 4px;
        background: #f8f9fa;
        cursor: pointer;
        aspect-ratio: 1/1;
        
        .file-icon {
          font-size: 24px;
          margin-bottom: 4px;
          color: #666;
        }
        
        .file-ext {
          font-size: 12px;
          color: #666;
        }
        
        &:hover {
          background: #e9ecef;
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
    }

    .topic-bottom {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .topic-tags {
        display: flex;

        .topic-tag {
          display: flex;
          justify-content: center;
          align-items: center;
          padding: 4px 10px;
          border-radius: 18px;
          background: var(--bg-color6);
          color: var(--text-color3);
          font-size: 13px;

          &:hover {
            color: var(--text-color3-hover);
            background: var(--bg-color6-hover);
          }

          img {
            display: block;
            width: 20px;
            height: 20px;
            margin: 0 4px 0 0;
            border-radius: 50%;
            object-fit: cover;
          }
        }
      }

      .topic-actions {
      display: flex;
      align-items: center;
      margin-top: 4px;
      font-size: 12px;
      user-select: none;

      .btn {
        color: var(--text-color3);
        cursor: pointer;
        display: flex;
        align-items: center;

        &:not(:last-child) {
          margin-right: 12px;
        }

        &:hover {
          color: var(--text-link-color);
        }

        i {
          margin-right: 3px;
          font-size: 16px;
          position: relative;
        }

        span {
          line-height: 24px;
          font-size: 15px;
        }

        &.liked {
          color: var(--color-red) !important;
        }
      }
    }
    }
  }
}

@media screen and (max-width: 768px) {
  .topic-list {
    .topic-item {
      padding: 8px 16px;

      &:after {
        left: 16px !important;
        right: 16px !important;
      }
      
      .topic-header {
        .topic-title-container {
          .topic-title {
            font-size: 14px;
          }
        }
      }
      
      .topic-content {
        margin-top: 3px;
        
        .topic-summary {
          margin-bottom: 3px;
        }
      }
      
      .topic-actions {
        margin-top: 3px;
        
        .topic-action-item {
          margin-right: 10px;
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
