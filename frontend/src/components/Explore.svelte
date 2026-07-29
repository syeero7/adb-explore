<script lang="ts">
  import type { main } from "@wails/go/models";
  import Header from "./Header.svelte";
  import {
    type InfoTitle,
    getEntries,
    directory,
    isStorageDir,
    sortBy,
    toDir,
    removePrefix,
  } from "@/lib/fs.svelte";
  import { svg, DOWN_ARROW, UP_ARROW, FILE, FOLDER } from "@/lib/svg";

  let selected: string[] = $state([]);

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
  {const isStorage = isStorageDir(directory.current)}
  {const path = isSymlink(entry.mode) ? removePrefix(entry.id, "1|") : entry.id}
  <!-- NOTE: paths prefix with "1|" points to a symlink or regular file -->

  <tr>
    <td>
      <input type="checkbox" value={path} bind:group={selected} disabled={isStorage} />
    </td>

    {const isFile = !entry.isDir && !isSymlink(entry.mode)}
    <td>
      {#if isFile}
        <div>
          {@render svg({ d: FILE, fileExt: entry.ext.length < 5 ? entry.ext : undefined })}
          <span>
            {entry.name}
          </span>
        </div>
      {:else}
        <button ondblclick={() => toDir(path)}>
          {@render svg({ d: FOLDER })}
          <span>
            {entry.name}
          </span>
        </button>
      {/if}
    </td>

    <td>{isFile ? entry.size : ""}</td>

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
