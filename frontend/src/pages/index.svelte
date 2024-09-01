<script lang="ts">
    import HTTP from '../components/HTTP.svelte';
    import WS from '../components/WS.svelte';
    import CURL from '../components/CURL.svelte';
    let isHTTP = true;
    let isWS = false;
    let isCURL = false;
    function enable(value: 'http' | 'ws' | 'curl') {
        switch (value) {
            case 'http': {
                isHTTP = true;
                isWS = false;
                isCURL = false;
                break;
            }
            case 'ws': {
                isWS = true;
                isHTTP = false;
                isCURL = false;
                break;
            }
            case 'curl': {
                isCURL = true;
                isHTTP = false;
                isWS = false;
                break;
            }
        }
    }
</script>

<div class="grid grid-cols-3 gap-x-2">
    <div class="flex flex-col">
        <button class="border-2 rounded-md" on:click={() => enable('http')}>HTTP</button>
        <div class="{isHTTP ? 'flex' : 'hidden'} flex-col w-full absolute left-0 top-8">
            <HTTP />
        </div>
    </div>
    <div class="flex flex-col">
        <button class="border-2 rounded-md" on:click={() => enable('ws')}>WS</button>
        <div class="{isWS ? 'flex' : 'hidden'} flex-col w-full absolute left-0 top-8">
            <WS />
        </div>
    </div>
    <div class="flex flex-col">
        <button class="border-2 rounded-md" on:click={() => enable('curl')}>CURL</button>
        <div class="{isCURL ? 'flex' : 'hidden'} flex-col w-full absolute left-0 top-8">
            <CURL />
        </div>
    </div>
</div>
