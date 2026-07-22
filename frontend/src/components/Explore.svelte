<script lang="ts">
  import { getEntries, directory } from "@/lib/fs.svelte";
  import Header from "./Header.svelte";
</script>

<Header />

{#await getEntries(directory)}
  <p>Loading...</p>
{:then data}
  {#if Array.isArray(data)}
    {#each data as entry}
      <pre>{JSON.stringify(entry)}</pre>
    {/each}
  {/if}
{:catch}
  <p>failed</p>
{/await}

<style>
</style>
