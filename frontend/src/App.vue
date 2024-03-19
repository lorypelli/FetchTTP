<script setup lang="ts">
import { ElNotification, ElTabPane, ElTabs } from 'element-plus';
import HTTP from './components/HTTP.vue';
import WS from './components/WS.vue';
import 'element-plus/dist/index.css';
import 'primevue/resources/themes/aura-light-green/theme.css';
import { ref } from 'vue';
import { handleTab as HTTPTab } from './components/HTTP.vue';
import { handleTab as WSTab } from './components/WS.vue';
defineOptions({
    name: 'App',
    components: {
        ElTabPane,
        ElTabs,
        HTTP,
        WS
    }
});
</script>

<script lang="ts">
const selectedTab = ref('HTTP');
export default {
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
                                }
                            },
                            {
                                label: 'WS',
                                onClick: () => {
                                    selectedTab.value = 'WS';
                                    WSTab(undefined, 'add');
                                }
                            }
                        ]
                    },
                    {
                        label: 'Open Devtools',
                        onClick: () => {
                            ElNotification({
                                title: 'Devtools',
                                message: 'To open devtools use CMD/CTRL + SHIFT + F12',
                                type: 'info',
                                position: 'bottom-right'
                            });
                        }
                    }
                ]
            });
        }
    }
};
</script>

<template>
    <ElTabs v-model="selectedTab" v-on:contextmenu="onContextMenu($event)">
        <ElTabPane label="HTTP" name="HTTP">
            <HTTP />
        </ElTabPane>
        <ElTabPane label="WS" name="WS">
            <WS />
        </ElTabPane>
    </ElTabs>
</template>