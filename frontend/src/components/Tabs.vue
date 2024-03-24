<script setup lang="ts">
import { ElTabs, ElTabPane, TabPaneName } from 'element-plus';
import { ref } from 'vue';
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Tabs',
    components: {
        ElTabs,
        ElTabPane
    }
});
const props = defineProps<{
    type: 'http' | 'ws'
}>();
</script>

<script lang="ts">
let httpTabIndex = 1;
const httpSelectedTab = ref('1');
const httpTab = ref([
    {
        name: '1',
        select: 'GET',
        input: ''
    }
]);
let wsTabIndex = 1;
const wsSelectedTab = ref('1');
const wsTab = ref([
    {
        name: '1',
        input: '',
        connected: false
    }
]);
export function httpTabHandle(targetName: TabPaneName | undefined, action: 'add' | 'remove') {
    switch (action) {
    case 'add': {
        const newTabIndex = `${++httpTabIndex}`;
        httpTab.value.push({
            name: newTabIndex,
            select: 'GET',
            input: ''
        });
        httpSelectedTab.value = newTabIndex;
        localStorage.setItem('httpTab', JSON.stringify(httpTab.value));
        break;
    }
    case 'remove': {
        const t = httpTab.value;
        if (t.length > 1) {
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
        }
        localStorage.setItem('httpTab', JSON.stringify(httpTab.value));
        break;
    }
    }
}
export function wsTabHandle(targetName: TabPaneName | undefined, action: 'add' | 'remove') {
    switch (action) {
    case 'add': {
        const newTabIndex = `${++wsTabIndex}`;
        wsTab.value.push({
            name: newTabIndex,
            input: '',
            connected: false
        });
        wsSelectedTab.value = newTabIndex;
        localStorage.setItem('wsTab', JSON.stringify(wsTab.value));
        break;
    }
    case 'remove': {
        const t = wsTab.value;
        if (t.length > 1) {
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
        }
        localStorage.setItem('wsTab', JSON.stringify(wsTab.value));
        break;
    }
    }
}
export interface CompleteItem {
    name: string,
    select: string,
    input: string,
    connected: boolean
}
export default {
    mounted() {
        const http = localStorage.getItem('httpTab');
        if (http) {
            const httpParse = JSON.parse(http);
            httpTab.value = httpParse;
            httpSelectedTab.value = httpParse[0].name;
            httpTabIndex = parseInt(httpSelectedTab.value);
        }
        else {
            httpTab.value = [
                {
                    name: '1',
                    select: 'GET',
                    input: ''
                }
            ];
        }
        const ws = localStorage.getItem('wsTab');
        if (ws) {
            const wsParse = JSON.parse(ws);
            wsTab.value = wsParse;
            wsSelectedTab.value = wsParse[0].name;
            wsTabIndex = parseInt(wsSelectedTab.value);
        }
        else {
            wsTab.value = [
                {
                    name: '1',
                    input: '',
                    connected: false
                }
            ];
        }
    },
    computed: {
        selectedTab: {
            get() {
                return this.type == 'http' ? httpSelectedTab.value : wsSelectedTab.value;
            },
            set(value: string) {
                if (this.type == 'http') {
                    httpSelectedTab.value = value;
                }
                else {
                    wsSelectedTab.value = value;
                }
                return this.type == 'http' ? httpSelectedTab.value : wsSelectedTab.value;
            }
        },
        tabHandle: function () {
            return this.type == 'http' ? httpTabHandle : wsTabHandle;
        }
    },
};
</script>

<template>
    <ElTabs v-model="selectedTab" tab-position="left" editable v-on:edit="tabHandle" v-on:keydown.alt="(e: KeyboardEvent) => {
        const tab = props.type == 'http' ? httpTab.length : props.type == 'ws' ? wsTab.length : null
        if (tab) {
            switch (props.type) {
                case 'http': {
                    for (let i = 0; i < tab; i++) {
                        if (e.key == httpTab[i].name) {
                            httpSelectedTab = e.key
                        }
                    }
                    break
                }
                case 'ws': {
                    for (let i = 0; i < tab; i++) {
                        if (e.key == wsTab[i].name) {
                            wsSelectedTab = e.key
                        }
                    }
                    break
                }
            }
        }
    }">
        <ElTabPane :label="item.name"
            v-for="(item, index) in (props.type == 'http' ? httpTab : props.type == 'ws' ? wsTab : null)"
            :name="item.name" :key="index">
            <slot :item="item as CompleteItem"></slot>
        </ElTabPane>
    </ElTabs>
</template>