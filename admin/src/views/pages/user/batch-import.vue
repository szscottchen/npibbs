<template>
  <div class="container">
    <div class="container-header">
      <h2>{{ $t('pages.user.batchImport.title') }}</h2>
    </div>

    <div class="container-main">
      <a-card>
        <!-- 使用说明 -->
        <a-alert type="info" :closable="false" style="margin-bottom: 20px">
          <template #icon><icon-info-circle /></template>
          <div>
            <p>{{ $t('pages.user.batchImport.step1') }}</p>
            <p>{{ $t('pages.user.batchImport.step2') }}</p>
            <p>{{ $t('pages.user.batchImport.step3') }}</p>
            <p>{{ $t('pages.user.batchImport.step4') }}</p>
          </div>
        </a-alert>

        <!-- 模板下载 -->
        <div style="margin-bottom: 30px">
          <a-button type="primary" @click="downloadTemplate">
            <template #icon><icon-download /></template>
            {{ $t('pages.user.batchImport.downloadTemplate') }}
          </a-button>
        </div>

        <!-- 文件上传 -->
        <a-upload
          :custom-request="handleUpload"
          :before-upload="beforeUpload"
          :show-file-list="true"
          :limit="1"
          :auto-upload="false"
          accept=".xlsx,.xls"
          draggable
          ref="uploadRef"
        >
          <template #upload-button>
            <div class="upload-area">
              <div>
                <icon-upload />
              </div>
              <div class="upload-text">
                {{ $t('pages.user.batchImport.dragFileHere') }}
                <span class="upload-link">{{ $t('pages.user.batchImport.clickToUpload') }}</span>
              </div>
              <div class="upload-tip">
                {{ $t('pages.user.batchImport.onlyExcelFile') }}
              </div>
            </div>
          </template>
        </a-upload>

        <!-- 导入按钮 -->
        <div style="margin-top: 20px; text-align: center">
          <a-button
            type="primary"
            size="large"
            :loading="importing"
            @click="submitUpload"
          >
            <template #icon><icon-upload /></template>
            {{ importing ? $t('pages.user.batchImport.importing') : $t('pages.user.batchImport.startImport') }}
          </a-button>
        </div>

        <!-- 进度条 -->
        <div v-if="importing" style="margin-top: 30px">
          <a-progress :percent="progress" :status="progressStatus" />
        </div>

        <!-- 导入结果 -->
        <div v-if="importResult" style="margin-top: 30px">
          <a-alert type="success" :closable="false" style="margin-bottom: 20px">
            <template #icon><icon-check-circle /></template>
            <div>
              <p><strong>{{ $t('pages.user.batchImport.importCompleted') }}</strong></p>
              <p>{{ $t('pages.user.batchImport.totalCount') }}: {{ importResult.totalCount }}</p>
              <p>{{ $t('pages.user.batchImport.successCount') }}: {{ importResult.successCount }}</p>
              <p>{{ $t('pages.user.batchImport.failCount') }}: {{ importResult.failCount }}</p>
            </div>
          </a-alert>

          <!-- 错误详情表格 -->
          <a-table
            v-if="importResult.errors && importResult.errors.length > 0"
            :data="importResult.errors"
            :pagination="false"
            :bordered="true"
          >
            <template #columns>
              <a-table-column
                :title="$t('pages.user.batchImport.rowNumber')"
                data-index="rowNumber"
                :width="100"
              />
              <a-table-column
                :title="$t('pages.user.batchImport.employeeId')"
                data-index="employeeId"
                :width="150"
              />
              <a-table-column
                :title="$t('pages.user.batchImport.nickname')"
                data-index="nickname"
                :width="150"
              />
              <a-table-column
                :title="$t('pages.user.batchImport.errorMessage')"
                data-index="errorMessage"
              />
            </template>
          </a-table>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useI18n } from 'vue-i18n'
import { downloadTemplate as downloadTemplateApi, batchImportUsers } from '@/api/user'
import {
  IconInfoCircle,
  IconDownload,
  IconUpload,
  IconCheckCircle,
} from '@arco-design/web-vue/es/icon'

