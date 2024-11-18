<script setup lang="ts">
import { VueMonacoEditor } from '@guolao/vue-monaco-editor';
import {
    ElButton,
    ElCheckbox,
    ElDatePicker,
    ElForm,
    ElFormItem,
    ElInput,
    ElMessageBox,
    ElOption,
    ElSelect,
    ElTabPane,
    ElTabs,
} from 'element-plus';
import { h, onMounted, reactive, ref } from 'vue';
import { EventsEmit } from '../../wailsjs/runtime/runtime';
import type { Header, Query } from '../types';
const props = defineProps<{
    name: string;
    type: 'http' | 'ws';
}>();
const emit = defineEmits<{
    headers: [value: Header[]];
    query: [value: Query[]];
    body: [value: string];
    message: [value: string];
}>();
const bodyType = ref('JSON');
let headers: Header[] = reactive([
    {
        enabled: true,
        name: 'User-Agent',
        value: 'FetchTTP',
    },
]);
let query: Query[] = reactive([]);
const body = ref('');
const message = ref('');
onMounted(() => {
    let localHeaders = localStorage.getItem(
        `${props.name}-headers-${props.type}`,
    );
    if (localHeaders) {
        headers = JSON.parse(localHeaders);
    }
    let localQuery = localStorage.getItem(`${props.name}-query-${props.type}`);
    if (localQuery) {
        query = JSON.parse(localQuery);
    }
    body.value = localStorage.getItem(`${props.name}-body`) || '';
    message.value = localStorage.getItem(`${props.name}-message`) || '';
});
function add(type: 'header' | 'query', index: number) {
    switch (type) {
        case 'header': {
            for (let i = 0; i < headers.length; i++) {
                if (i == index) {
                    headers.splice(i + 1, 0, {
                        enabled: false,
                        name: '',
                        value: '',
                    });
                }
            }
            break;
        }
        case 'query': {
            for (let i = 0; i < query.length; i++) {
                if (i == index) {
                    query.splice(i + 1, 0, {
                        enabled: false,
                        name: '',
                        value: '',
                    });
                }
            }
            break;
        }
    }
}
function remove(type: 'header' | 'query', index: number) {
    switch (type) {
        case 'header': {
            for (let i = 0; i < headers.length; i++) {
                if (i == index) {
                    if (i == 0) {
                        headers[0] = {
                            enabled: false,
                            name: '',
                            value: '',
                        };
                        localStorage.removeItem(
                            `${props.name}-headers-${props.type}`,
                        );
                    } else {
                        headers.splice(i, 1);
                        localStorage.setItem(
                            `${props.name}-headers-${props.type}`,
                            JSON.stringify(headers),
                        );
                    }
                }
            }
            break;
        }
        case 'query': {
            for (let i = 0; i < query.length; i++) {
                if (i == index) {
                    if (i == 0) {
                        query[0] = {
                            enabled: false,
                            name: '',
                            value: '',
                        };
                        localStorage.removeItem(
                            `${props.name}-query-${props.type}`,
                        );
                    } else {
                        query.splice(i, 1);
                        localStorage.setItem(
                            `${props.name}-query-${props.type}`,
                            JSON.stringify(query),
                        );
                    }
                }
            }
            break;
        }
    }
}
function getCorrectType(s: string) {
    switch (s) {
        case 'JSON':
            return 'application/json; charset=utf-8';
        case 'YAML':
            return 'application/yaml; charset=utf-8';
        case 'XML':
            return 'text/xml; charset=utf-8';
        case 'TXT':
            return 'text/plain; charset=utf-8';
    }
}
function setCorrectType(s: string) {
    let found = false;
    headers.forEach((h) => {
        if (h.name == 'Content-Type') {
            h.value = getCorrectType(s) || 'text/plain';
            found = true;
        }
    });
    if (!found) {
        headers.push({
            name: 'Content-Type',
            value: getCorrectType(s) || 'text/plain',
            enabled: true,
        });
    }
}
function sendHeader() {
    localStorage.setItem(
        `${props.name}-headers-${props.type}`,
        JSON.stringify(headers),
    );
    emit('headers', headers);
}
function sendQuery() {
    localStorage.setItem(
        `${props.name}-query-${props.type}`,
        JSON.stringify(query),
    );
    emit('query', query);
}
function sendBody() {
    localStorage.setItem(`${props.name}-body`, body.value);
    emit('body', body.value);
}
function sendMessage() {
    localStorage.setItem(`${props.name}-message`, message.value);
    emit('message', message.value);
}
</script>

