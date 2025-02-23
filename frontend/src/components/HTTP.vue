<script setup lang="ts">
import {
    ElButton,
    ElInput,
    ElNotification,
    ElOption,
    ElSelect,
} from 'element-plus';
import { HTTP } from '../../wailsjs/go/main/App.js';
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
let body = '';
let previousRequestTime = 0;
let requestTime = Date.now();
export default {
    data() {
        return {
            status: '',
            header: {},
            response: '',
            url: '',
        };
    },
    methods: {
        getRequest() {
            return previousRequestTime - requestTime;
        },
        handleSelect(item: CompleteItem) {
            localStorage.setItem(`${item.name}-select`, item.select);
        },
        handleInput(item: CompleteItem) {
            localStorage.setItem(`${item.name}-input-http`, item.input);
        },
        handleHeader(h: Header[]) {
            headers = h;
        },
        handleQuery(q: Query[]) {
            query = q;
        },
        handleBody(b: string) {
            body = b;
        },
        update(res: Response) {
            this.status = res.Status;
            this.header = res.Header;
            this.response = res.Body;
            this.url = res.URL;
        },
        sendRequest(item: CompleteItem) {
            if (item.input) {
                if (!/https?/.test(item.input)) {
                    item.input = 'https://' + item.input;
                }
            } else {
                item.input = 'https://echo.zuplo.io';
            }
            try {
                // @ts-ignore Types aren't correct!
                HTTP(item.select, item.input, headers, query, body).then(
                    (res: any) => {
                        if (res.Error) {
                            ElNotification({
                                title: 'Something went wrong!',
                                message: res.Error,
                                type: 'error',
                                position: 'bottom-right',
                            });
                            return;
                        }
                        this.update(res);
                        requestTime = Date.now();
                    },
                );
            } catch {
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
    <Tabs type="http">
        <template #default="{ item }">
            <div class="flex p-1 space-x-1">
                <ElSelect
                    v-model="item.select"
                    class="w-32"
                    @change="handleSelect(item)"
                >
                    <ElOption value="GET" />
                    <ElOption value="HEAD" />
                    <ElOption value="POST" />
                    <ElOption value="PUT" />
                    <ElOption value="DELETE" />
                    <ElOption value="CONNECT" />
                    <ElOption value="OPTIONS" />
                    <ElOption value="TRACE" />
                    <ElOption value="PATCH" />
                </ElSelect>
                <ElInput
                    v-model="item.input"
                    placeholder="echo.zuplo.io"
                    @input="handleInput(item)"
                    @keydown.enter="
                        () => {
                            previousRequestTime = Date.now();
                            const res = getRequest();
                            if (res <= 175) {
                                requestTime = Date.now();
                                return;
                            } else {
                                sendRequest(item);
                            }
                        }
                    "
                />
                <ElButton
                    class="w-20"
                    @click="
                        () => {
                            previousRequestTime = Date.now();
                            const res = getRequest();
                            if (res <= 175) {
                                requestTime = Date.now();
                                return;
                            } else {
                                sendRequest(item);
                            }
                        }
                    "
                >
                    Send
                </ElButton>
            </div>
            <Split
                :name="item.name"
                :url="url"
                :status="status"
                :header="header"
                :response="response"
                type="http"
                @headers="handleHeader"
                @query="handleQuery"
                @body="handleBody"
            />
        </template>
    </Tabs>
</template>
