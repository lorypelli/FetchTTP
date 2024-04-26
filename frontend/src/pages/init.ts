import type { App } from 'vue';
import 'element-plus/theme-chalk/dark/css-vars.css';
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css';
import ContextMenu from '@imengyu/vue3-context-menu';

export function init(app: App) {
    app.use(ContextMenu)
}