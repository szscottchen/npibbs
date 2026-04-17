/**
 * 企业微信相关 Composable
 */
export const useWeCom = () => {
  // 确保只在客户端访问 navigator
  const getUserAgent = () => {
    if (process.client && typeof navigator !== 'undefined') {
      return navigator.userAgent.toLowerCase()
    }
    return ''
  }

  /**
   * 检测是否在企业微信环境
   */
  const isWeCom = computed(() => {
    const ua = getUserAgent()
    return ua.includes('wxwork') || ua.includes('micromessenger')
  })

  /**
   * 检测是否是移动端
   */
  const isMobile = computed(() => {
    const ua = process.client && typeof navigator !== 'undefined' ? navigator.userAgent : ''
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(ua)
  })

  /**
   * 检测是否需要移动版布局
   * 注意：企业微信电脑端不应使用移动端布局
   */
  const needMobileLayout = computed(() => {
    // 只有在真正的移动设备（非企业微信电脑端）才使用移动端布局
    return isMobile.value
  })

  return {
    isWeCom,
    isMobile,
    needMobileLayout
  }
}
