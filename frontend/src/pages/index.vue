<script setup lang="ts">
import ContextMenu from '@imengyu/vue3-context-menu';
import { ElNotification, ElTabPane, ElTabs } from 'element-plus';
import 'element-plus/dist/index.css';
import 'primevue/resources/themes/aura-light-green/theme.css';
import { onMounted, ref } from 'vue';
import CURL from '../components/CURL.vue';
import HTTP from '../components/HTTP.vue';
import {
    httpTabHandle as HTTPTab,
    wsTabHandle as WSTab,
} from '../components/Tabs.vue';
import Updater from '../components/Updater.vue';
import WS from '../components/WS.vue';
import { checkUpdates } from '../functions/checkUpdates';
defineOptions({
    name: 'App',
    components: {
        ElTabPane,
        ElTabs,
        HTTP,
        WS,
        CURL,
        Updater,
    },
});
const selectedTab = ref('HTTP');
onMounted(() => {
    import.meta.env.PROD && checkUpdates('load');
});
function onContextMenu(e: MouseEvent) {
    e.preventDefault();
    ContextMenu.showContextMenu({
        x: e.x,
        y: e.y,
        theme: 'dark',
        items: [
            {
                label: 'New',
                children: [
                    {
                        label: 'HTTP',
                        onClick: () => {
                            selectedTab.value = 'HTTP';
                            HTTPTab(undefined, 'add');
                        },
                    },
                    {
                        label: 'WS',
                        onClick: () => {
                            selectedTab.value = 'WS';
                            WSTab(undefined, 'add');
                        },
                    },
                    {
                        label: 'CURL',
                        onClick: () => (selectedTab.value = 'CURL'),
                    },
                ],
            },
            {
                label: 'Check for Updates',
                onClick: () => checkUpdates('menu'),
            },
            {
                label: 'Open Devtools',
                onClick: () =>
                    ElNotification({
                        title: 'Devtools',
                        message: 'To open devtools use CMD/CTRL + SHIFT + F12',
                        type: 'info',
                        position: 'bottom-right',
                    }),
            },
        ],
    });
}
</script>

<template>
    <ElTabs
        v-model="selectedTab"
        @contextmenu="onContextMenu($event)"
        @keydown.alt="
            (e: KeyboardEvent) => {
                switch (e.key.toUpperCase()) {
                    case 'H': {
                        selectedTab = 'HTTP';
                        break;
                    }
                    case 'W': {
                        selectedTab = 'WS';
                        break;
                    }
                }
            }
        "
    >
        <ElTabPane label="HTTP" name="HTTP">
            <HTTP />
        </ElTabPane>
        <ElTabPane label="WS" name="WS">
            <WS />
        </ElTabPane>
        <ElTabPane
            class="flex justify-center items-center h-max place-items-center scale-150"
            label="CURL"
            name="CURL"
        >
            <CURL />
        </ElTabPane>
    </ElTabs>
</template>
