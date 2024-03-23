<script setup lang="ts">
import { ElButton, ElInput, ElNotification, ElOption, ElSelect } from 'element-plus';
import { HTTP } from '../../wailsjs/go/main/App.js';
import Split from './Split.vue';
import Tabs, { CompleteItem } from './Tabs.vue';
defineOptions({
    name: 'HTTP',
    components: {
        ElButton,
        ElInput,
        ElOption,
        ElSelect,
        Split,
        Tabs,
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
    URL: string,
    Status: string,
    Header: [],
    Body: string
    Error: string
}
let headers: Header[] = [
    {
        enabled: true,
        name: 'User-Agent',
        value: 'FetchTTP'
    }
];
let query: Query[] = [];
let body = '';
export default {
    data() {
        return {
            status: '',
            header: {},
            response: '',
            url: ''
        };
    },
    methods: {
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
                if (!item.input.startsWith('http://') && !item.input.startsWith('https://')) {
                    item.input = 'https://' + item.input;
                }
            }
            else {
                item.input = 'https://echo.zuplo.io';
            }
            try {
                HTTP(item.select, item.input, headers, query, body).then((res) => {
                    if (res.Error) {
                        ElNotification({
                            title: 'Something went wrong!',
                            message: res.Error,
                            type: 'error',
                            position: 'bottom-right'
                        });
                        return;
                    }
                    this.update(res);
                });
            }
            catch {
                ElNotification({
                    title: 'Something went wrong!',
                    type: 'error',
                    position: 'bottom-right'
                });
            }
        }
    }
};
</script>

<template>
    <Tabs type="http">
        <template #default="{ item }">
            <div class="flex p-1 space-x-1">
                <ElSelect class="w-32" v-model="item.select">
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
                <ElInput v-model="item.input" placeholder="echo.zuplo.io" v-on:keydown.enter="sendRequest(item)"></ElInput>
                <ElButton class="w-20" v-on:click="sendRequest(item)">Send</ElButton>
            </div>
            <Split :name="item.name" :url="url" :status="status" :header="header" :response="response" type='http'
                v-on:headers="handleHeader" v-on:query="handleQuery" v-on:body="handleBody" />
        </template>
    </Tabs>
</template>