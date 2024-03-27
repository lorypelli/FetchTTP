<script setup lang="ts">
import { ElButton, ElForm, ElFormItem, ElInput, ElMessageBox } from 'element-plus';
import { CURL } from '../../wailsjs/go/main/App.js';
import { ref, h } from 'vue';
import { httpTabHandle as HTTPTab } from './Tabs.vue';
defineOptions({
    name: 'CURL',
    components: {
        ElButton,
        ElForm,
        ElFormItem,
        ElInput
    }
});
const curl = ref('');
</script>

<template>
  <ElButton
    @click="() => {
      ElMessageBox({
        title: 'cURL Request',
        message: () => h(ElForm, {
          labelPosition: 'top'
        }, [
          h(ElFormItem, {
            label: 'URL'
          }, [
            h(ElInput, {
              modelValue: curl,
              'onUpdate:modelValue': (val) => {
                curl = val
              },
              type: 'textarea',
              resize: 'none'
            })
          ])
        ])
      })
        .then(() => {
          let url = curl.replace('curl', '')
          if (!url.startsWith('http://') && !url.startsWith('https://')) {
            url = 'https://' + url
          }
          CURL(url).then(() => {
            HTTPTab(undefined, 'add')
          })
        })
        .catch(() => { })
    }"
  >
    Open cURL
  </ElButton>
</template>