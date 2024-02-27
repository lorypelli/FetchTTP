<script setup lang="ts">
import { ElTabPane, ElTabs } from 'element-plus';
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Split',
    components: {
        Splitter,
        SplitterPanel,
        ElTabPane,
        ElTabs
    }
});
const props = defineProps<{
    type: 'http' | 'ws'
}>();
</script>
<template>
    <Splitter>
        <SplitterPanel :min-size="20">
            <ElTabs>
                <ElTabPane label="Headers"></ElTabPane>
                <ElTabPane label="Query"></ElTabPane>
                <ElTabPane label="Cookies"></ElTabPane>
                <ElTabPane label="Body" v-if="props.type == 'http'"></ElTabPane>
                <ElTabPane label="Message" v-if="props.type == 'ws'"></ElTabPane>
            </ElTabs>
        </SplitterPanel>
        <SplitterPanel :min-size="20">
            <ElTabs>
                <ElTabPane label="Headers"></ElTabPane>
                <ElTabPane label="Response"></ElTabPane>
            </ElTabs>
        </SplitterPanel>
    </Splitter>
</template>