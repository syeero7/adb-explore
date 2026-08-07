<script lang="ts">
  import Logs from "./Logs.svelte";
  import { directory, isStorageDir, refresh, toParentDir, toStorageDir } from "@/lib/fs.svelte";
  import { svg, RELOAD, UP_ARROW, STORAGE, SEARCH, CLOSE } from "@/lib/svg";

  let isSearching = $state(false);
  const isStorage = $derived(isStorageDir(directory.current));
  let timeout: number | undefined;

  const openSearch = async () => {
    isSearching = true;
  };

  const refreshDir = () => {
    clearTimeout(timeout);
    refresh();
  };

  const closeSearch = () => {
    clearTimeout(timeout);
    isSearching = false;
    directory.query = "";
  };

  const search = (e: Event) => {
    clearTimeout(timeout);
    const query = (e.target as HTMLInputElement).value.trim();
    timeout = setTimeout(() => {
      directory.query = query;
    }, 300);
  };
</script>

<header>
  <button title="refresh current directory" onclick={refreshDir}>
    {@render svg({ d: RELOAD })}
  </button>

  <button title="go to parent directory" onclick={toParentDir} disabled={isStorage}>
    {@render svg({ d: UP_ARROW })}
  </button>

  <button title="go to storage directory" onclick={toStorageDir} disabled={isStorage}>
    {@render svg({ d: STORAGE })}
  </button>

  {#if isSearching}
    <!-- svelte-ignore a11y_autofocus -->
    <input type="text" autofocus oninput={search} aria-label="search query" />
    <button title="close" onclick={closeSearch}>
      {@render svg({ d: CLOSE })}
    </button>
  {:else}
    <input
      type="text"
      bind:value={directory.current}
      readonly
      aria-label="current directory path" />
    <button title="search entry" onclick={openSearch}>
      {@render svg({ d: SEARCH })}
    </button>
  {/if}

  <Logs />
</header>

<style>
  header {
    position: fixed;
    top: 0;
    display: flex;
    align-items: center;
    gap: 0.5em;
    background: var(--background);
    min-width: var(--explorer-min-width);
    min-height: var(--explorer-header-height);
    margin-top: var(--explorer-header-margin);

    input {
      flex: 1;
      min-width: 25ch;
      cursor: text;
      outline: none;
      box-shadow: unset;
      max-height: var(--btn-size);
    }

    button {
      flex-shrink: 0;
      height: var(--btn-size);
      width: var(--btn-size);
      box-shadow: var(--shadow);
      border-radius: var(--radius);
    }
  }
</style>
