<script lang="ts">
    export let type: 'http' | 'ws';
    let tabs = [1];
</script>

<div class="flex">
    {#each tabs as t}
        <div>
            {t}
            <slot />
        </div>
    {/each}
</div>
<button
    class="mt-0.5"
    on:click={() => (tabs = [...tabs, tabs.length + 1])}>+</button>
