<template>
    <ElTabs
        v-model="selectedTab"
        tab-position="left"
        editable
        @edit="tabHandle"
        @keyup="
            () => {
                key = '';
            }
        "
        @keydown.alt="keyHandle"
    >
        <ElTabPane
            v-for="(item, index) in props.type == 'http' ? httpTab : wsTab"
            :index="index"
            :label="item.name"
            :name="item.name"
        >
            <slot
                :item="item as CompleteItem"
                :index="parseInt(item.name) - 1"
            />
        </ElTabPane>
    </ElTabs>
</template>

<script setup lang="ts">
import { ElTabPane, ElTabs, TabPaneName } from 'element-plus';
import { computed, ref } from 'vue';
import type { CompleteItem } from '../types';
const props = defineProps<{
    type: 'http' | 'ws';
}>();
let httpTabIndex = 1;
const httpSelectedTab = ref('1');
const httpTab = ref([
    {
        name: '1',
        select: 'GET',
        input: '',
    },
]);
let wsTabIndex = 1;
const wsSelectedTab = ref('1');
const wsTab = ref([
    {
        name: '1',
        input: '',
        connected: false,
    },
]);
let key = '';
function httpTabHandle(
    targetName: TabPaneName | undefined,
    action: 'add' | 'remove',
) {
    switch (action) {
        case 'add': {
            let newTabIndex = ++httpTabIndex;
            const http = localStorage.getItem('httpTab');
            const tabs = http ? JSON.parse(http) : [];
            tabs.forEach((t: { name: string }) => {
                if (t.name == newTabIndex.toString()) {
                    newTabIndex++;
                }
            });
            httpTab.value.push({
                name: newTabIndex.toString(),
                select: 'GET',
                input: '',
            });
            httpSelectedTab.value = newTabIndex.toString();
            localStorage.setItem('httpTab', JSON.stringify(httpTab.value));
            break;
        }
        case 'remove': {
            const t = httpTab.value;
            let activeTab = httpSelectedTab.value;
            if (activeTab == targetName) {
                t.forEach((tab, index) => {
                    if (tab.name == targetName) {
                        const nextTab = t[index + 1] || t[index - 1];
                        if (nextTab) {
                            activeTab = nextTab.name;
                        }
                    }
                });
            }
            httpSelectedTab.value = activeTab;
            httpTab.value = t.filter((tab) => tab.name != targetName);
            localStorage.setItem('httpTab', JSON.stringify(httpTab.value));
            if (t.length == 1) {
                localStorage.clear();
                httpTabIndex = 0;
                httpTabHandle(undefined, 'add');
            }
            break;
        }
    }
}
function wsTabHandle(
    targetName: TabPaneName | undefined,
    action: 'add' | 'remove',
) {
    switch (action) {
        case 'add': {
            let newTabIndex = ++wsTabIndex;
            const ws = localStorage.getItem('wsTab');
            const tabs = ws ? JSON.parse(ws) : [];
            tabs.forEach((t: { name: string }) => {
                if (t.name == newTabIndex.toString()) {
                    newTabIndex++;
                }
            });
            wsTab.value.push({
                name: newTabIndex.toString(),
                input: '',
                connected: false,
            });
            wsSelectedTab.value = newTabIndex.toString();
            localStorage.setItem('wsTab', JSON.stringify(wsTab.value));
            break;
        }
        case 'remove': {
            const t = wsTab.value;
            let activeTab = wsSelectedTab.value;
            if (activeTab == targetName) {
                t.forEach((tab, index) => {
                    if (tab.name == targetName) {
                        const nextTab = t[index + 1] || t[index - 1];
                        if (nextTab) {
                            activeTab = nextTab.name;
                        }
                    }
                });
            }
            wsSelectedTab.value = activeTab;
            wsTab.value = t.filter((tab) => tab.name != targetName);
            localStorage.setItem('wsTab', JSON.stringify(wsTab.value));
            if (t.length == 1) {
                localStorage.clear();
                wsTabIndex = 0;
                wsTabHandle(undefined, 'add');
            }
            break;
        }
    }
}
const selectedTab = computed({
    get() {
        return props.type == 'http'
            ? httpSelectedTab.value
            : wsSelectedTab.value;
    },
    set(value: string) {
        if (props.type == 'http') {
            httpSelectedTab.value = value;
        } else {
            wsSelectedTab.value = value;
        }
        return props.type == 'http'
            ? httpSelectedTab.value
            : wsSelectedTab.value;
    },
});
const tabHandle = computed(() =>
    props.type == 'http' ? httpTabHandle : wsTabHandle,
);
function keyHandle(e: KeyboardEvent) {
    setTimeout(() => {
        if (/^[0-9]$/.test(e.key)) {
            key += e.key;
        }
    }, 100);
    const tab =
        props.type == 'http' ? httpTab.value.length : wsTab.value.length;
    switch (props.type) {
        case 'http': {
            for (let i = 0; i < tab; i++) {
                if (key == httpTab.value[i].name) {
                    httpSelectedTab.value = key;
                }
            }
            break;
        }
        case 'ws': {
            for (let i = 0; i < tab; i++) {
                if (key == wsTab.value[i].name) {
                    wsSelectedTab.value = key;
                }
            }
            break;
        }
    }
}
</script>
