<template>
  <div class="container">
    <div class="container-header">
      <h2>{{ isCreate ? '新增用户' : '编辑用户' }}</h2>
    </div>
    
    <div class="container-main">
      <a-card>
        <a-form ref="formRef" :model="form" :rules="rules" :label-col="{ span: 4 }" :wrapper-col="{ span: 14 }">
          <a-form-item label="类型" field="type">
            <a-select v-model="form.type" placeholder="用户类型">
              <a-option :value="0" label="用户" />
              <a-option :value="1" label="员工" />
            </a-select>
          </a-form-item>

          <a-form-item label="昵称" field="nickname">
            <a-input v-model="form.nickname" />
          </a-form-item>

          <a-form-item label="头像" field="avatar">
            <image-upload v-model="form.avatar" />
          </a-form-item>

          <a-form-item label="用户名" field="username">
            <a-input v-model="form.username" />
          </a-form-item>

          <a-form-item label="邮箱" field="email">
            <a-input v-model="form.email" />
          </a-form-item>

          <a-form-item label="性别" field="gender">
            <a-select v-model="form.gender">
              <a-option value="Male" label="男" />
              <a-option value="Female" label="女" />
            </a-select>
          </a-form-item>

          <a-form-item label="主页" field="homePage">
            <a-input v-model="form.homePage" />
          </a-form-item>

          <a-form-item label="描述" field="description">
            <a-input v-model="form.description" />
          </a-form-item>

          <a-form-item label="角色" field="roles">
            <a-select v-model="form.roleIds" multiple placeholder="请选择角色">
              <a-option
                v-for="role in roles"
                :key="role.id"
                :value="role.id"
                :label="role.name"
              />
            </a-select>
          </a-form-item>

          <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
            <a-space>
              <a-button type="primary" @click="handleSubmit">
                提交
              </a-button>
              <a-button @click="handleCancel">
                取消
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import ImageUpload from '@/components/ImageUpload.vue';

const router = useRouter();
const route = useRoute();

const formRef = ref();
interface Role {
  id: number;
  name: string;
}

const roles = ref<Role[]>([]);
const isCreate = ref(!route.params.id);

const form = reactive({
  type: 0,
  username: undefined,
  email: undefined,
  nickname: undefined,
  avatar: undefined,
  gender: undefined,
  homePage: undefined,
  description: undefined,
  roleIds: undefined,
});

const rules = {
  type: [{ required: true, message: '请选择用户类型' }],
  nickname: [{ required: true, message: '请输入用户昵称' }],
};

const loadRoles = async () => {
  try {
    roles.value = await axios.get('/api/admin/role/roles');
  } catch (e: any) {
    useHandleError(e);
  }
};

const loadUser = async (id: string) => {
  try {
    const userData = await axios.get(`/api/admin/user/${id}`);
    Object.assign(form, userData);
  } catch (e: any) {
    useHandleError(e);
  }
};

const handleSubmit = async () => {
  const validateErr = await formRef.value.validate();
  if (validateErr) {
    return;
  }
  
  try {
    const url = isCreate.value
      ? '/api/admin/user/create'
      : '/api/admin/user/update';
    await axios.postForm<any>(url, jsonToFormData(form));
    Message.success('提交成功');
    router.push('/user');
  } catch (e: any) {
    useHandleError(e);
  }
};

const handleCancel = () => {
  router.push('/user');
};

onMounted(async () => {
  await loadRoles();
  if (!isCreate.value && route.params.id) {
    await loadUser(route.params.id as string);
  }
});
</script>

<style lang="less" scoped>
.container {
  padding: 0 20px 20px 20px;
}

.container-header {
  margin-bottom: 20px;
  
  h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 500;
  }
}

.container-main {
  background: var(--color-bg-1);
  border-radius: 4px;
  padding: 20px;
}
</style>