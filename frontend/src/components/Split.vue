<script setup lang="ts">
import { ElButton, ElCheckbox, ElInput, ElTabPane, ElTabs, ElText, ElDivider, ElForm, ElFormItem, ElDatePicker, ElMessageBox, ElSelect, ElOption, ElEmpty } from 'element-plus';
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
import { reactive, h } from 'vue';
import { EventsEmit } from '../../wailsjs/runtime/runtime';
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Split',
    components: {
        Splitter,
        SplitterPanel,
        ElTabPane,
        ElTabs,
        ElInput,
        ElButton,
        ElCheckbox,
        ElText,
        ElDivider,
        ElForm,
        ElFormItem,
        ElDatePicker,
        ElSelect,
        ElOption,
        ElEmpty
    }
});
const props = defineProps<{
    name: string,
    url?: string,
    status: string,
    header: object,
    response: string,
    type: 'http' | 'ws'
}>();
</script>

<script lang="ts">
export default {
    data() {
        return {
            width: 0,
            headers: [{
                enabled: true,
                name: 'User-Agent',
                value: 'FetchTTP'
            }],
            query: [{
                enabled: false,
                name: '',
                value: ''
            }],
            body: '',
            message: ''
        };
    },
    mounted() {
        let headers = localStorage.getItem(`${this.name}-headers`);
        if (headers) {
            this.headers = JSON.parse(headers);
        }
        else {
            this.headers = [{
                enabled: false,
                name: '',
                value: ''
            }];
        }
        let query = localStorage.getItem(`${this.name}-query`);
        if (query) {
            this.query = JSON.parse(query);
        }
        else {
            this.query = [{
                enabled: false,
                name: '',
                value: ''
            }];
        }
        this.body = localStorage.getItem(`${this.name}-body`) || '';
        this.message = localStorage.getItem(`${this.name}-message`) || '';
        this.update();
        window.addEventListener('resize', this.update);
    },
    beforeUnmount() {
        window.removeEventListener('resize', this.update);
    },
    methods: {
        update() {
            this.width = window.innerWidth;
        },
        add(type: 'header' | 'query', index: number) {
            switch (type) {
            case 'header': {
                for (let i = 0; i < this.headers.length; i++) {
                    if (i == index) {
                        this.headers.splice(i + 1, 0, {
                            enabled: false,
                            name: '',
                            value: ''
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
                            value: ''
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
                                value: ''
                            };
                        }
                        else {
                            this.headers.splice(i, 1);
                        }
                    }
                }
                break;
            }
            case 'query': {
                for (let i = 0; i < this.headers.length; i++) {
                    if (i == index) {
                        if (i == 0) {
                            this.query[0] = {
                                enabled: false,
                                name: '',
                                value: ''
                            };
                        }
                        else {
                            this.query.splice(i, 1);
                        }
                    }
                }
                break;
            }
            }
        },
        sendHeader() {
            localStorage.setItem(`${this.name}-headers`, JSON.stringify(this.headers));
            this.$emit('headers', this.headers);
        },
        sendQuery() {
            localStorage.setItem(`${this.name}-query`, JSON.stringify(this.query));
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
        isText(h: object) {
            return !this.isImage(h) && !this.isVideo(h) && !this.isAudio(h);
        },
        isImage(h: object) {
            const regex = /image\/*/;
            return Object.entries(h).filter(([k, v]) => {
                return k == 'Content-Type' && regex.test(v[0]);
            }).length > 0;
        },
        isVideo(h: object) {
            const regex = /video\/*/;
            return Object.entries(h).filter(([k, v]) => {
                return k == 'Content-Type' && regex.test(v[0]);
            }).length > 0;
        },
        isAudio(h: object) {
            const regex = /audio\/*/;
            return Object.entries(h).filter(([k, v]) => {
                return k == 'Content-Type' && regex.test(v[0]);
            }).length > 0;
        }
    }
};
</script>

<template>
    <Splitter :layout="width < 900 ? 'vertical' : 'horizontal'">
        <SplitterPanel :min-size="20">
            <ElTabs class="pr-2">
                <ElTabPane :label="`Headers (${headers.filter((h) => {
        return h.enabled
    }).length})`">
                    <div class="flex space-x-1 pr-2 pt-2" v-for="(item, index) in headers" :key="index">
                        <ElCheckbox v-model="item.enabled" v-on:change="() => {
        sendHeader()
        if (!item.name || !item.value) {
            item.enabled = false
        }
    }" class="w-10" />
                        <ElInput v-model="item.name" v-on:input="() => {
        sendHeader()
        if (item.name && item.value) {
            item.enabled = true
        }
        else {
            item.enabled = false
        }
    }" />
                        <ElInput v-model="item.value" v-on:input="() => {
        sendHeader()
        if (item.name && item.value) {
            item.enabled = true
        }
        else {
            item.enabled = false
        }
    }" />
                        <ElButton v-on:click="add('header', index)">+</ElButton>
                        <ElButton v-on:click="remove('header', index)">-</ElButton>
                    </div>
                </ElTabPane>
                <ElTabPane :label="`Query (${query.filter((q) => {
        return q.enabled
    }).length})`">
                    <div class="flex space-x-1 pr-2 pt-2" v-for="(item, index) in query" :key="index">
                        <ElCheckbox v-model="item.enabled" v-on:change="() => {
        sendQuery()
        if (!item.name || !item.value) {
            item.enabled = false
        }
    }" class="w-10" />
                        <ElInput v-model="item.name" v-on:input="() => {
        sendQuery()
        if (item.name && item.value) {
            item.enabled = true
        }
        else {
            item.enabled = false
        }
    }" />
                        <ElInput v-model="item.value" v-on:input="() => {
        sendQuery()
        if (item.name && item.value) {
            item.enabled = true
        }
        else {
            item.enabled = false
        }
    }" />
                        <ElButton v-on:click="add('query', index)">+</ElButton>
                        <ElButton v-on:click="remove('query', index)">-</ElButton>
                    </div>
                </ElTabPane>
                <ElTabPane label="Cookies" class="flex justify-center">
                    <ElButton v-on:click="() => {
        const form = reactive({
            key: '',
            value: '',
            domain: '',
            path: '',
            expires: '',
            same_site: 'Lax',
            secure: false,
            http_only: false
        });
        ElMessageBox({
            title: 'New Cookie',
            message: h(ElForm, {
                labelPosition: 'top'
            }, [
                h(ElFormItem, {
                    label: 'Key'
                }, [
                    h(ElInput, {
                        modelValue: form.key,
                        'onUpdate:modelValue': (val) => {
                            form.key = val
                        }
                    })
                ]),
                h(ElFormItem, {
                    label: 'Value'
                }, [
                    h(ElInput, {
                        modelValue: form.value,
                        'onUpdate:modelValue': (val) => {
                            form.value = val
                        }
                    })
                ]),
                h(ElFormItem, {
                    label: 'Domain'
                }, [
                    h(ElInput, {
                        modelValue: form.domain,
                        'onUpdate:modelValue': (val) => {
                            form.domain = val
                        }
                    })
                ]),
                h(ElFormItem, {
                    label: 'Path'
                }, [
                    h(ElInput, {
                        modelValue: form.path,
                        'onUpdate:modelValue': (val) => {
                            form.path = val
                        }
                    })
                ]),
                h(ElFormItem, {
                    label: 'Expires'
                }, [
                    h('div', {
                        class: 'w-full'
                    }, [
                        h(ElDatePicker, {
                            modelValue: form.expires,
                            type: 'datetime',
                            valueFormat: 'x',
                            'onUpdate:modelValue': (val) => {
                                form.expires = val
                            }
                        })
                    ])
                ]),
                h(ElFormItem, {
                    label: 'SameSite'
                }, [
                    h(ElSelect, {
                        modelValue: form.same_site,
                        'onUpdate:modelValue': (val) => {
                            form.same_site = val
                        }
                    }, [
                        h(ElOption, {
                            value: 'None'
                        }),
                        h(ElOption, {
                            value: 'Lax'
                        }),
                        h(ElOption, {
                            value: 'Strict'
                        })
                    ])
                ]),
                h('div', {
                    class: 'flex justify-around'
                }, [
                    h(ElFormItem, {
                        label: 'Secure'
                    }, [
                        h(ElCheckbox, {
                            modelValue: form.secure,
                            'onUpdate:modelValue': (val) => {
                                form.secure = !!val
                            }
                        })
                    ]),
                    h(ElFormItem, {
                        label: 'HTTPOnly'
                    }, [
                        h(ElCheckbox, {
                            modelValue: form.http_only,
                            'onUpdate:modelValue': (val) => {
                                form.http_only = !!val
                            }
                        })
                    ])
                ])
            ]),
            confirmButtonText: 'Add',
            closeOnClickModal: false
        })
            .then(() => {
                headers.push({
                    enabled: true,
                    name: 'Set-Cookie',
                    value: `${form.key}=${form.value};${form.domain ? `Domain=${form.domain};` : ''}${form.path ? `Path=${form.path};` : ''}${form.expires ? `Expires=${new Date(form.expires).toUTCString()};` : ''}${form.same_site ? `SameSite=${form.same_site};` : ''}${form.secure ? 'Secure;' : ''}${form.http_only ? 'HTTPOnly;' : ''}`
                })
            })
            .catch(() => { })
    }">Add Cookie</ElButton>
                </ElTabPane>
                <ElTabPane label="Body" v-if="props.type == 'http'">
                    <ElInput v-model="body" type="textarea" resize="none" />
                </ElTabPane>
                <ElTabPane label="Message" v-if="props.type == 'ws'">
                    <div class="flex flex-col space-y-1">
                        <ElInput v-model="message" type="textarea" resize="none" />
                        <ElButton v-on:click="() => {
        EventsEmit('message', message)
    }">Send</ElButton>
                    </div>
                </ElTabPane>
            </ElTabs>
        </SplitterPanel>
        <SplitterPanel :min-size="20">
            <ElTabs class="pl-2">
                <ElTabPane label="Headers">
                    <ElEmpty
                        :class="`${Object.keys(props.header).length == 0 ? 'flex' : 'hidden'} justify-center h-scrollbar-partial`"
                        description="Nothing to display here..." />
                    <ElText class="flex justify-center">{{ props.status }}</ElText>
                    <ElDivider v-if="props.status" />
                    <ElText v-if="Object.keys(props.header).length > 0">{{ props.header }}</ElText>
                </ElTabPane>
                <ElTabPane label="Response">
                    <ElEmpty
                        :class="`${['', 'null'].includes(props.response.trim()) ? 'flex' : 'hidden'} justify-center h-scrollbar-partial`"
                        description="Nothing to display here..." />
                    <ElText v-if="isText(props.header) && !['', 'null'].includes(props.response.trim())">{{
        props.response }}</ElText>
                    <div
                        :class="`${isText(props.header) ? 'hidden' : 'flex'} justify-center items-center h-scrollbar-full`">
                        <img v-if="isImage(props.header)" :src="props.url" />
                        <audio v-if="isAudio(props.header)" :src="props.url" controls></audio>
                        <video v-if="isVideo(props.header)" :src="props.url" controls></video>
                    </div>
                </ElTabPane>
            </ElTabs>
        </SplitterPanel>
    </Splitter>
</template>