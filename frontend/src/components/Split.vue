<script setup lang="ts">
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
import Request from './Request.vue';
import Response from './Response.vue';
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Split',
    components: {
        Splitter,
        SplitterPanel,
        Request,
        Response
    }
});
const props = defineProps<{
    name: string,
    url?: string,
    status: string,
    header: Header,
    response: string,
    type: 'http' | 'ws'
}>();
</script>

<script lang="ts">
interface Header {
    [x: string]: string[]
}
interface Query {
    enabled: boolean,
    name: string,
    value: string
}
export default {
    emits: {
        'query': null,
        'body': null,
        'headers': null,
        'message': null
    },
    data() {
        return {
            width: 0
        };
    },
    mounted() {
        let headers = localStorage.getItem(`${this.name}-headers-${this.type}`);
        if (headers) {
            this.handleHeader(JSON.parse(headers));
        }
        let query = localStorage.getItem(`${this.name}-query-${this.type}`);
        if (query) {
            this.handleQuery(JSON.parse(query));
        }
        this.handleBody(localStorage.getItem(`${this.name}-body`) || '');
        this.handleMessage(localStorage.getItem(`${this.name}-message`) || '');
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
        handleHeader(h: Header) {
            this.$emit('headers', h);
        },
        handleQuery(q: Query) {
            this.$emit('query', q);
        },
        handleBody(b: string) {
            this.$emit('body', b);
        },
        handleMessage(m: string) {
            this.$emit('message', m);
        }
    }
};
</script>

<template>
  <Splitter :layout="width < 900 ? 'vertical' : 'horizontal'">
    <SplitterPanel :min-size="20">
      <Request
        :name="props.name"
        :type="props.type"
        @headers="handleHeader"
        @query="handleQuery"
        @body="handleBody"
        @message="handleMessage"
      />
    </SplitterPanel>
    <SplitterPanel :min-size="20">
      <Response
        :url="props.url"
        :status="props.status"
        :header="props.header"
        :response="props.response"
      />
    </SplitterPanel>
  </Splitter>
</template>