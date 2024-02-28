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
const headerName = ref('');
const headerValue = ref('');
const queryName = ref('');
const queryValue = ref('');
</script>
<template>
    <Splitter>
        <SplitterPanel :min-size="20">
            <ElTabs class="pr-2">
                <ElTabPane label="Headers">
                    <div class="flex space-x-2 pr-2">
                        <ElCheckbox />
                        <ElInput v-model="headerName"></ElInput>
                        <ElInput v-model="headerValue"></ElInput>
                    </div>
                    <div class="flex pt-2">
                        <ElButton>+</ElButton>
                        <ElButton>-</ElButton>
                    </div>
                </ElTabPane>
                <ElTabPane label="Query">
                    <div class="flex space-x-2 pr-2">
                        <ElCheckbox />
                        <ElInput v-model="queryName"></ElInput>
                        <ElInput v-model="queryValue"></ElInput>
                    </div>
                    <div class="flex pt-2">
                        <ElButton>+</ElButton>
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