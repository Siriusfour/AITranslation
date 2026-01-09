import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],

  server: {
    host:"0.0.0.0",
    proxy: {   //代理信息
      '/noteApi': {  //当你在前端请求 /kapi/xxx 时，Vite 会把请求 转发 到http://localhost:3008/xxx
        target: 'http://localhost:3008',
        rewrite: (path) => path.replace(/^\/noteApi/, ''),
        changeOrigin: true,
        secure: false,
        configure: (proxy, _options) => {
          proxy.on('proxyReq', (proxyReq, req) => {
            console.log('[VITE PROXY ->]',
                proxy.options.target?.toString(),
                proxyReq.path,
                '| from', req.url
            );
          });

          proxy.on('proxyRes', (proxyRes, req) => {
            console.log('[VITE PROXY RES]', proxyRes.statusCode, req.url);
          });
          proxy.on('error', (err, req) => {
            console.error('[VITE PROXY ERR]', req.url, err.message);
          });
        },
      }
    }
  },
})