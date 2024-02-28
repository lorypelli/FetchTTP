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
            headers: [{ [headerName.value]: headerValue }],
            query: [{ [queryName.value]: queryValue }]
        };
    },
    methods: {
        addHeader() {
            this.headers.push({});
        },
        addQuery() {
            this.query.push({});
        }
    }
};
</script>
<template>
    <Splitter>
        <SplitterPanel :min-size="20">
            <ElTabs class="pr-2">
                <ElTabPane label="Headers">
                    <div class="flex space-x-2 pr-2 pt-2" v-for="(item, index) in headers" :key="index">
                        <ElCheckbox />
                        <ElInput v-model="item.headerName.value"></ElInput>
                        <ElInput v-model="item.headerValue.value"></ElInput>
                    </div>
                    <div class="flex pt-2">
                        <ElButton v-on:click="addHeader">+</ElButton>
                        <ElButton>-</ElButton>
                    </div>
                </ElTabPane>
                <ElTabPane label="Query">
                    <div class="flex space-x-2 pr-2 pt-2" v-for="(item, index) in query" :key="index">
                        <ElCheckbox />
                        <ElInput v-model="item.queryName.value"></ElInput>
                        <ElInput v-model="item.queryValue.value"></ElInput>
                    </div>
                    <div class="flex pt-2">
                        <ElButton v-on:click="addQuery">+</ElButton>
                        <ElButton>-</ElButton>
                    </div>
                </ElTabPane>
                <ElTabPane label="Cookies"></ElTabPane>
                <ElTabPane label="Body" v-if="props.type == 'http'"></ElTabPane>
                <ElTabPane label="Message" v-if="props.type == 'ws'"></ElTabPane>
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