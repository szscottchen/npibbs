export default defineNuxtPlugin((nuxtApp) => {
  // 检测是否在企业微信环境
  const isWeCom = () => {
    const ua = navigator.userAgent.toLowerCase()
    return ua.includes('wxwork') || ua.includes('micromessenger')
  }
  
  // 检测是否为移动设备
  const isMobile = () => {
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
  }
  
  nuxtApp.provide('isWeCom', isWeCom)
  nuxtApp.provide('isMobile', isMobile)
  
  // 如果是企业微信移动端，重定向到移动端页面
  // 注意：跳过 /api、/admin、/wecom-、/install 等特殊路径
  if (process.client && isWeCom() && isMobile()) {
    const currentPath = window.location.pathname
    console.log('当前路径:', currentPath)

    // 跳过特殊路径
    const skipPaths = ['/api', '/admin', '/install', '/wecom-', '/_nuxt', '/static']
    const shouldSkip = skipPaths.some(prefix => currentPath.startsWith(prefix))

    if (!shouldSkip && !currentPath.startsWith('/mobile')) {
      console.log('重定向到移动端:', currentPath)
      const mobilePath = currentPath === '/' ? '/mobile' : `/mobile${currentPath}`
      window.location.replace(mobilePath)
    } else {
      console.log('跳过重定向:', currentPath)
    }
  }
})