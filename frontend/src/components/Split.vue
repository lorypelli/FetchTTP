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
interface Header {
    [x: string]: string[]
}
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
export default {
    data() {
        return {
            width: 0,
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