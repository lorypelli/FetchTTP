<script setup lang="ts">
import { ElButton, ElInput, ElOption, ElSelect } from 'element-plus';
import { MakeRequest } from '../../wailsjs/go/main/App.js';
import { ref } from 'vue';
import Split from './Split.vue';
defineOptions({
    name: 'HTTP',
    components: {
        ElButton,
        ElInput,
        ElOption,
        ElSelect,
        Split
    }
});
const select = ref('GET');
const input = ref('');
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
    Body: string
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
            response: ''
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
        }
    }
};
</script>

<template>
    <div class="flex p-1 space-x-1">
        <ElSelect class="w-32" v-model="select">
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
        <ElInput v-model="input" placeholder="echo.zuplo.io"></ElInput>
        <ElButton class="w-20" v-on:click="() => {
            if (input) {
                if (!input.startsWith('http://') && !input.startsWith('https://')) {
                    input = 'https://' + input
                }
            }
            else {
                input = 'https://echo.zuplo.io'
            }
            MakeRequest(select, input, headers, query, body).then((res) => {
                update(res)
            })
        }">Send</ElButton>
    </div>
    <Split :status=status :header=header :response=response type='http' v-on:headers="handleHeader" v-on:query="handleQuery" v-on:body="handleBody" />
</template>