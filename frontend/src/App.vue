<script lang="ts">
import { ElTabPane, ElTabs } from 'element-plus';
import HTTP from './components/HTTP.vue';
import WS from './components/WS.vue';
import 'element-plus/dist/index.css';
export default {
    name: 'App',
    components: {
        ElTabPane,
        ElTabs,
        HTTP,
        WS
    }
};
</script>
<template>
    <div class="px-1">
        <ElTabs>
            <ElTabPane label="http"><HTTP /></ElTabPane>
            <ElTabPane label="ws"><WS /></ElTabPane>
        </ElTabs>
    </div>
</template>