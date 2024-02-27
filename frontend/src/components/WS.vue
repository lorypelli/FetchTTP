<script setup lang="ts">
import { ElButton, ElInput } from 'element-plus';
import 'element-plus/dist/index.css';
import { ref } from 'vue';
defineOptions({
    name: 'WS',
    components: {
        ElButton,
        ElInput
    }
});
const input = ref('');
</script>
<template>
    <div class="flex p-1 space-x-1">
        <ElInput v-model="input" placeholder="echo.websocket.org"></ElInput>
        <ElButton class="w-28">Connect</ElButton>
    </div>
</template>