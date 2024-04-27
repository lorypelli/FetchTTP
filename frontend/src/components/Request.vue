<script setup lang="ts">
import {
    ElButton,
    ElCheckbox,
    ElInput,
    ElTabPane,
    ElTabs,
    ElForm,
    ElFormItem,
    ElDatePicker,
    ElMessageBox,
    ElSelect,
    ElOption,
} from 'element-plus';
import { reactive, h } from 'vue';
import { EventsEmit } from '../../wailsjs/runtime/runtime';
defineOptions({
    name: 'Request',
    components: {
        ElButton,
        ElCheckbox,
        ElInput,
        ElTabPane,
        ElTabs,
        ElForm,
        ElFormItem,
        ElDatePicker,
        ElSelect,
        ElOption,
    },
});
const props = defineProps<{
    name: string;
    type: 'http' | 'ws';
}>();
</script>

<script lang="ts">
export default {
    emits: {
        query: null,
        body: null,
        headers: null,
        message: null,
    },
    data() {
        return {
            headers: [
                {
                    enabled: true,
                    name: 'User-Agent',
                    value: 'FetchTTP',
                },
            ],
            query: [
                {
                    enabled: false,
                    name: '',
                    value: '',
                },
            ],
            body: '',
            message: '',
        };
    },
    mounted() {
        let headers = localStorage.getItem(`${this.name}-headers-${this.type}`);
        if (headers) {
            this.headers = JSON.parse(headers);
        } else {
            this.headers = [
                {
                    enabled: true,
                    name: 'User-Agent',
                    value: 'FetchTTP',
                },
            ];
        }
        let query = localStorage.getItem(`${this.name}-query-${this.type}`);
        if (query) {
            this.query = JSON.parse(query);
        } else {
            this.query = [
                {
                    enabled: false,
                    name: '',
                    value: '',
                },
            ];
        }
        this.body = localStorage.getItem(`${this.name}-body`) || '';
        this.message = localStorage.getItem(`${this.name}-message`) || '';
    },
    methods: {
        add(type: 'header' | 'query', index: number) {
            switch (type) {
                case 'header': {
                    for (let i = 0; i < this.headers.length; i++) {
                        if (i == index) {
                            this.headers.splice(i + 1, 0, {
                                enabled: false,
                                name: '',
                                value: '',
                            });
                        }
                    }
                    break;
                }
                case 'query': {
                    for (let i = 0; i < this.query.length; i++) {
                        if (i == index) {
                            this.query.splice(i + 1, 0, {
                                enabled: false,
                                name: '',
                                value: '',
                            });
                        }
                    }
                    break;
                }
            }
        },
        remove(type: 'header' | 'query', index: number) {
            switch (type) {
                case 'header': {
                    for (let i = 0; i < this.headers.length; i++) {
                        if (i == index) {
                            if (i == 0) {
                                this.headers[0] = {
                                    enabled: false,
                                    name: '',
                                    value: '',
                                };
                                localStorage.removeItem(
                                    `${this.name}-headers-${this.type}`,
                                );
                            } else {
                                this.headers.splice(i, 1);
                                localStorage.setItem(
                                    `${this.name}-headers-${this.type}`,
                                    JSON.stringify(this.headers),
                                );
                            }
                        }
                    }
                    break;
                }
                case 'query': {
                    for (let i = 0; i < this.query.length; i++) {
                        if (i == index) {
                            if (i == 0) {
                                this.query[0] = {
                                    enabled: false,
                                    name: '',
                                    value: '',
                                };
                                localStorage.removeItem(
                                    `${this.name}-query-${this.type}`,
                                );
                            } else {
                                this.query.splice(i, 1);
                                localStorage.setItem(
                                    `${this.name}-query-${this.type}`,
                                    JSON.stringify(this.query),
                                );
                            }
                        }
                    }
                    break;
                }
            }
        },
        sendHeader() {
            localStorage.setItem(
                `${this.name}-headers-${this.type}`,
                JSON.stringify(this.headers),
            );
            this.$emit('headers', this.headers);
        },
        sendQuery() {
            localStorage.setItem(
                `${this.name}-query-${this.type}`,
                JSON.stringify(this.query),
            );
            this.$emit('query', this.query);
        },
        sendBody() {
            localStorage.setItem(`${this.name}-body`, this.body);
            this.$emit('body', this.body);
        },
        sendMessage() {
            localStorage.setItem(`${this.name}-message`, this.message);
            this.$emit('message', this.message);
        },
    },
};
</script>

<template>
    <ElTabs class="pr-2">
        <ElTabPane :label="`Headers (${headers.filter((h) => {
            return h.enabled;
        }).length
            })`">
            <div v-for="(item, index) in headers" :key="index" class="flex space-x-1 pr-2 pt-2">
                <ElCheckbox v-model="item.enabled" class="w-10" @change="() => {
                    if (!item.name || !item.value) {
                        item.enabled = false;
                    }
                    sendHeader();
                }
                    " />
                <ElInput v-model="item.name" @input="() => {
                    if (item.name && item.value) {
                        item.enabled = true;
                    } else {
                        item.enabled = false;
                    }
                    sendHeader();
                }
                    " />
                <ElInput v-model="item.value" @input="() => {
                    if (item.name && item.value) {
                        item.enabled = true;
                    } else {
                        item.enabled = false;
                    }
                    sendHeader();
                }
                    " />
                <ElButton @click="add('header', index)"> + </ElButton>
                <ElButton @click="remove('header', index)"> - </ElButton>
            </div>
        </ElTabPane>
        <ElTabPane :label="`Query (${query.filter((q) => {
            return q.enabled;
        }).length
            })`">
            <div v-for="(item, index) in query" :key="index" class="flex space-x-1 pr-2 pt-2">
                <ElCheckbox v-model="item.enabled" class="w-10" @change="() => {
                    if (!item.name || !item.value) {
                        item.enabled = false;
                    }
                    sendQuery();
                }
                    " />
                <ElInput v-model="item.name" @input="() => {
                    if (item.name && item.value) {
                        item.enabled = true;
                    } else {
                        item.enabled = false;
                    }
                    sendQuery();
                }
                    " />
                <ElInput v-model="item.value" @input="() => {
                    if (item.name && item.value) {
                        item.enabled = true;
                    } else {
                        item.enabled = false;
                    }
                    sendQuery();
                }
                    " />
                <ElButton @click="add('query', index)"> + </ElButton>
                <ElButton @click="remove('query', index)"> - </ElButton>
            </div>
        </ElTabPane>
        <ElTabPane label="Cookies" class="flex justify-center items-center h-full">
            <ElButton class="scale-150" @click="() => {
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
                    .catch(() => { });
            }
                ">
                Add Cookie
            </ElButton>
        </ElTabPane>
        <ElTabPane v-if="props.type == 'http'" label="Body">
            <ElInput v-model="body" type="textarea" resize="none" @input="sendBody" />
        </ElTabPane>
        <ElTabPane v-if="props.type == 'ws'" label="Message">
            <div class="flex flex-col space-y-1">
                <ElInput v-model="message" type="textarea" resize="none" @input="sendMessage" />
                <ElButton @click="() => {
                    EventsEmit('message', message);
                }
                    ">
                    Send
                </ElButton>
            </div>
        </ElTabPane>
    </ElTabs>
</template>
