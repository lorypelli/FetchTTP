<script setup lang="ts">
import { ElButton, ElCheckbox, ElInput, ElTabPane, ElTabs, ElText, ElDivider, ElScrollbar, ElForm, ElFormItem, ElDatePicker, ElMessageBox, ElSelect, ElOption, ElEmpty } from 'element-plus';
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
import { reactive, h } from 'vue';
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
        addHeader(index: number) {
            for (let i = 0; i < this.headers.length; i++) {
                if (i == index) {
                    if (!this.headers[i + 1]) {
                        this.headers.push({
                            enabled: false,
                            name: '',
                            value: ''
                        });
                    }
                    else {
                        this.headers.push({
                            enabled: this.headers[i + 1].enabled,
                            name: this.headers[i + 1].name,
                            value: this.headers[i + 1].value
                        });
                        this.headers[i + 1] = {
                            enabled: false,
                            name: '',
                            value: ''
                        };
                    }
                }
            }
        },
        removeHeader(index: number) {
            for (let i = this.headers.length; i >= 0; i--) {
                if (i == index) {
                    if (!this.headers[i]) {
                        this.headers[i] = {
                            enabled: false,
                            name: '',
                            value: ''
                        };
                    }
                    else {
                        this.headers[i] = {
                            enabled: false,
                            name: '',
                            value: ''
                        };
                        if (this.headers.length > 1) {
                            this.headers.pop();
                        }
                    }
                }
            }
        },
        addQuery(index: number) {
            for (let i = 0; i < this.query.length; i++) {
                if (i == index) {
                    if (!this.query[i + 1]) {
                        this.query.push({
                            enabled: false,
                            name: '',
                            value: ''
                        });
                    }
                    else {
                        this.query.push({
                            enabled: this.query[i + 1].enabled,
                            name: this.query[i + 1].name,
                            value: this.query[i + 1].value
                        });
                        this.query[i + 1] = {
                            enabled: false,
                            name: '',
                            value: ''
                        };
                    }
                }
            }
        },
        removeQuery(index: number) {
            for (let i = this.query.length; i >= 0; i--) {
                if (i == index) {
                    if (!this.query[i]) {
                        this.query[i] = {
                            enabled: false,
                            name: '',
                            value: ''
                        };
                    }
                    else {
                        this.query[i] = {
                            enabled: false,
                            name: '',
                            value: ''
                        };
                        if (this.query.length > 1) {
                            this.query.pop();
                        }
                    }
                }
            }
        },
        sendHeader() {
            this.$emit('headers', this.headers);
        },
        sendQuery() {
            this.$emit('query', this.query);
        },
        sendBody() {
            this.$emit('body', this.body);
        },
        sendMessage() {
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
                    <ElScrollbar height="83.5vh">
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
                            <ElButton v-on:click="addHeader(index)">+</ElButton>
                            <ElButton v-on:click="removeHeader(index)">-</ElButton>
                        </div>
                    </ElScrollbar>
                </ElTabPane>
                <ElTabPane :label="`Query (${query.filter((q) => {
        return q.enabled
    }).length})`">
                    <ElScrollbar height="83.5vh">
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
                            <ElButton v-on:click="addQuery(index)">+</ElButton>
                            <ElButton v-on:click="removeQuery(index)">-</ElButton>
                        </div>
                    </ElScrollbar>
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
            message: () => h(ElForm, {
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
                    value: ''
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
                        <ElButton>Send</ElButton>
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
                    <ElScrollbar class="h-scrollbar-partial">
                        <ElText v-if="Object.keys(props.header).length > 0">{{ props.header }}</ElText>
                    </ElScrollbar>
                </ElTabPane>
                <ElTabPane label="Response">
                    <ElEmpty
                        :class="`${Object.keys(props.header).length == 0 ? 'flex' : 'hidden'} justify-center h-scrollbar-partial`"
                        description="Nothing to display here..." />
                    <ElScrollbar class="h-scrollbar-full">
                        <ElText v-if="isText(props.header)">{{ props.response }}</ElText>
                        <div
                            :class="`${isText(props.header) ? 'hidden' : 'flex'} justify-center items-center h-scrollbar-full`">
                            <img v-if="isImage(props.header)" :src="props.url" />
                            <audio v-if="isAudio(props.header)" :src="props.url" controls></audio>
                            <video v-if="isVideo(props.header)" :src="props.url" controls></video>
                        </div>
                    </ElScrollbar>
                </ElTabPane>
            </ElTabs>
        </SplitterPanel>
    </Splitter>
</template>