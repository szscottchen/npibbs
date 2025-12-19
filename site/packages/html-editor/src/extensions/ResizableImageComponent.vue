<template>
  <NodeViewWrapper class="resizable-image-wrapper" :class="{ 'is-selected': selected }">
    <!-- 判断是否为图片文件 -->
    <template v-if="isImageFile">
      <img
        ref="imageRef"
        :src="src"
        :alt="alt"
        :title="title"
        :width="currentWidth"
        :height="currentHeight"
        class="editor-image resizable"
        @load="onImageLoad"
        @click="selectImage"
      />
    </template>
    
    <!-- 非图片文件显示文件图标和名称 -->
    <template v-else>
      <div class="file-placeholder" @click="selectImage">
        <div class="file-icon">
          <i class="iconfont icon-file"></i>
        </div>
        <div class="file-info">
          <div class="file-name">{{ fileName }}</div>
          <div class="file-size">{{ fileSize }}</div>
        </div>
      </div>
    </template>
    
    <!-- 选择框和调整大小的手柄（仅图片文件） -->
    <template v-if="selected && isImageFile && imageLoaded">
      <!-- 选择框边框 -->
      <div class="selection-border">
        <div class="border-line border-top"></div>
        <div class="border-line border-right"></div>
        <div class="border-line border-bottom"></div>
        <div class="border-line border-left"></div>
      </div>
      
      <!-- 调整大小的手柄 -->
      <div
        v-for="handle in resizeHandles"
        :key="handle"
        :class="`resize-handle resize-handle-${handle}`"
        @mousedown="(e) => startResize(e, handle)"
      ></div>
    </template>
    
    <!-- 选择框（仅非图片文件） -->
    <template v-else-if="selected && !isImageFile">
      <div class="selection-border">
        <div class="border-line border-top"></div>
        <div class="border-line border-right"></div>
        <div class="border-line border-bottom"></div>
        <div class="border-line border-left"></div>
      </div>
    </template>
  </NodeViewWrapper>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { NodeViewWrapper } from '@tiptap/vue-3'

interface Props {
  editor: any
  node: any
  decorations: any
  selected: boolean
  extension: any
  getPos: () => number
  updateAttributes: (attributes: Record<string, any>) => void
  deleteNode: () => void
}

const props = defineProps<Props>()

const imageRef = ref<HTMLImageElement | null>(null)
const imageLoaded = ref(false)
const currentWidth = ref<number | null>(null)
const currentHeight = ref<number | null>(null)
const originalAspectRatio = ref<number>(1)

// 调整大小的手柄
const resizeHandles = ['nw', 'n', 'ne', 'e', 'se', 's', 'sw', 'w']

// 获取属性
const src = computed(() => props.node.attrs.src)
const alt = computed(() => props.node.attrs.alt)
const title = computed(() => props.node.attrs.title)

// 判断是否为图片文件
const isImageFile = computed(() => {
  // 检查src是否为base64图片数据
  if (src.value && src.value.startsWith('data:image/')) {
    return true
  }
  // 检查节点属性中的isImage字段（由服务器返回）
  if (props.node.attrs.isImage !== undefined) {
    return props.node.attrs.isImage
  }
  // 检查文件扩展名（备用方案）
  if (title.value) {
    const ext = title.value.toLowerCase().split('.').pop()
    const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg']
    return imageExts.includes(ext || '')
  }
  return false
})

// 获取文件名
const fileName = computed(() => {
  // 优先使用节点属性中的title，如果没有则尝试从src中提取文件名
  if (title.value) {
    return title.value
  }
  
  // 尝试从src中提取文件名（假设src是URL格式）
  if (src.value) {
    try {
      const url = new URL(src.value)
      const pathname = url.pathname
      const filename = pathname.split('/').pop()
      if (filename) {
        return filename
      }
    } catch (e) {
      // 如果src不是有效的URL，则直接使用src作为文件名
      return src.value
    }
  }
  
  return '未知文件'
})

