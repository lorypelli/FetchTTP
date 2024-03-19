<script setup lang="ts">
import { ElButton, ElInput, ElNotification, ElOption, ElSelect, ElTabPane, ElTabs, TabPaneName } from 'element-plus';
import { HTTP } from '../../wailsjs/go/main/App.js';
import Split from './Split.vue';
import { ref } from 'vue';
defineOptions({
    name: 'HTTP',
    components: {
        ElButton,
        ElInput,
        ElOption,
        ElSelect,
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
        }
    }
};
let tabIndex = 1;
const selectedTab = ref('1');
const tabs = ref([
    {
        name: '1',
        select: 'GET',
        input: ''
    }
]);
export function handleTab(targetName: TabPaneName | undefined, action: 'add' | 'remove') {
    switch (action) {
    case 'add': {
        const newTabIndex = `${++tabIndex}`;
        tabs.value.push({
            name: newTabIndex,
            select: 'GET',
            input: ''
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

<template>
    <ElTabs v-model="selectedTab" tab-position="left" editable v-on:edit="handleTab">
        <ElTabPane :label="item.name" v-for="(item, index) in tabs" :name="item.name" :key="index">
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
                <ElInput v-model="item.input" placeholder="echo.zuplo.io"></ElInput>
                <ElButton class="w-20" v-on:click="() => {
        if (item.input) {
            if (!item.input.startsWith('http://') && !item.input.startsWith('https://')) {
                item.input = 'https://' + item.input
            }
        }
        else {
            item.input = 'https://echo.zuplo.io'
        }
        try {
            HTTP(item.select, item.input, headers, query, body).then((res) => {
                if (res.Error) {
                    ElNotification({
                        title: 'Something went wrong!',
                        message: res.Error,
                        type: 'error',
                        position: 'bottom-right'
                    })
                    return
                }
                update(res)
            })
        }
        catch {
            ElNotification({
                title: 'Something went wrong!',
                type: 'error',
                position: 'bottom-right'
            })
        }
    }">Send</ElButton>
            </div>
            <Split :url="url" :status="status" :header="header" :response="response" type='http'
                v-on:headers="handleHeader" v-on:query="handleQuery" v-on:body="handleBody" />
        </ElTabPane>
    </ElTabs>
</template>