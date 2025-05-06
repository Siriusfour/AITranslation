import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { crx } from '@crxjs/vite-plugin'
// import manifest from './manifest.json' // Node 14 & 16
import manifest from './manifest.json' assert { type: 'json' } // Node >=17

export default defineConfig({
  plugins: [
    vue(),
    crx({ manifest }),
  ],

  server: {
    hmr: {
      host: 'localhost', // 确保 HMR 服务器地址可从插件环境访问
      port: 5173,
      protocol: 'ws' // 明确指定协议为 ws
    },
    port: 5173, // 确保端口号与错误信息中的一致
    strictPort: true, // 建议添加，确保端口不被占用而改变
    cors: true 

  },
})