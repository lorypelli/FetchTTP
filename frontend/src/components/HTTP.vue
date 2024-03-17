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
let tabIndex = 1;
const selectedTab = ref(`${tabIndex}`);
const t = ref([
    {
        name: `${selectedTab.value}`,
        select: 'GET',
        input: ''
    }
]);
export default {
    data() {
        return {
            status: '',
            header: {},
            response: '',
            url: '',
            tabs: [{
                name: `${selectedTab.value}`,
                select: 'GET',
                input: ''
            }]
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
        handleTab(targetName: TabPaneName | undefined, action: 'add' | 'remove') {
            switch (action) {
            case 'add': {
                const newTabIndex = `${++tabIndex}`;
                this.tabs.push({
                    name: newTabIndex,
                    select: 'GET',
                    input: ''
                });
                selectedTab.value = newTabIndex;
                break;
            }
            case 'remove': {
                const tabs = t.value;
                let activeTab = selectedTab.value;
                if (activeTab == targetName) {
                    tabs.forEach((tab, index) => {
                        if (tab.name == targetName) {
                            const nextTab = tabs[index + 1] || tabs[index - 1];
                            if (nextTab) {
                                activeTab = nextTab.name;
                            }
                        }
                    });
                }
                selectedTab.value = activeTab;
                t.value = tabs.filter((tab) => tab.name != targetName);
                break;
            }
            }
        },
        update(res: Response) {
            this.status = res.Status;
            this.header = res.Header;
            this.response = res.Body;
            this.url = res.URL;
        }
    }
};
</script>

<template>
    <ElTabs tab-position="left" editable v-on:edit="handleTab">
        <ElTabPane :label="item.select" v-for="(item, index) in tabs" :name="item.name" :key="index">
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