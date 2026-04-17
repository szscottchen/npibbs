<template>
  <div class="comments">
    <load-more-async
      ref="loadMore"
      v-slot="{ results }"
      :params="{ entityType, entityId }"
      url="/api/comment/comments"
    >
      <div v-for="comment in processComments(results)" :key="comment.id" class="comment">
        <div class="comment-item-left">
          <my-avatar :user="comment.user" :size="30" has-border />
        </div>
        <div class="comment-item-main">
          <div class="comment-meta">
            <nuxt-link
              :to="`/user/${comment.user.id}`"
              class="comment-nickname"
            >
              {{ comment.user.nickname }}
            </nuxt-link>
            <div class="comment-meta-right">
              <time class="comment-time">{{
                usePrettyDate(comment.createTime)
              }}</time>
              <span v-if="comment.ipLocation" class="comment-ip-area"
                >{{ t('component.comment.list.ipLocation') }}{{ comment.ipLocation }}</span
              >
            </div>
          </div>
          <div class="comment-content-wrapper">
            <template v-if="comment.content">
              <div
                v-if="comment.contentType === 'text'"
                class="comment-content content"
                v-text="comment.content"
              />
              <div
                v-else
                class="comment-content content"
                v-html="comment.content"
              />
            </template>
            <div
              v-if="comment.imageList && comment.imageList.length"
              class="comment-image-list"
            >
              <img
                v-for="(image, imageIndex) in comment.imageList"
                :key="imageIndex"
                :src="image.preview"
              />
            </div>
          </div>
          <div class="comment-actions">
            <div
              class="comment-action-item"
              :class="{ active: comment.liked }"
              @click="like(comment)"
            >
              <i class="iconfont icon-like" />
              <span>{{ comment.liked ? t('component.comment.list.liked') : t('component.comment.list.like') }}</span>
              <span v-if="comment.likeCount > 0">{{ comment.likeCount }}</span>
            </div>
            <div
              class="comment-action-item"
              :class="{ active: reply.commentId === comment.id }"
              @click="switchShowReply(comment)"
            >
              <i class="iconfont icon-comment" />
              <span>{{ reply.commentId === comment.id ? t('component.comment.list.cancelReply') : t('component.comment.list.reply') }}</span>
            </div>
            <div
              v-if="userStore.user && userStore.user.id == topicUserId && valueTypes.length > 0"
              class="comment-action-item"
              @click="showValue(comment)"
            >
              <i class="iconfont icon-star" />
              <span>价值</span>
            </div>
          </div>
          <div v-if="reply.commentId === comment.id" class="comment-reply-form">
            <text-editor
              :ref="setEditorRef(comment.id)"
              v-model:content="reply.value.content"
              v-model:imageList="reply.value.imageList"
              :height="80"
              @submit="submitReply(comment)"
            />
          </div>
          <CommentSubList
            v-if="
              comment.replies &&
              comment.replies.results &&
              comment.replies.results.length
            "
            :comment-id="comment.id"
            :data="comment.replies"
            :topic-user-id="topicUserId"
            :entity-type="entityType"
            :entity-id="entityId"
            @reply="onReply(comment, $event)"
          />
        </div>
      </div>
    </load-more-async>
    
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
  </div>
</template>

<script setup>
const { t } = useI18n();

const props = defineProps({
  entityType: {
    type: String,
    default: "",
    required: true,
  },
  entityId: {
    type: Number,
    default: 0,
    required: true,
  },
});
const reply = reactive({
  commentId: 0,
  value: {
    content: "",
    imageList: [],
  },
});

const showValueDialog = ref(false);
const currentComment = ref(null);
const valueTypes = ref([]);
const topicUserId = ref(0);

const editorRefs = ref({});
const userStore = useUserStore();
const loadMore = ref(null);

const append = (data) => {
  if (loadMore.value) {
    loadMore.value.refresh();
  }
};

const like = async (comment) => {
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
};

