<script setup lang="ts">
import { ElInput, ElText } from 'element-plus';
import { ref } from 'vue';
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Updater',
    components: {
        ElText,
        ElInput
    }
});
const props = defineProps<{
    version: string,
    description: string
}>();
const updateVersion = ref(props.version);
const updateDescription = ref(props.description);
</script>

<template>
    <div class="space-y-1">
        <ElInput class="flex flex-col items-center" readonly v-model="updateVersion" />
        <ElInput class="flex flex-col" type="textarea" resize="none" readonly v-model="updateDescription" />
    </div>
</template>