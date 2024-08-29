<script lang="ts">
    import { Tabs } from '@svelteuidev/core';
    import HTTP from '../components/HTTP.svelte';
    import WS from '../components/WS.svelte';
    import CURL from '../components/CURL.svelte';
</script>

<Tabs>
    <Tabs.Tab label="HTTP">
        <HTTP />
    </Tabs.Tab>
    <Tabs.Tab label="WS">
        <WS />
    </Tabs.Tab>
    <Tabs.Tab label="CURL">
        <CURL />
    </Tabs.Tab>
</Tabs>