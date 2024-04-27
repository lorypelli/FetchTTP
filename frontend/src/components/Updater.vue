<script setup lang="ts">
import { ElInput, ElText } from 'element-plus';
import { ref } from 'vue';
defineOptions({
    name: 'Updater',
    components: {
        ElText,
        ElInput,
    },
});
const props = defineProps<{
    version: string;
    description: string;
}>();
const updateVersion = ref(props.version);
const updateDescription = ref(props.description);
</script>

<template>
    <div class="space-y-1">
        <ElInput v-model="updateVersion" class="flex flex-col items-center" readonly />
        <ElInput v-model="updateDescription" class="flex flex-col" type="textarea" resize="none" readonly />
    </div>
</template>
