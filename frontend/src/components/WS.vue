<script setup lang="ts">
import { ElButton, ElInput, ElNotification, ElTabs, ElTabPane } from 'element-plus';
import { WS } from '../../wailsjs/go/main/App';
import Split from './Split.vue';
import { EventsEmit, EventsOff, EventsOn } from '../../wailsjs/runtime/runtime';
import Tabs from './Tabs.vue';
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
        update(res: Response) {
            this.status = res.Status;
            this.header = res.Header;
            this.response += res.Message + '\n';
        }
    }
};
</script>

<template>
    <Tabs type="ws">
        <template #default="{ item }">
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
                v-on:query="handleQuery" />
        </template>
    </Tabs>
</template>