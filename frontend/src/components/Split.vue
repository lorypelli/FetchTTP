<script setup lang="ts">
import { ElButton, ElCheckbox, ElInput, ElTabPane, ElTabs } from 'element-plus';
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
import { ref } from 'vue';
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
        ElCheckbox
    }
});
const props = defineProps<{
    type: 'http' | 'ws'
}>();
</script>
<script lang="ts">
const headerName = ref('');
const headerValue = ref('');
const queryName = ref('');
const queryValue = ref('');
export default {
    data() {
        return {
            headers: [{ [headerName.value]: headerValue.value }],
            query: [{ [queryName.value]: queryValue.value }]
        };
    },
    methods: {
        addHeader() {
            this.headers.push({});
        },
        removeHeader() {
            if (this.headers.length == 1) {
                return;
            }
            this.headers.pop();
        },
        addQuery() {
            this.query.push({});
        },
        removeQuery() {
            if (this.query.length == 1) {
                return;
            }
            this.query.pop();
        },
    }
};
</script>
<template>
    <Splitter>
        <SplitterPanel :min-size="20">
            <ElTabs class="pr-2">
                <ElTabPane label="Headers">
                    <div class="flex space-x-2 pr-2 pt-2" v-for="(item, index) in headers" :key="index">
                        <ElCheckbox class="w-10" />
                        <ElInput />
                        <ElInput />
                    </div>
                    <div class="flex flex-col pt-2 space-y-2">
                        <ElButton v-on:click="addHeader">+</ElButton>
                        <ElButton v-on:click="removeHeader">-</ElButton>
                    </div>
                </ElTabPane>
                <ElTabPane label="Query">
                    <div class="flex space-x-2 pr-2 pt-2" v-for="(item, index) in query" :key="index">
                        <ElCheckbox class="w-10" />
                        <ElInput />
                        <ElInput />
                    </div>
                    <div class="flex flex-col pt-2 space-y-2">
                        <ElButton v-on:click="addQuery">+</ElButton>
                        <ElButton v-on:click="removeQuery">-</ElButton>
                    </div>
                </ElTabPane>
                <ElTabPane label="Cookies"></ElTabPane>
                <ElTabPane label="Body" v-if="props.type == 'http'">
                    <ElInput type="textarea" resize="none" />
                </ElTabPane>
                <ElTabPane label="Message" v-if="props.type == 'ws'">
                    <div class="flex flex-col space-y-1">
                        <ElInput type="textarea" resize="none" />
                        <ElButton>Send</ElButton>
                    </div>
                </ElTabPane>
            </ElTabs>
        </SplitterPanel>
        <SplitterPanel :min-size="20">
            <ElTabs class="pl-2">
                <ElTabPane label="Headers"></ElTabPane>
                <ElTabPane label="Response"></ElTabPane>
            </ElTabs>
        </SplitterPanel>
    </Splitter>
</template>