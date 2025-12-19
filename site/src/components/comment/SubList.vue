<template>
  <div class="replies">
    <div v-for="comment in replies.results" :key="comment.id" class="comment">
      <div class="comment-item-left">
        <my-avatar :user="comment.user" :size="24" has-border />
      </div>
      <div class="comment-item-main">
        <div class="comment-meta">
          <div>
            <nuxt-link
              :to="`/user/${comment.user.id}`"
              class="comment-nickname"
            >
              {{ comment.user.nickname }}
            </nuxt-link>
            <template v-if="comment.quote">
              &nbsp;<span>{{ $t('component.comment.subList.replyTo') }}</span>&nbsp;
              <nuxt-link
                :to="`/user/${comment.quote.user.id}`"
                class="comment-nickname"
              >
                {{ comment.quote.user.nickname }}
              </nuxt-link>
            </template>
          </div>
          <time class="comment-time">{{
            usePrettyDate(comment.createTime)
          }}</time>
        </div>
        <div class="comment-content-wrapper">
          <div v-if="comment.content" class="comment-content content">
            <div v-text="comment.content" />
          </div>
          <div
            v-if="comment.imageList && comment.imageList.length"
            class="comment-image-list"
          >
            <img
              v-for="(image, imageIndex) in comment.imageList"
              :key="imageIndex"
              :src="image.url"
            />
          </div>

          <div v-if="comment.quote" class="comment-quote">
            <div
              class="comment-quote-content content"
              v-html="comment.quote.content"
            />
            <div
              v-if="comment.quote.imageList && comment.quote.imageList.length"
              class="comment-quote-image-list"
            >
              <img
                v-for="(image, imageIndex) in comment.quote.imageList"
                :key="imageIndex"
                :src="image.preview"
              />
            </div>
          </div>
        </div>
        <div class="comment-actions">
          <div
            class="comment-action-item"
            :class="{ active: comment.liked }"
            @click="like(comment)"
          >
            <i class="iconfont icon-like" />
            <span>{{ comment.liked ? $t('component.comment.subList.liked') : $t('component.comment.subList.like') }}</span>
            <span v-if="comment.likeCount > 0">{{ comment.likeCount }}</span>
          </div>
          <div
            class="comment-action-item"
            :class="{ active: reply.quoteId === comment.id }"
            @click="switchShowReply(comment)"
          >
            <i class="iconfont icon-comment" />
            <span>{{
              reply.quoteId === comment.id ? $t('component.comment.subList.cancelReply') : $t('component.comment.subList.reply')
            }}</span>
          </div>
          <div
            v-if="user && user.id === localTopicUserId && valueTypes.length > 0"
            class="comment-action-item"
            @click="showValue(comment)"
          >
            <i class="iconfont icon-star" />
            <span>价值</span>
          </div>
        </div>
        <div v-if="reply.quoteId === comment.id" class="comment-reply-form">
          <text-editor
            :ref="`editor${comment.id}`"
            v-model:content="reply.value.content"
            v-model:imageList="reply.value.imageList"
            :height="80"
            @submit="submitReply(comment)"
          />
        </div>
      </div>
    </div>
    <div v-if="replies.hasMore === true" class="comment-more">
      <a @click="loadMore">
        <span>{{ $t('component.comment.subList.loadMore') }}</span>
        <i class="iconfont icon-right" />
      </a>
    </div>
  </div>
  
  <!-- 价值评价对话框 -->
  <el-dialog
    v-model="showValueDialog"
    title="价值评价"
    width="400px"
    :close-on-click-modal="false"
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
</template>

<script>
import { ref } from 'vue';

