<script setup lang="ts">
import { ElButton, ElForm, ElFormItem, ElInput, ElMessageBox } from 'element-plus';
import { h } from 'vue';
defineOptions({
    name: 'CURL',
    components: {
        ElButton,
        ElForm,
        ElFormItem,
        ElInput
    }
});
</script>

<template>
    <ElButton v-on:click="() => {
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