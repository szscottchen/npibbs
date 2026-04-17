import { mergeConfig } from 'vite';
import baseConfig from './vite.config.base';

export default mergeConfig(
  {
    mode: 'development',
    base: '/',
    server: {
      port: 8080,
      open: true,
      //开发环境允许所有主机访问
      allowedHosts: true,
      fs: {
        strict: true,
      },
      proxy: {
        '/api': {
          target: 'http://localhost:8082',
          changeOrigin: true,
          rewrite: (path: string) => path
        }
      }
    },
    plugins: [],
  },
  baseConfig
);
