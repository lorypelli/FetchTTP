import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
export default defineConfig({
    plugins: [vue()],
    optimizeDeps: {
        include: ['vue'],
    },
    server: {
        host: '127.0.0.1',
    },
});
