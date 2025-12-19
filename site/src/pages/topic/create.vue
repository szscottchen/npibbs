<template>
  <section class="main">
    <div class="container">
      <article v-if="isNeedEmailVerify" class="message is-warning">
        <div class="message-header">
          <p>{{ $t("pages.topic.create.needEmailTitle") }}</p>
        </div>
        <div class="message-body">
          {{ $t("pages.topic.create.needEmailBody") }}
          <strong>
            <nuxt-link
              to="/user/profile/account"
              style="color: var(--text-link-color)"
              >{{ $t("pages.topic.create.goVerify") }}</nuxt-link
            ></strong
          >
        </div>
      </article>
      <div v-else class="publish-form">
        <div class="form-title">
          <div class="form-title-left">
            <div class="form-title-name">
              {{
                postForm.type === 0
                  ? $t("pages.topic.create.post")
                  : $t("pages.topic.create.tweet")
              }}
            </div>
            <div v-if="postForm.type === 0" class="form-title-support">
              <label class="checkbox">
                <input 
                  v-model="postForm.needAHand" 
                  type="checkbox"
                  :true-value="1"
                  :false-value="0"
                />
                <span class="support-label">呼叫支援</span>
              </label>
            </div>
          </div>
          <div
            v-if="postForm.type === 0"
            class="form-title-switch"
            @click="switchEditor"
          >
            <div v-if="postForm.contentType === 'markdown'" class="editor-type">
              <img src="~/assets/images/markdown.svg" />
              <span>Markdown</span>
            </div>
            <div v-else class="editor-type">
              <img src="~/assets/images/html.svg" />
              <span>HTML</span>
            </div>
            <i class="iconfont icon-switch" />
          </div>
        </div>

        <div class="topic-tags">
          <div
            v-for="node in nodes"
            :key="node.id"
            class="topic-tag"
            :class="{ selected: postForm.nodeId === node.id }"
            @click="postForm.nodeId = node.id"
          >
            <span>{{ node.name }}</span>
          </div>
        </div>

        <div v-if="postForm.type === 0" class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input topic-title"
              type="text"
              :placeholder="$t('pages.topic.create.titlePlaceholder')"
            />
          </div>
        </div>

        <div v-if="postForm.type === 0" class="field">
          <div class="control">
            <markdown-editor
              v-if="postForm.contentType === 'markdown'"
              v-model="postForm.content"
              :placeholder="$t('pages.topic.create.contentPlaceholder')"
            />
            <MEditor
              v-else
              v-model="postForm.content"
              :uploadImage="uploadImageForEditor"
            />
          </div>
        </div>

        <div v-if="postForm.type === 0 && isEnableHideContent" class="field">
          <div class="control">
            <MEditor v-model="postForm.hideContent" height="200px" />
          </div>
        </div>

        <div v-if="postForm.type === 1" class="field">
          <div class="control">
            <simple-editor
              ref="simpleEditorComponent"
              :placeholder="$t('pages.topic.create.contentPlaceholder')"
              v-model:content="postForm.content"
              v-model:imageList="postForm.imageList"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="form-footer">
          <a
            :class="{ 'is-loading': publishing }"
            class="button is-primary btn-publish"
            @click="publish"
            >{{
              postForm.type === 1
                ? $t("pages.topic.create.tweetBtn")
                : $t("pages.topic.create.postBtn")
            }}</a
          >
        </div>
      </div>
    </div>

    <CaptchaDialog ref="captchaDialog" />
  </section>
</template>

<script setup>
import { useI18n } from "vue-i18n";

definePageMeta({
  middleware: "auth",
});

const userStore = useUserStore();
const configStore = useConfigStore();
const route = useRoute();
const router = useRouter();

// 为MEditor创建兼容的上传函数
const uploadImageForEditor = async (file) => {
  const response = await useUploadImage(file);
  return {
    url: response.url,
    name: response.filename || file.name,
    isImage: response.isImage
  };
};

const type = Number.parseInt(route.query.type) || 0;
const nodeId =
  parseInt(route.query.nodeId) || configStore.config.defaultNodeId || 0;

if (type === 1 && !configStore.config.modules.tweet) {
  showError("😱 Tweet module is not enabled");
}
if (type === 0 && !configStore.config.modules.topic) {
  showError("😱 Topic module is not enabled");
}

const postForm = ref({
  type: type,
  nodeId: nodeId,
  title: "",
  tags: [],
  contentType: "html",
  content: "",
  hideContent: "",
  imageList: [],
  needAHand: false,

  captchaId: "",
  captchaCode: "",
  captchaProtocol: 2,
});
const publishing = ref(false);
const simpleEditorComponent = ref(null);
const captchaDialog = ref(null);

