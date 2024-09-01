<script lang="ts">
    export let type: 'http' | 'ws';
    let tabs = [1];
</script>

{#each tabs as t}
    <div class="flex">
        <div class="grid grid-cols-1">
            <button>{t}</button>
        </div>
        <slot />
    </div>
{/each}
<button
    class="fixed bottom-0 mt-0.5"
    on:click={() => (tabs = [...tabs, tabs.length + 1])}>+</button
>