const { t } = useI18n()

const uploadRef = ref(null)
const importing = ref(false)
const progress = ref(0)
const progressStatus = ref('normal')
const importResult = ref(null)

// 下载模板
const downloadTemplate = () => {
  downloadTemplateApi().then((res) => {
    // 添加调试信息
    console.log('Download template response:', res);
    
    // 检查响应是否为Blob对象
    if (res instanceof Blob) {
      // 直接使用响应作为Blob
      const blob = res;
      
      // 检查Blob大小
      console.log('Blob size:', blob.size);
      
      if (blob.size === 0) {
        Message.error('下载失败: 文件内容为空');
        return;
      }
      
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', '用户导入模板.xlsx');
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } else if (res && res.data instanceof Blob) {
      // 如果响应中有data属性且为Blob
      const blob = res.data;
      
      // 检查Blob大小
      console.log('Blob size:', blob.size);
      
      if (blob.size === 0) {
        Message.error('下载失败: 文件内容为空');
        return;
      }
      
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', '用户导入模板.xlsx');
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } else {
      Message.error('下载失败: 响应格式不正确');
      console.error('Unexpected response format:', res);
    }
  }).catch((error) => {
    console.error('Download template error:', error);
    Message.error('下载失败: ' + (error.message || '未知错误'));
  });
}

// 上传前验证
const beforeUpload = (file) => {
  const isExcel =
    file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
    file.type === 'application/vnd.ms-excel' ||
    file.name.endsWith('.xlsx') ||
    file.name.endsWith('.xls')
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isExcel) {
    Message.error(t('pages.user.batchImport.onlyExcelFile'))
    return false
  }
  if (!isLt10M) {
    Message.error(t('pages.user.batchImport.fileSizeLimit'))
    return false
  }
  return true
}

// 提交上传
const submitUpload = () => {
  uploadRef.value?.submit()
}

// 自定义上传
const handleUpload = async (options) => {
  importing.value = true
  progress.value = 0
  progressStatus.value = 'normal'
  importResult.value = null

  const formData = new FormData()
  formData.append('file', options.fileItem.file)

  // 模拟进度
  const progressTimer = setInterval(() => {
    if (progress.value < 90) {
      progress.value += 10
    }
  }, 200)

  try {
    const result = await batchImportUsers(formData)

    clearInterval(progressTimer)
    progress.value = 100
    progressStatus.value = 'success'
    importResult.value = result

    // 显示完成提示
    Message.success(t('pages.user.batchImport.importCompleted'))

    // 2秒后隐藏进度条
    setTimeout(() => {
      importing.value = false
    }, 2000)

    // 清空文件列表 - 使用更兼容的方法
    if (uploadRef.value) {
      // 尝试使用clearFiles方法
      if (typeof uploadRef.value.clearFiles === 'function') {
        uploadRef.value.clearFiles();
      } 
      // 如果clearFiles不存在，尝试使用内部的fileList
      else if (uploadRef.value.fileList && uploadRef.value.fileList.length > 0) {
        uploadRef.value.fileList = [];
      }
    }

    // 调用成功回调
    options.onSuccess()
  } catch (error) {
    clearInterval(progressTimer)
    progressStatus.value = 'danger'
    Message.error(error.message || t('pages.user.batchImport.importFailed'))
    importing.value = false

    // 调用失败回调
    options.onError()
  }
}
</script>

<style scoped>
.container {
  padding: 20px;
}

.container-header h2 {
  margin: 0 0 20px 0;
  font-size: 20px;
  font-weight: 600;
}

.upload-area {
  padding: 40px;
  text-align: center;
}

.upload-text {
  margin-top: 10px;
  font-size: 14px;
  color: var(--color-text-2);
}

.upload-link {
  color: rgb(var(--primary-6));
  cursor: pointer;
}

.upload-tip {
  margin-top: 8px;
  font-size: 12px;
  color: var(--color-text-3);
}
</style>
