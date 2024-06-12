import ContextMenu from '@imengyu/vue3-context-menu';
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css';
import 'element-plus/theme-chalk/dark/css-vars.css';
import type { App } from 'vue';

export function init(app: App) {
    app.use(ContextMenu);
}
