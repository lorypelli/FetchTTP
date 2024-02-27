<script setup lang="ts">
import { ElTabPane, ElTabs } from 'element-plus';
import HTTP from './components/HTTP.vue';
import WS from './components/WS.vue';
import 'element-plus/dist/index.css';
defineOptions({
    name: 'App',
    components: {
        ElTabPane,
        ElTabs,
        HTTP,
        WS
    }
});
</script>
<template>
    <ElTabs>
        <ElTabPane label="http">
            <HTTP />
        </ElTabPane>
        <ElTabPane label="ws">
            <WS />
        </ElTabPane>
    </ElTabs>
</template>