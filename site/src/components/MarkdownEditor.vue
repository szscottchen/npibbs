<template>
  <client-only>
    <div class="markdown-editor-wrapper">
      <div class="upload-tips">
        <i class="iconfont icon-info" />
        <span>支持上传图片、文档、压缩包等任意文件类型（最大10MB）</span>
      </div>
      <MdEditor
        v-model="value"
        :theme="$colorMode.preference"
        @onChange="change"
        @onUploadImg="uploadFiles"
        :toolbars="toolbars"
        :style="{ height: height }"
        :placeholder="placeholder"
        :preview="true"
        :language="language"
        :footers="[]"
      />
    </div>
  </client-only>
</template>

<script setup>
import { MdEditor } from "md-editor-v3";
import "md-editor-v3/lib/style.css";

const language = computed(() => {
  const { locale } = useI18n();
  return locale.value;
});

const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },
  height: {
    type: String,
    default: "400px",
  },
  placeholder: {
    type: String,
    default: "",
  },
});

const emits = defineEmits(["update:modelValue"]);

const value = ref(props.modelValue);

const toolbars = ref([
  "bold",
  "underline",
  "italic",
  "strikeThrough",
  "-",
  "title",
  "sub",
  "sup",
  "quote",
  "unorderedList",
  "orderedList",
  "task",
  "-",
  "codeRow",
  "code",
  "link",
  "image",
  "table",
  // "mermaid",
  // "katex",
  "-",
  "revoke",
  "next",
  // "save",
  "-",
  // "pageFullscreen",
  "preview",
  // "htmlPreview",
  "catalog",
  "=",
  "fullscreen",
]);

function change(v) {
  emits("update:modelValue", v);
}

async function uploadFiles(files, callback) {
  console.log('📤 上传文件:', files);
  
  try {
    const res = await Promise.all(
      files.map((file) => {
        return useUploadImage(file);
      })
    );
    
    console.log('📥 上传结果:', res);
    
    // 分离图片和非图片文件
    const imageUrls = [];
    const fileMarkdownLinks = [];
    
    res.forEach((item) => {
      if (item.isImage) {
        // 图片文件：直接返回URL给md-editor-v3处理
        console.log('🖼️ 图片文件:', item.filename, '->', item.url);
        imageUrls.push(item.url);
      } else {
        // 非图片文件：构建带图标的HTML链接
        console.log('📎 非图片文件:', item.filename, '->', item.url);
        const fileSize = formatFileSize(item.size || 0);
        const iconClass = getFileIconClass(item.fileExt);
        const htmlLink = `<div class="file-attachment"><i class="${iconClass}"></i> <a href="${item.url}" target="_blank">${item.filename}</a> <span class="file-size">(${fileSize})</span></div>`;
        fileMarkdownLinks.push(htmlLink);
      }
    });
    
    // 先处理图片文件（通过callback）
    if (imageUrls.length > 0) {
      console.log('🖼️ 处理图片:', imageUrls);
      callback(imageUrls);
    }
    
    // 再手动插入非图片文件的markdown链接
    if (fileMarkdownLinks.length > 0) {
      console.log('📎 处理非图片文件:', fileMarkdownLinks);
      
      // 获取当前编辑器的内容
      const currentContent = value.value || '';
      
      // 在内容末尾添加文件链接（添加两个换行符作为分隔）
      const separator = currentContent ? '\n\n' : '';
      const newContent = currentContent + separator + fileMarkdownLinks.join('\n');
      
      // 更新编辑器内容
      value.value = newContent;
      emits('update:modelValue', newContent);
      
      console.log('✅ 已添加文件链接到编辑器:', fileMarkdownLinks);
    }
    
    console.log('✅ 上传完成');
  } catch (error) {
    console.error('❌ 上传失败:', error);
    // 发生错误时返回空数组，避免编辑器卡住
    callback([]);
  }
}

// 获取文件图标类名（参照Topic系统）
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

// 格式化文件大小函数
function formatFileSize(bytes) {
  if (bytes === 0) return '0 B';
  
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 添加样式
const uploadTipsStyle = `
.markdown-editor-wrapper {
  position: relative;
}

.upload-tips {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  margin-bottom: 10px;
  font-size: 14px;
  color: #6c757d;
}

.upload-tips .iconfont {
  margin-right: 8px;
  color: #17a2b8;
}
`;

// 动态添加样式
if (typeof document !== 'undefined') {
  const style = document.createElement('style');
  style.textContent = uploadTipsStyle;
  document.head.appendChild(style);
}
</script>