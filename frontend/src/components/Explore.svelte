<script lang="ts">
  import type { main } from "@wails/go/models";
  import Header from "./Header.svelte";
  import ContextMenu from "./ContextMenu.svelte";
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
  const isStorage = $derived(isStorageDir(directory.current));

  const isSymlink = (mode: number) => (mode & 0xf000) === 0xa000;
</script>

<Header />
<ContextMenu {selected} />
<!-- TODO: style loging and failed states -->

<section class="explorer">
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
</section>

{#snippet row(entry: main.Entry)}
  {const path = isSymlink(entry.mode) ? removePrefix(entry.id, "1|") : entry.id}
  <!-- NOTE: paths prefix with "1|" points to a symlink or regular file -->

  <tr>
    <td>
      <input type="checkbox" value={path} bind:group={selected} disabled={isStorage} />
    </td>

    {const isFile = !entry.isDir && !isSymlink(entry.mode)}
    <td>
      {#if isFile}
        <div class="file">
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
  :root {
    --explorer-header-height: 2.5em;
    --explorer-header-margin: 1.25em;
    --explorer-min-width: min(48em, 100vw);
    --btn-flex-gap: 0.5em;

    --tr-name: 2;
  }

  :global body:has(div section.explorer) {
    align-content: unset;
  }

  section {
    margin-top: calc(var(--explorer-header-height) + (var(--explorer-header-margin) * 1.5));
    max-height: calc(100vh - var(--explorer-header-height) - (var(--explorer-header-margin) * 2.5));
    min-width: var(--explorer-min-width);
    overflow-y: scroll;
    background: var(--background);
    box-shadow: 0 1px 1px var(--background-a30);
  }

  table {
    table-layout: fixed;
    border-spacing: 0;
    min-width: 100%;
    min-height: 100%;
  }

  th,
  td {
    border: 1px solid var(--background-a30);
  }

  thead {
    --title-width: 0ch;

    th:not(th:nth-child(1)) {
      --btn-padding: 1em;
      min-width: calc(
        var(--title-width) + var(--btn-flex-gap) + var(--svg-size) + var(--btn-padding)
      );
    }

    th:nth-child(2) {
      --title-width: 4ch;
    }

    th:nth-child(3) {
      --title-width: 4ch;
    }

    th:nth-child(4) {
      --title-width: 13ch;
    }

    th:first-of-type {
      border-color: transparent;
      border-bottom-color: var(--background-a30);
    }

    tr {
      background: var(--background);
      position: sticky;
      top: 0;
      z-index: 10;
    }

    button {
      text-transform: capitalize;
      justify-content: center;
    }
  }

  .file,
  button {
    display: flex;
    align-items: center;
    gap: 0.5em;
    padding: 0.5em;
    min-height: 1.5em;
    min-width: 100%;
    background: transparent;
  }

  tbody {
    --td-input-width: 2.5em;

    td:has(input) {
      width: var(--td-input-width);
    }

    td:nth-child(3),
    td:nth-child(4) {
      text-align: center;
      padding-inline: 0.5em;
    }

    input {
      width: 1em;
      height: 1em;
      display: block;
      margin-inline: auto;
      accent-color: var(--success);
    }

    tr {
      --row-cursor: pointer;

      position: relative;
      transform: rotate(0deg);
      transition-property: color, background-color;
      transition-duration: 300ms;

      &:hover,
      &:focus-visible {
        color: var(--accent);
        background: var(--background-a10);
      }

      .file {
        --row-cursor: default;
      }

      button,
      .file {
        &:hover,
        &:focus-visible {
          border-color: transparent;
        }

        &::after {
          content: "";
          position: fixed;
          inset: 0;
          z-index: 5;
          cursor: var(--row-cursor);
          left: var(--td-input-width);
        }
      }
    }
  }
</style>
