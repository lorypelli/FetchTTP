<template>
    <ElSwitch
        v-model="readable"
        class="pt-2 absolute right-1 z-1"
        inactive-text="Raw"
        active-text="Human Readable"
    />
    <ElTabs class="pl-2">
        <ElTabPane label="Headers">
            <ElEmpty
                v-if="Object.keys(props.header).length == 0"
                class="flex justify-center h-1/2 lg:h-full"
                description="Nothing to show here..."
            />
            <ElText
                class="flex justify-center sticky top-0 bg-primary"
                :style="`color: ${getColor(props.status)};`"
            >
                {{ props.status }}
            </ElText>
            <ElText v-if="Object.keys(props.header).length > 0 && !readable">
                {{ props.header }}
            </ElText>
            <table v-if="Object.keys(props.header).length > 0 && readable">
                <tr
                    v-for="(item, index) in Object.keys(props.header)"
                    :index="index"
                >
                    <td class="break-all w-1/2 pl-5">
                        <ElText>{{ item }}</ElText>
                    </td>
                    <td class="break-all w-1/2 pr-5">
                        <ElText>{{ props.header[item].join('') }}</ElText>
                    </td>
                </tr>
            </table>
        </ElTabPane>
        <ElTabPane label="Response">
            <ElEmpty
                v-if="['', 'null'].includes(props.response.trim())"
                class="flex justify-center h-1/2 lg:h-full"
                description="Nothing to show here..."
            />
            <VueMonacoEditor
                v-if="
                    (isText(props.header) &&
                        !['', 'null'].includes(props.response.trim())) ||
                    !readable
                "
                :language="
                    props.header['Content-Type']
                        ? props.header['Content-Type'][0]
                              .split(';')[0]
                              .split('/')[1]
                        : undefined
                "
                :options="{
                    automaticLayout: true,
                    minimap: { enabled: false },
                    maxTokenizationLineLength: Infinity,
                    readOnly: true,
                }"
                :value="props.response"
                theme="vs-dark"
            />
            <div
                v-if="!isText(props.header) && readable"
                class="flex justify-center items-center h-1/2 lg:h-full"
            >
                <img v-if="isImage(props.header)" :src="props.url" />
                <audio v-if="isAudio(props.header)" :src="props.url" controls />
                <video v-if="isVideo(props.header)" :src="props.url" controls />
                <iframe
                    v-if="isPage(props.header)"
                    :srcdoc="baseURL(props.response, props.url)"
                    class="w-full rounded-2xl h-1/2 lg:h-full"
                    sandbox="allow-scripts allow-forms"
                />
                <embed
                    v-if="isPDF(props.header)"
                    :src="props.url"
                    class="w-full rounded-2xl h-1/2 lg:h-full"
                />
            </div>
        </ElTabPane>
    </ElTabs>
</template>

<script setup lang="ts">
import { VueMonacoEditor } from '@guolao/vue-monaco-editor';
import { ElEmpty, ElSwitch, ElTabPane, ElTabs, ElText } from 'element-plus';
import { ref } from 'vue';
import type { GenericHeader } from '../types';
const props = defineProps<{
    url?: string;
    status: string;
    header: GenericHeader;
    response: string;
}>();
const readable = ref(true);
function getColor(s: string) {
    const status = parseInt(s.split(' ')[0]);
    if (status >= 100 && status < 200) {
        return 'blue';
    } else if (status >= 200 && status < 300) {
        return 'green';
    } else if (status >= 300 && status < 400) {
        return 'yellow';
    } else if (status >= 400 && status < 500) {
        return 'orange';
    } else if (status >= 500 && status < 600) {
        return 'red';
    }
}
function isText(h: GenericHeader) {
    return !isPage(h) && !isPDF(h) && !isImage(h) && !isVideo(h) && !isAudio(h);
}
function isPDF(h: GenericHeader) {
    return (
        Object.entries(h).filter(([k, v]) => {
            return k == 'Content-Type' && v[0].includes('application/pdf');
        }).length > 0
    );
}
function isPage(h: GenericHeader) {
    return (
        Object.entries(h).filter(([k, v]) => {
            return k == 'Content-Type' && v[0].includes('text/html');
        }).length > 0
    );
}
function isImage(h: GenericHeader) {
    const regex = /image\/*/;
    return (
        Object.entries(h).filter(([k, v]) => {
            return k == 'Content-Type' && regex.test(v[0]);
        }).length > 0
    );
}
function isVideo(h: GenericHeader) {
    const regex = /video\/*/;
    return (
        Object.entries(h).filter(([k, v]) => {
            return k == 'Content-Type' && regex.test(v[0]);
        }).length > 0
    );
}
function isAudio(h: GenericHeader) {
    const regex = /audio\/*/;
    return (
        Object.entries(h).filter(([k, v]) => {
            return k == 'Content-Type' && regex.test(v[0]);
        }).length > 0
    );
}
function baseURL(r: string, u: string | undefined) {
    if (u) {
        return r.replace('<head>', `<head><base href="${u}">`);
    }
    return r;
}
</script>
