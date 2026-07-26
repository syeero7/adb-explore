<script lang="ts">
  import type { main } from "@wails/go/models";
  import Header from "./Header.svelte";
  import { type InfoTitle, getEntries, directory, useIsStorageDir, sortBy } from "@/lib/fs.svelte";
  import { svg, DOWN_ARROW, UP_ARROW } from "@/lib/svg";

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
          {@render th("name")}
          {@render th("size")}
          {@render th("date modified")}
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

{#snippet th(title: InfoTitle)}
  {const { isActive, isAsc, handler } = sortBy(title, directory.sortBy)}
  <th>
    <button onclick={handler}>
      <span>
        {title}
      </span>
      {#if isActive}
        {#if isAsc}
          {@render svg({ d: UP_ARROW })}
        {:else}
          {@render svg({ d: DOWN_ARROW })}
        {/if}
      {/if}
    </button>
  </th>
{/snippet}

<style>
</style>
