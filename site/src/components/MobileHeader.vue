<template>
  <div class="mobile-header">
    <div v-if="showBack" class="header-left" @click="handleBack">
      <i class="iconfont icon-back"></i>
    </div>
    <div v-else class="header-left header-logo">
      <img
        v-if="config.siteLogo"
        :alt="config.siteTitle"
        :src="config.siteLogo"
        class="logo-img"
      />
      <img v-else :alt="config.siteTitle" src="~/assets/images/logo.png" class="logo-img" />
      <span class="logo-title">{{ config.siteTitle }}</span>
    </div>
    <div class="header-title">{{ title }}</div>
    <div class="header-right">
      <slot name="create-topic"></slot>
      <slot name="right"></slot>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  showBack: {
    type: Boolean,
    default: false
  }
})

const router = useRouter()
const configStore = useConfigStore()
const { config } = storeToRefs(configStore)

const handleBack = () => {
  router.back()
}
</script>

<style scoped>
.mobile-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 44px;
  padding: 0 16px;
  background: white;
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  width: 40px;
  display: flex;
  align-items: center;
  cursor: pointer;
}

.header-left .iconfont {
  font-size: 20px;
  color: #333;
}

.header-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.logo-img {
  height: 28px;
  width: auto;
  object-fit: contain;
}

.logo-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.header-title {
  font-size: 17px;
  font-weight: 600;
  color: #333;
  flex: 1;
  text-align: center;
}

.header-right {
  width: 40px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}
</style>
