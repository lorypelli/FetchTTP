<template>
    <Tabs type="ws">
        <template #default="{ item, index }">
            <div class="flex p-1 space-x-1">
                <ElInput
                    v-model="item.input"
                    :disabled="item.connected"
                    placeholder="echo.websocket.org"
                    @input="handleInput(item, index)"
                    @keydown.enter="sendWebsocket(item)"
                />
                <ElButton
                    v-model="item.connected"
                    class="w-36"
                    @click="sendWebsocket(item)"
                >
                    {{ item.connected ? 'Disconnect' : 'Connect' }}
                </ElButton>
            </div>
            <Split
                :name="item.name"
                :index="index"
                :status="status"
                :headers="responseHeaders"
                :response="response"
                type="ws"
                @headers="handleHeader"
                @query="handleQuery"
            />
        </template>
    </Tabs>
</template>

<script setup lang="ts">
import { ElButton, ElInput, ElNotification } from 'element-plus';
import { reactive, ref } from 'vue';
import { WS } from '../../wailsjs/go/main/App';
import { EventsEmit, EventsOff, EventsOn } from '../../wailsjs/runtime/runtime';
import { wsTabItem } from '../functions/useStorage';
import type {
    CompleteItem,
    GenericHeader,
    Header,
    Query,
    Response,
} from '../types';
import Split from './Split.vue';
import Tabs from './Tabs.vue';
let headers: Header[] = reactive([
    {
        enabled: true,
        name: 'User-Agent',
        value: 'FetchTTP',
    },
]);
let status = ref('');
let responseHeaders: GenericHeader[] = reactive([]);
let response = ref('');
let query: Query[] = reactive([]);
function handleInput(item: CompleteItem, key: number) {
    if (!wsTabItem.value[key]) {
        wsTabItem.value[key] = {};
    }
    wsTabItem.value[key].url = item.input;
}
function handleHeader(h: Header[]) {
    headers = h;
}
function handleQuery(q: Query[]) {
    query = q;
}
function update(res: Response) {
    status.value = res.Status;
    responseHeaders = res.Header;
    response.value += res.Message + '\n';
}
async function sendWebsocket(item: CompleteItem) {
    if (item.input) {
        if (!/ws?s/.test(item.input)) {
            item.input = 'wss://' + item.input;
        }
    } else {
        item.input = 'wss://echo.websocket.org';
    }
    item.connected = !item.connected;
    try {
        EventsEmit('connected', item.connected);
        let error = await WS(item.input, headers, query, item.connected);
        if (error) {
            item.connected = false;
            ElNotification({
                title: 'Something went wrong!',
                message: error,
                type: 'error',
                position: 'bottom-right',
            });
        }
        if (item.connected) {
            response.value = '';
            EventsOn('websocket', (data: Response) => {
                update(data);
            });
        } else {
            EventsOff('websocket');
        }
    } catch (err: any) {
        item.connected = false;
        ElNotification({
            title: 'Something went wrong!',
            message: err,
            type: 'error',
            position: 'bottom-right',
        });
    }
}
</script>
