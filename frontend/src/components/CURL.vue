<script setup lang="ts">
import { ElButton, ElForm, ElFormItem, ElInput, ElMessageBox, ElNotification } from 'element-plus';
import { CURL } from '../../wailsjs/go/main/App.js';
import { ref, h } from 'vue';
import { default as R } from './Response.vue';
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

<script lang="ts">
interface Response {
    URL: string,
    Status: string,
    Header: [],
    Body: string
    Error: string
}
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
          try {
            CURL(url).then((res: Response) => {
              if (res.Error) {
                ElNotification({
                  title: 'Something went wrong!',
                  message: res.Error,
                  type: 'error',
                  position: 'bottom-right'
                });
                return;
              }
              ElMessageBox({
                title: 'Response',
                message: () => h(R, {
                  url: res.URL,
                  status: res.Status,
                  header: res.Header,
                  response: res.Body
                })
              })
                .catch(() => { })
            })
          }
          catch {
            ElNotification({
              title: 'Something went wrong!',
              type: 'error',
              position: 'bottom-right'
            });  
          }
        })
        .catch(() => { })
    }"
  >
    Open cURL
  </ElButton>
</template>