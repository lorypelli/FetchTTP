<script setup lang="ts">
defineOptions({
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'Split',
});
defineProps<{
    type: 'http' | 'ws'
}>();
</script>