// 获取文件大小
const fileSize = computed(() => {
  // 从节点属性中获取文件大小
  const size = props.node.attrs.size
  if (size !== undefined && size !== null) {
    return formatFileSize(size)
  }
  return '未知大小'
})

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 调整大小相关的状态
let isResizing = false
let startX = 0
let startY = 0
let startWidth = 0
let startHeight = 0
let currentHandle = ''

onMounted(() => {
  // 初始化尺寸
  currentWidth.value = props.node.attrs.width
  currentHeight.value = props.node.attrs.height
  
  // 添加全局事件监听
  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
})

onBeforeUnmount(() => {
  document.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseup', handleMouseUp)
})

/**
 * 图片加载完成时的处理
 */
function onImageLoad() {
  imageLoaded.value = true
  
  if (!imageRef.value) return
  
  const img = imageRef.value
  originalAspectRatio.value = img.naturalWidth / img.naturalHeight
  
  // 如果没有设置初始尺寸，使用图片的原始尺寸（限制最大宽度）
  if (!currentWidth.value && !currentHeight.value) {
    const maxWidth = 600 // 最大宽度
    const naturalWidth = img.naturalWidth
    const naturalHeight = img.naturalHeight
    
    if (naturalWidth > maxWidth) {
      currentWidth.value = maxWidth
      currentHeight.value = Math.round(maxWidth / originalAspectRatio.value)
    } else {
      currentWidth.value = naturalWidth
      currentHeight.value = naturalHeight
    }
    
    // 更新节点属性
    props.updateAttributes({
      width: currentWidth.value,
      height: currentHeight.value,
    })
  }
}

/**
 * 选择图片
 */
function selectImage() {
  const pos = props.getPos()
  props.editor.commands.setNodeSelection(pos)
}

/**
 * 开始调整大小
 */
function startResize(event: MouseEvent, handle: string) {
  event.preventDefault()
  event.stopPropagation()
  
  if (!imageRef.value) return
  
  isResizing = true
  currentHandle = handle
  startX = event.clientX
  startY = event.clientY
  startWidth = currentWidth.value || imageRef.value.offsetWidth
  startHeight = currentHeight.value || imageRef.value.offsetHeight
  
  // 添加调整大小时的样式
  document.body.style.cursor = getResizeCursor(handle)
  props.editor.view.dom.classList.add('resizing-image')
}

/**
 * 处理鼠标移动
 */
function handleMouseMove(event: MouseEvent) {
  if (!isResizing || !imageRef.value) return
  
  event.preventDefault()
  
  const deltaX = event.clientX - startX
  const deltaY = event.clientY - startY
  
  let newWidth = startWidth
  let newHeight = startHeight
  
  // 根据手柄方向计算新尺寸
  switch (currentHandle) {
    case 'se': // 右下角
    case 'e': // 右边
    case 's': // 下边
      newWidth = startWidth + deltaX
      break
    case 'sw': // 左下角
    case 'w': // 左边
      newWidth = startWidth - deltaX
      break
    case 'ne': // 右上角
      newWidth = startWidth + deltaX
      break
    case 'nw': // 左上角
    case 'n': // 上边
      newWidth = startWidth - deltaX
      break
  }
  
  // 限制最小和最大尺寸
  newWidth = Math.max(50, Math.min(800, newWidth))
  
  // 根据宽高比计算新高度
  newHeight = Math.round(newWidth / originalAspectRatio.value)
  
  // 更新当前尺寸
  currentWidth.value = newWidth
  currentHeight.value = newHeight
}

/**
 * 处理鼠标释放
 */
function handleMouseUp() {
  if (!isResizing) return
  
  isResizing = false
  document.body.style.cursor = ''
  props.editor.view.dom.classList.remove('resizing-image')
  
  // 更新节点属性
  props.updateAttributes({
    width: currentWidth.value,
    height: currentHeight.value,
  })
}

