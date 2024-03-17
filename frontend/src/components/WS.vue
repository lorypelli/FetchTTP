<script setup lang="ts">
import { ElButton, ElInput, ElNotification, ElTabs, ElTabPane, TabPaneName } from 'element-plus';
import 'element-plus/dist/index.css';
import { ref } from 'vue';
import Split from './Split.vue';
import { WS } from '../../wailsjs/go/main/App';
import { EventsEmit, EventsOff, EventsOn } from '../../wailsjs/runtime/runtime';
defineOptions({
    name: 'WS',
    components: {
        ElButton,
        ElInput,
        Split,
        ElTabs,
        ElTabPane
    }
});
let tabIndex = 1;
const selectedTab = ref('1');
const tabs = ref([
    {
        name: '1',
        input: '',
        connected: false
    }
]);
function handleTab(targetName: TabPaneName | undefined, action: 'add' | 'remove') {
    switch (action) {
    case 'add': {
        const newTabIndex = `${++tabIndex}`;
        tabs.value.push({
            name: newTabIndex,
            input: '',
            connected: false
        });
        selectedTab.value = newTabIndex;
        break;
    }
    case 'remove': {
        const t = tabs.value;
        if (t.length > 1) {
            let activeTab = selectedTab.value;
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
            selectedTab.value = activeTab;
            tabs.value = t.filter((tab) => tab.name != targetName);
        }
        break;
    }
    }
}
</script>

<script lang="ts">
interface Header {
    enabled: boolean,
    name: string,
    value: string
}
interface Query {
    enabled: boolean,
    name: string,
    value: string
}
interface Response {
    Status: string,
    Header: [],
    Message: string
}
let headers: Header[] = [
    {
        enabled: true,
        name: 'User-Agent',
        value: 'FetchTTP'
    }
];
let query: Query[] = [];
let message = '';
export default {
    data() {
        return {
            status: '',
            header: {},
            response: '',
        };
    },
    methods: {
        handleHeader(h: Header[]) {
            headers = h;
        },
        handleQuery(q: Query[]) {
            query = q;
        },
        handleMessage(m: string) {
            message = m;
        },
        update(res: Response) {
            this.status = res.Status;
            this.header = res.Header;
            this.response += res.Message + '\n';
        }
    }
};
</script>

<template>
    <ElTabs v-model="selectedTab" tab-position="left" editable v-on:edit="handleTab">
        <ElTabPane :label="item.name" v-for="(item, index) in tabs" :name="item.name" :key="index">
            <div class="flex p-1 space-x-1">
                <ElInput v-model="item.input" :disabled="item.connected" placeholder="echo.websocket.org"></ElInput>
                <ElButton v-model="item.connected" class="w-36" v-on:click="() => {
        item.connected = !item.connected
        if (item.connected) {
            response = ''
        }
        if (item.input) {
            if (!item.input.startsWith('ws://') && !item.input.startsWith('wss://')) {
                item.input = 'wss://' + item.input
            }
        }
        else {
            item.input = 'wss://echo.websocket.org'
        }
        try {
            EventsEmit('connected', item.connected)
            WS(item.input, headers, query, item.connected)
            if (item.connected) {
                EventsOn('websocket', (data: Response) => {
                    update(data)
                })
            }
            else {
                EventsOff('websocket')
            }
        }
        catch {
            item.connected = false
            ElNotification({
                title: 'Something went wrong!',
                type: 'error',
                position: 'bottom-right'
            })
        }
    }">{{ item.connected ? 'Disconnect' : 'Connect' }}</ElButton>
            </div>
            <Split :status="status" :header="header" :response="response" type='ws' v-on:headers="handleHeader"
                v-on:query="handleQuery" v-on:message="handleMessage" />
        </ElTabPane>
    </ElTabs>
</template>