// 获取价值类型配置
const loadValueTypes = async () => {
  try {
    const config = await useHttpGet("/api/sys-config/configs");
    if (config && config.valueTypes && config.valueTypes.length > 0) {
      valueTypes.value = config.valueTypes;
    }
  } catch (e) {
    console.error("Failed to load value types:", e);
  }
};

// 获取话题用户信息
const loadTopicUser = async () => {
  if (props.entityType === 'topic' && props.entityId) {
    try {
      const topic = await useHttpGet(`/api/topic/${props.entityId}`);
      if (topic && topic.user) {
        topicUserId.value = topic.user.id;
      }
    } catch (e) {
      console.error("Failed to load topic user:", e);
    }
  }
};

// 处理评论数据，将子评论挂载到父评论下
const processComments = (results) => {
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
};

// 显示价值评价对话框
const showValue = (comment) => {
  if (!userStore.user) {
    useMsgSignIn();
    return;
  }
  
  // 只有主贴发布人可以评价价值
  if (userStore.user.id != topicUserId.value) {
    useMsgError("只有主贴发布人可以评价价值");
    return;
  }
  
  currentComment.value = comment;
  showValueDialog.value = true;
};

// 提交价值评价
const submitValue = async (valueType, index) => {
  try {
    await useHttpPost(
      "/api/comment/value",
      useJsonToForm({
        commentId: currentComment.value.id,
        valueType: index,
      })
    );
    useMsgSuccess(`评价成功：${valueType.label} (+${valueType.score}积分)`);
    showValueDialog.value = false;
    currentComment.value = null;
  } catch (e) {
    useCatchError(e);
  }
};

const switchShowReply = (comment) => {
  if (!userStore.user) {
    useMsgSignIn();
    return;
  }

  if (reply.commentId === comment.id) {
    hideReply(comment);
  } else {
    reply.commentId = comment.id;

    setTimeout(() => {
      const editorRef = editorRefs.value[`editor${comment.id}`];
      if (editorRef && typeof editorRef.focus === "function") {
        editorRef.focus();
      } else {
        console.warn(`Editor with id editor${comment.id} not found.`);
      }
    }, 100);
  }
};

const setEditorRef = (id) => {
  return (el) => {
    if (el) {
      editorRefs.value[`editor${id}`] = el;
    }
  };
};

const hideReply = (comment) => {
  reply.commentId = 0;
  reply.value.content = "";
  reply.value.imageList = [];
};

const submitReply = async (parent) => {
  try {
    const ret = await useHttpPost(
      "/api/comment/create",
      useJsonToForm({
        entityType: "comment",
        entityId: parent.id,
        content: reply.value.content,
        imageList:
          reply.value.imageList && reply.value.imageList.length
            ? JSON.stringify(reply.value.imageList)
            : "",
      })
    );
    hideReply();
    appendReply(parent, ret);
    useMsgSuccess(t('component.comment.list.publishSuccess'));
  } catch (e) {
    useCatchError(e);
  }
};

const onReply = (parent, comment) => {
  appendReply(parent, comment);
};

const appendReply = (parent, comment) => {
  if (parent.replies && parent.replies.results) {
    parent.replies.results.push(comment);
  } else {
    parent.replies = {
      results: [comment],
    };
  }
};

defineExpose({
  append,
});

onMounted(() => {
  loadValueTypes();
  loadTopicUser();
});
</script>

<style scoped lang="scss">
.comments {
  font-size: 14px;

  .comment {
    display: flex;
    padding: 10px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color4);
    }

    .comment-item-main {
      flex: 1 1 auto;
      margin-left: 10px;

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

        .comment-meta-right {
          .comment-time {
            font-size: 12px;
            color: var(--text-color3);
          }
          .comment-ip-area {
            font-size: 12px;
            color: var(--text-color3);
            margin-left: 10px;
          }
        }
      }

      .comment-content-wrapper {
        .comment-content {
          margin-top: 10px;
          margin-bottom: 0;
          color: var(--text-color);
          white-space: pre-wrap;
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

      .comment-replies {
        margin-top: 10px;
        // padding: 10px;
        background-color: var(--bg-color2);
      }
    }
  }

  .reply {
    display: flex;
  }

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
</style>
