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

<div class="flex justify-between">
    <div class="flex flex-col">
        <button on:click={() => enable('http')}>HTTP</button>
        <div class="{isHTTP ? 'flex' : 'hidden'} flex-col">
            <HTTP />
        </div>
    </div>
    <div class="flex flex-col">
        <button on:click={() => enable('ws')}>WS</button>
        <div class="{isWS ? 'flex' : 'hidden'} flex-col">
            <WS />
        </div>
    </div>
    <div class="flex flex-col">
        <button on:click={() => enable('curl')}>CURL</button>
        <div class="{isCURL ? 'flex' : 'hidden'} flex-col">
            <CURL />
        </div>
    </div>
</div>
