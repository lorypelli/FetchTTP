<script lang="ts">
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
    <div class="flex px-8 w-full gap-1 justify-center">
        <select class="w-24">
            {#each http_methods as method}
                <option>{method}</option>
            {/each}
        </select>
        <input class="w-full" />
        <button>Send</button>
    </div>
</Tabs>
