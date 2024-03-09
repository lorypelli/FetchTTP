<script setup lang="ts">
import { ElButton, ElCheckbox, ElInput, ElTabPane, ElTabs, ElText, ElDivider, ElScrollbar } from 'element-plus';
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Split',
    components: {
        Splitter,
        SplitterPanel,
        ElTabPane,
        ElTabs,
        ElInput,
        ElButton,
        ElCheckbox,
        ElText,
        ElDivider
    }
});
const props = defineProps<{
    url: string,
    status: string,
    header: object,
    response: string,
    type: 'http' | 'ws'
}>();
</script>

<script lang="ts">
export default {
    data() {
        return {
            headers: [{
                enabled: true,
                name: 'User-Agent',
                value: 'FetchTTP'
            }],
            query: [{
                enabled: false,
                name: '',
                value: ''
            }],
            body: '',
            message: ''
        };
    },
    methods: {
        addHeader() {
            this.headers.push({
                enabled: false,
                name: '',
                value: ''
            });
        },
        removeHeader() {
            if (this.headers.length == 1) {
                return;
            }
            this.headers.pop();
        },
        addQuery() {
            this.query.push({
                enabled: false,
                name: '',
                value: ''
            });
        },
        removeQuery() {
            if (this.query.length == 1) {
                return;
            }
            this.query.pop();
        },
        sendHeader() {
            this.$emit('headers', this.headers);
        },
        sendQuery() {
            this.$emit('query', this.query);
        },
        sendBody() {
            this.$emit('body', this.body);
        },
        sendMessage() {
            this.$emit('message', this.message);
        },
        isText(h: object) {
            return !this.isImage(h);
        },
        isImage(h: object) {
            const regex = /image\/*/;
            return Object.entries(h).filter(([k, v]) => {
                return k == 'Content-Type' && regex.test(v[0]);
            }).length > 0;
        }
    }
};
</script>

<template>
    <Splitter>
        <SplitterPanel :min-size="20">
            <ElTabs class="pr-2">
                <ElTabPane :label="`Headers (${headers.filter((h) => {
                    return h.enabled
                }).length})`">
                    <ElScrollbar height="83.5vh">
                        <div class="flex space-x-2 pr-2 pt-2" v-for="(item, index) in headers" :key="index">
                            <ElCheckbox v-model="item.enabled" v-on:change="() => {
                                sendHeader()
                                if (!item.name || !item.value) {
                                    item.enabled = false
                                }
                            }" class="w-10" />
                            <ElInput v-model="item.name" v-on:input="() => {
                                sendHeader()
                                if (item.name && item.value) {
                                    item.enabled = true
                                }
                                else {
                                    item.enabled = false
                                }
                            }" />
                            <ElInput v-model="item.value" v-on:input="() => {
                                sendHeader()
                                if (item.name && item.value) {
                                    item.enabled = true
                                }
                                else {
                                    item.enabled = false
                                }
                            }" />
                        </div>
                        <div class="flex flex-col pt-2 space-y-2">
                            <ElButton v-on:click="addHeader">+</ElButton>
                            <ElButton v-on:click="removeHeader">-</ElButton>
                        </div>
                    </ElScrollbar>
                </ElTabPane>
                <ElTabPane :label="`Query (${query.filter((q) => {
                    return q.enabled
                }).length})`">
                    <ElScrollbar height="83.5vh">
                        <div class="flex space-x-2 pr-2 pt-2" v-for="(item, index) in query" :key="index">
                            <ElCheckbox v-model="item.enabled" v-on:change="() => {
                                sendQuery()
                                if (!item.name || !item.value) {
                                    item.enabled = false
                                }
                            }" class="w-10" />
                            <ElInput v-model="item.name" v-on:input="() => {
                                sendQuery()
                                if (item.name && item.value) {
                                    item.enabled = true
                                }
                                else {
                                    item.enabled = false
                                }
                            }" />
                            <ElInput v-model="item.value" v-on:input="() => {
                                sendQuery()
                                if (item.name && item.value) {
                                    item.enabled = true
                                }
                                else {
                                    item.enabled = false
                                }
                            }" />
                        </div>
                        <div class="flex flex-col pt-2 space-y-2">
                            <ElButton v-on:click="addQuery">+</ElButton>
                            <ElButton v-on:click="removeQuery">-</ElButton>
                        </div>
                    </ElScrollbar>
                </ElTabPane>
                <ElTabPane label="Cookies"></ElTabPane>
                <ElTabPane label="Body" v-if="props.type == 'http'">
                    <ElInput v-model="body" type="textarea" resize="none" />
                </ElTabPane>
                <ElTabPane label="Message" v-if="props.type == 'ws'">
                    <div class="flex flex-col space-y-1">
                        <ElInput v-model="message" type="textarea" resize="none" />
                        <ElButton>Send</ElButton>
                    </div>
                </ElTabPane>
            </ElTabs>
        </SplitterPanel>
        <SplitterPanel :min-size="20">
            <ElTabs class="pl-2">
                <ElTabPane label="Headers">
                    <ElText class="flex justify-center">{{ props.status }}</ElText>
                    <ElDivider v-if="props.status" />
                    <ElScrollbar height="75vh">
                        <ElText v-if="Object.keys(props.header).length > 0">{{ props.header }}</ElText>
                    </ElScrollbar>
                </ElTabPane>
                <ElTabPane label="Response">
                    <ElScrollbar height="83.5vh">
                        <ElText v-if="isText(props.header)">{{ props.response }}</ElText>
                        <div class="flex justify-center">
                            <img v-if="isImage(props.header)" :src="props.url" />
                        </div>
                    </ElScrollbar>
                </ElTabPane>
            </ElTabs>
        </SplitterPanel>
    </Splitter>
</template>