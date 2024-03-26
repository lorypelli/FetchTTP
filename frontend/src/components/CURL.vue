<script setup lang="ts">
import { ElButton, ElForm, ElFormItem, ElInput, ElMessageBox } from 'element-plus';
import { ref, h } from 'vue';
defineOptions({
    name: 'CURL',
    components: {
        ElButton,
        ElForm,
        ElFormItem,
        ElInput
    }
});
const curl = ref('');
</script>

<template>
    <ElButton v-model="curl" v-on:click="() => {
        ElMessageBox({
            title: 'cURL Request',
            message: h(ElForm, {
                labelPosition: 'top'
            }, [
                h(ElFormItem, {
                    label: 'URL'
                }, [
                    h(ElInput, {
                        type: 'textarea',
                        resize: 'none'
                    })
                ])
            ])
        })
            .catch(() => { })
    }">Open cURL</ElButton>
</template>