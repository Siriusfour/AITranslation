import path from 'path';

export default {
  base: './',
  build: {
    outDir: path.resolve(__dirname, 'dist'),
    assetsDir: 'assets',
    cssCodeSplit: false,
    rollupOptions: {
      input: {
        popup: './Src/popup.html',
      },
      output: {
        entryFileNames: 'js/[name].[hash].js',
        chunkFileNames: 'js/[name].[hash].js',
        assetFileNames: 'assets/[name].[hash][extname]'
      }
    },
  }
}
