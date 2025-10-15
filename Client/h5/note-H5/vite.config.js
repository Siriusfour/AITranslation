import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],

  server: {
    host:"0.0.0.0",
    proxy: {
      '/kapi': {
        target: 'http://localhost:3008',
        rewrite: (path) => path.replace(/^\/kapi/, ''),
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
