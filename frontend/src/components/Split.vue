<script setup lang="ts">
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
import type { GenericHeader as Header, Query } from '../types';
import Request from './Request.vue';
import Response from './Response.vue';
import { onBeforeUnmount, onMounted, ref } from 'vue';
const props = defineProps<{
    name: string;
    url?: string;
    status: string;
    header: Header;
    response: string;
    type: 'http' | 'ws';
}>();
const emit = defineEmits<{
    header: [value: Header[]];
    query: [value: Query[]];
    body: [value: string];
    message: [value: string];
}>();
const width = ref(0);
onMounted(() => {
    const headers = localStorage.getItem(`${props.name}-headers-${props.type}`);
    if (headers) {
        handleHeader(JSON.parse(headers));
    }
    const query = localStorage.getItem(`${props.name}-query-${props.type}`);
    if (query) {
        handleQuery(JSON.parse(query));
    }
    handleBody(localStorage.getItem(`${props.name}-body`) || '');
    handleMessage(localStorage.getItem(`${props.name}-message`) || '');
    update();
    window.addEventListener('resize', update);
})
onBeforeUnmount(() => {
    window.removeEventListener('resize', update);
})
function update() {
    width.value = window.innerWidth;
}
function handleHeader(h: Header[]) {
    emit('header', h);
}
function handleQuery(q: Query[]) {
    emit('query', q);
}
function handleBody(b: string) {
    emit('body', b);
}
function handleMessage(m: string) {
    emit('message', m);
}
</script>

<template>
    <Splitter :layout="width <= 1024 ? 'vertical' : 'horizontal'">
        <SplitterPanel :min-size="25">
            <Request :name="props.name" :type="props.type" @headers="handleHeader" @query="handleQuery"
                @body="handleBody" @message="handleMessage" />
        </SplitterPanel>
        <SplitterPanel :min-size="25">
            <Response :url="props.url" :status="props.status" :header="props.header" :response="props.response" />
        </SplitterPanel>
    </Splitter>
</template>