/**
 * 获取调整大小时的鼠标样式
 */
function getResizeCursor(handle: string): string {
  const cursors: Record<string, string> = {
    'nw': 'nw-resize',
    'n': 'n-resize',
    'ne': 'ne-resize',
    'e': 'e-resize',
    'se': 'se-resize',
    's': 's-resize',
    'sw': 'sw-resize',
    'w': 'w-resize',
  }
  return cursors[handle] || 'default'
}
</script>

<style lang="scss">
.resizable-image-wrapper {
  position: relative;
  display: inline-block;
  margin: 0 !important;
  
  &.is-selected {
    .resize-handle {
      opacity: 1;
      pointer-events: auto;
    }
    
    .selection-border {
      opacity: 1;
    }
  }
  
  .selection-border {
    position: absolute;
    top: -2px;
    left: -2px;
    right: -2px;
    bottom: -2px;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
    
    .border-line {
      position: absolute;
      background-color: #3b82f6;
      
      &.border-top {
        top: 0;
        left: 0;
        right: 0;
        height: 2px;
      }
      
      &.border-right {
        top: 0;
        right: 0;
        bottom: 0;
        width: 2px;
      }
      
      &.border-bottom {
        bottom: 0;
        left: 0;
        right: 0;
        height: 2px;
      }
      
      &.border-left {
        top: 0;
        left: 0;
        bottom: 0;
        width: 2px;
      }
    }
  }
  
  .editor-image.resizable {
    border: none !important;
    outline: none !important;
    box-shadow: none !important;
    background: transparent !important;
    border-radius: 4px;
    transition: all 0.2s ease;
    display: block;
  }
  
  .resize-handle {
    position: absolute;
    width: 8px;
    height: 8px;
    background-color: #3b82f6;
    border: 1px solid white;
    border-radius: 50%;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
    z-index: 100;
    box-shadow: 0 0 2px rgba(0, 0, 0, 0.3);
    
    &:hover {
      background-color: #2563eb;
      transform: scale(1.2);
    }
  }
  
  // 角落手柄
  .resize-handle-nw {
    top: -6px;
    left: -6px;
    cursor: nw-resize;
  }
  
  .resize-handle-ne {
    top: -6px;
    right: -6px;
    cursor: ne-resize;
  }
  
  .resize-handle-sw {
    bottom: -6px;
    left: -6px;
    cursor: sw-resize;
  }
  
  .resize-handle-se {
    bottom: -6px;
    right: -6px;
    cursor: se-resize;
  }
  
  // 边缘中间的手柄
  .resize-handle-n {
    top: -6px;
    left: 50%;
    transform: translateX(-50%);
    cursor: n-resize;
  }
  
  .resize-handle-s {
    bottom: -6px;
    left: 50%;
    transform: translateX(-50%);
    cursor: s-resize;
  }
  
  .resize-handle-e {
    right: -6px;
    top: 50%;
    transform: translateY(-50%);
    cursor: e-resize;
  }
  
  .resize-handle-w {
    left: -6px;
    top: 50%;
    transform: translateY(-50%);
    cursor: w-resize;
  }
  
  // 非图片文件的占位符样式
  .file-placeholder {
    display: flex;
    align-items: center;
    padding: 12px;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    background-color: #f9fafb;
    cursor: pointer;
    min-width: 200px;
    
    .file-icon {
      font-size: 32px;
      margin-right: 12px;
      color: #6b7280;
    }
    
    .file-info {
      flex: 1;
      
      .file-name {
        font-weight: 500;
        color: #1f2937;
        margin-bottom: 4px;
        word-break: break-all;
      }
      
      .file-size {
        font-size: 12px;
        color: #6b7280;
      }
    }
  }
}

// 全局样式：调整大小时禁用文本选择
:global(.ProseMirror.resizing-image) {
  user-select: none !important;
}
</style>
