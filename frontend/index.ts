import ContextMenu from '@imengyu/vue3-context-menu';
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css';
import 'element-plus/theme-chalk/dark/css-vars.css';
import { createApp } from 'vue';
import App from './src/pages/index.vue';
import './src/styles/style.css';

const app = createApp(App);
app.use(ContextMenu);
app.mount('#app');
