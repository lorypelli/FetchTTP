import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import { multip } from 'vite-plugin-multip';
export default defineConfig({
    plugins: [vue(), multip()],
    optimizeDeps: {
        include: ['vue'],
    },
    server: {
        host: '127.0.0.1',
    },
});
