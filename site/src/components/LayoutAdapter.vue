<template>
  <component :is="layoutComponent">
    <slot />
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useWeCom } from '~/composables/useWeCom'

const { needMobileLayout } = useWeCom()

// 根据环境选择不同的布局组件
const layoutComponent = computed(() => {
  return needMobileLayout.value ? 'MobileLayout' : 'DesktopLayout'
})
</script>
