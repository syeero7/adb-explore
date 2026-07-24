<script lang="ts">
  import type { main } from "@wails/go/models";
  import { getEntries, directory, useIsStorageDir } from "@/lib/fs.svelte";
  import Header from "./Header.svelte";

  let selected: string[] = $state([]);
  const isStorageDir = useIsStorageDir().value;

  const isSymlink = (mode: number) => (mode & 0xf000) === 0xa000;
</script>

<Header />

{#await getEntries(directory)}
  <p>Loading...</p>
{:then data}
  {#if Array.isArray(data)}
    <table>
      <thead>
        <tr>
          <th></th>
          <th>Name</th>
          <th>Size</th>
          <th>Date Modified</th>
        </tr>
      </thead>

      <tbody>
        {#each data as entry}
          {@render row(entry)}
        {/each}
      </tbody>
    </table>
  {/if}
{:catch}
  <p>Failed</p>
{/await}

{#snippet row(entry: main.Entry)}
  <tr>
    <td>
      <input type="checkbox" value={entry.path} bind:group={selected} disabled={isStorageDir} />
    </td>
    <td>{entry.name}</td>
    <td>{entry.isDir || isSymlink(entry.mode) ? "" : entry.size}</td>
    {const modified = new Date(entry.lastModified)}
    <td>{modified.toLocaleDateString()}</td>
  </tr>
{/snippet}

<style>
</style>