const { t } = useI18n();

const isNeedEmailVerify = computed(() => {
  return (
    configStore.config.createTopicEmailVerified && !userStore.user.emailVerified
  );
});

const isEnableHideContent = computed(() => {
  return configStore.config.enableHideContent;
});

const topicCaptchaEnabled = computed(() => {
  return configStore.config.topicCaptcha;
});

const { data: nodes } = await useMyFetch("/api/topic/nodes");

watch(
  () => route.query,
  (newQuery, oldQuery) => {
    init();
  },
  { deep: true }
);

const init = () => {
  postForm.value.type = Number.parseInt(route.query.type) || 0;

  let contentType = route.query.contentType;
  if (!contentType) {
    contentType = "html";
  }
  postForm.value.contentType = "contentType";

  useHead({
    title: useSiteTitle(
      type === 0 ? t("pages.topic.create.post") : t("pages.topic.create.tweet")
    ),
  });
};

init();

const switchEditor = () => {
  useConfirm(t("pages.topic.create.switchEditorConfirm"))
    .then(() => {
      postForm.value.content = "";
      if (postForm.value.contentType === "markdown") {
        postForm.value.contentType = "html";
      } else {
        postForm.value.contentType = "markdown";
      }
    })
    .catch(() => {});
};

const publish = () => {
  if (publishing.value) {
    return;
  }

  if (topicCaptchaEnabled.value) {
    captchaDialog.value.show().then((captcha) => {
      publishSubmit(captcha);
    });
  } else {
    publishSubmit();
  }
};

const publishSubmit = async (captcha) => {
  if (publishing.value) {
    return;
  }

  if (postForm.value.type === 1) {
    if (simpleEditorComponent.value.loading) {
      useMsgWarning(t("pages.topic.create.imageUploading"));
      return;
    }
  }

  if (captcha) {
    postForm.value.captchaId = captcha.captchaId;
    postForm.value.captchaCode = captcha.captchaCode;
    postForm.value.captchaProtocol = 2;
  }

  publishing.value = true;
  try {
    const submitData = {
      ...postForm.value,
      needAHand: postForm.value.needAHand ? 1 : 0
    };
    const topic = await useHttpPost("/api/topic/create", submitData);
    router.push(`/topic/${topic.id}`);
  } catch (e) {
    useMsgError(e.message || e);
    publishing.value = false;
  }
};
</script>

<style lang="scss" scoped>
.main {
  padding: 20px 0;

  .publish-form {
    background: var(--bg-color);
    border-radius: 8px;
    padding: 20px;
    box-shadow: 1px 1px 4px 0 rgba(0, 0, 0, 0.1);

    .form-title {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 20px;

      .form-title-left {
        display: flex;
        align-items: center;

        .form-title-name {
          font-size: 20px;
          font-weight: bold;
          color: var(--text-color);
        }

        .form-title-support {
          display: flex;
          align-items: center;
          margin-left: 15px;
          
          .checkbox {
            display: flex;
            align-items: center;
            cursor: pointer;
            
            input[type="checkbox"] {
              margin-right: 6px;
              cursor: pointer;
            }
            
            .support-label {
              color: var(--text-color);
              font-size: 14px;
              font-weight: normal;
            }
          }
        }
      }

      .form-title-switch {
        display: flex;
        align-items: center;
        cursor: pointer;
        color: var(--text-color4);
        font-size: 14px;

        .editor-type {
          display: flex;
          align-items: center;

          img {
            width: 16px;
            height: 16px;
            margin-right: 5px;
          }
        }

        .icon-switch {
          margin-left: 5px;
          font-size: 12px;
        }
      }
    }

    .topic-tags {
      display: flex;
      flex-wrap: wrap;
      margin-bottom: 20px;

      .topic-tag {
        padding: 4px 10px;
        background: var(--bg-color3);
        border-radius: 3px;
        margin-right: 10px;
        margin-bottom: 10px;
        cursor: pointer;
        font-size: 14px;

        &.selected {
          background: var(--primary-color);
          color: white;
        }
      }
    }

    .topic-title {
      font-size: 16px;
      height: 40px;
      border: 0;
      outline: 0;
      background: var(--bg-color2);
      color: var(--text-color);
    }

    .form-footer {
      margin-top: 20px;
      text-align: right;

      .btn-publish {
        padding: 0 30px;
        height: 36px;
        display: inline-flex;
        align-items: center;
      }
    }
  }
}
</style>
