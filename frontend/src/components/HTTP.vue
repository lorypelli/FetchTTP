<script setup lang="ts">
import { ElButton, ElInput, ElOption, ElSelect } from 'element-plus';
//import { MakeRequest } from '../../wailsjs/go/main/App.js';
import { ref } from 'vue';
import Split from './Split.vue';
defineOptions({
    name: 'HTTP',
    components: {
        ElButton,
        ElInput,
        ElOption,
        ElSelect,
        Split
    }
});
const select = ref('GET');
const input = ref('');
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
            if (input) {
                if (!input.startsWith('http://') && !input.startsWith('https://')) {
                    input = 'https://' + input
                }
            }
            else {
                input = 'https://echo.zuplo.io'
            }
            //MakeRequest(select, input).then((res) => {})
        }">Send</ElButton>
    </div>
    <Split type='http' ref="data" />
</template>