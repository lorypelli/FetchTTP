<script setup lang="ts">
import { ElButton, ElInput } from 'element-plus';
import 'element-plus/dist/index.css';
import { ref } from 'vue';
import Split from './Split.vue';
defineOptions({
    name: 'WS',
    components: {
        ElButton,
        ElInput,
        Split
    }
});
const input = ref('');
</script>

<template>
    <div class="flex p-1 space-x-1">
        <ElInput v-model="input" placeholder="echo.websocket.org"></ElInput>
        <ElButton class="w-28">Connect</ElButton>
    </div>
    <Split type='ws' />
</template>