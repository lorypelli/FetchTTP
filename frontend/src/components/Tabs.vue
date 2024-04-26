<script setup lang="ts">
import { ElTabs, ElTabPane, TabPaneName } from 'element-plus';
import { ref } from 'vue';
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Tabs',
    components: {
        ElTabs,
        ElTabPane,
    },
});
const props = defineProps<{
    type: 'http' | 'ws';
}>();
</script>

<script lang="ts">
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
export function httpTabHandle(
    targetName: TabPaneName | undefined,
    action: 'add' | 'remove',
) {
    switch (action) {
        case 'add': {
            const newTabIndex = `${++httpTabIndex}`;
            httpTab.value.push({
                name: newTabIndex,
                select: 'GET',
                input: '',
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
export function wsTabHandle(
    targetName: TabPaneName | undefined,
    action: 'add' | 'remove',
) {
    switch (action) {
        case 'add': {
            const newTabIndex = `${++wsTabIndex}`;
            wsTab.value.push({
                name: newTabIndex,
                input: '',
                connected: false,
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
    name: string;
    select: string;
    input: string;
    connected: boolean;
}
export default {
    computed: {
        selectedTab: {
            get() {
                return this.type == 'http'
                    ? httpSelectedTab.value
                    : wsSelectedTab.value;
            },
            set(value: string) {
                if (this.type == 'http') {
                    httpSelectedTab.value = value;
                } else {
                    wsSelectedTab.value = value;
                }
                return this.type == 'http'
                    ? httpSelectedTab.value
                    : wsSelectedTab.value;
            },
        },
        tabHandle: function () {
            return this.type == 'http' ? httpTabHandle : wsTabHandle;
        },
    },
    mounted() {
        let http = localStorage.getItem('httpTab');
        let httpLen = 1;
        if (http) {
            let httpJson = JSON.parse(http);
            httpLen = httpJson.length;
            httpTab.value = httpJson;
            httpTabIndex = httpJson[0].name;
            httpSelectedTab.value = httpJson[0].name.toString();
        }
        for (let i = 0; i < httpLen; i++) {
            const select = localStorage.getItem(
                `${httpTab.value[i].name}-select`,
            );
            if (select) {
                httpTab.value[i].select = select;
            }
            const httpInput = localStorage.getItem(
                `${httpTab.value[i].name}-input-http`,
            );
            if (httpInput) {
                httpTab.value[i].input = httpInput;
            }
        }
        let ws = localStorage.getItem('wsTab');
        let wsLen = 1;
        if (ws) {
            let wsJson = JSON.parse(ws);
            wsLen = wsJson.length;
            wsTab.value = wsJson;
            wsTabIndex = wsJson[0].name;
            wsSelectedTab.value = wsJson[0].name.toString();
        }
        for (let i = 0; i < wsLen; i++) {
            const wsInput = localStorage.getItem(
                `${wsTab.value[i].name}-input-ws`,
            );
            if (wsInput) {
                wsTab.value[i].input = wsInput;
            }
        }
    },
    methods: {
        keyHandle(e: KeyboardEvent) {
            setTimeout(() => {
                if (
                    ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'].includes(
                        e.key,
                    )
                ) {
                    key += e.key;
                }
            }, 100);
            const tab =
                this.type == 'http'
                    ? httpTab.value.length
                    : this.type == 'ws'
                        ? wsTab.value.length
                        : null;
            if (tab) {
                switch (this.type) {
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
        },
    },
};
</script>

<template>
    <ElTabs v-model="selectedTab" tab-position="left" editable @edit="tabHandle" @keyup="() => {
            key = '';
        }
        " @keydown.alt="keyHandle">
        <ElTabPane v-for="(item, index) in props.type == 'http'
            ? httpTab
            : props.type == 'ws'
                ? wsTab
                : null" :key="index" :label="item.name" :name="item.name">
            <slot :item="item as CompleteItem" />
        </ElTabPane>
    </ElTabs>
</template>
