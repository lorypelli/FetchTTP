<script lang="ts">
    import { Tabs, Button } from '@svelteuidev/core';
    export let type: 'http' | 'ws';
    let tabs = [
        1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
        21,
    ];
</script>

{#key tabs}
    <Tabs orientation="vertical">
        {#each tabs as t}
            <Tabs.Tab label={t}><slot /></Tabs.Tab>
        {/each}
    </Tabs>
{/key}
<Button
    class="mt-2"
    on:click={() => (tabs = [...tabs, tabs.length + 1])}
    variant="outline"
    color="gray">+</Button
>
