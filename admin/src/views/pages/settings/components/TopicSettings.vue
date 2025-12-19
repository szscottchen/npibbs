<template>
  <div class="topic-settings">
    <a-card :title="$t('pages.settings.topic.valueTypes.title')">
      <a-form :model="form" :label-col-props="{ span: 6 }" :wrapper-col-props="{ span: 18 }">
        <a-form-item
          v-for="(item, index) in form.valueTypes"
          :key="index"
          :label="$t('pages.settings.topic.valueTypes.item') + (index + 1)"
        >
          <a-row :gutter="16">
            <a-col :span="10">
              <a-input
                v-model="item.label"
                :placeholder="$t('pages.settings.topic.valueTypes.placeholder.label')"
              />
            </a-col>
            <a-col :span="10">
              <a-input-number
                v-model="item.score"
                :min="1"
                :max="100"
                :placeholder="$t('pages.settings.topic.valueTypes.placeholder.score')"
                style="width: 100%"
              />
            </a-col>
            <a-col :span="4">
              <a-button
                type="outline"
                status="danger"
                @click="removeValueType(index)"
              >
                <icon-delete />
              </a-button>
            </a-col>
          </a-row>
        </a-form-item>
        
        <a-form-item>
          <a-button type="outline" @click="addValueType">
            <icon-plus />
            {{ $t('pages.settings.topic.valueTypes.addType') }}
          </a-button>
        </a-form-item>
        
        <a-form-item>
          <a-button type="primary" @click="submit" :loading="loading">
            {{ $t('pages.settings.topic.submit') }}
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, getCurrentInstance } from 'vue';
import axios from 'axios';

const { proxy } = getCurrentInstance() || {};
const loading = ref(false);

const form = reactive({
  valueTypes: []
});

// 获取配置
const loadConfig = async () => {
  try {
    const res = await axios.get('/api/admin/sys-config/configs');
    console.log('获取配置响应:', res);
    if (res && res.valueTypes) {
      form.valueTypes = res.valueTypes;
      console.log('加载到价值类型配置:', form.valueTypes);
    } else {
      console.log('未获取到价值类型配置，使用默认值');
      // 默认配置
      form.valueTypes = [
        { label: '值得讨论', score: 2 },
        { label: '深有启发', score: 3 },
        { label: '可以采纳', score: 4 },
        { label: '价值超高', score: 5 }
      ];
    }
  } catch (error) {
    console.error('加载配置失败:', error);
    proxy.$message.error('加载配置失败: ' + (error.response?.data?.message || error.message || '未知错误'));
    // 出错时也使用默认配置
    form.valueTypes = [
      { label: '值得讨论', score: 2 },
      { label: '深有启发', score: 3 },
      { label: '可以采纳', score: 4 },
      { label: '价值超高', score: 5 }
    ];
  }
};

// 添加价值类型
const addValueType = () => {
  form.valueTypes.push({ label: '', score: 1 });
};

// 删除价值类型
const removeValueType = (index) => {
  form.valueTypes.splice(index, 1);
};

// 提交配置
const submit = async () => {
  // 验证数据
  for (let i = 0; i < form.valueTypes.length; i++) {
    const item = form.valueTypes[i];
    if (!item.label || item.label.trim() === '') {
      proxy.$message.error(`第${i + 1}个价值类型标签不能为空`);
      return;
    }
    if (!item.score || item.score <= 0) {
      proxy.$message.error(`第${i + 1}个价值类型积分必须大于0`);
      return;
    }
  }

  loading.value = true;
  try {
    await axios.post('/api/admin/sys-config/save', {
      valueTypes: JSON.stringify(form.valueTypes)
    });
    proxy.$message.success('保存成功');
  } catch (error) {
    console.error('保存失败:', error);
    proxy.$message.error('保存失败: ' + (error.response?.data?.message || error.message || '未知错误'));
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadConfig();
});
</script>

<style scoped lang="less">
.topic-settings {
  padding: 20px;
  min-height: 400px;
}
</style>