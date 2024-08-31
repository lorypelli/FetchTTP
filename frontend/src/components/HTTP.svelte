<script lang="ts">
    import { NativeSelect } from '@svelteuidev/core';
    import Tabs from './Tabs.svelte';
    const http_methods = [
        'GET',
        'HEAD',
        'POST',
        'PUT',
        'DELETE',
        'CONNECT',
        'OPTIONS',
        'TRACE',
        'PATH',
    ];
</script>

<Tabs type="http">
    <NativeSelect class="w-20" data={http_methods} />
</Tabs>
