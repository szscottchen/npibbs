<template>
  <div class="image-upload-button">
    <ToolbarButton
      title="上传文件"
      @click="handleImageUpload"
      :disabled="!editor"
    >
      <LucideUpload :size="TOOLBAR_ICON_SIZE" />
    </ToolbarButton>
    
    <!-- 文件上传进度提示 -->
    <div v-if="uploading" class="upload-progress">
      <div class="upload-spinner"></div>
      <span>上传中...</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Editor } from '@tiptap/core'
import { LucideUpload } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'
import type { UploadImageFunction } from "../../utils/imageUtils";

import {
  createFileInput, 
  isValidFileSize,
  formatFileSize,
  MAX_FILE_SIZE
} from '../../utils/imageUtils'

const props = defineProps<{
  editor: Editor | null | undefined
  uploadImage: UploadImageFunction
}>()

const uploading = ref(false)

/**
 * 处理文件上传（默认显示所有文件类型）
 */
async function handleImageUpload() {
  if (!props.editor) return
  
  // 创建支持所有文件类型的文件选择器（默认显示所有文件）
  const input = createFileInput(true)
  
  input.onchange = async (event) => {
    const files = (event.target as HTMLInputElement).files
    if (!files || files.length === 0) return
    
    const file = files[0]
    
    try {
      uploading.value = true
      
      if (!isValidFileSize(file)) {
        alert(`文件大小超过限制！最大支持 ${formatFileSize(MAX_FILE_SIZE)}`)
        return
      }
      
      // 上传文件（支持所有文件类型）
      const resp = await props.uploadImage(file)
      
      // 插入文件到编辑器，传递服务器返回的 isImage 字段
      props.editor?.chain().focus().setResizableImage({ 
        src: resp.url,
        alt: resp.name || '',
        title: resp.name || '',
        isImage: resp.isImage !== undefined ? resp.isImage : undefined,
      }).run()
      
    } catch (error) {
      console.error('文件上传失败:', error)
      alert(error instanceof Error ? error.message : '文件上传失败，请重试。')
    } finally {
      uploading.value = false
      // 清理input
      document.body.removeChild(input)
    }
  }
  
  // 添加到DOM并触发点击
  document.body.appendChild(input)
  input.click()
}
</script>

<style lang="scss" scoped>
.image-upload-button {
  position: relative;
  display: inline-block;
}

.upload-progress {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background: var(--editor-bg);
  border: 1px solid var(--editor-border);
  border-radius: 4px;
  padding: 8px 12px;
  font-size: 12px;
  color: var(--editor-text);
  white-space: nowrap;
  z-index: 1000;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  gap: 6px;
}

.upload-spinner {
  width: 12px;
  height: 12px;
  border: 2px solid var(--editor-border);
  border-top: 2px solid var(--editor-focus);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>