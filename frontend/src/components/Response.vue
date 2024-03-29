<script setup lang="ts">
import { ElTabPane, ElTabs, ElEmpty, ElText, ElDivider } from 'element-plus';
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Response',
    components: {
        ElTabPane,
        ElTabs,
        ElEmpty,
        ElText,
        ElDivider
    }
});
interface Header {
    [x: string]: string[]
}
const props = defineProps<{
    url?: string,
    status: string,
    header: Header,
    response: string
}>();
</script>

<script lang="ts">
export default {
    methods: {
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
  <ElTabs class="pl-2">
    <ElTabPane label="Headers">
      <ElEmpty
        :class="`${Object.keys(props.header).length == 0 ? 'flex' : 'hidden'} justify-center h-full`"
        description="Nothing to display here..."
      />
      <ElText class="flex justify-center">
        {{ props.status }}
      </ElText>
      <ElDivider v-if="props.status" />
      <table v-if="Object.keys(props.header).length > 0">
        <tbody>
          <tr
            v-for="(item, index) in Object.keys(props.header)"
            :key="index"
          >
            <th class="break-all w-1/2">
              <ElText>{{ item }}</ElText>
            </th>
            <td class="break-all w-1/2">
              <ElText>{{ props.header[item].join("") }}</ElText>
            </td>
          </tr>
        </tbody>
      </table>
    </ElTabPane>
    <ElTabPane label="Response">
      <ElEmpty
        :class="`${['', 'null'].includes(props.response.trim()) ? 'flex' : 'hidden'} justify-center h-full`"
        description="Nothing to display here..."
      />
      <ElText v-if="isText(props.header) && !['', 'null'].includes(props.response.trim())">
        {{
          props.response }}
      </ElText>
      <div :class="`${isText(props.header) ? 'hidden' : 'flex'} justify-center items-center h-full`">
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
      </div>
    </ElTabPane>
  </ElTabs>
</template>