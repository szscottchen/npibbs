import { Extension } from '@tiptap/core'
import { Plugin, PluginKey } from '@tiptap/pm/state'
import { UploadImageFunction } from '../utils/imageUtils'

export interface PasteImageOptions {
  /**
   * 是否启用粘贴图片功能
   */
  enabled: boolean
  /**
   * 自定义图片上传函数
   */
  uploadImage: UploadImageFunction
}

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 粘贴图片扩展
 * 支持粘贴剪贴板中的截图和复制的磁盘图片文件
 */
export const PasteImage = Extension.create<PasteImageOptions>({
  name: 'pasteImage',

  addOptions() {
    return {
      enabled: true,
      uploadImage: undefined,
    }
  },

  onCreate() {
  },

  addProseMirrorPlugins() {
    return [
      new Plugin({
        key: new PluginKey('pasteImage'),
        props: {
          handlePaste: (view, event, slice) => {
            console.log('📋 接收到粘贴事件', { enabled: this.options.enabled, event })

            if (!this.options.enabled) {
              return false
            }

            // 优先使用 files，如果没有则使用 items
            let imageFiles: File[] = []

            // 直接从 files 获取所有文件（支持所有文件类型）
            const files = Array.from(event.clipboardData?.files || [])
            imageFiles = files

            // 如果 files 中没有文件，再尝试从 items 中获取
            if (imageFiles.length === 0) {
              const items = Array.from(event.clipboardData?.items || [])
              imageFiles = items
                .filter(item => item.kind === 'file')
                .map(item => item.getAsFile())
                .filter((file): file is File => file !== null)
            }

            if (imageFiles.length === 0) {
              console.log('📋 粘贴事件中没有找到图片文件')
              return false
            }

            console.log('📋 找到文件:', imageFiles.map(f => ({ name: f.name, type: f.type, size: f.size })))

            // 阻止默认粘贴行为
            event.preventDefault()

            // 处理图片文件
            imageFiles.forEach((file, index) => {
              setTimeout(async () => {
                try {
                  console.log('开始处理粘贴的图片:', file.name, file.type, file.size)

                  // 使用配置中的uploadImage函数或默认的uploadImage函数
                  const uploadImageFn = this.options.uploadImage
                  const resp = await uploadImageFn(file)

                  // 获取当前光标位置
                  const currentState = view.state
                  const pos = currentState.selection.from + index // 为每个图片偏移位置

                  // 根据文件类型创建不同的节点
                  let node
                  if (resp.isImage) {
                    // 图片文件：创建resizableImage节点
                    node = currentState.schema.nodes.resizableImage.create({
                      src: resp.url,
                      alt: resp.name || '',
                      title: resp.name || '',
                    })
                  } else {
                    // 非图片文件：创建普通链接
                    const linkText = `${resp.name || '文件'}${resp.size ? ` (${formatFileSize(resp.size)})` : ''}`
                    node = currentState.schema.text(linkText, [
                      currentState.schema.marks.link.create({
                        href: resp.url,
                        title: resp.name || '',
                        target: '_blank',
                      })
                    ])
                  }

                  // 插入节点
                  const tr = currentState.tr.replaceWith(pos, pos, node)
                  view.dispatch(tr)

                  console.log('文件粘贴成功')
                } catch (error) {
                  console.error('粘贴图片失败:', error)
                  alert('文件粘贴失败: ' + (error instanceof Error ? error.message : '未知错误'))
                }
              }, index * 10) // 轻微延迟以确保顺序
            })

            return true
          },

          handleDrop: (view, event, slice, moved) => {
            console.log('🖱️ 接收到拖拽事件', { enabled: this.options.enabled, event })

            if (!this.options.enabled) {
              return false
            }

            const files = Array.from(event.dataTransfer?.files || [])
            const imageFiles = files

            if (imageFiles.length === 0) {
              console.log('🖱️ 拖拽事件中没有找到图片文件')
              return false
            }

            console.log('🖱️ 找到拖拽图片文件:', imageFiles.map(f => ({ name: f.name, type: f.type, size: f.size })))

            // 阻止默认拖拽行为
            event.preventDefault()

            // 获取拖拽位置
            const coordinates = view.posAtCoords({
              left: event.clientX,
              top: event.clientY,
            })

            if (!coordinates) {
              return false
            }

            // 处理拖拽的图片文件
            imageFiles.forEach((file, index) => {
              setTimeout(async () => {
                try {
                  console.log('开始处理拖拽的文件:', file.name, file.type, file.size)

                  // 使用配置中的uploadImage函数或默认的uploadImage函数
                  const uploadImageFn = this.options.uploadImage
                  const resp = await uploadImageFn(file)

                  // 为每个文件计算插入位置（避免重叠）
                  const insertPos = coordinates.pos + index

                  // 根据文件类型创建不同的节点
                  let node
                  if (resp.isImage) {
                    // 图片文件：创建resizableImage节点
                    node = view.state.schema.nodes.resizableImage.create({
                      src: resp.url,
                      alt: resp.name || '',
                      title: resp.name || '',
                    })
                  } else {
                    // 非图片文件：创建普通链接
                    const linkText = `${resp.name || '文件'}${resp.size ? ` (${formatFileSize(resp.size)})` : ''}`
                    node = view.state.schema.text(linkText, [
                      view.state.schema.marks.link.create({
                        href: resp.url,
                        title: resp.name || '',
                        target: '_blank',
                      })
                    ])
                  }

                  // 在拖拽位置插入节点
                  const tr = view.state.tr.replaceWith(insertPos, insertPos, node)
                  view.dispatch(tr)

                  console.log('文件拖拽成功')
                } catch (error) {
                  console.error('拖拽图片失败:', error)
                  alert('文件拖拽失败: ' + (error instanceof Error ? error.message : '未知错误'))
                }
              }, index * 10) // 轻微延迟以确保顺序
            })

            return true
          },
        },
      }),
    ]
  },
})