<template>
    <ElTabs class="pr-2">
        <ElTabPane
            :label="`Headers (${
                headers.filter((h) => {
                    return h.enabled;
                }).length
            })`"
        >
            <div
                v-for="(item, index) in headers"
                :key="index"
                class="flex space-x-1 pr-2 pt-2"
            >
                <ElCheckbox
                    v-model="item.enabled"
                    class="w-10"
                    @change="
                        () => {
                            if (!item.name || !item.value) {
                                item.enabled = false;
                            }
                            sendHeader();
                        }
                    "
                />
                <ElInput
                    v-model="item.name"
                    @input="
                        () => {
                            if (item.name && item.value) {
                                item.enabled = true;
                            } else {
                                item.enabled = false;
                            }
                            sendHeader();
                        }
                    "
                />
                <ElInput
                    v-model="item.value"
                    @input="
                        () => {
                            if (item.name && item.value) {
                                item.enabled = true;
                            } else {
                                item.enabled = false;
                            }
                            sendHeader();
                        }
                    "
                />
                <ElButton @click="add('header', index)"> + </ElButton>
                <ElButton @click="remove('header', index)"> - </ElButton>
            </div>
        </ElTabPane>
        <ElTabPane
            :label="`Query (${
                query.filter((q) => {
                    return q.enabled;
                }).length
            })`"
        >
            <div
                v-for="(item, index) in query"
                :key="index"
                class="flex space-x-1 pr-2 pt-2"
            >
                <ElCheckbox
                    v-model="item.enabled"
                    class="w-10"
                    @change="
                        () => {
                            if (!item.name || !item.value) {
                                item.enabled = false;
                            }
                            sendQuery();
                        }
                    "
                />
                <ElInput
                    v-model="item.name"
                    @input="
                        () => {
                            if (item.name && item.value) {
                                item.enabled = true;
                            } else {
                                item.enabled = false;
                            }
                            sendQuery();
                        }
                    "
                />
                <ElInput
                    v-model="item.value"
                    @input="
                        () => {
                            if (item.name && item.value) {
                                item.enabled = true;
                            } else {
                                item.enabled = false;
                            }
                            sendQuery();
                        }
                    "
                />
                <ElButton @click="add('query', index)"> + </ElButton>
                <ElButton @click="remove('query', index)"> - </ElButton>
            </div>
        </ElTabPane>
        <ElTabPane
            label="Cookies"
            class="flex justify-center items-center h-1/4 lg:h-full"
        >
            <ElButton
                class="scale-150"
                @click="
                    () => {
                        const form = reactive({
                            key: '',
                            value: '',
                            domain: '',
                            path: '',
                            expires: '',
                            same_site: 'Lax',
                            secure: false,
                            http_only: false,
                        });
                        ElMessageBox({
                            title: 'New Cookie',
                            message: () =>
                                h(
                                    ElForm,
                                    {
                                        labelPosition: 'top',
                                    },
                                    [
                                        h(
                                            ElFormItem,
                                            {
                                                label: 'Key',
                                            },
                                            [
                                                h(ElInput, {
                                                    modelValue: form.key,
                                                    'onUpdate:modelValue': (
                                                        val,
                                                    ) => {
                                                        form.key = val;
                                                    },
                                                }),
                                            ],
                                        ),
                                        h(
                                            ElFormItem,
                                            {
                                                label: 'Value',
                                            },
                                            [
                                                h(ElInput, {
                                                    modelValue: form.value,
                                                    'onUpdate:modelValue': (
                                                        val,
                                                    ) => {
                                                        form.value = val;
                                                    },
                                                }),
                                            ],
                                        ),
                                        h(
                                            ElFormItem,
                                            {
                                                label: 'Domain',
                                            },
                                            [
                                                h(ElInput, {
                                                    modelValue: form.domain,
                                                    'onUpdate:modelValue': (
                                                        val,
                                                    ) => {
                                                        form.domain = val;
                                                    },
                                                }),
                                            ],
                                        ),
                                        h(
                                            ElFormItem,
                                            {
                                                label: 'Path',
                                            },
                                            [
                                                h(ElInput, {
                                                    modelValue: form.path,
                                                    'onUpdate:modelValue': (
                                                        val,
                                                    ) => {
                                                        form.path = val;
                                                    },
                                                }),
                                            ],
                                        ),
                                        h(
                                            ElFormItem,
                                            {
                                                label: 'Expires',
                                            },
                                            [
                                                h(
                                                    'div',
                                                    {
                                                        class: 'w-full',
                                                    },
                                                    [
                                                        h(ElDatePicker, {
                                                            modelValue:
                                                                form.expires,
                                                            type: 'datetime',
                                                            valueFormat: 'x',
                                                            'onUpdate:modelValue':
                                                                (val) => {
                                                                    form.expires =
                                                                        val;
                                                                },
                                                        }),
                                                    ],
                                                ),
                                            ],
                                        ),
                                        h(
                                            ElFormItem,
                                            {
                                                label: 'SameSite',
                                            },
                                            [
                                                h(
                                                    ElSelect,
                                                    {
                                                        modelValue:
                                                            form.same_site,
                                                        'onUpdate:modelValue': (
                                                            val,
                                                        ) => {
                                                            form.same_site =
                                                                val;
                                                        },
                                                    },
                                                    [
                                                        h(ElOption, {
                                                            value: 'None',
                                                        }),
                                                        h(ElOption, {
                                                            value: 'Lax',
                                                        }),
                                                        h(ElOption, {
                                                            value: 'Strict',
                                                        }),
                                                    ],
                                                ),
                                            ],
                                        ),
                                        h(
                                            'div',
                                            {
                                                class: 'flex justify-around',
                                            },
                                            [
                                                h(
                                                    ElFormItem,
                                                    {
                                                        label: 'Secure',
                                                    },
                                                    [
                                                        h(ElCheckbox, {
                                                            modelValue:
                                                                form.secure,
                                                            'onUpdate:modelValue':
                                                                (val) => {
                                                                    form.secure =
                                                                        !!val;
                                                                },
                                                        }),
                                                    ],
                                                ),
                                                h(
                                                    ElFormItem,
                                                    {
                                                        label: 'HTTPOnly',
                                                    },
                                                    [
                                                        h(ElCheckbox, {
                                                            modelValue:
                                                                form.http_only,
                                                            'onUpdate:modelValue':
                                                                (val) => {
                                                                    form.http_only =
                                                                        !!val;
                                                                },
                                                        }),
                                                    ],
                                                ),
                                            ],
                                        ),
                                    ],
                                ),
                            confirmButtonText: 'Add',
                            closeOnClickModal: false,
                        })
                            .then(() => {
                                headers.push({
                                    enabled: true,
                                    name: 'Cookie',
                                    value: `${form.key}=${form.value};${form.domain ? `Domain=${form.domain};` : ''}${form.path ? `Path=${form.path};` : ''}${form.expires ? `Expires=${new Date(form.expires).toUTCString()};` : ''}${form.same_site ? `SameSite=${form.same_site};` : ''}${form.secure ? 'Secure;' : ''}${form.http_only ? 'HTTPOnly;' : ''}`,
                                });
                            })
                            .catch(() => {});
                    }
                "
            >
                Add Cookie
            </ElButton>
        </ElTabPane>
        <ElTabPane v-if="props.type == 'http'" label="Body">
            <VueMonacoEditor
                class="editor"
                v-model:value="body"
                :language="bodyType.toLowerCase()"
                :options="{
                    automaticLayout: true,
                    minimap: { enabled: false },
                    maxTokenizationLineLength: Infinity,
                }"
                theme="vs-dark"
                @change="
                    () => {
                        setCorrectType(bodyType);
                        sendBody();
                    }
                "
            />
            <ElSelect v-model="bodyType" class="pt-2" @change="setCorrectType">
                <ElOption value="JSON" />
                <ElOption value="YAML" />
                <ElOption value="XML" />
                <ElOption value="TXT" />
            </ElSelect>
        </ElTabPane>
        <ElTabPane v-if="props.type == 'ws'" label="Message">
            <div class="flex flex-col space-y-1">
                <VueMonacoEditor
                    v-model:value="message"
                    style="height: 70vh"
                    :options="{
                        automaticLayout: true,
                        minimap: { enabled: false },
                        maxTokenizationLineLength: Infinity,
                    }"
                    theme="vs-dark"
                    @change="sendMessage"
                />
                <ElButton
                    @click="
                        () => {
                            EventsEmit('message', message);
                        }
                    "
                >
                    Send
                </ElButton>
            </div>
        </ElTabPane>
    </ElTabs>
</template>

<style>
.editor {
    height: 70vh !important;
}

@media only screen and (max-width: 1024px) {
    .editor {
        height: 20vh !important;
    }
}
</style>
