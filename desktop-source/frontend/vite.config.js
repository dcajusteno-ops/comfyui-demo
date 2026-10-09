import path from 'node:path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// 图片、视频以及 __variant__ 变体这些路径由 Wails 后端（main.go 的 AssetServer.Handler → ServeImage）负责。
// dev 模式下前端由 vite dev server 提供，它的 SPA fallback 会把找不到的路径统统返回 index.html。
// 而 Wails 的 dev 代理只在收到 404 时才会回退到后端 Handler，所以这里必须显式返回 404，
// 否则所有 <img> 拿到的都是 HTML，表现为整页图片全裂。
const MEDIA_EXTENSION = /\.(png|jpe?g|webp|gif|bmp|avif|svg|ico|mp4|webm|mov|mkv)$/i

const wailsBackendAssets = {
  name: 'wails-backend-assets',
  configureServer(server) {
    server.middlewares.use((req, res, next) => {
      const url = (req.url || '').split('?')[0]
      if (url.startsWith('/__variant__/') || MEDIA_EXTENSION.test(url)) {
        res.statusCode = 404
        res.end()
        return
      }
      next()
    })
  },
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    wailsBackendAssets,
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    sourcemap: true
  }
})
