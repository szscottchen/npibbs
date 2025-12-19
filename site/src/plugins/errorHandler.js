export default defineNuxtPlugin((nuxtApp) => {
  // 捕获Vue组件渲染和观察期间的错误
  nuxtApp.vueApp.config.errorHandler = (err, instance, info) => {
    console.error('Vue Error Handler:', err, info);
    // 在生产环境中，您可能想要将错误上报到错误跟踪服务
    // 例如：reportErrorToService(err, info)
  };
  
  // 捕获未处理的Promise拒绝
  if (import.meta.client) {
    window.addEventListener('unhandledrejection', event => {
      console.error('Unhandled promise rejection:', event.reason);
      // 可以阻止默认的错误处理
      // event.preventDefault();
    });
  }
});