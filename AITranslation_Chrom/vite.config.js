import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { crx } from '@crxjs/vite-plugin'
// import manifest from './manifest.json' // Node 14 & 16
import manifest from './manifest.json' assert { type: 'json' } // Node >=17
import { resolve } from 'path';

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

  build: {
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'index.html'), 
        AITranslation: resolve(__dirname, 'src/scripts/AITranslation.js'),
      },
      output: {
        entryFileNames: (chunkInfo) => {
          if (chunkInfo.name === 'AITranslation') {
            return 'assets/AITranslation.js'; 
          }
          // 对于其他入口，您可能希望保留哈希或者有不同的命名规则
          return 'assets/[name]-[hash].js';
        },
        // chunkFileNames: 'assets/[name]-[hash].js',
        // assetFileNames: 'assets/[name]-[hash].[ext]',
      }
    },
   outDir: 'dist', // 确保输出目录是正确的
    emptyOutDir: true, // 每次构建前清空输出目录
  }
})