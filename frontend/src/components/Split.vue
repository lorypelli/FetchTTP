<script setup lang="ts">
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
import { onBeforeUnmount, onMounted, ref } from 'vue';
import type { GenericHeader, Header, Query } from '../types';
import Request from './Request.vue';
import Response from './Response.vue';
const props = defineProps<{
    name: string;
    url?: string;
    status: string;
    header: GenericHeader;
    response: string;
    type: 'http' | 'ws';
}>();
const emit = defineEmits<{
    headers: [value: Header[]];
    query: [value: Query[]];
    body: [value: string];
    message: [value: string];
}>();
const width = ref(0);
onMounted(() => {
    const headers = localStorage.getItem(`${props.name}-headers-${props.type}`);
    if (headers) {
        handleHeaders(JSON.parse(headers));
    }
    const query = localStorage.getItem(`${props.name}-query-${props.type}`);
    if (query) {
        handleQuery(JSON.parse(query));
    }
    handleBody(localStorage.getItem(`${props.name}-body`) || '');
    handleMessage(localStorage.getItem(`${props.name}-message`) || '');
    update();
    window.addEventListener('resize', update);
});
onBeforeUnmount(() => window.removeEventListener('resize', update));
function update() {
    width.value = window.innerWidth;
}
function handleHeaders(h: Header[]) {
    emit('headers', h);
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
            <Request
                :name="props.name"
                :type="props.type"
                @headers="handleHeaders"
                @query="handleQuery"
                @body="handleBody"
                @message="handleMessage"
            />
        </SplitterPanel>
        <SplitterPanel :min-size="25">
            <Response
                :url="props.url"
                :status="props.status"
                :header="props.header"
                :response="props.response"
            />
        </SplitterPanel>
    </Splitter>
</template>
