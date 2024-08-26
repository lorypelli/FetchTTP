<script setup lang="ts">
import { ElNotification, ElTabPane, ElTabs } from 'element-plus';
import 'element-plus/dist/index.css';
import 'primevue/resources/themes/aura-light-green/theme.css';
import { ref } from 'vue';
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
</script>

<script lang="ts">
const selectedTab = ref('HTTP');
export default {
    mounted() {
        import.meta.env.PROD && checkUpdates('load');
    },
    methods: {
        onContextMenu(e: MouseEvent) {
            e.preventDefault();
            this.$contextmenu({
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
                        label: 'Duplicate',
                        onClick: () => {
                            selectedTab.value == 'HTTP'
                                ? HTTPTab(undefined, 'add')
                                : WSTab(undefined, 'add');
                            localStorage.getItem(`${selectedTab}`);
                        },
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
                                message:
                                    'To open devtools use CMD/CTRL + SHIFT + F12',
                                type: 'info',
                                position: 'bottom-right',
                            }),
                    },
                ],
            });
        },
    },
};
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
