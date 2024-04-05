<script setup lang="ts">
import { ElTabPane, ElTabs, ElEmpty, ElText, ElDivider, ElSwitch } from 'element-plus';
import { ref } from 'vue';
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Response',
    components: {
        ElTabPane,
        ElTabs,
        ElEmpty,
        ElText,
        ElDivider,
        ElSwitch
    }
});
const props = defineProps<{
  url?: string,
  status: string,
  header: Header,
  response: string
}>();
const readable = ref(true);
</script>

<script lang="ts">
interface Header {
  [x: string]: string[]
}
export default {
    methods: {
        isText(h: Header) {
            return !this.isPage(h) && !this.isImage(h) && !this.isVideo(h) && !this.isAudio(h);
        },
        isPage(h: Header) {
            return Object.entries(h).filter(([k, v]) => {
                return k == 'Content-Type' && v[0].includes('text/html');
            }).length > 0;
        },
        isImage(h: Header) {
            const regex = /image\/*/;
            return Object.entries(h).filter(([k, v]) => {
                return k == 'Content-Type' && regex.test(v[0]);
            }).length > 0;
        },
        isVideo(h: Header) {
            const regex = /video\/*/;
            return Object.entries(h).filter(([k, v]) => {
                return k == 'Content-Type' && regex.test(v[0]);
            }).length > 0;
        },
        isAudio(h: Header) {
            const regex = /audio\/*/;
            return Object.entries(h).filter(([k, v]) => {
                return k == 'Content-Type' && regex.test(v[0]);
            }).length > 0;
        },
        baseURL(r: string, u: string | undefined) {
            if (u) {
                return r.replace('<head>', `<head><base href="${u}">`);
            }
            return r;
        }
    }
};
</script>

<template>
  <ElTabs class="pl-2">
    <ElTabPane label="Headers">
      <ElEmpty
        v-if="Object.keys(props.header).length == 0"
        class="flex justify-center h-full"
        description="Nothing to display here..."
      />
      <ElText class="flex justify-center">
        {{ props.status }}
      </ElText>
      <ElDivider v-if="props.status" />
      <ElText v-if="Object.keys(props.header).length > 0 && !readable">
        {{ props.header }}
      </ElText>
      <table v-if="Object.keys(props.header).length > 0 && readable">
        <tr
          v-for="(item, index) in Object.keys(props.header)"
          :key="index"
        >
          <td class="break-all w-1/2 pl-5">
            <ElText>{{ item }}</ElText>
          </td>
          <td class="break-all w-1/2 pr-5">
            <ElText>{{ props.header[item].join("") }}</ElText>
          </td>
        </tr>
      </table>
    </ElTabPane>
    <ElTabPane label="Response">
      <ElEmpty
        v-if="['', 'null'].includes(props.response.trim())"
        class="flex justify-center h-full"
        description="Nothing to display here..."
      />
      <ElText v-if="isText(props.header) && !['', 'null'].includes(props.response.trim()) || !readable">
        {{
          props.response }}
      </ElText>
      <div
        v-if="!isText(props.header) && readable"
        class="flex justify-center items-center h-full"
      >
        <img
          v-if="isImage(props.header)"
          :src="props.url"
        >
        <audio
          v-if="isAudio(props.header)"
          :src="props.url"
          controls
        />
        <video
          v-if="isVideo(props.header)"
          :src="props.url"
          controls
        />
        <iframe
          v-if="isPage(props.header)"
          :srcdoc="baseURL(props.response, props.url)"
          class="w-full h-full rounded-2xl"
          sandbox="allow-scripts allow-forms"
        />
      </div>
    </ElTabPane>
  </ElTabs>
  <ElSwitch
    v-model="readable"
    class="pl-2"
    inactive-text="Raw"
    active-text="Human Readable"
  />
</template>