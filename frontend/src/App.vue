<script setup lang="ts">
import { ElTabPane, ElTabs } from 'element-plus';
import HTTP from './components/HTTP.vue';
import WS from './components/WS.vue';
import 'element-plus/dist/index.css';
import 'primevue/resources/themes/aura-light-green/theme.css';
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
        <ElTabPane label="HTTP">
            <HTTP />
        </ElTabPane>
        <ElTabPane label="WS">
            <WS />
        </ElTabPane>
    </ElTabs>
</template>