import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import path from 'path'
import { fileURLToPath } from 'url'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'

// 获取 __dirname (ESM 模块需要)
const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    uni(),
    // 自动导入 Vue 相关函数
    AutoImport({
      imports: [
        'vue',
        'pinia',
        // uni-app 相关不需要自动导入，因为是全局的
      ],
      dts: 'src/auto-imports.d.ts', // 生成类型声明文件
      eslintrc: {
        enabled: false,
      },
    }),
    // 自动导入组件
    Components({
      dts: 'src/components.d.ts',
      // uni-app 组件不需要配置，会自动识别
    }),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
})
