<template>
  <div class="image-uploads">
    <div
      v-for="(file, index) in previewFiles"
      :key="index"
      class="preview-item"
      :class="{ deleted: file.deleted }"
      :style="{ width: size, height: size }"
    >
      <img v-if="file.isImage" :src="file.url" class="image-item" />
      <div v-else class="file-icon">
        <i :class="getFileIconClass(file.fileExt)"></i>
        <div class="file-ext">{{ file.fileExt.replace('.', '') }}</div>
      </div>
      <el-progress
        v-show="file.progress < 100"
        :percentage="file.progress"
        color="#25A9F6"
        :show-text="false"
        class="progress"
      />
      <div v-show="file.progress < 100" class="cover">
        {{ $t("component.imageUpload.uploading") }}
      </div>
      <div
        class="upload-delete"
        :class="{
          'show-delete': file.progress === 100,
        }"
        @click="removeItem(index)"
      >
        <i class="iconfont icon-delete" />
      </div>
      <div v-if="!file.isImage && file.progress === 100" class="file-name">
        {{ file.filename }}
      </div>
    </div>
    <div
      v-show="previewFiles.length < limit"
      class="add-image-btn"
      :style="{ width: size, height: size }"
      @click="onClick($event)"
    >
      <input
        ref="currentInput"
        :accept="accept"
        type="file"
        multiple
        @input="onInput"
      />
      <div class="add-image-btn-wrapper">
        <slot name="add-image-button">
          <i class="iconfont icon-add" style="font-size: 15px; opacity: 0.4" />
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue';
const { t } = useI18n();

const props = defineProps({
  modelValue: {
    type: Array,
    default() {
      return [];
    },
  },
  accept: {
    type: String,
    default: "image/jpeg,image/png,image/gif,image/bmp,image/webp,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.pdf,.txt,.zip,.rar",
  },
  limit: {
    type: Number,
    default: 9,
  },
  sizeLimit: {
    type: Number,
    default: 1024 * 1024 * 20,
  },
  size: {
    type: String,
    default: "94px",
  },
});

const emit = defineEmits(["update:modelValue"]);

const fileList = ref(props.modelValue || []);
const previewFiles = ref([]);
const currentInput = ref(null);
const loading = ref(false);

watch(
  () => props.modelValue,
  (newModelValue) => {
    // 更新 fileList
    fileList.value = Array.isArray(newModelValue) ? [...newModelValue] : [];

    // 更新 previewFiles
    previewFiles.value = fileList.value.map((item) => ({
      name: item.name || "unknown", // 假设每个文件对象有 name 属性
      url: item.url, // 假设每个文件对象有 url 属性
      progress: 100, // 已上传完成
      deleted: false,
      size: item.size || 0, // 假设每个文件对象有 size 属性
    }));
  },
  { immediate: true } // 立即执行一次，确保初始化时同步
);

const onClick = () => {
  if (currentInput.value) {
    currentInput.value.dispatchEvent(new MouseEvent("click"));
  }
};

const onInput = (e) => {
  const files = e.target.files;
  addFiles(files);
};

const addFiles = (files) => {
  if (!files || !files.length) return; // 没有文件
  if (!checkSizeLimit(files)) return; // 文件大小检查
  if (!checkLengthLimit(files)) return; // 文件数量检查

  const fileArray = [];
  for (let i = 0; i < files.length; i++) {
    const url = getObjectURL(files[i]);
    const ext = '.' + files[i].name.split('.').pop().toLowerCase();
    const isImage = /^\.(jpe?g|png|gif|bmp|webp)$/i.test(ext);
    previewFiles.value.push({
      name: files[i].name,
      url,
      progress: 0,
      deleted: false,
      size: files[i].size,
      filename: files[i].name,
      isImage,
      fileExt: ext,
    });
    fileArray.push(files[i]);
  }
  const promiseList = fileArray.reduce((result, file, index, array) => {
    result.push(uploadFile(file, index, array.length));
    return result;
  }, []);
  uploadFiles(promiseList);
};

