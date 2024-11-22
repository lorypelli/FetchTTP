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
                :headers="responseHeaders"
                :response="response"
                type="http"
                @headers="handleHeaders"
                @query="handleQuery"
                @body="handleBody"
            />
        </template>
    </Tabs>
</template>

<script setup lang="ts">
import {
    ElButton,
    ElInput,
    ElNotification,
    ElOption,
    ElSelect,
} from 'element-plus';
import { reactive, ref } from 'vue';
import { HTTP } from '../../wailsjs/go/main/App.js';
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
const status = ref('');
let responseHeaders: GenericHeader[] = reactive([]);
const response = ref('');
const url = ref('');
let query: Query[] = reactive([]);
const body = ref('');
const previousRequestTime = ref(0);
const requestTime = ref(Date.now());
function getRequest() {
    return previousRequestTime.value - requestTime.value;
}
function handleSelect(item: CompleteItem) {
    localStorage.setItem(`${item.name}-select`, item.select);
}
function handleInput(item: CompleteItem) {
    localStorage.setItem(`${item.name}-input-http`, item.input);
}
function handleHeaders(h: Header[]) {
    headers = h;
}
function handleQuery(q: Query[]) {
    query = q;
}
function handleBody(b: string) {
    body.value = b;
}
function update(res: Response) {
    status.value = res.Status;
    responseHeaders = res.Header;
    response.value = res.Body;
    url.value = res.URL;
}
function sendRequest(item: CompleteItem) {
    if (item.input) {
        if (!/http?s/.test(item.input)) {
            item.input = 'https://' + item.input;
        }
    } else {
        item.input = 'https://echo.zuplo.io';
    }
    try {
        HTTP(item.select, item.input, headers, query, body).then((res: any) => {
            if (res.Error) {
                ElNotification({
                    title: 'Something went wrong!',
                    message: res.Error,
                    type: 'error',
                    position: 'bottom-right',
                });
                return;
            }
            update(res);
            requestTime.value = Date.now();
        });
    } catch {
        ElNotification({
            title: 'Something went wrong!',
            type: 'error',
            position: 'bottom-right',
        });
    }
}
</script>
