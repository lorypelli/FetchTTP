<script setup lang="ts">
import { ElButton, ElInput, ElNotification } from 'element-plus';
import { WS } from '../../wailsjs/go/main/App';
import { EventsEmit, EventsOff, EventsOn } from '../../wailsjs/runtime/runtime';
import type { Header, Query, Response } from '../types';
import Split from './Split.vue';
import Tabs, { CompleteItem } from './Tabs.vue';
</script>

<script lang="ts">
let headers: Header[] = [
    {
        enabled: true,
        name: 'User-Agent',
        value: 'FetchTTP',
    },
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
        handleInput(item: CompleteItem) {
            localStorage.setItem(`${item.name}-input-ws`, item.input);
        },
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
        },
        async sendWebsocket(item: CompleteItem) {
            if (item.input) {
                if (!/wss?/.test(item.input)) {
                    item.input = 'wss://' + item.input;
                }
            } else {
                item.input = 'wss://echo.websocket.org';
            }
            item.connected = !item.connected;
            try {
                EventsEmit('connected', item.connected);
                // @ts-ignore Types aren't correct!
                let error = await WS(
                    item.input,
                    headers,
                    query,
                    item.connected,
                );
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
                    this.response = '';
                    EventsOn('websocket', (data: Response) => {
                        this.update(data);
                    });
                } else {
                    EventsOff('websocket');
                }
            } catch {
                item.connected = false;
                ElNotification({
                    title: 'Something went wrong!',
                    type: 'error',
                    position: 'bottom-right',
                });
            }
        },
    },
};
</script>

<template>
    <Tabs type="ws">
        <template #default="{ item }">
            <div class="flex p-1 space-x-1">
                <ElInput
                    v-model="item.input"
                    :disabled="item.connected"
                    placeholder="echo.websocket.org"
                    @input="handleInput(item)"
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
                :status="status"
                :header="header"
                :response="response"
                type="ws"
                @headers="handleHeader"
                @query="handleQuery"
            />
        </template>
    </Tabs>
</template>