const uploadFile = (file, index, length) => {
  const formData = new FormData();
  formData.append("file", file, file.name);
  return useHttp("/api/upload", {
    method: "POST",
    body: formData,
  });
};
const uploadFiles = (promiseList) => {
  loading.value = true;

  Promise.all(promiseList).then(
    (resList) => {
      // 请求响应后，更新到 100%
      // 使用新的数组替换previewFiles.value，而不是直接修改每个项
      previewFiles.value = previewFiles.value.map(item => ({
        ...item,
        progress: 100
      }));
      
      // 添加上传结果到fileList
      resList.forEach((item) => {
        fileList.value.push({
          url: item.url,
          preview: item.url,
          filename: item.filename,
          isImage: item.isImage,
          fileExt: item.fileExt,
        });
      });
      
      if (currentInput.value) {
        currentInput.value.value = "";
      }
      loading.value = false;
      
      // 使用nextTick延迟触发更新，避免无限递归
      nextTick(() => {
        emit("update:modelValue", [...fileList.value]);
      });
    },
    (e) => {
      useMsgError(e.message || e);

      if (currentInput.value) {
        currentInput.value.value = "";
      }

      // 失败的时候取消对应的预览照片
      const length = promiseList.length;
      previewFiles.value.splice(previewFiles.value.length - length, length);

      loading.value = false;
    }
  );
};
const removeItem = (index) => {
  ElMessageBox.confirm(
    t("component.imageUpload.confirmDelete"),
    t("component.imageUpload.deleteTitle"),
    {
      confirmButtonText: t("component.imageUpload.confirmButton"),
      cancelButtonText: t("component.imageUpload.cancelButton"),
      type: "warning",
    }
  ).then(
    () => {
      previewFiles.value[index].deleted = true; // 删除动画
      fileList.value.splice(index, 1);
      emit("update:modelValue", fileList.value); // 避免和回显冲突，先修改 fileList
      setTimeout(() => {
        previewFiles.value.splice(index, 1);
      }, 900);
    },
    () => console.log("取消删除")
  );
};
const checkSizeLimit = (files) => {
  let pass = true;
  for (let i = 0; i < files.length; i++) {
    if (files[i].size > props.sizeLimit) {
      pass = false;
    }
  }
  if (!pass)
    useMsgError(
      t("component.imageUpload.sizeLimitError", {
        size: props.sizeLimit / 1024 / 1024,
      })
    );
  return pass;
};
const checkLengthLimit = (files) => {
  if (previewFiles.value.length + files.length > props.limit) {
    useMsgWarning(
      t("component.imageUpload.countLimitError", { limit: props.limit })
    );
    return false;
  } else {
    return true;
  }
};
const getObjectURL = (file) => {
  let url = null;
  if (window.createObjectURL) {
    // basic
    url = window.createObjectURL(file);
  } else if (window.URL) {
    // mozilla(firefox)
    url = window.URL.createObjectURL(file);
  } else if (window.webkitURL) {
    // webkit or chrome
    url = window.webkitURL.createObjectURL(file);
  }
  return url;
};
const clear = () => {
  fileList.value = [];
  previewFiles.value = [];
};

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

defineExpose({
  onClick,
  addFiles,
  clear,
  loading,
  getFileIconClass,
});
</script>

<style lang="scss" scoped>
.image-uploads {
  display: flex;
  column-gap: 10px;
  row-gap: 10px;

  .preview-item {
    position: relative;
    border: 1px solid var(--border-color);
    border-radius: 2px;

    &.deleted {
      transition: 1s all;
      transform: translateY(-100%);
      opacity: 0;
    }

    .image-item {
      cursor: pointer;
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .file-icon {
      width: 100%;
      height: 100%;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      background: #f0f0f0;
    }

    .file-icon i {
      font-size: 32px;
      margin-bottom: 4px;
      color: #666;
    }

    .file-icon .file-ext {
      font-size: 12px;
      color: #666;
    }

    .file-name {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      background: rgba(0, 0, 0, 0.7);
      color: white;
      font-size: 10px;
      padding: 2px;
      text-align: center;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .progress {
      position: absolute;
      top: 80px;
      width: 100%;
      height: 6px;
      padding: 0 10px;
    }

    .cover {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      color: var(--text-color2);
      background: rgba(255, 255, 255, 0.5);
      font-size: 12px;
      display: flex;
      justify-content: center;
      align-items: center;
    }

    .upload-delete {
      cursor: pointer;
      position: absolute;
      left: 0;
      bottom: 0;
      height: 20px;
      width: 100%;
      display: none;
      justify-content: center;
      align-items: center;
      background: rgba(0, 0, 0, 0.3);
      text-align: center;
      vertical-align: middle;
      line-height: 20px;
      z-index: 10;

      i.iconfont {
        font-size: 14px;
        fill: white;
        color: var(--text-color5);
        font-weight: 700;
      }
    }

    &:hover {
      .upload-delete.show-delete {
        display: flex;
      }
    }
  }

  .add-image-btn {
    cursor: pointer;
    border: 1px solid var(--border-color);
    border-radius: 2px;
    position: relative;

    input[type="file"] {
      cursor: pointer;
      display: none;
    }

    .add-image-btn-wrapper {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
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
  color: #333;
}

.icon-file-zip {
  color: #ffb900;
}

.icon-file {
  color: #666;
}
</style>
