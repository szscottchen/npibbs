/**
 * 图片上传相关工具函数
 */

export type UploadImageResponse = {
  url: string;
  name?: string;
  width?: number;
  height?: number;
  isImage?: boolean;
  size?: number;
};

export type UploadImageFunction = (file: File) => Promise<UploadImageResponse>;

// 支持的图片格式
export const SUPPORTED_IMAGE_TYPES = [
  'image/jpeg',
  'image/jpg',
  'image/png',
  'image/gif',
  'image/webp',
  'image/svg+xml'
]

// 最大文件大小 (10MB)
export const MAX_FILE_SIZE = 10 * 1024 * 1024

/**
 * 检查文件大小是否符合要求
 */
export function isValidFileSize(file: File): boolean {
  return file.size <= MAX_FILE_SIZE
}

/**
 * 将文件转换为Base64格式
 */
export function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      if (typeof reader.result === 'string') {
        resolve(reader.result)
      } else {
        reject(new Error('Failed to convert file to base64'))
      }
    }
    reader.onerror = () => reject(new Error('Failed to read file'))
    reader.readAsDataURL(file)
  })
}

/**
 * 创建文件选择器
 */
export function createFileInput(acceptAllFiles: boolean = false): HTMLInputElement {
  const input = document.createElement('input')
  input.type = 'file'
  // 如果acceptAllFiles为true，接受所有文件类型，否则只接受图片
  input.accept = acceptAllFiles ? '*' : SUPPORTED_IMAGE_TYPES.join(',')
  input.style.display = 'none'
  return input
}

/**
 * 格式化文件大小显示
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 图片上传处理函数
 * 这里使用 base64 作为示例，实际项目中应该上传到服务器
 */
export async function uploadImage(file: File): Promise<UploadImageResponse> {
  if (!isValidFileSize(file)) {
    throw new Error(`文件大小超过限制。最大支持：${formatFileSize(MAX_FILE_SIZE)}`)
  }
  try {
    // 这里使用 base64 作为演示
    const base64 = await fileToBase64(file)
    return {
      url: base64,
      name: file.name,
    }
  } catch (error) {
    throw new Error('图片上传失败：' + (error as Error).message)
  }
}

/**
 * 通用文件上传处理函数
 * 支持所有文件类型
 */
export async function uploadFile(file: File): Promise<UploadImageResponse> {
  // 构造表单数据
  const formData = new FormData();
  formData.append('file', file, file.name);
  
  try {
    // 使用原生 fetch API 调用上传接口
    const response = await fetch('/api/upload', {
      method: 'POST',
      body: formData,
    });
    
    // 检查响应状态
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    // 解析响应数据
    const result: any = await response.json();
    
    // 根据服务器返回的数据构建响应
    return {
      url: result.url,
      name: result.filename,
      isImage: result.isImage,
      size: result.size || file.size, // 优先使用服务器返回的size，如果没有则使用原始文件大小
      // 如果是图片且有宽高信息，可以设置宽高
      // width: result.isImage && result.width ? result.width : undefined,
      // height: result.isImage && result.height ? result.height : undefined,
    };
  } catch (error) {
    throw new Error('文件上传失败：' + (error as Error).message);
  }
}
