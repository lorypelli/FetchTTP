<script setup lang="ts">
import { ElMessageBox, ElNotification, ElTabPane, ElTabs } from 'element-plus';
import HTTP from '../components/HTTP.vue';
import WS from '../components/WS.vue';
import CURL from '../components/CURL.vue';
import {
    httpTabHandle as HTTPTab,
    wsTabHandle as WSTab,
} from '../components/Tabs.vue';
import Updater from '../components/Updater.vue';
import 'element-plus/dist/index.css';
import 'primevue/resources/themes/aura-light-green/theme.css';
import { h, ref } from 'vue';
import { CheckUpdates, Update as U } from '../../wailsjs/go/main/App.js';
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
interface Update {
    IsLatest: boolean;
    Version: string;
    Description: string;
    Error: string;
}
const selectedTab = ref('HTTP');
export default {
    mounted() {
        try {
            CheckUpdates().then((res: Update) => {
                if (res.Error) {
                    ElNotification({
                        title: 'Error while checking for updates!',
                        message: res.Error,
                        type: 'error',
                        position: 'bottom-right',
                    });
                    return;
                }
                if (!res.IsLatest) {
                    ElMessageBox({
                        title: 'A new version is avaible!\nDo you want to update?',
                        message: h(Updater, {
                            version: res.Version,
                            description: res.Description,
                        }),
                        confirmButtonText: 'Yes',
                        showCancelButton: true,
                        cancelButtonText: 'No',
                        showClose: false,
                        closeOnClickModal: false,
                        closeOnHashChange: false,
                        closeOnPressEscape: false,
                        center: true,
                    })
                        .then(() => {
                            ElMessageBox({
                                title: 'Warning!',
                                message:
                                    'The app will now exit and will be re-opened automatically',
                                type: 'warning',
                                showClose: false,
                                closeOnClickModal: false,
                                closeOnHashChange: false,
                                closeOnPressEscape: false,
                                center: true,
                            }).then(() => {
                                U();
                            });
                        })
                        .catch(() => { });
                }
            });
        } catch {
            ElNotification({
                title: 'Error while checking for updates!',
                type: 'error',
                position: 'bottom-right',
            });
        }
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
                        ],
                    },
                    {
                        label: 'Check for Updates',
                        onClick: () => {
                            try {
                                CheckUpdates().then((res: Update) => {
                                    if (res.Error) {
                                        ElNotification({
                                            title: 'Error while checking for updates!',
                                            message: res.Error,
                                            type: 'error',
                                            position: 'bottom-right',
                                        });
                                        return;
                                    }
                                    if (res.IsLatest) {
                                        ElMessageBox({
                                            message:
                                                'You are using latest version',
                                        }).catch(() => { });
                                    } else {
                                        ElMessageBox({
                                            title: 'A new version is avaible!\nDo you want to update?',
                                            message: h(Updater, {
                                                version: res.Version,
                                                description: res.Description,
                                            }),
                                            confirmButtonText: 'Yes',
                                            showCancelButton: true,
                                            cancelButtonText: 'No',
                                            showClose: false,
                                            closeOnClickModal: false,
                                            closeOnHashChange: false,
                                            closeOnPressEscape: false,
                                            center: true,
                                        })
                                            .then(() => {
                                                ElMessageBox({
                                                    title: 'Warning!',
                                                    message:
                                                        'The app will now exit and will be re-opened automatically',
                                                    type: 'warning',
                                                    showClose: false,
                                                    closeOnClickModal: false,
                                                    closeOnHashChange: false,
                                                    closeOnPressEscape: false,
                                                    center: true,
                                                }).then(() => {
                                                    U();
                                                });
                                            })
                                            .catch(() => { });
                                    }
                                });
                            } catch {
                                ElNotification({
                                    title: 'Error while checking for updates!',
                                    type: 'error',
                                    position: 'bottom-right',
                                });
                            }
                        },
                    },
                    {
                        label: 'Open Devtools',
                        onClick: () => {
                            ElNotification({
                                title: 'Devtools',
                                message:
                                    'To open devtools use CMD/CTRL + SHIFT + F12',
                                type: 'info',
                                position: 'bottom-right',
                            });
                        },
                    },
                ],
            });
        },
    },
};
</script>

<template>
    <ElTabs v-model="selectedTab" @contextmenu="onContextMenu($event)" @keydown.alt="(e: KeyboardEvent) => {
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
        ">
        <ElTabPane label="HTTP" name="HTTP">
            <HTTP />
        </ElTabPane>
        <ElTabPane label="WS" name="WS">
            <WS />
        </ElTabPane>
        <ElTabPane class="flex justify-center items-center h-max place-items-center scale-150" label="CURL" name="CURL">
            <CURL />
        </ElTabPane>
    </ElTabs>
</template>