export default {
  props: {
    commentId: {
      type: Number,
      required: true,
    },
    data: {
      type: Object,
      required: true,
    },
    topicUserId: {
      type: Number,
      default: 0,
    },
    entityType: {
      type: String,
      default: "",
    },
    entityId: {
      type: Number,
      default: 0,
    },
  },
  emits: ['reply'],
  data() {
    return {
      replies: this.data,
      reply: {
        quoteId: 0,
        value: {
          content: "",
          imageList: [],
        },
      },
      showValueDialog: false,
      currentComment: null,
      valueTypes: [],
      localTopicUserId: this.topicUserId,
    };
  },
  computed: {
    user() {
      const userStore = useUserStore();
      return userStore.user;
    },
  },
  mounted() {
    this.loadValueTypes();
    this.loadTopicUser();
  },
  methods: {
    async loadTopicUser() {
      // 如果props中提供了topicUserId，直接使用
      if (this.localTopicUserId) {
        return;
      }
      // 如果没有提供，尝试从实体加载
      if (this.entityType === 'topic' && this.entityId) {
        try {
          const topic = await useHttpGet(`/api/topic/${this.entityId}`);
          if (topic && topic.user) {
            this.localTopicUserId = topic.user.id;
          }
        } catch (e) {
          console.error("Failed to load topic user:", e);
        }
      }
    },
    async loadMore() {
      const ret = await useHttpGet("/api/comment/replies", {
        params: {
          commentId: this.commentId,
          cursor: this.replies.cursor,
        },
      });
      this.replies.cursor = ret.cursor;
      this.replies.hasMore = ret.hasMore;
      this.replies.results.push(...ret.results);
    },
    async like(comment) {
      try {
        if (comment.liked) {
          await useHttpPost(
            "/api/like/unlike",
            useJsonToForm({
              entityType: "comment",
              entityId: comment.id,
            })
          );
          comment.liked = false;
          comment.likeCount = comment.likeCount > 0 ? comment.likeCount - 1 : 0;
        } else {
          await useHttpPost(
            "/api/like/like",
            useJsonToForm({
              entityType: "comment",
              entityId: comment.id,
            })
          );
          comment.liked = true;
          comment.likeCount = comment.likeCount + 1;
        }
      } catch (e) {
        useCatchError(e);
      }
    },
    switchShowReply(comment) {
      if (!this.user) {
        useMsgSignIn();
        return;
      }

      if (this.reply.quoteId === comment.id) {
        this.hideReply(comment);
      } else {
        this.reply.quoteId = comment.id;
        setTimeout(() => {
          const refs = this.$refs[`editor${comment.id}`];
          if (refs && refs.length > 0) {
            refs[0].focus();
          }
        }, 100);
      }
    },
    hideReply(comment) {
      this.reply.quoteId = 0;
      this.reply.value.content = "";
      this.reply.value.imageList = [];
    },
    async submitReply(parent) {
      try {
        const ret = await useHttpPost(
          "/api/comment/create",
          useJsonToForm({
            entityType: "comment",
            entityId: this.commentId,
            quoteId: this.reply.quoteId,
            content: this.reply.value.content,
            imageList:
              this.reply.value.imageList && this.reply.value.imageList.length
                ? JSON.stringify(this.reply.value.imageList)
                : "",
          })
        );
        this.hideReply(parent);
        this.$emit("reply", ret);
      } catch (e) {
        useCatchError(e);
      }
    },
    async loadValueTypes() {
      // 加载价值类型配置
      try {
        const config = await useHttpGet("/api/sys-config/configs");
        if (config && config.valueTypes && config.valueTypes.length > 0) {
          this.valueTypes = config.valueTypes;
        }
      } catch (e) {
        console.error("Failed to load value types:", e);
      }
    },
    showValue(comment) {
      // 显示价值评价对话框
      if (!this.user) {
        useMsgSignIn();
        return;
      }
      
      // 只有主贴发布人可以评价
      if (this.user.id !== this.localTopicUserId) {
        useMsgError("只有主贴发布人可以评价价值");
        return;
      }
      
      // 不能评价自己的评论
      if (this.user.id === comment.user.id) {
        useMsgError("不能评价自己的评论");
        return;
      }
      
      this.currentComment = comment;
      this.showValueDialog = true;
    },
    async submitValue(valueType, index) {
      // 提交价值评价
      try {
        await useHttpPost("/api/comment/value", useJsonToForm({
          commentId: this.currentComment.id,
          valueType: index
        }));
        
        useMsgSuccess(`评价成功：${valueType.label} (+${valueType.score}积分)`);
        this.showValueDialog = false;
        
        // 更新评论的价值显示
        if (!this.currentComment.values) {
          this.currentComment.values = [];
        }
        this.currentComment.values.push({
          type: valueType.type,
          score: valueType.score,
          userId: this.user.id
        });
      } catch (e) {
        useMsgError("价值评价失败:" + e.message);
      }
    },
  },
};
</script>

<style scoped lang="scss">
.replies {
  background-color: var(--bg-color2);
  padding: 10px;
  border-radius: 3px;

  .comment {
    display: flex;
    padding: 10px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color4);
    }

    .comment-item-left {
      margin-right: 10px;
    }

    .comment-item-main {
      flex: 1 1 auto;

      .comment-meta {
        display: flex;
        justify-content: space-between;
        .comment-nickname {
          font-size: 14px;
          color: var(--text-color2);

          &:hover {
            color: var(--text-link-color);
          }
        }

        .comment-time {
          font-size: 12px;
          color: var(--text-color3);
        }
      }

      .comment-content-wrapper {
        .comment-content {
          margin-top: 10px;
          margin-bottom: 0;
          color: var(--text-color);
          white-space: pre-wrap;
          font-size: 14px;
        }

        .comment-quote {
          margin-top: 10px;
          padding: 10px;
          background-color: var(--bg-color);
          border-left: 3px solid var(--border-color);

          .comment-quote-content {
            font-size: 14px;
            color: var(--text-color3);
          }
        }

        .comment-image-list {
          margin-top: 10px;

          img {
            width: 72px;
            height: 72px;
            line-height: 72px;
            cursor: pointer;
            &:not(:last-child) {
              margin-right: 8px;
            }

            object-fit: cover;
            transition: all 0.5s ease-out 0.1s;

            &:hover {
              transform: matrix(1.04, 0, 0, 1.04, 0, 0);
              backface-visibility: hidden;
            }
          }
        }
      }

      .comment-actions {
        margin-top: 10px;
        display: flex;
        align-items: center;
        column-gap: 10px;

        .comment-action-item {
          font-size: 12px;
          cursor: pointer;
          color: var(--text-color3);
          user-select: none;
          display: flex;
          align-items: center;
          column-gap: 2px;

          &:hover {
            color: var(--text-link-color);
          }

          &.active {
            color: var(--text-link-color);
            font-weight: 500;
          }

          i {
            font-size: 12px;
          }
        }
      }

      .comment-reply-form {
        margin-top: 10px;
      }
    }
  }

  .comment-more {
    text-align: center;
    margin-top: 10px;
    font-size: 14px;

    a {
      color: var(--text-color3);
      cursor: pointer;

      &:hover {
        color: var(--text-link-color);
      }

      i {
        font-size: 12px;
      }
    }
  }
}

// 价值评价对话框样式
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
</style>