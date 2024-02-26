<script lang="ts">
import { ElButton, ElInput, ElOption, ElSelect } from 'element-plus';
import 'element-plus/dist/index.css';
import { MakeRequest } from '../../wailsjs/go/main/App.js';
export default {
    name: 'HTTP',
    components: {
        ElButton,
        ElInput,
        ElOption,
        ElSelect
    }
};
</script>
<template>
    <div class="flex p-1 space-x-1">
        <ElSelect class="w-32" v-model="select">
            <ElOption value="GET" />
            <ElOption value="HEAD" />
            <ElOption value="POST" />
            <ElOption value="PUT" />
            <ElOption value="DELETE" />
            <ElOption value="CONNECT" />
            <ElOption value="OPTIONS" />
            <ElOption value="TRACE" />
            <ElOption value="PATCH" />
        </ElSelect>
        <ElInput v-model="input" placeholder="echo.zuplo.io"></ElInput>
        <ElButton class="w-20" v-on:click="() => {
            MakeRequest(select, input ? (input.startsWith('https://') ? input : (input.startsWith('http://') ? input : 'https' + input)) : 'https://echo.zuplo.io').then((res) => {
                console.log(res)
                return
            })
        }">Send</ElButton>
    </div>
</template>
<script setup lang="ts">
import { ref } from 'vue';
const select = ref('GET');
const input = ref('');
</script>