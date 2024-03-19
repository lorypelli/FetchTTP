import { createApp } from 'vue';
import App from './App.vue';
import './style.css';
import 'element-plus/theme-chalk/dark/css-vars.css';
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css';
import ContextMenu from '@imengyu/vue3-context-menu';
const app = createApp(App);
app.use(ContextMenu);
app.mount('#app');