<script setup lang="ts">
import { ElButton, ElInput } from 'element-plus';
import 'element-plus/dist/index.css';
import { ref } from 'vue';
import Split from './Split.vue';
import { WS } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';
defineOptions({
    name: 'WS',
    components: {
        ElButton,
        ElInput,
        Split
    }
});
const input = ref('');
let connected = false;
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
    <div class="flex p-1 space-x-1">
        <ElInput v-model="input" :disabled="connected" placeholder="echo.websocket.org"></ElInput>
        <ElButton v-model="connected" class="w-36" v-on:click="() => {
            if (connected) {
                connected = false
            }
            else {
                connected = true
            }
            if (input) {
                if (!input.startsWith('ws://') && !input.startsWith('wss://')) {
                    input = 'wss://' + input
                }
            }
            else {
                input = 'wss://echo.websocket.org'
            }
            WS(input, headers, query, connected)
            if (connected) {
                EventsOn('websocket', (data: Response) => {
                    update(data)
                })
            }
        }">{{ connected ? 'Disconnect' : 'Connect' }}</ElButton>
    </div>
    <Split :status="status" :header="header" :response="response" type='ws' v-on:headers="handleHeader" v-on:query="handleQuery" v-on:message="handleMessage" />
</template>